package records

import "bytes"

type TrailProvider struct {
	chapters []ChapterInfo
}

func NewTrailProvider(chapters []ChapterInfo) TrailProvider {
	return TrailProvider{
		chapters: chapters,
	}
}

func (tp *TrailProvider) Get(from, to int) TrailingData {
	strands := TrailingData{Multibyte: 0}
	for i, chap := range tp.chapters {
		end := chap.Start + chap.Length
		if chap.Start <= from && end >= to {
			atExactBoundary := chap.Start == from || end == to
			strands.Strands = &StrandData{
				Index:        i,
				FlagTBSType:  8,
				FlagDoesSpan: !atExactBoundary,
			}
			break
		}

		if chap.Start >= from && chap.Start <= to || end >= from && end <= to {
			if strands.Strands == nil {
				strands.Strands = &StrandData{
					Index:        i,
					FlagTBSType:  8,
					FlagDoesSpan: false,
				}
			}
			strands.Strands.FlagNumSiblings++
		}
	}

	return strands
}

// TrailingData represents is the trailing entries that are appended
// to every text record as indicated by the extra data bitflags in the
// MOBI header. This implementation only supports flags 0b11 meaning
// entries for multibyte overlap and indexing data.
type TrailingData struct {
	Multibyte byte
	Strands   *StrandData
}

// StrandData is the indexing data that represents one hierarchy of
// chapters and sub-chapters in the trailing byte sequence of a text
// record. This data would normally be separated into multiple strands
// and sequences but this implementation does not support sub-chapters
// meaning a strand must contain exactly one sequence and the indexing
// data must consist of either one or zero strands.
type StrandData SequenceData

// SequenceData is the indexing data that represents one chapter of an
// arbitrary level in the trailing byte sequence of a text record. If
// the chapter and its siblings do not have sub-chapters, its siblings
// are also combined into the sequence.
type SequenceData struct {
	Index                     int
	FlagFirstOfNotFirstStrand bool
	FlagTBSType               int
	FlagNumSiblings           byte
	FlagDoesSpan              bool
}

func (td TrailingData) Encode() []byte {
	b := bytes.NewBuffer([]byte{td.Multibyte})
	if td.Strands != nil {
		b.Write(td.Strands.Encode())
	}

	return encodeTrailingBytes(b.Bytes())
}

func (sd StrandData) Encode() []byte {
	value := sd.Index << 3
	if sd.FlagDoesSpan {
		value |= 0b0001
	}
	if sd.FlagTBSType != 0 {
		value |= 0b0010
	}
	if sd.FlagNumSiblings > 1 {
		value |= 0b0100
	}
	if sd.FlagFirstOfNotFirstStrand {
		value |= 0b1000
	}

	b := bytes.NewBuffer(nil)
	b.Write(encodeVwi(value))
	if sd.FlagTBSType != 0 {
		b.Write(encodeVwi(sd.FlagTBSType))
	}
	if sd.FlagNumSiblings > 1 {
		b.WriteByte(sd.FlagNumSiblings)
	}
	if sd.FlagDoesSpan {
		b.Write(encodeVwi(0))
	}

	return b.Bytes()
}

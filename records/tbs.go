package records

import "bytes"

type StrandProvider struct {
	chapters []ChapterInfo
}

func NewStrandProvider(chapters []ChapterInfo) StrandProvider {
	return StrandProvider{
		chapters: chapters,
	}
}

func (s *StrandProvider) Get(from, to int) StrandsData {
	strands := StrandsData{StrandData: nil}
	for i, chap := range s.chapters {
		end := chap.Start + chap.Length
		if chap.Start < from && end > to {
			strands.StrandData = &StrandData{
				Index:        i,
				FlagTBSType:  8,
				FlagDoesSpan: true,
			}
			break
		}

		if chap.Start >= from && chap.Start <= to || end >= from && end <= to {
			if strands.StrandData == nil {
				strands.StrandData = &StrandData{
					Index:        i,
					FlagTBSType:  8,
					FlagDoesSpan: false,
				}
			}
			strands.FlagNumSiblings++
		}
	}

	return strands
}

// StrandsData is the indexing data that represents all chapters of
// all levels in the trailing byte sequence of a text record.  This
// data would normally be separated into strands, which would further
// be separated into sequences. However, this implementation does not
// support sub-chapters meaning the trailing byte sequence may contain
// only no strands or one strand which may contain only one sequence.
type StrandsData struct {
	*StrandData
}

// StrandData is the indexing data that represents one hierarchy of
// chapters and sub-chapters in the trailing byte sequence of a text
// record. This data would normally be separated into sequences.
// However, this implementation does not support sub-chapters meaning
// a strand may contain only one sequence.
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

func (td StrandsData) Encode() []byte {
	if td.StrandData == nil {
		return encodeTrailingBytes([]byte{0})
	}

	value := td.Index << 3
	if td.FlagDoesSpan != false {
		value |= 0b0001
	}
	if td.FlagTBSType != 0 {
		value |= 0b0010
	}
	if td.FlagNumSiblings > 1 {
		value |= 0b0100
	}
	if td.FlagFirstOfNotFirstStrand != false {
		value |= 0b1000
		panic("unreachable")
	}

	b := bytes.NewBuffer([]byte{0})
	b.Write(encodeVwi(value))
	if td.FlagTBSType != 0 {
		b.Write(encodeVwi(td.FlagTBSType))
	}
	if td.FlagNumSiblings > 1 {
		b.WriteByte(td.FlagNumSiblings)
	}
	if td.FlagDoesSpan != false {
		b.Write(encodeVwi(0))
	}

	return encodeTrailingBytes(b.Bytes())
}

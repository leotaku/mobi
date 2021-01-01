package records

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
)

type ImageRecord struct {
	img     *image.Image
	encoded []byte
	error   error
}

func NewImageRecord(img image.Image) ImageRecord {
	return ImageRecord{
		img: &img,
	}
}

func (r ImageRecord) Length() int {
	r.maybeEncode()
	return len(r.encoded)
}

func (r ImageRecord) Write(w io.Writer) error {
	r.maybeEncode()
	if r.error != nil {
		return r.error
	}

	_, err := w.Write(r.encoded)
	return err
}

func (r *ImageRecord) maybeEncode() {
	if r.img != nil {
		r.encoded, r.error = encodeJFIF(*r.img, nil)
		r.img = nil
	}
}

var naiveJFIFHeader = []byte{
	0xFF, 0xD8, // SOI
	0xFF, 0xE0, // APP0 Marker
	0x00, 0x10, // Length
	0x4A, 0x46, 0x49, 0x46, 0x00, // JFIF\0
	0x01, 0x02, // 1.02
	0x00,       // Density type
	0x00, 0x01, // X Density
	0x00, 0x01, // Y Density
	0x00, 0x00, // No Thumbnail
}

func encodeJFIF(img image.Image, o *jpeg.Options) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := jpeg.Encode(buf, img, o)
	if err != nil {
		return nil, err
	}

	// Connect header and body
	body := buf.Bytes()[2:]
	data := append(naiveJFIFHeader, body...)

	return data, err
}

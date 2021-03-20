package records

import (
	"image"
	"io"

	"github.com/leotaku/mobi/jfif"
)

type ImageRecord struct {
	img image.Image
}

func NewImageRecord(img image.Image) ImageRecord {
	return ImageRecord{
		img: img,
	}
}

func (r ImageRecord) Write(w io.Writer) error {
	return jfif.Encode(w, r.img, nil)
}

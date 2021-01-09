// Package jfif implements writing JPEG images with fixed JFIF header.
package jfif

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
)

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

func Encode(w io.Writer, img image.Image, o *jpeg.Options) error {
	buf := bytes.NewBuffer(nil)
	err := jpeg.Encode(buf, img, o)
	if err != nil {
		return err
	}

	// Connect header and body
	body := buf.Bytes()[2:]
	w.Write(naiveJFIFHeader)
	w.Write(body)

	return nil
}

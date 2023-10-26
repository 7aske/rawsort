package main

import (
	"github.com/rwcarlsen/goexif/exif"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ExifData contains the relevant data extracted from a file's EXIF data.
type ExifData struct {
	// Make Camera make (Nikon, Canon, ...)
	Make string
	// Model Camera model (D750, 5D Mark IV, ...)
	Model string
	// DateTime Date and time
	DateTime time.Time
	// Ext File extension (.jpg, .png, .raw, .nef, ...)
	Ext string
}

// ReadExifData reads the EXIF data from a file and returns an ExifData struct.
func ReadExifData(path string) (ExifData, error) {
	file, err := os.Open(path)
	if err != nil {
		return ExifData{}, err
	}
	defer file.Close()

	x, err := exif.Decode(file)
	if err != nil {
		return ExifData{}, err
	}

	camMake, _ := x.Get(exif.Make)
	camMakeStr, _ := camMake.StringVal()

	camModel, _ := x.Get(exif.Model)
	camModelStr, _ := camModel.StringVal()

	tm, _ := x.DateTime()

	return ExifData{
		Make:     camMakeStr,
		Model:    camModelStr,
		DateTime: tm,
		Ext:      strings.ToUpper(filepath.Ext(path)),
	}, nil
}

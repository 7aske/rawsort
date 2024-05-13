package exif

import (
	"github.com/rwcarlsen/goexif/exif"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var MakeMap = map[string]string{
	"NIKON CORPORATION":           "Nikon",
	"NIKON":                       "Nikon",
	"FUJIFILM":                    "Fujifilm",
	"Sony Ericsson":               "Sony",
	"OLYMPUS IMAGING CORP.":       "Olympus",
	"SIGMA":                       "Sigma",
	"LEICA":                       "Leica",
	"RICOH IMAGING COMPANY, LTD.": "Ricoh",
	"KODAK":                       "Kodak",
	"LG Electronics":              "LG",
	"samsung":                     "Samsung",
}

var ModelMap = map[string]string{
	"FinePix X100": "X100",
}

var OnlyAlphaRegex = regexp.MustCompile(`\W+`)
var WhiteSpaceRegex = regexp.MustCompile(`\s+`)

// Data contains the relevant data extracted from a file's EXIF data.
type Data struct {
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
func ReadExifData(path string) (*Data, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	x, err := exif.Decode(file)
	if err != nil {
		return nil, err
	}

	dateTime := time.UnixMilli(0)
	if tm, err := x.DateTime(); err == nil {
		dateTime = tm
	}

	return &Data{
		Make:     adaptMakeName(getStringValue(x, exif.Make)),
		Model:    adaptModelName(getStringValue(x, exif.Model)),
		DateTime: dateTime,
		Ext:      strings.ToUpper(filepath.Ext(path)),
	}, nil
}

func getStringValue(x *exif.Exif, tag exif.FieldName) string {
	val, err := x.Get(tag)
	if err != nil {
		return "Unknown-" + string(tag)
	}
	str, err := val.StringVal()
	if err != nil {
		return "Unknown-" + string(tag)
	}
	return str
}

func adaptModelName(model string) string {
	if val, ok := ModelMap[model]; ok {
		return val
	}

	model = strings.TrimSpace(model)
	model = WhiteSpaceRegex.ReplaceAllString(model, "-")
	model = OnlyAlphaRegex.ReplaceAllString(model, "")
	return model
}

func adaptMakeName(make string) string {
	if val, ok := MakeMap[make]; ok {
		return val
	}
	make = strings.TrimSpace(make)
	make = WhiteSpaceRegex.ReplaceAllString(make, "-")
	make = OnlyAlphaRegex.ReplaceAllString(make, "")
	return make
}

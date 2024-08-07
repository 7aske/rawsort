package util

import (
	"fmt"
	"github.com/7aske/rawsort/internal/exif"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// FormatFilename formats a filename based on the given format string and ExifData
// extracted from a file.
//
// Format:
//
// %K - Camera Make (Nikon, Canon, ...)
//
// %L - Camera Model (D750, 5D Mark IV, ...)
//
// %D - Date (yyyy-mm-dd)
//
// %t - Time (HH:MM:SS)
//
// %y - Year (yyyy)
//
// %m - Month (mm)
//
// %d - Day (dd)
//
// %e - Extension (.jpg, .png, .raw, .nef, ...)
func FormatFilename(format string, data *exif.Data) string {
	var sb strings.Builder
	isFormat := false
	for _, c := range format {
		if c == '%' {
			isFormat = true
		} else if isFormat {
			isFormat = false
			switch c {
			case 'K':
				sb.WriteString(data.Make)
			case 'L':
				sb.WriteString(data.Model)
			case 'D':
				sb.WriteString(data.DateTime.Format(time.DateOnly))
			case 't':
				// Colons cause problems in smb mounted drives
				formatted := strings.ReplaceAll(data.DateTime.Format(time.TimeOnly), ":", "")
				sb.WriteString(formatted)
			case 'y':
				sb.WriteString(strconv.Itoa(data.DateTime.Year()))
			case 'm':
				sb.WriteString(monthToNumber(data.DateTime.Month()))
			case 'd':
				sb.WriteString(strconv.Itoa(data.DateTime.Day()))
			case 'e':
				sb.WriteString(data.Ext)
			}
		} else {
			sb.WriteRune(c)
		}
	}

	return sb.String()
}

// monthToNumber converts a time.Month to a string representation of the month number with a leading zero.
func monthToNumber(month time.Month) string {
	switch month {
	case time.January:
		return "01"
	case time.February:
		return "02"
	case time.March:
		return "03"
	case time.April:
		return "04"
	case time.May:
		return "05"
	case time.June:
		return "06"
	case time.July:
		return "07"
	case time.August:
		return "08"
	case time.September:
		return "09"
	case time.October:
		return "10"
	case time.November:
		return "11"
	case time.December:
		return "12"
	}
	return ""
}

// Copy copies the contents of the file named src to the file named by dst.
func Copy(src, dst string) (size int64, err error) {
	in, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	written, err := io.Copy(out, in)
	if err != nil {
		return 0, err
	}
	err = out.Sync()
	return written, err
}

// RenameFile renames a file adding a number suffix before the extensions
func RenameFile(path string) string {
	dir := filepath.Dir(path)
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(filepath.Base(path), ext)
	for i := 1; ; i++ {
		newPath := filepath.Join(dir, fmt.Sprintf("%s-%d%s", base, i, ext))
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
	}
}

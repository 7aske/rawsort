package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"log"
	"os"
	"path/filepath"
)

// Args Command line arguments
type Args struct {
	// Src    string Source folder
	Src string
	// Dest   string Destination folder
	Dest string
	// Format string Filename format
	Format string
}

// DefaultFormat Default format for file sorting:
// Manufacturer/Model/Date/Date_Time_Manufacturer_Model.extension
// FUJIFULM/X100/2023-01-01/2023-01-01_12-00-00_FUJIFILM_X100.RAF
const DefaultFormat = "%K/%L/%D/%D_%t_%K_%L%e"

func parseArgs() Args {
	parser := argparse.NewParser("rawsort",
		"Sorts raw files into folders by camera make, model and date")

	s := parser.String("s", "src", &argparse.Options{
		Required: true,
		Help:     "Source folder",
	})
	d := parser.String("d", "dest", &argparse.Options{
		Required: true,
		Help:     "Destination folder",
	})
	f := parser.String("f", "format", &argparse.Options{
		Required: false,
		Help: "Filename format\n" +
			"  Filename format options:\n" +
			"    %%D - Date\n" +
			"    %%t - Time\n" +
			"    %%y - Year\n" +
			"    %%m - Month (mm)\n" +
			"    %%d - Day\n" +
			"    %%K - Make\n" +
			"    %%L - Model\n" +
			"    %%e - Extension\n",
	})

	// If the default value is specified in the parser.String the help output
	// is rather ugly, so we set it here instead
	if *f == "" {
		*f = DefaultFormat
	}

	err := parser.Parse(os.Args)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, parser.Usage(err))
		os.Exit(1)
	}

	return Args{
		Src:    *s,
		Dest:   *d,
		Format: *f,
	}
}

func main() {
	args := parseArgs()

	err := filepath.Walk(args.Src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		data, err := ReadExifData(path)
		if err != nil {
			log.Println(err)
			return nil
		}

		fileName := FormatFilename(args.Format, data)

		destPath := filepath.Join(args.Dest, fileName)
		err = os.MkdirAll(filepath.Dir(destPath), 0755)
		if err != nil {
			return err
		}

		destInfo, err := os.Stat(destPath)
		// File exists and is the same file
		if err != nil && (destInfo != nil && destInfo.Size() == info.Size()) {
			return nil
		}

		// File exists and is not the same file
		if destInfo != nil && destInfo.Size() != info.Size() {
			// TODO: rename dest file
			return nil
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s -> %s\n", path, destPath)
		err = CopyFileContents(path, destPath)

		return nil
	})

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

}

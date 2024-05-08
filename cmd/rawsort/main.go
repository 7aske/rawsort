package main

import (
	"fmt"
	"github.com/7aske/rawsort/internal/exif"
	"github.com/7aske/rawsort/internal/util"
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
	// Verbose verbosity of the output
	Verbose bool
	// Interactive asks for user input before renaming files
	Interactive bool
}

// DefaultFormat Default format for file sorting:
// Manufacturer/Model/Date/Date_Time_Manufacturer_Model.extension
// FUJIFILM/X100/2023-01-01/2023-01-01_12-00-00_FUJIFILM_X100.RAF
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
			"    %D - Date\n" +
			"    %t - Time\n" +
			"    %y - Year\n" +
			"    %m - Month (mm)\n" +
			"    %d - Day\n" +
			"    %K - Make\n" +
			"    %L - Model\n" +
			"    %e - Extension\n",
	})
	v := parser.Flag("v", "verbose", &argparse.Options{
		Required: false,
		Help:     "Verbose output",
	})
	i := parser.Flag("i", "interactive", &argparse.Options{
		Required: false,
		Help:     "Interactive mode, ask user on conflicts",
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
		Src:         *s,
		Dest:        *d,
		Format:      *f,
		Verbose:     *v,
		Interactive: *i,
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

		data, err := exif.ReadExifData(path)
		if err != nil {
			log.Println(err)
			return nil
		}

		fileName := util.FormatFilename(args.Format, data)

		destPath := filepath.Join(args.Dest, fileName)
		err = os.MkdirAll(filepath.Dir(destPath), 0755)
		if err != nil {
			return err
		}

		destInfo, err := os.Stat(destPath)
		// File exists and is the same file
		if err != nil && (destInfo != nil && destInfo.Size() == info.Size()) {
			if args.Verbose {
				_, _ = fmt.Fprintf(os.Stderr, "Duplicate file found: %s\n", destPath)
			}
			return nil
		}

		// File exists and is not the same file
		if destInfo != nil && destInfo.Size() != info.Size() {
			// skip, overwrite or rename
			// ask user for input
			if args.Interactive {
				for {
					_, _ = fmt.Fprintf(os.Stderr, "File %s already exists (a)bort, (s)kip, (o)verwrite or (r)ename? (a,s,o,r): ", destPath)
					var input string
					_, err := fmt.Scanln(&input)
					if err != nil {
						return err
					}

					switch input {
					case "a":
						os.Exit(1)
					case "s":
						return nil
					case "o":
						break
					case "r":
						destPath = util.RenameFile(destPath)
					}
				}
			} else {
				if args.Verbose {
					_, _ = fmt.Fprintf(os.Stderr, "File %s already exists\n", destPath)
				}
				destPath = util.RenameFile(destPath)
			}
		}

		if args.Verbose {
			_, _ = fmt.Fprintf(os.Stderr, "%s -> %s\n", path, destPath)
		}
		err = util.CopyFileContents(path, destPath)

		return nil
	})

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

}

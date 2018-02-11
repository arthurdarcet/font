// font is a utility that can parse and print information about font files.
package main

import (
	"fmt"
	"os"

	"github.com/ConradIrwin/font/sfnt"
)

func usage() {
	fmt.Println(`
<<<<<<< HEAD
Usage: font [features|info|metrics|scrub|stats] font.[otf,ttf,woff,woff2] ...
=======
Usage: font [cmap|features|info|metrics|scrub|stats] font.[otf,ttf,woff]
>>>>>>> 956c6b9 (Add initial support for parsing the cmap table.)

cmap: prints out Character To Glyph mappings
features: prints the gpos/gsub tables (contains font features)
info: prints the name table (contains metadata)
metrics: prints the hhea table (contains font metrics)
scrub: remove the name table (saves significant space)
stats: prints each table and the amount of space used`)
}

func main() {
	command := "help"
	if len(os.Args) > 1 {
		command = os.Args[1]
		os.Args = os.Args[1:]
	}

	cmds := map[string]func(*sfnt.Font) error{
		"scrub":    Scrub,
		"info":     Info,
		"stats":    Stats,
		"metrics":  Metrics,
		"features": Features,
		"cmap":     Cmap,
	}
	if _, found := cmds[command]; !found {
		usage()
		return
	}

	if len(os.Args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: font %s <font file> ...\n", command)
		os.Exit(1)
	}

	exitCode := 0
	for _, filename := range os.Args[1:] {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open font: %s\n", err)
			exitCode = 1
			continue
		}
		defer file.Close()

		font, err := sfnt.Parse(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse font: %s\n", err)
			exitCode = 1
			continue
		}

		if len(os.Args[1:]) > 1 {
			fmt.Println("==>", filename, "<==")
		}
		if err := cmds[command](font); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			exitCode = 1
			continue
		}
	}
	os.Exit(exitCode)
}

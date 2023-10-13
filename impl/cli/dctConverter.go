package main

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/zocker-160/Berlin-Taxi-DCT/impl/dct"
)

const UsageMsg = `
Usage:
    dctConverter <imagefile> <destination folder>
`

func dct2png(filename, name, dstFolder string) error {
	fmt.Println("Converting", name, "to PNG")

	dct, err := dct.DecodeFile(filename)
	if err != nil {
		return err
	}

	outFile, err := os.Create(filepath.Join(dstFolder, name+".png"))
	if err != nil {
		return err
	}
	defer outFile.Close()

	return png.Encode(outFile, dct)
}

func png2dct(filename, name, dstFolder string) error {
	fmt.Println("Converting", name, "to DCT")

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	image, err := png.Decode(file)
	if err != nil {
		return err
	}

	outFile, err := os.Create(filepath.Join(dstFolder, name+".dct"))
	if err != nil {
		return err
	}
	defer outFile.Close()

	return dct.Encode(outFile, image)
}

func main() {
	fmt.Println("DCT converter by zocker_160")

	if len(os.Args) != 3 {
		fmt.Print(UsageMsg)
		return
	}

	sourceFile := os.Args[1]
	dstFolder := os.Args[2]

	basename := filepath.Base(sourceFile)
	name := strings.TrimSuffix(basename, filepath.Ext(basename))

	switch ext := filepath.Ext(sourceFile); strings.ToLower(ext) {
	case ".dct":
		err := dct2png(sourceFile, name, dstFolder)
		if err != nil {
			panic(err)
		}
	case ".png":
		err := png2dct(sourceFile, name, dstFolder)
		if err != nil {
			panic(err)
		}
	default:
		fmt.Println("ERROR: unsupported file type", ext)
		return
	}
}

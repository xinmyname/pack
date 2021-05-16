package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"pack/lzss"
	"path/filepath"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: pack <input file name> [output file name]")
		fmt.Println()
		fmt.Println("       If no output file is specified, the output file name will be the same")
		fmt.Println("       as the input file name with the .lzss extension.")
		return
	}

	inFname := os.Args[1]

	_, err := os.Stat(inFname)

	if os.IsNotExist(err) {
		log.Fatalf("%v does not exist!\n", inFname)
	}

	var outFname string

	if len(os.Args) > 2 {
		outFname = os.Args[2]
	} else {
		outFname = strings.TrimSuffix(inFname, filepath.Ext(inFname)) + ".lzss"
	}

	inFile, err := os.Open(inFname)
	if err != nil {
		log.Fatal(err)
	}

	outFile, err := os.Create(outFname)
	if err != nil {
		log.Fatal(err)
	}

	outCompressed := lzss.NewWriter(outFile)

	_, err = io.Copy(outCompressed, inFile)
	if err != nil {
		log.Fatal(err)
	}

	outCompressed.Close()
}

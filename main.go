package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"pack/lzss"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: pack [-u] <input file name>")
		fmt.Println()
		fmt.Println("    -u  Unpack file")
		return
	}

	var srcPath string
	unpack := false

	args := os.Args[1:]

	for ; len(args) > 0; args = args[1:] {
		if args[0] == "-u" {
			unpack = true
		} else {
			srcPath = args[0]
		}
	}

	_, err := os.Stat(srcPath)

	if os.IsNotExist(err) {
		log.Fatalf("%v does not exist!\n", srcPath)
	}

	inFile, err := os.Open(srcPath)
	if err != nil {
		log.Fatal(err)
	}

	if unpack {

		outUncompressed := lzss.NewReader(inFile)

		_, err := io.Copy(os.Stdout, outUncompressed)
		if err != nil {
			log.Fatal(err)
		}

		outUncompressed.Close()

	} else {

		outCompressed := lzss.NewWriter(os.Stdout)

		_, err = io.Copy(outCompressed, inFile)
		if err != nil {
			log.Fatal(err)
		}

		outCompressed.Close()
	}
}

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

	inFileInfo, err := os.Stat(inFname)

	if os.IsNotExist(err) {
		log.Fatalf("%v does not exist!\n", inFname)
	}

	fmt.Printf("%v - %v bytes in", inFname, inFileInfo.Size())

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

	/*
		data, err := ioutil.ReadAll(inFile)
		if err != nil {
			log.Fatal(err)
		}

		compressed := compress(data)
		encoded := encode(compressed)
	*/
	outFile, err := os.Open(outFname)
	if err != nil {
		log.Fatal(err)
	}

	outCompressed := lzss.NewWriter(outFile, 4096)

	wrote, err := io.Copy(outCompressed, inFile)
	if err != nil {
		log.Fatal(err)
	}

	outFile.Close()

	fmt.Printf("%v - %v bytes out", outFname, wrote)
}

/*
func compress(data []byte) []int {

	compressed := []int{}

	for _, b := range data {
		compressed = append(compressed, int(b))
	}

	return compressed
}

func encode(compressed []int) []byte {

	encoded := []byte{}

	for _, i := range compressed {
		encoded = append(encoded, byte(i))
	}

	return encoded
}
*/

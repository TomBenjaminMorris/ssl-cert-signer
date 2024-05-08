package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	ca, _ := loadCA()
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <csr>\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// fileData, err := getFileData()
	fileData, err := getCSR()

	if err != nil {
		exitGracefully(err)
	}

	if _, err := checkIfValidFile(fileData.filepath); err != nil {
		exitGracefully(err)
	}

	// fmt.Println(fileData)
	csr := parseCSR(fileData)
	cert := signCSR(csr, ca)
	writeCert(cert)
	// fmt.Println(csr.csr)
}

// changePDFmeta changes the META field from PDF file to new field.
// This tool use pdftk (pdftk must be installed in you system).
// The pdf file names and new fields are being taken from the csv-file (filename;new title)
// The csv file and field name to change are taken from args.

package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {

	newPdfDirName := "newpdf"

	if len(os.Args) < 3 {
		fmt.Println(os.Args[0] + " — change the PDF meta field on the files from a csv-list.")
		fmt.Println("CSV-file format: (filename;new title)")
		fmt.Println("\nSelect the csv file name and the meta-field name.")
		fmt.Println(os.Args[0] + " filename.csv" + " metaname")
		os.Exit(-1)
	}

	filename := os.Args[1]
	metaname := os.Args[2]

	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't open "+filename)
		os.Exit(1)
	}

	r := csv.NewReader(f)

	// Set csv-delimiter
	r.Comma = ';'

	s, err := r.ReadAll()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't read or parse "+filename)
		os.Exit(2)
	}

	// Try to make a dir to the new pdf files
	if _, err := os.Stat(newPdfDirName); os.IsNotExist(err) {
		if err := os.Mkdir(newPdfDirName, os.ModeDir|os.ModePerm); err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, "Can't create the directory: "+newPdfDirName)
			os.Exit(5)
		}
	}

	for _, line := range s {
		if len(line) >= 2 {
			changePDFmeta(line[0], line[1], metaname, newPdfDirName)
		} else {
			fmt.Fprintln(os.Stderr, "Format of the csv-file line is not correct:")
			fmt.Fprintln(os.Stderr, line)
		}
	}
}

// changePDFmeta sets the META metaname properties to str in the pdffilename
func changePDFmeta(pdfFilename string, str string, metaname string, newPdfDir string) {

	rand.Seed(time.Now().UTC().UnixNano())

	tempDir := os.TempDir() + "/" + "changePDFmeta"
	tempFilename := tempDir + "/" + metaname + strconv.Itoa(rand.Int())

	var fileContent []byte

	fileContent = []byte("InfoKey: " + metaname + "\nInfoValue: " + str)

	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		if err := os.Mkdir(tempDir, os.ModeDir|os.ModePerm); err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, "Can't create the temporary directory: "+tempDir)
			os.Exit(3)
		}
	}

	// Make tempFile to pdftk
	if err := ioutil.WriteFile(tempFilename, fileContent, os.ModePerm); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "Can't write to the file: "+tempFilename)
		os.Exit(4)
	}

	c := exec.Command("pdftk", pdfFilename, "update_info_utf8", tempFilename, "output", newPdfDir+"/"+pdfFilename)
	fmt.Printf("Run: %s — ", c.Args)
	err := c.Run()

	if err != nil {
		fmt.Println("Error")
	} else {
		fmt.Println("Success")
	}

}

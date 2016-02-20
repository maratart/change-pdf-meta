package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"
)

func main() {

	newPdfDirName := "new-pdf"

	if len(os.Args) < 4 {
		execname := path.Base(os.Args[0])
		fmt.Println(execname + " — change the PDF meta field on the files from a csv-list.")
		fmt.Println("CSV-file format: (filename;new title)")
		fmt.Println("\nUsage:")
		fmt.Println("Select the csv file name and the meta-field name.")
		fmt.Println("./" + execname + " filename.csv" + " metaname" + " pdfdir")
		os.Exit(0)
	}

	filename := os.Args[1]
	metaname := os.Args[2]
	pdfdir := os.Args[3]

	newPdfDirName = pdfdir + "/" + newPdfDirName

	if wd, err := os.Getwd(); err == nil {
		fmt.Println("Work dir: " + wd)
	}

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
			changePDFmeta(line[0], line[1], metaname, newPdfDirName, pdfdir)
		} else {
			fmt.Fprintln(os.Stderr, "Format of the csv-file line is not correct:")
			fmt.Fprintln(os.Stderr, line)
		}
	}
}

// changePDFmeta sets the META metaname properties to str in the pdffilename
func changePDFmeta(pdfFilename, str, metaname, newPdfDir, pdfdir string) {

	rand.Seed(time.Now().UTC().UnixNano())

	tempDir := os.TempDir() + "/" + "change-pdf-meta"
	tempFilename := tempDir + "/" + strconv.Itoa(rand.Int())

	var fileContent []byte

	fileContent = []byte("InfoKey: " + metaname + "\nInfoValue: " + str)

	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		if err := os.Mkdir(tempDir, os.ModeDir|os.ModePerm); err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, "Can't create the temporary directory: "+tempDir)
			os.Exit(3)
		}
	}

	// Make tempFile for pdftk
	if err := ioutil.WriteFile(tempFilename, fileContent, os.ModePerm); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "Can't write to the file: "+tempFilename)
		os.Exit(4)
	}

	//args := []string{"'" + pdfdir + "/" + pdfFilename + "'", "update_info_utf8", tempFilename, "output", "'" + newPdfDir + "/" + pdfFilename + "'"}

	args := []string{pdfdir + "/" + pdfFilename, "update_info_utf8",
		tempFilename, "output", newPdfDir + "/" + pdfFilename}

	c := exec.Command("pdftk", args...)

	fmt.Printf("Run: %s — ", c.Args)

	err := c.Run()

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Success")
	}
}

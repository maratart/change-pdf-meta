# ChangePDFmeta changes the META field from PDF file to new field.

This tool use pdftk (pdftk must be installed in you system).

The pdf file names and new fields are being taken from the csv-file (filename;new title)
The csv file and field name to change are taken from args.

## Compile
You must compile this tool from the sourse by running:
```
go build changePDFmeta.go
```


##USAGE: 
```
changePDFmeta table.csv Title
```

table.csv
---------
file1.pdf;New title for file1   
file2.pdf;New title for file2
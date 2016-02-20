# Change-pdf-meta changes the content of a META field of a PDF file to the new field content.

This tool use pdftk (pdftk must be installed in you system).

The pdf file names and new fields are being taken from the csv-file (filename;new title)
The csv file and field name to change are taken from args.

## Compile
You must compile this tool from the source by running:  
```
go build change-pdf-meta.go
```


##USAGE: 
```
./change-pdf-meta table.csv Title
```

table.csv
---------
file1.pdf;New title for file1   
file2.pdf;New title for file2
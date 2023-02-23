package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (app *application) ImageHandler(w http.ResponseWriter, r *http.Request) { 
 
	r.ParseMultipartForm(10 * 1024 * 1024) 
	
	file, handler, err := r.FormFile("myfile") 
	
	if err != nil { 
		fmt.Println(err) 
		return 
	} 
	
	defer file.Close() 
	
	fileBytes, err := ioutil.ReadAll(file) 
	if err != nil { 
		fmt.Println(err) 
		return 
	} 
	handler.Header.Get("Content-Type") 
	
	fileType := http.DetectContentType(fileBytes) 
	
	switch fileType { 
	case "image/jpeg", "image/jpg": 
		tempFile, err := ioutil.TempFile("images", "upload-*.jpg") 
		if err != nil { 
			fmt.Println(err) 
			return 
	 	} 
		defer tempFile.Close() 
		tempFile.Write(fileBytes) 
		fmt.Println("JPEG file uploaded") 
	case "application/pdf": 
		tempFile, err := ioutil.TempFile("pdfs", "upload-*.pdf") 
		if err != nil { 
			fmt.Println(err) 
			return 
		} 
		defer tempFile.Close() 
		tempFile.Write(fileBytes) 
		fmt.Println("PDF file uploaded") 
	case "application/docx": 
	 tempFile, err := ioutil.TempFile("docs", "upload-*.docx") 
		if err != nil { 
			fmt.Println(err) 
			return 
		} 
		defer tempFile.Close() 
		tempFile.Write(fileBytes) 
		fmt.Println("Word file uploaded") 
	default: 
		fmt.Println("Unsupported file type") 
		return 
	} 
   }
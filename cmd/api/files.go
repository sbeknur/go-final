package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (app *application) fileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024)

	file, handler, err := r.FormFile("myfile")
	if err != nil {      fmt.Println(err)
	   return   
	}
 
	defer file.Close()
	fmt.Println("File info")
	fmt.Println("File name:", handler.Filename)   
	fmt.Println("File size:", handler.Size)
	fmt.Println("File type:", handler.Header.Get("Content-Type"))

	tempFile, err2 := ioutil.TempFile("files/image", "upload-*.jpg")
	if err2 != nil {      fmt.Println(err2)
	   return   }
	defer tempFile.Close()

	tempFile, err4 := ioutil.TempFile("files/pdf", "upload-*.pdf")   
	if err2 != nil {
	   fmt.Println(err4)
	   return   }
	defer tempFile.Close()

	fileBytes, err3 := ioutil.ReadAll(file)   
	if err3 != nil {
	   fmt.Println(err2)      
	   return
	}   
	
	tempFile.Write(fileBytes)

	fmt.Println("File sended!")
 }
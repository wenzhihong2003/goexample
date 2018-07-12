package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	f1()
}

// client 上传单个文件
func f1() {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("files", "file.txt")
	file, _ := os.Open("file.txt")
	defer file.Close()
	io.Copy(fileWriter, file)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	resp, _ := http.Post("http://localhost:8080/upload", contentType, bodyBuffer)
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	log.Println(resp.Status)
	log.Println(string(respBody))
}

// client 上传多个文件
func f2() {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	// file1
	fileWriter1, _ := bodyWriter.CreateFormFile("files", "file1.txt")
	file1, _ := os.Open("file1.txt")
	defer file1.Close()
	io.Copy(fileWriter1, file1)

	// file2
	fileWriter2, _ := bodyWriter.CreateFormFile("files", "file2.txt")
	file2, _ := os.Open("fil2.txt")
	defer file2.Close()
	io.Copy(fileWriter2, file2)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, _ := http.Post("http://localhost:8080/upload", contentType, bodyBuffer)
	resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	log.Println(resp.Status)
	log.Println(string(respBody))
}

// 上传其它form数据
func f3() {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	// file1
	fileWriter1, _ := bodyWriter.CreateFormFile("files", "file1.txt")
	file1, _ := os.Open("file1.txt")
	defer file1.Close()
	io.Copy(fileWriter1, file1)

	// file2
	fileWriter2, _ := bodyWriter.CreateFormFile("files", "file2.txt")
	file2, _ := os.Open("fil2.txt")
	defer file2.Close()
	io.Copy(fileWriter2, file2)

	// other form data
	extraParams := map[string]string{
		"title": "My document",
		"auth":  "Matt aimonetti",
	}
	for k, v := range extraParams {
		_ = bodyWriter.WriteField(k, v)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, _ := http.Post("http://localhost:8080/upload", contentType, bodyBuffer)
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	log.Println(resp.Status)
	log.Println(string(respBody))
}

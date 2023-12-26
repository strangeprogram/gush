package main

import (
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

var hardfilesFlag bool

func init() {
	flag.BoolVar(&hardfilesFlag, "hf", false, "Push the file to Hardfiles.org")
	flag.Parse()
}

func UploadToHardFiles(filename string) error {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Prepare the file for upload
	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	writer.Close()

	// Perform the upload to hardfiles.org
	request, err := http.NewRequest("POST", "https://hardfiles.org/upload", strings.NewReader(body.String()))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body using io.ReadAll (replacement for ioutil.ReadAll)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Handle the response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with status: %s", resp.Status)
	}

	// TODO: Parse the response to extract the file URL.
	// The actual implementation depends on how hardfiles.org structures its response.
	// Here's a placeholder for where you'd parse the URL:
	uploadedFileURL := string(respBody) // Replace with actual parsing logic
	fmt.Printf("File uploaded successfully to HardFiles.org! File URL: %s\n", uploadedFileURL)
	return nil
}

func main() {
	if flag.NArg() != 1 {
		fmt.Println("Usage: gush -hf filename")
		return
	}

	if hardfilesFlag {
		filename := flag.Arg(0)
		err := UploadToHardFiles(filename)
		if err != nil {
			fmt.Printf("Failed to upload file: %s\n", err)
			return
		}
	} else {
		fmt.Println("Please specify an associated flag for a service to push the file")
	}
}

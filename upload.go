package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
)

type UploadResponse struct {
	Files []struct {
		URL string `json:"url"`
	} `json:"files"`
}

func detectContentType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the file's magic numbers
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading file:", err)
		return "", err
	}

	// Detect the content type
	contentType := http.DetectContentType(buffer)
	fmt.Println("Detected Content-Type:", contentType)
	if contentType == "text/plain; charset=utf-8" {
		ext := filepath.Ext(filePath)
		mimeType := mime.TypeByExtension(ext)
		if mimeType != "" {
			contentType = mimeType
			fmt.Println("Fallback Content-Type:", contentType)
		}
	}
	return contentType, nil
}

func uploadFile(serverName, token, filePath string) (string, error) {
	parsedURL, err := url.Parse(serverName)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return "", err
	}

	hostname := parsedURL.Hostname()

	if hostname == "" {
		// If parsing as a full URL fails, treat serverName as just the hostname
		hostname = serverName
	}

	if hostname == "" {
		return "", fmt.Errorf("invalid servername in config file: %s", serverName)
	}

	contentType, err := detectContentType(filePath)
	if err != nil {
		fmt.Println("Error detecting content type:", err)
		fmt.Println("Content-Type set to: application/octet-stream")
		return "", err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Reset the file reader to the beginning of the file
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println("Error seeking to start of file:", err)
		return "", err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			"file", filepath.Base(filePath)))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	var uploadURL string
	if config.HTTPInsecure {
		uploadURL = "http://" + hostname + "/api/upload"
	} else {
		uploadURL = "https://" + hostname + "/api/upload"
	}

	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var uploadResponse UploadResponse
	err = json.Unmarshal(respBody, &uploadResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON: %w, body: %s", err, string(respBody))
	}

	if len(uploadResponse.Files) > 0 {
		return uploadResponse.Files[0].URL, nil
	}

	return "", fmt.Errorf("no files found in response: %s", string(respBody))
}

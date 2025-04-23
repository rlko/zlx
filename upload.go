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

func prepareMultipartBody(filePath string) (*bytes.Buffer, string, string, error) {
	contentType, err := detectContentType(filePath)
	if err != nil {
		fmt.Println("Error detecting content type:", err)
		fmt.Println("Content-Type set to: application/octet-stream")
		return nil, "", "", err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", "", err
	}
	defer file.Close()

	// Reset file reader
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println("Error seeking to start of file:", err)
		return nil, "", "", err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filepath.Base(filePath)))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, "", "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, "", "", err
	}

	err = writer.Close()
	if err != nil {
		return nil, "", "", err
	}

	return body, writer.FormDataContentType(), contentType, nil
}

func getFullURL(serverName string, pathName string) (string, error) {
	var scheme string

	parsedURL, err := url.Parse(serverName)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return "", err
	}

	scheme = parsedURL.Scheme
	serverName = parsedURL.Hostname()

	if scheme == "" {
		scheme = "https://"
	}

	if serverName == "" {
		serverName = config.ServerName
		if serverName == "" {
			return "", fmt.Errorf("empty or invalid \"servername\" in config file")
		}
	}

	if pathName == "" {
		pathName = "/api/upload"
	}

	return scheme + "://" + serverName + pathName, nil
}

func uploadFile(config Config, filePath string) (string, error) {

	//	fmt.Println("Flags: config.Upload.MaxViews:", config.Upload.MaxViews)
	//	fmt.Println("Flags: config.Upload.OriginalName:", config.Upload.OriginalName)
	//	fmt.Println("Flags: config.Upload.Clipboard:", config.Upload.Clipboard)

	url, err := getFullURL(config.ServerName, config.PathName)
	if err != nil {
		fmt.Println("Error getting full URL:", err)
		return "", err
	}

	body, contentTypeHeader, _, err := prepareMultipartBody(filePath)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", config.Token)
	req.Header.Set("Content-Type", contentTypeHeader)

	if config.Upload.MaxViews > 0 {
		req.Header.Set("x-zipline-max-views", fmt.Sprintf("%d", config.Upload.MaxViews))
	}
	req.Header.Set("x-zipline-original-name", fmt.Sprintf("%t", config.Upload.OriginalName))

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

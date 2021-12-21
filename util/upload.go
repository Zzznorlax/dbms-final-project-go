package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

type ImgurUploadData struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Datetime    int    `json:"datetime"`
	Type        string `json:"type"`
	Animated    bool   `json:"animated"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Link        string `json:"link"`
}

type ImgurUploadResp struct {
	Data    ImgurUploadData `json:"data"`
	Success bool            `json:"success"`
	Status  int             `json:"status"`
}

func UploadFile(url string, headers map[string]string, filename string, content []byte, target interface{}) error {

	payload := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(payload)

	fileWriter, err := bodyWriter.CreateFormFile("image", filename)
	if err != nil {
		return fmt.Errorf("creating formfile: %w", err)
	}

	fileWriter.Write(content)
	bodyWriter.Close()

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return fmt.Errorf("creating post request: %w", err)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	return err
}

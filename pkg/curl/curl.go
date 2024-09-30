package curl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func Get(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			var responseBody bytes.Buffer
			_, err = responseBody.ReadFrom(resp.Body)
			if err != nil {
				return "", err
			}
			return responseBody.String(), nil
		}
		return "", fmt.Errorf("GET request failed with status code: %d", resp.StatusCode)
	}

	var responseBody bytes.Buffer
	_, err = responseBody.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	return responseBody.String(), nil
}

func Post(url string, headers map[string]string, body interface{}) (string, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			var responseBody bytes.Buffer
			_, err = responseBody.ReadFrom(resp.Body)
			if err != nil {
				return "", err
			}
			return responseBody.String(), nil
		}
		return "", fmt.Errorf("POST request failed with status code: %d", resp.StatusCode)
	}

	var responseBody bytes.Buffer
	_, err = responseBody.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	return responseBody.String(), nil
}

func Put(url string, headers map[string]string, body map[string]interface{}) (string, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("PUT request failed with status code: %d", resp.StatusCode)
	}

	var responseBody bytes.Buffer
	_, err = responseBody.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	return responseBody.String(), nil
}

func Delete(url string, headers map[string]string, body map[string]interface{}) (string, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("DELETE request failed with status code: %d", resp.StatusCode)
	}

	var responseBody bytes.Buffer
	_, err = responseBody.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	return responseBody.String(), nil
}

type PostFormInput struct {
	Headers map[string]string
	Texts   map[string]string
	Files   map[string]*multipart.FileHeader
}

func PostForm(url string, i PostFormInput) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// The text fields
	for key, value := range i.Texts {
		err := writer.WriteField(key, fmt.Sprintf("%v", value))
		if err != nil {
			return "", err
		}
	}

	// The files
	for key, fileHeader := range i.Files {
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Println("Error opening file: ", err)
			return "", err
		}
		defer file.Close()

		part, err := writer.CreateFormFile(key, fileHeader.Filename)
		if err != nil {
			fmt.Println("pkg/curl: PostForm: CreateFormFile: ", err)
			return "", err
		}

		if _, err := io.Copy(part, file); err != nil {
			fmt.Println("pkg/curl: PostForm: Copy: ", err)
			return "", err
		}
	}

	if err := writer.Close(); err != nil {
		fmt.Println("pkg/curl: PostForm: Close: ", err)
		return "", err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	for key, value := range i.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("POST request failed with status code: %d", resp.StatusCode)
	}

	var responseBody bytes.Buffer
	_, err = responseBody.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	return responseBody.String(), nil
}

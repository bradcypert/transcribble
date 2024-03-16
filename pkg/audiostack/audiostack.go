package audiostack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type AudioStack struct {
	apiKey string
}

func MakeAudioStack(apiKey string) AudioStack {
	return AudioStack{
		apiKey,
	}
}

type UploadBody struct {
	FilePath string `json:"filePath"`
	FileType string `json:"fileType"`
}

type UploadResponseData struct {
	Field         string `json:"field"`
	FileUploadURL string `json:"fileUploadUrl"`
}

type UploadResponse struct {
	Data UploadResponseData `json:"data"`
}

func (u AudioStack) GetFileUploadUrl(file *os.File) string {
	name := file.Name()
	nameParts := strings.Split(name, "/")
	fileName := nameParts[len(nameParts)-1]

	url := "https://v2.api.audio/content/file/create-upload-url"
	fmt.Println(fileName)
	payload, _ := json.Marshal(&UploadBody{
		FilePath: fileName,
		FileType: "video",
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-api-key", u.apiKey)

	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()

	uploadResponse := UploadResponse{}
	json.Unmarshal(body, &uploadResponse)

	// Upload the file to the returned URL
	fileUploadUrl := uploadResponse.Data.FileUploadURL
	fmt.Println(fileUploadUrl)
	return fileUploadUrl
}

func (u AudioStack) UploadFile(uploadUrl string, file *os.File) {
	writePayload := &bytes.Buffer{}
	writer := multipart.NewWriter(writePayload)
	part3, errFile3 := writer.CreateFormFile("file", "out.mp3")
	_, errFile3 = io.Copy(part3, file)
	if errFile3 != nil {
		fmt.Println(errFile3)
		return
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("PUT", uploadUrl, writePayload)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.Status)

	fmt.Println(string(body))
}

package main

import (
	"fmt"
	"os"
)

const API_KEY = "10f5e762-1f78-4943-9dda-79a00723b079"

func main() {
	if len(os.Args) != 2 {
		panic("Please pass a file path to use this CLI")
	}
	filePath := os.Args[1]
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic(fmt.Sprintf("Unable to find the specified file at %s", filePath))
	}

	const API_KEY = "10f5e762-1f78-4943-9dda-79a00723b079"

	audioStack := MakeAudioStack(API_KEY)
	uploadUrl := audioStack.GetFileUploadUrl(file)
	audioStack.UploadFile(uploadUrl, file)
}

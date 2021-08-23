# Golang Yandex Speech Kit [![Test And Linter](https://github.com/zemags/golang-yandex-speech-kit/actions/workflows/pipeline.yml/badge.svg?branch=master)](https://github.com/zemags/golang-yandex-speech-kit/actions/workflows/pipeline.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/zemags/golang-yandex-speech-kit)](https://goreportcard.com/report/github.com/zemags/golang-yandex-speech-kit)


Small simple SDK to convert text to audio by Yandex Speech Kit Service.

SDK uses only premium voices.

Get Yandex Cloud service account Api-Key: https://cloud.yandex.com/en/docs/speechkit/concepts/auth

Example usage:
```go
package main

import (
    "log"
    "os"
    "github.com/pkg/errors"
    speechkit "https://github.com/zemags/golang-yandex-speech-kit"

)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Can't find .env file")
	}
	API_KEY, exist := os.LookupEnv("API_KEY")
	if !exist {
		error.New("Yandex cloud service account Api-Key not provided")
	}
	
	client := &http.Client{Timeout: 5 * time.Second}
	apiParams := speechkit.APIParams{APIKey: API_KEY, Client: client}

        // define folder for mp3 audio
        currentDir, _ := os.Getwd()
        pathToFiles := path.Join(currentDir, "temp_folder")

	speechParams := speechkit.SpeechParams{
        	Voice: "male",
		Speed: 1.0,
        	PathToFiles: pathToFiles,
        }

        client := speechkit.NewSpeechKitClient(apiParams, speechParams)

        exampleTextForAudio := "Lorem Ipsum is simply dummy."
        err := client.CreateAudio(exampleTextForAudio)
        if err != nil {
           	log.Println(err)
	}
}

```
Define Yandex Cloud service accout Api-Key:
| Param   | Type   | Defenition                                               |
| ------- | ------ | -------------------------------------------------------- |
| API_KEY | string | https://cloud.yandex.com/en/docs/speechkit/concepts/auth |


You can define SpeechParams by this table:
| Param       | Type    | Variables             | Default        | Defenition                                                           |
| ----------- | ------- | --------------------- | -------------- | -------------------------------------------------------------------- |
| Voice       | string  | "male" <br/> "female" | "male"         | "male" provided voice 'filipp' <br/> "female" provided voice 'alena' |
| Speed       | float32 | from 0.1 to 3.0       | 1.0            | speed of synthesized speech                                          |
| PathToFiles | string  | NEED TO DEFINE        | NEED TO DEFINE | path to synthesized audio mp3                                        |

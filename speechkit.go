// Simple SDK to sythesize russian voice from text with
// Yandex Speech Kit Service premium voices

package speechkit

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
)

// URL provided remote host to yandex speech kit api
const URL = "https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize"

var (
	speechSpeed    = float32(1.0)
	speechLanguage = "ru-RU"
	speechFormat   = "oggopus"
	speechEmotion  = "neutral"
	textMaxLen     = 2000
	output         = "output.txt"
)

// SpeechKitClient main client for generating audio
type SpeechKitClient struct { //nolint
	APIParams
	SpeechParams
}

// APIParams define how access remote yandex api endpoint
type APIParams struct {
	Client *http.Client
	APIKey string
}

// SpeechParams user options for audio
type SpeechParams struct {
	emotion     string
	voice       string
	speed       float32
	pathToFiles string
}

// NewSpeechKitClient create new client
func NewSpeechKitClient(apiParams APIParams, speechParams SpeechParams) *SpeechKitClient {
	return &SpeechKitClient{
		APIParams:    apiParams,
		SpeechParams: speechParams,
	}
}

// CreateAudio receive user text and generate audio
func (c *SpeechKitClient) CreateAudio(text string) error {
	output, err := c.createFile()
	if err != nil {
		return errors.Wrap(err, "error: while creating output.txt file")
	}
	defer output.Close()

	textParts, err := splitTextToParts(text)
	if err != nil {
		return errors.Wrap(err, "error: occurred while splitting the text")
	}

	for fileIndex, textPart := range textParts {
		fileName := fmt.Sprintf("%v.ogg", fileIndex)
		fmt.Println()
		err := c.doRequest(textPart, fileName)

		if err != nil {
			return err
		}
		_, err = output.WriteString(fmt.Sprintf("file '%s'\n", fileName))
		if err != nil {
			return err
		}
	}

	if err := c.convertToMP3(text); err != nil {
		return err
	}

	return nil
}

// createFile output file by audio parts
func (c *SpeechKitClient) createFile() (*os.File, error) {
	// check if file exists
	output := path.Join(c.pathToFiles, output)
	var _, err = os.Stat(output)
	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(output)
		if err != nil {
			return nil, err
		}
		return file, err
	}
	return nil, errors.New("error: file already exists")
}

// generateURL prepare url with voice opts
func (c *SpeechKitClient) generateURL(text string) string {
	if c.SpeechParams.speed == 0.0 {
		c.SpeechParams.speed = speechSpeed
	}

	// define sex for russian premium voice
	if c.SpeechParams.voice == "female" {
		c.SpeechParams.voice = "alena"
	} else if c.SpeechParams.voice == "male" {
		c.SpeechParams.voice = "filipp"
	} else {
		// set default
		c.SpeechParams.voice = "filipp"
	}

	if c.SpeechParams.emotion == "" {
		c.SpeechParams.emotion = speechEmotion
	}

	v := url.Values{}
	v.Add("text", text)
	v.Add("speed", fmt.Sprintf("%.2f", c.SpeechParams.speed))
	v.Add("emotion", c.SpeechParams.emotion)
	v.Add("voice", c.SpeechParams.voice)
	v.Add("lang", speechLanguage)
	v.Add("format", speechFormat)
	return v.Encode()
}

// doRequest make request and save content in 'oggopus' format
func (c *SpeechKitClient) doRequest(text, fileName string) error {
	body := strings.NewReader(c.generateURL(text))
	req, err := http.NewRequest(http.MethodPost, URL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Add("Authorization", fmt.Sprintf("Api-Key %s", c.APIParams.APIKey))

	response, err := c.APIParams.Client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("error: api occurred with status: %v", response.StatusCode))
	}

	fullFilePath := path.Join(c.pathToFiles, fileName)
	outputFile, err := os.Create(fullFilePath)
	if err != nil {
		return errors.Wrap(err, "error: occurred while creating audio file")
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, response.Body)
	if err != nil {
		return errors.Wrap(err, "error: occurred while copying response to file")
	}
	return nil
}

func (c *SpeechKitClient) convertToMP3(text string) error {
	var bound int
	pathToOutFile := path.Join(c.pathToFiles, output)

	if len(text) > 20 {
		bound = 20
	} else {
		bound = len(text)
	}

	mp3FileName := strings.Map(removeNonUTF, fmt.Sprintf("%s.mp3", text[:bound]))
	pathToMP3 := path.Join(c.pathToFiles, mp3FileName)

	fmt.Println(pathToOutFile, pathToMP3)

	cmd := exec.Command(
		"ffmpeg", "-y", "-f", "concat", "-safe", "0", "-i", pathToOutFile, "-vn", "-ar", "44100", "-ac", "2", "-ab", "128k", "-f", "mp3", pathToMP3,
	)

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "error: occurred while generating mp3 from ogg parts")
	}
	return nil
}

// splitTextToParts split text into small parts with size textMaxLen
func splitTextToParts(text string) ([]string, error) {
	if text == "" {
		return nil, errors.New("error: empty txt provided")
	}

	var (
		start, stop int
		textSlice   []string
	)
	for start, stop = 0, textMaxLen; stop < len(text); start, stop = stop, stop+textMaxLen {
		var finalString string
		for _, word := range strings.Split(text[start:stop], " ") {
			s := strings.Map(removeNonUTF, word)
			finalString = finalString + s + " "
		}
		textSlice = append(textSlice, finalString)
	}
	textSlice = append(textSlice, text[start:])
	return textSlice, nil
}

// removeNonUTF remove from string non UTF characters
func removeNonUTF(r rune) rune {
	if r == utf8.RuneError {
		return -1
	}
	return r
}

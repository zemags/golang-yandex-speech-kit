package speechkit

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestUrl(t *testing.T) {
	assert.Equal(
		t,
		"https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize",
		URL,
	)
}

func TestNewSpeechKitClient(t *testing.T) {
	actual := NewSpeechKitClient(APIParams{}, SpeechParams{})
	expected := &SpeechKitClient{APIParams{}, SpeechParams{}}
	assert.Equal(t, actual, expected)
}

func TestRemoveNonUTF(t *testing.T) {
	testString := "Lorem Ipsum is simply dum�y"
	actual := strings.Map(removeNonUTF, testString)
	expected := "Lorem Ipsum is simply dumy"
	assert.Equal(t, actual, expected)
}

func TestSplitTextToParts(t *testing.T) {
	actual, err := splitTextToParts("")
	assert.Nil(t, actual)
	assert.Error(t, err)

	textMaxLen = 10
	actual, err = splitTextToParts("Lorem Ipsum is simply dummy.")
	assert.NoError(t, err)
	assert.Equal(t, []string{"Lorem Ipsu ", "m is simpl ", "y dummy."}, actual)
}

func TestConvertToMP3(t *testing.T) {
	currentDir, _ := os.Getwd()
	pathToMp3 := path.Join(currentDir, "data")
	text := "Мгновенно воцарилась глубокая тишина"
	client := SpeechKitClient{
		APIParams{},
		SpeechParams{
			pathToFiles: pathToMp3,
		},
	}
	err := client.convertToMP3(text)
	assert.NoError(t, err)

	mp3FileName := strings.Map(removeNonUTF, fmt.Sprintf("%s.mp3", text[:20]))
	assert.FileExists(t, path.Join(pathToMp3, mp3FileName))

	client = SpeechKitClient{
		APIParams{},
		SpeechParams{
			pathToFiles: path.Join(currentDir, "not_exist_folder"),
		},
	}
	err = client.convertToMP3(text)
	assert.Error(t, err)
}

func TestGenerateURL(t *testing.T) {
	text := "Lorem Ipsum is simply dummy."
	client := SpeechKitClient{
		APIParams{},
		SpeechParams{
			speed:   0.0,
			emotion: "neutral",
			voice:   "female",
		},
	}
	actual := client.generateURL(text)

	expected := "emotion=neutral&format=oggopus&lang=ru-RU&speed=1.00&text=Lorem+Ipsum+is+simply+dummy.&voice=alena"

	assert.Equal(t, actual, expected)

	client.SpeechParams.voice = "male"
	actual = client.generateURL(text)
	expected = "emotion=neutral&format=oggopus&lang=ru-RU&speed=1.00&text=Lorem+Ipsum+is+simply+dummy.&voice=filipp"
	assert.Equal(t, actual, expected)

	client.SpeechParams.voice = ""
	actual = client.generateURL(text)
	expected = "emotion=neutral&format=oggopus&lang=ru-RU&speed=1.00&text=Lorem+Ipsum+is+simply+dummy.&voice=filipp"
	assert.Equal(t, actual, expected)

}

func TestCreateFile(t *testing.T) {
	currentDir, _ := os.Getwd()
	pathToExistFile := path.Join(currentDir, "data")
	client := SpeechKitClient{
		APIParams{},
		SpeechParams{
			pathToFiles: pathToExistFile,
		},
	}
	file, err := client.createFile()
	assert.Error(t, err)
	assert.Nil(t, file)

	output = "new_test_file.txt"
	pathToNewFile := path.Join(currentDir, "data")
	client = SpeechKitClient{
		APIParams{},
		SpeechParams{
			pathToFiles: pathToNewFile,
		},
	}
	file, err = client.createFile()
	assert.NoError(t, err)
	assert.NotNil(t, file)
	os.Remove(path.Join(pathToNewFile, output))
}

func TestDoRequest(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	APIKey, exists := os.LookupEnv("API_KEY")
	if !exists {
		APIKey = os.Getenv("API_KEY")
	}

	currentDir, _ := os.Getwd()
	pathToFiles := path.Join(currentDir, "temp")
	_ = os.Mkdir(pathToFiles, 0755)
	text := "Мгновенно воцарилась глубокая тишина"
	client := SpeechKitClient{
		APIParams{
			Client: &http.Client{},
			APIKey: APIKey,
		},
		SpeechParams{
			pathToFiles: pathToFiles,
		},
	}
	actual := client.doRequest(text, "1.ogg")
	assert.Nil(t, actual)
	assert.FileExists(t, path.Join(pathToFiles, "1.ogg"))
	os.RemoveAll(pathToFiles)

	// test errors
	client = SpeechKitClient{
		APIParams{
			Client: &http.Client{},
			APIKey: "invalid-api-key",
		},
		SpeechParams{
			pathToFiles: pathToFiles,
		},
	}
	err := client.doRequest(text, "1.ogg")
	assert.EqualError(
		t, err, "error: api occurred with status: 401",
	)

	client.pathToFiles = "invalid path to folder"
	client.APIKey = APIKey
	err = client.doRequest(text, "1.ogg")
	assert.EqualError(
		t, err, "error: occurred while creating audio file: open invalid path to folder/1.ogg: no such file or directory",
	)
}

func TestCreateAudio(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	APIKey, exists := os.LookupEnv("API_KEY")
	if !exists {
		APIKey = os.Getenv("API_KEY")
	}

	textMaxLen = 2000

	currentDir, _ := os.Getwd()
	pathToFiles := path.Join(currentDir, "temp")
	_ = os.Mkdir(pathToFiles, 0755)

	client := SpeechKitClient{
		APIParams{
			Client: &http.Client{},
			APIKey: APIKey,
		},
		SpeechParams{
			pathToFiles: pathToFiles,
		},
	}

	pathToSampleText := path.Join(currentDir, "data", "sample_text.txt")
	text, err := ioutil.ReadFile(pathToSampleText)
	if err != nil {
		log.Fatal("Error occurred while opening sample_text.txt")
	}

	err = client.CreateAudio(string(text))
	assert.NoError(t, err)

	os.RemoveAll(pathToFiles)

	err = client.CreateAudio("")
	assert.Error(t, err)

}

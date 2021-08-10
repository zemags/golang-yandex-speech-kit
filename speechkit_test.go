package speechkit

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

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
			text:        text,
		},
	}
	err := client.convertToMP3()
	assert.NoError(t, err)

	mp3FileName := strings.Map(removeNonUTF, fmt.Sprintf("%s.mp3", text[:20]))
	assert.FileExists(t, path.Join(pathToMp3, mp3FileName))

	client = SpeechKitClient{
		APIParams{},
		SpeechParams{
			pathToFiles: path.Join(currentDir, "not_exist_folder"),
			text:        text,
		},
	}
	err = client.convertToMP3()
	assert.Error(t, err)
}

func TestGenerateURL(t *testing.T) {
	client := SpeechKitClient{
		APIParams{},
		SpeechParams{
			text:    "Lorem Ipsum is simply dummy.",
			speed:   0.0,
			emotion: "neutral",
			voice:   "female",
		},
	}
	actual := client.generateURL("Lorem Ipsum is simply dummy.")

	expected := "emotion=neutral&format=oggopus&lang=ru-RU&speed=1.00&text=Lorem+Ipsum+is+simply+dummy.&voice=alena"

	assert.Equal(t, actual, expected)

	client.SpeechParams.voice = "male"
	actual = client.generateURL("Lorem Ipsum is simply dummy.")
	expected = "emotion=neutral&format=oggopus&lang=ru-RU&speed=1.00&text=Lorem+Ipsum+is+simply+dummy.&voice=filipp"
	assert.Equal(t, actual, expected)

	client.SpeechParams.voice = ""
	actual = client.generateURL("Lorem Ipsum is simply dummy.")
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

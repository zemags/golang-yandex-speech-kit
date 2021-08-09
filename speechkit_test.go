package speechkit

import (
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

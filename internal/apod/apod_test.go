package apod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatStatus_NormalTitle(t *testing.T) {
	a := &APOD{Title: "Pillars of Creation"}
	assert.Equal(t, "🔭 APOD: \"Pillars of Creation\"", a.FormatStatus())
}

func TestFormatStatus_LongTitle(t *testing.T) {
	a := &APOD{Title: "The Magnificent and Extraordinary Tail of Comet McNaught"}
	s := a.FormatStatus()
	assert.Contains(t, s, "…")
	assert.LessOrEqual(t, len([]rune(s)), 60)
}

func TestFormatStatus_EmptyTitle(t *testing.T) {
	a := &APOD{Title: ""}
	assert.Equal(t, "", a.FormatStatus())
}

func TestFormatStatus_Nil(t *testing.T) {
	var a *APOD
	assert.Equal(t, "", a.FormatStatus())
}

func TestParseResponse_Valid(t *testing.T) {
	data := []byte(`{"title":"Test Image","date":"2026-04-03","media_type":"image"}`)
	result, err := parseResponse(data)
	assert.NoError(t, err)
	assert.Equal(t, "Test Image", result.Title)
	assert.Equal(t, "2026-04-03", result.Date)
	assert.Equal(t, "image", result.MediaType)
}

func TestParseResponse_InvalidJSON(t *testing.T) {
	_, err := parseResponse([]byte("not json"))
	assert.Error(t, err)
}

func TestFormatStatus_ExactlyFortyChars(t *testing.T) {
	a := &APOD{Title: "1234567890123456789012345678901234567890"}
	assert.Equal(t, "🔭 APOD: \"1234567890123456789012345678901234567890\"", a.FormatStatus())
}

func TestFormatStatus_FortyOneCharsTruncated(t *testing.T) {
	a := &APOD{Title: "12345678901234567890123456789012345678901"}
	assert.Contains(t, a.FormatStatus(), "…")
}

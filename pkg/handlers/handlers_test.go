package handlers

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertTopic_InvalidPayload(t *testing.T) {
	body := "{'name': 'foo'}"
	request, err := http.NewRequest(http.MethodGet, "", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatal(err)
	}
	response := InsertTopic(request)
	expected := Response{
		StatusCode: 400,
		Body: ErrorBody{ErrorMsg: ErrInvalidPayload},
	}
	assert.Equal(t, expected, response)
}

func TestInsertDeleteTopic_ValidPayload(t *testing.T) {
	body := `{"userId": "1d7ee7f0-36f5-4e33-a766-26981e62d9cf", "title": "testTopic_InsertValid"}`
	request, err := http.NewRequest(http.MethodGet, "", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatal(err)
	}
	response := InsertTopic(request)
	expected := Response{
		StatusCode: 200,
		Body: nil,
	}
	assert.Equal(t, expected, response)

	body = `{"userId": "1d7ee7f0-36f5-4e33-a766-26981e62d9cf", "title": "testTopic_InsertValid"}`
	request, err = http.NewRequest(http.MethodDelete, "", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatal(err)
	}
	response = DeleteTopic(request)
	expected = Response{
		StatusCode: 200,
		Body: nil,
	}
	assert.Equal(t, expected, response)
}

func TestInsertNote(t *testing.T) {
	body := `{"userId": "1d7ee7f0-36f5-4e33-a766-26981e62d9cf", "title": "testInsert"}`
	request, err := http.NewRequest(http.MethodGet, "", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatal(err)
	}
	response := InsertTopic(request)
	expected := Response{
		StatusCode: 200,
		Body: nil,
	}
	assert.Equal(t, expected, response)

	body = `{"userId": "1d7ee7f0-36f5-4e33-a766-26981e62d9cf", "title": "testInsert", "note": {"title": "note_title", "content": "some test content"}}`
	request, err = http.NewRequest(http.MethodGet, "", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatal(err)
	}
	response = InsertNote(request)
	expected = Response{
		StatusCode: 200,
		Body: nil,
	}
	assert.Equal(t, expected, response)
}
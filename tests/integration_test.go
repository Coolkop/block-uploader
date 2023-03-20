//go:build integration

package tests_test

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"file-storage/internal/container"
)

func TestIntegration(t *testing.T) {
	fileData := []byte(gofakeit.Sentence(1000))

	c, err := container.New()
	if !assert.NoError(t, err) {
		return
	}

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		part, err := writer.CreateFormFile("file", "somefile")
		if err != nil {
			t.Error(err)
		}

		_, err = part.Write(fileData)
		if err != nil {
			t.Error(err)
		}
	}()

	request := httptest.NewRequest("PUT", "/", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()

	// PUT file
	c.HttpHandler.ServeHTTP(response, request)
	if !assert.Equal(t, 200, response.Code) {
		return
	}

	s := struct {
		ID string `json:"id"`
	}{}

	err = json.Unmarshal(response.Body.Bytes(), &s)
	if !assert.NoError(t, err) {
		return
	}

	createdFileID, err := uuid.Parse(s.ID)
	if !assert.NoError(t, err) {
		return
	}

	request = httptest.NewRequest("GET", fmt.Sprintf("/%s", createdFileID.String()), nil)
	response = httptest.NewRecorder()

	// GET file
	c.HttpHandler.ServeHTTP(response, request)
	if !assert.Equal(t, 200, response.Code) {
		return
	}

	assert.Equal(t, fileData, response.Body.Bytes())
}

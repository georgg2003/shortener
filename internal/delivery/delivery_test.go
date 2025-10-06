package delivery_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/georgg2003/shortener/internal/delivery"
	"github.com/georgg2003/shortener/internal/repository"
	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	name                      string
	testURL                   string
	shortenMethod             string
	shortenExpectedStatusCode int
	processMethod             string
	processExpectedStatusCode int
	wrongShortID              bool
}

func TestDelivery(t *testing.T) {
	repo := repository.New()
	usecase := usecase.New(repo)
	delivery := delivery.New(usecase)

	testCases := []TestCase{
		{
			name:                      "all ok",
			testURL:                   "https://calendar.mail.ru",
			shortenMethod:             http.MethodPost,
			shortenExpectedStatusCode: http.StatusCreated,
			processMethod:             http.MethodGet,
			processExpectedStatusCode: http.StatusTemporaryRedirect,
			wrongShortID:              false,
		},
		{
			name:                      "wrong shorten method",
			testURL:                   "https://calendar.mail.ru",
			shortenMethod:             http.MethodGet,
			shortenExpectedStatusCode: http.StatusMethodNotAllowed,
			processMethod:             http.MethodGet,
			processExpectedStatusCode: http.StatusTemporaryRedirect,
			wrongShortID:              false,
		},
		{
			name:                      "wrong proccess method",
			testURL:                   "https://calendar.mail.ru",
			shortenMethod:             http.MethodPost,
			shortenExpectedStatusCode: http.StatusCreated,
			processMethod:             http.MethodPost,
			processExpectedStatusCode: http.StatusMethodNotAllowed,
			wrongShortID:              false,
		},
		{
			name:                      "wrong short id",
			testURL:                   "https://calendar.mail.ru",
			shortenMethod:             http.MethodPost,
			shortenExpectedStatusCode: http.StatusCreated,
			processMethod:             http.MethodGet,
			processExpectedStatusCode: http.StatusNotFound,
			wrongShortID:              true,
		},
	}

	for _, test := range testCases {
		shortURL := ""
		t.Run(test.name, func(t *testing.T) {
			t.Run("shorten", func(t *testing.T) {
				body := bytes.Buffer{}
				body.WriteString(test.testURL)

				request := httptest.NewRequest(test.shortenMethod, "/", &body)
				recorder := httptest.NewRecorder()

				delivery.ShortenURL(recorder, request)

				res := recorder.Result()
				defer res.Body.Close()
				assert.Equal(t, test.shortenExpectedStatusCode, res.StatusCode)
				assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))

				if res.StatusCode == http.StatusCreated {
					b, err := io.ReadAll(res.Body)
					require.Nil(t, err)
					shortURL = string(b)
					t.Log(shortURL)

					require.NotEmpty(t, shortURL)
				}
			})

			t.Run("processShortURL", func(t *testing.T) {
				if shortURL == "" {
					t.Skip("Skipping test because shortURL was not get")
				}
				if test.wrongShortID {
					shortURL = "http://localhost:8080/wrong_id"
				}

				t.Logf("Making request to %v", shortURL)
				request := httptest.NewRequest(test.processMethod, shortURL, nil)
				shortID := strings.Split(shortURL, "/")[3]
				request.SetPathValue("id", shortID)
				recorder := httptest.NewRecorder()

				delivery.ProcessShortURL(recorder, request)

				res := recorder.Result()
				defer res.Body.Close()
				assert.Equal(t, test.processExpectedStatusCode, res.StatusCode)
				if res.StatusCode == http.StatusTemporaryRedirect {
					assert.Equal(t, test.testURL, res.Header.Get("Location"))
				}
			})
		})
	}
}

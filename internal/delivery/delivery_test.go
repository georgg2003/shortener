package delivery_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/delivery"
	"github.com/georgg2003/shortener/internal/repository"
	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/go-resty/resty/v2"
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

func testRequest(
	t *testing.T,
	ts *httptest.Server,
	method, url string,
	body interface{},
) *resty.Response {
	req := resty.New().SetRedirectPolicy(resty.NoRedirectPolicy()).R()
	req.Method = method
	req.URL = url
	req.SetBody(body)

	resp, err := req.Send()

	if !errors.Is(err, resty.ErrAutoRedirectDisabled) {
		require.NoError(t, err, "error making HTTP request")
	}

	return resp
}

func TestDelivery(t *testing.T) {
	ts := httptest.NewServer(nil)

	conf := &config.Config{
		BaseURL:    ts.URL,
		ListenAddr: ts.URL,
	}
	repo := repository.New()
	usecase := usecase.New(repo, conf)
	delivery := delivery.New(usecase)

	r := delivery.GetNewRouter()

	ts.Config.Handler = r

	defer ts.Close()

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
				res := testRequest(t, ts, test.shortenMethod, ts.URL+"/", test.testURL)
				statusCode := res.StatusCode()
				assert.Equal(t, test.shortenExpectedStatusCode, statusCode)

				if statusCode == http.StatusCreated {
					b := res.Body()
					shortURL = string(b)

					require.NotEmpty(t, shortURL)
				}
			})

			t.Run("processShortURL", func(t *testing.T) {
				if shortURL == "" {
					t.Skip("Skipping test because shortURL was not get")
				}
				if test.wrongShortID {
					shortURL = "https://google.com/wrong_id"
				}

				t.Logf("Making request to %v", shortURL)

				res := testRequest(t, ts, test.processMethod, shortURL, nil)
				statusCode := res.StatusCode()

				assert.Equal(t, test.processExpectedStatusCode, statusCode)
				if statusCode == http.StatusTemporaryRedirect {
					assert.Equal(t, test.testURL, res.Header().Get("Location"))
				}
			})
		})
	}
}

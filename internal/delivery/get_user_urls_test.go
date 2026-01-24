package delivery_test

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/delivery"
	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/internal/repository/db/mock"
	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const testURL = "https://calendar.mail.ru"
const errHTTPReqText = "error making HTTP request"

func newTestServer(t *testing.T) (*httptest.Server, *mock.MockRepository) {
	ts := httptest.NewServer(nil)
	logger := logrus.New()

	conf := &config.Config{
		BaseURL:     ts.URL,
		ListenAddr:  ts.URL,
		DataBaseDSN: "127.0.0.1:5432",
	}

	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)

	usecase := usecase.New(repo, conf, logger)
	delivery := delivery.New(usecase, logger, conf)

	r := delivery.GetNewRouter()

	ts.Config.Handler = r

	return ts, repo
}

func TestGetUserUrls(t *testing.T) {
	ts, repo := newTestServer(t)
	defer ts.Close()

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := resty.New()
	client.SetCookieJar(jar)

	userID := int64(1)
	repo.EXPECT().NewUser(gomock.Any()).Return(userID, nil)
	repo.EXPECT().NewShortURL(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	res, err := client.R().SetBody(testURL).Post(ts.URL + "/")
	require.NoError(t, err, errHTTPReqText)
	assert.Equal(t, http.StatusCreated, res.StatusCode())

	hasAccessToken := false
	for _, cookie := range res.Cookies() {
		if cookie.Name == "access_token" {
			hasAccessToken = true
			break
		}
	}
	require.True(t, hasAccessToken)

	repo.EXPECT().GetUserURLs(gomock.Any(), userID).Return([]models.URLEntity{{
		OriginalURL: testURL,
		ShortID:     "adxzdsadsa",
	}}, nil)

	req, err := client.R().Get(ts.URL + "/api/user/urls")
	require.NoError(t, err, errHTTPReqText)

	var jsonBody models.APIUserURLsResponse
	json.Unmarshal(req.Body(), &jsonBody)

	assert.Len(t, jsonBody, 1)
	url := jsonBody[0]
	assert.Equal(t, url.OriginalURL, testURL)
}

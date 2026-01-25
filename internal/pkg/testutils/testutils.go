package testutils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/delivery"
	"github.com/georgg2003/shortener/internal/repository/db/mock"
	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/georgg2003/shortener/pkg/jwthelper"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const TestUserID = 1

type DeliveryTestCase struct {
	Name       string
	Method     string
	Path       string
	Body       any
	StatusCode int
	Response   []byte
	MockFunc   func()
}

type TestServer struct {
	ts          *httptest.Server
	accessToken string

	Repo *mock.MockRepository
}

func (ts *TestServer) MakeRequest(t *testing.T, method, path string, body any) *resty.Response {
	req := resty.New().
		SetRedirectPolicy(resty.NoRedirectPolicy()).
		R().
		SetCookie(&http.Cookie{
			Name:  "access_token",
			Value: ts.accessToken,
		})
	req.Method = method
	req.URL = ts.ts.URL + path
	req.SetBody(body)

	resp, err := req.Send()

	if !errors.Is(err, resty.ErrAutoRedirectDisabled) {
		require.NoError(t, err, "error making HTTP request")
	}

	return resp
}

func (ts *TestServer) RunTestCase(tc DeliveryTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		if tc.MockFunc != nil {
			tc.MockFunc()
		}
		resp := ts.MakeRequest(t, tc.Method, tc.Path, tc.Body)

		assert.Equal(t, tc.StatusCode, resp.StatusCode())
		assert.Equal(t, tc.Response, resp.Body())
	}
}

func (ts *TestServer) Close() {
	ts.ts.Close()
}

func NewTestServer(t *testing.T) *TestServer {
	ts := httptest.NewServer(nil)
	logger := logrus.New()

	cfg := config.New()
	cfg.BaseURL = ts.URL
	cfg.ListenAddr = ts.URL
	cfg.DataBaseDSN = "some fake dsn"

	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)

	usecase := usecase.New(repo, cfg, logger)
	delivery := delivery.New(usecase, logger, cfg)

	r := delivery.GetNewRouter()

	accessToken, err := jwthelper.New([]byte(cfg.JWTSecretKey)).NewAccessToken(TestUserID)
	require.NoError(t, err)

	ts.Config.Handler = r

	return &TestServer{
		ts:          ts,
		accessToken: accessToken,
		Repo:        repo,
	}
}

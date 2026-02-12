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
	"github.com/georgg2003/shortener/pkg/utils"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	TestUserID      = 1
	TestShortID     = "01234567"
	TestOriginalURL = "https://calendar.mail.ru"
)

var ErrSomeError = errors.New("some error")

type TestReporter interface {
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	FailNow()
}

type Common struct {
	Name     string
	Method   string
	Path     string
	Body     any
	MockFunc func(t TestReporter, ts *TestServer)
}

type DeliveryTestCase struct {
	Common
	StatusCode   int
	Response     []byte
	CheckReponse func(t *testing.T, resp *resty.Response)
}

type DeliveryBenchmark struct {
	Common
}

type TestServer struct {
	ts          *httptest.Server
	accessToken string

	Repo *mock.MockRepository
}

func (ts *TestServer) makeRequest(t TestReporter, method, path string, body any) *resty.Response {
	req := resty.New().
		SetRedirectPolicy(resty.NoRedirectPolicy()).
		R().
		SetCookie(&http.Cookie{
			Name:  "access_token",
			Value: ts.accessToken,
		})
	req.Method = method
	req.URL = ts.MakeAbsoluteURL(path)
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
			tc.MockFunc(t, ts)
		}
		resp := ts.makeRequest(t, tc.Method, tc.Path, tc.Body)

		assert.Equal(t, tc.StatusCode, resp.StatusCode())
		assert.Equal(t, tc.Response, resp.Body())

		if tc.CheckReponse != nil {
			tc.CheckReponse(t, resp)
		}
	}
}

func (ts *TestServer) RunBenchmark(tc DeliveryBenchmark) func(b *testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if tc.MockFunc != nil {
				tc.MockFunc(b, ts)
			}
			ts.makeRequest(b, tc.Method, tc.Path, tc.Body)
		}
	}
}

func (ts *TestServer) Close() {
	ts.ts.Close()
}

func (ts *TestServer) MakeAbsoluteURL(path string) string {
	return ts.ts.URL + path
}

func NewTestServer(t TestReporter) *TestServer {
	ts := httptest.NewServer(nil)
	logger := logrus.New()

	cfg := config.New()
	cfg.BaseURL = ts.URL
	cfg.ListenAddr = ts.URL
	cfg.DataBaseDSN = "some fake dsn"

	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)

	gen := utils.StubGenerator{}
	usecase := usecase.New(repo, cfg, logger, usecase.WithBase62Generator(gen))
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

package delivery_test

import (
	"encoding/json"
	"net/http"
	"slices"
	"testing"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/internal/pkg/testutils"
	"github.com/georgg2003/shortener/pkg/postgres"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const shortenAPIBatchPath = "/api/shorten/batch"
const shortenAPIPath = "/api/shorten"
const shortenPath = "/"

var (
	successCommon = testutils.Common{
		Name:   "success",
		Method: http.MethodPost,
		Path:   shortenAPIBatchPath,
		Body: models.APIShortenURLBatchRequest{
			models.APIShortenURLBatchRequestRecord{
				CorrelationID: "1",
				OriginalURL:   testutils.TestOriginalURL,
			},
		},
		MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
			shortURL := ts.MakeAbsoluteURL("/" + testutils.TestShortID)

			entity := &models.URLEntity{
				CorrelationID: "1",
				OriginalURL:   testutils.TestOriginalURL,
				ShortID:       testutils.TestShortID,
				ShortURL:      shortURL,
			}

			ts.Repo.EXPECT().
				NewShortURLBatch(gomock.Any(), gomock.Eq([]*models.URLEntity{entity}))
		},
	}
)

func TestShortenURL(t *testing.T) {
	server := testutils.NewTestServer(t)
	defer server.Close()

	shortURL := server.MakeAbsoluteURL("/" + testutils.TestShortID)

	testCases := []testutils.DeliveryTestCase{
		{
			Common: testutils.Common{
				Name:   "success",
				Method: http.MethodPost,
				Path:   shortenPath,
				Body:   testutils.TestOriginalURL,
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					ts.Repo.EXPECT().
						NewShortURL(gomock.Any(), testutils.TestOriginalURL, testutils.TestShortID).
						Return(nil)
				},
			},
			StatusCode: http.StatusCreated,
			Response:   []byte(shortURL),
		},
		{
			Common: testutils.Common{
				Name:   "invalid body",
				Method: http.MethodPost,
				Path:   shortenPath,
			},
			StatusCode: http.StatusBadRequest,
			Response:   []byte("Invalid request body\n"),
		},
		{
			Common: testutils.Common{
				Name:   "invalid scheme",
				Method: http.MethodPost,
				Path:   shortenPath,
				Body: models.APIShortenURLRequest{
					URL: "invalid scheme",
				},
			},
			StatusCode: http.StatusBadRequest,
			Response:   []byte("Invalid url\n"),
		},
		{
			Common: testutils.Common{
				Name:   "parse error",
				Method: http.MethodPost,
				Path:   shortenPath,
				Body: models.APIShortenURLRequest{
					URL: "https://dsa dsa dsa",
				},
			},
			StatusCode: http.StatusBadRequest,
			Response:   []byte("Invalid url\n"),
		},
		{
			Common: testutils.Common{
				Name:   "invalid hostname",
				Method: http.MethodPost,
				Path:   shortenPath,
				Body: models.APIShortenURLRequest{
					URL: "https://",
				},
			},
			StatusCode: http.StatusBadRequest,
			Response:   []byte("Invalid url\n"),
		},
		{
			Common: testutils.Common{
				Name:   "internal error",
				Method: http.MethodPost,
				Path:   shortenPath,
				Body:   testutils.TestOriginalURL,
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					ts.Repo.EXPECT().
						NewShortURL(gomock.Any(), testutils.TestOriginalURL, testutils.TestShortID).
						Return(testutils.ErrSomeError)
				},
			},
			StatusCode: http.StatusInternalServerError,
			Response:   []byte("Internal Error\n"),
		},
		{
			Common: testutils.Common{
				Name:   "url already shortened",
				Method: http.MethodPost,
				Path:   shortenPath,
				Body:   testutils.TestOriginalURL,
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					call := ts.Repo.EXPECT().
						NewShortURL(gomock.Any(), testutils.TestOriginalURL, testutils.TestShortID).
						Return(postgres.ErrUniqueViolation)

					ts.Repo.EXPECT().GetShortID(gomock.Any(), testutils.TestOriginalURL).
						After(call).Return(testutils.TestShortID, nil)
				},
			},
			StatusCode: http.StatusConflict,
			Response:   []byte(shortURL),
		},
		{
			Common: testutils.Common{
				Name:   "url already shortened, get short id error",
				Method: http.MethodPost,
				Path:   shortenPath,
				Body:   testutils.TestOriginalURL,
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					call := ts.Repo.EXPECT().
						NewShortURL(gomock.Any(), testutils.TestOriginalURL, testutils.TestShortID).
						Return(postgres.ErrUniqueViolation)

					ts.Repo.EXPECT().GetShortID(gomock.Any(), testutils.TestOriginalURL).
						After(call).Return("", testutils.ErrSomeError)
				},
			},
			StatusCode: http.StatusInternalServerError,
			Response:   []byte("Internal Error\n"),
		},
	}
	for tc := range slices.Values(testCases) {
		t.Run(tc.Name, server.RunTestCase(tc))
	}
}

func TestAPIShortenURL(t *testing.T) {
	server := testutils.NewTestServer(t)
	defer server.Close()

	shortURL := server.MakeAbsoluteURL("/" + testutils.TestShortID)
	succResp := models.APIShortenURLResponse{
		Result: shortURL,
	}

	bytesResp, err := json.Marshal(succResp)
	bytesResp = append(bytesResp, '\n')
	require.NoError(t, err)

	testCases := []testutils.DeliveryTestCase{
		{
			Common: testutils.Common{
				Name:   "success",
				Method: http.MethodPost,
				Path:   shortenAPIPath,
				Body: models.APIShortenURLRequest{
					URL: testutils.TestOriginalURL,
				},
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					ts.Repo.EXPECT().
						NewShortURL(gomock.Any(), testutils.TestOriginalURL, testutils.TestShortID).
						Return(nil)
				},
			},
			StatusCode: http.StatusCreated,
			Response:   bytesResp,
		},
		{
			Common: testutils.Common{
				Name:   "invalid json",
				Method: http.MethodPost,
				Path:   shortenAPIPath,
				Body:   "invalid json",
			},
			StatusCode: http.StatusBadRequest,
			Response:   []byte("failed to decode body: invalid character 'i' looking for beginning of value\n"),
		},
		{
			Common: testutils.Common{
				Name:   "invalid scheme",
				Method: http.MethodPost,
				Path:   shortenAPIPath,
				Body: models.APIShortenURLRequest{
					URL: "invalid scheme",
				},
			},
			StatusCode: http.StatusBadRequest,
			Response:   []byte("invalid scheme\n"),
		},
		{
			Common: testutils.Common{
				Name:   "parse error",
				Method: http.MethodPost,
				Path:   shortenAPIPath,
				Body: models.APIShortenURLRequest{
					URL: "https://dsa dsa dsa",
				},
			},
			StatusCode: http.StatusBadRequest,
			Response:   []byte(`invalid url: parse "https://dsa dsa dsa": invalid character " " in host name` + "\n"),
		},
		{
			Common: testutils.Common{
				Name:   "invalid hostname",
				Method: http.MethodPost,
				Path:   shortenAPIPath,
				Body: models.APIShortenURLRequest{
					URL: "https://",
				},
			},
			StatusCode: http.StatusBadRequest,
			Response:   []byte("invalid hostname\n"),
		},
		{
			Common: testutils.Common{
				Name:   "internal error",
				Method: http.MethodPost,
				Path:   shortenAPIPath,
				Body: models.APIShortenURLRequest{
					URL: testutils.TestOriginalURL,
				},
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					ts.Repo.EXPECT().
						NewShortURL(gomock.Any(), testutils.TestOriginalURL, testutils.TestShortID).
						Return(testutils.ErrSomeError)
				},
			},
			StatusCode: http.StatusInternalServerError,
			Response:   []byte("Internal error\n"),
		},
		{
			Common: testutils.Common{
				Name:   "url already shortened",
				Method: http.MethodPost,
				Path:   shortenAPIPath,
				Body: models.APIShortenURLRequest{
					URL: testutils.TestOriginalURL,
				},
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					call := ts.Repo.EXPECT().
						NewShortURL(gomock.Any(), testutils.TestOriginalURL, testutils.TestShortID).
						Return(postgres.ErrUniqueViolation)

					ts.Repo.EXPECT().GetShortID(gomock.Any(), testutils.TestOriginalURL).
						After(call).Return(testutils.TestShortID, nil)
				},
			},
			StatusCode: http.StatusConflict,
			Response:   bytesResp,
		},
		{
			Common: testutils.Common{
				Name:   "url already shortened, get short id error",
				Method: http.MethodPost,
				Path:   shortenAPIPath,
				Body: models.APIShortenURLRequest{
					URL: testutils.TestOriginalURL,
				},
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					call := ts.Repo.EXPECT().
						NewShortURL(gomock.Any(), testutils.TestOriginalURL, testutils.TestShortID).
						Return(postgres.ErrUniqueViolation)

					ts.Repo.EXPECT().GetShortID(gomock.Any(), testutils.TestOriginalURL).
						After(call).Return("", testutils.ErrSomeError)
				},
			},
			StatusCode: http.StatusInternalServerError,
			Response:   []byte("Internal error\n"),
		},
	}
	for tc := range slices.Values(testCases) {
		t.Run(tc.Name, server.RunTestCase(tc))
	}
}

func TestAPIShortenURLBatch(t *testing.T) {
	server := testutils.NewTestServer(t)
	defer server.Close()

	shortURL := server.MakeAbsoluteURL("/" + testutils.TestShortID)
	succResp := models.APIShortenURLBatchResponse{
		models.APIShortenURLBatchResponseRecord{
			CorrelationID: "1",
			ShortURL:      shortURL,
		},
	}
	entity := &models.URLEntity{
		CorrelationID: "1",
		OriginalURL:   testutils.TestOriginalURL,
		ShortID:       testutils.TestShortID,
		ShortURL:      shortURL,
	}
	bytesResp, err := json.Marshal(succResp)
	bytesResp = append(bytesResp, '\n')
	require.NoError(t, err)

	testCases := []testutils.DeliveryTestCase{
		{
			Common:     successCommon,
			StatusCode: http.StatusCreated,
			Response:   bytesResp,
		},
		{
			Common: testutils.Common{
				Name:   "bad request",
				Method: http.MethodPost,
				Path:   shortenAPIBatchPath,
				Body:   "12321",
			},
			StatusCode: http.StatusBadRequest,
			Response:   []byte("failed to decode body\n"),
		},
		{
			Common: testutils.Common{
				Name:   "with existing urls",
				Method: http.MethodPost,
				Path:   shortenAPIBatchPath,
				Body: models.APIShortenURLBatchRequest{
					models.APIShortenURLBatchRequestRecord{
						CorrelationID: "1",
						OriginalURL:   testutils.TestOriginalURL,
					},
				},
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					exp1 := ts.Repo.EXPECT().
						NewShortURLBatch(gomock.Any(), gomock.Eq([]*models.URLEntity{entity})).
						Return(postgres.ErrUniqueViolation)
					ts.Repo.EXPECT().
						GetShortIDsBatch(gomock.Any(), gomock.Eq([]*models.URLEntity{entity})).
						After(exp1).
						Return(map[string]string{testutils.TestOriginalURL: testutils.TestShortID}, nil)
				},
			},
			StatusCode: http.StatusConflict,
			Response:   bytesResp,
		},
	}
	for tc := range slices.Values(testCases) {
		t.Run(tc.Name, server.RunTestCase(tc))
	}
}

func BenchmarkAPIShortenURLBatch(b *testing.B) {
	server := testutils.NewTestServer(b)
	defer server.Close()

	benchmarks := []testutils.DeliveryBenchmark{
		{
			Common: successCommon,
		},
	}

	for bench := range slices.Values(benchmarks) {
		b.Run(bench.Name, server.RunBenchmark(bench))
	}
}

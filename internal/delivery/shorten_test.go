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
			Response:   []byte("failed to decode body"),
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
					exp1 := server.Repo.EXPECT().
						NewShortURLBatch(gomock.Any(), gomock.Eq([]*models.URLEntity{entity})).
						Return(postgres.ErrUniqueViolation)
					server.Repo.EXPECT().
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

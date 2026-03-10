package delivery_test

import (
	"encoding/json"
	"net/http"
	"slices"
	"testing"

	"github.com/georgg2003/shortener/internal/delivery"
	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/internal/pkg/testutils"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const statsPath = "/api/internal/stats"

func TestInternalStats(t *testing.T) {
	server := testutils.NewTestServer(t)
	defer server.Close()

	successBody := models.APIInternalStatsResponse{
		URLsCount:  10,
		UsersCount: 2,
	}

	successBodyBytes, err := json.Marshal(successBody)
	successBodyBytes = append(successBodyBytes, '\n')
	require.NoError(t, err)

	testCases := []testutils.DeliveryTestCase{
		{
			Common: testutils.Common{
				Name:   "success",
				Method: http.MethodGet,
				Path:   statsPath,
				ModifyRequstFunc: func(req *resty.Request) {
					req.SetHeader(delivery.RealIPHeaderName, "127.0.0.1")
				},
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					server.Repo.EXPECT().GetStats(gomock.Any()).Return(models.Stats{
						UsersCount: 2,
						URLsCount:  10,
					}, nil)
				},
			},
			Response:   successBodyBytes,
			StatusCode: http.StatusOK,
		},
		{
			Common: testutils.Common{
				Name:   "forbidden",
				Method: http.MethodGet,
				Path:   statsPath,
				ModifyRequstFunc: func(req *resty.Request) {
					req.SetHeader(delivery.RealIPHeaderName, "127.32.19.15")
				},
			},
			Response:   []byte(""),
			StatusCode: http.StatusForbidden,
		},
		{
			Common: testutils.Common{
				Name:   "db error",
				Method: http.MethodGet,
				Path:   statsPath,
				ModifyRequstFunc: func(req *resty.Request) {
					req.SetHeader(delivery.RealIPHeaderName, "127.0.0.1")
				},
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					server.Repo.EXPECT().GetStats(gomock.Any()).Return(models.Stats{}, testutils.ErrSomeError)
				},
			},
			Response:   []byte("some error\n"),
			StatusCode: http.StatusInternalServerError,
		},
	}
	for tc := range slices.Values(testCases) {
		t.Run(tc.Name, server.RunTestCase(tc))
	}
}

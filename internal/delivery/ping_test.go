package delivery_test

import (
	"net/http"
	"slices"
	"testing"

	"github.com/georgg2003/shortener/internal/pkg/testutils"
	"go.uber.org/mock/gomock"
)

const pingPath = "/ping"

func TestPing(t *testing.T) {
	server := testutils.NewTestServer(t)
	defer server.Close()

	testCases := []testutils.DeliveryTestCase{
		{
			Common: testutils.Common{
				Name:   "success",
				Method: http.MethodGet,
				Path:   pingPath,
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					server.Repo.EXPECT().Ping(gomock.Any())
				},
			},
			StatusCode: http.StatusOK,
			Response:   []byte(""),
		},
		{
			Common: testutils.Common{
				Name:   "fail",
				Method: http.MethodGet,
				Path:   pingPath,
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					server.Repo.EXPECT().Ping(gomock.Any()).Return(testutils.ErrSomeError)
				},
			},
			StatusCode: http.StatusInternalServerError,
			Response:   []byte("ping failed"),
		},
	}
	for tc := range slices.Values(testCases) {
		t.Run(tc.Name, server.RunTestCase(tc))
	}
}

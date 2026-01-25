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
			Name:       "success",
			Method:     http.MethodGet,
			Path:       pingPath,
			StatusCode: http.StatusOK,
			MockFunc: func(t *testing.T) {
				server.Repo.EXPECT().Ping(gomock.Any())
			},
			Response: []byte(""),
		},
		{
			Name:       "fail",
			Method:     http.MethodGet,
			Path:       pingPath,
			StatusCode: http.StatusInternalServerError,
			MockFunc: func(t *testing.T) {
				server.Repo.EXPECT().Ping(gomock.Any()).Return(testutils.ErrSomeError)
			},
			Response: []byte("ping failed\n"),
		},
	}
	for tc := range slices.Values(testCases) {
		t.Run(tc.Name, server.RunTestCase(tc))
	}
}

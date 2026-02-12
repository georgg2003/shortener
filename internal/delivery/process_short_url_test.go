package delivery_test

import (
	"net/http"
	"slices"
	"testing"

	"github.com/georgg2003/shortener/internal/pkg/testutils"
	"github.com/georgg2003/shortener/internal/repository/db"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProcessShortURL(t *testing.T) {
	server := testutils.NewTestServer(t)
	defer server.Close()

	path := "/" + testutils.TestShortID

	testCases := []testutils.DeliveryTestCase{
		{
			Common: testutils.Common{
				Name:   "method not allowed",
				Method: http.MethodPut,
				Path:   path,
			},
			StatusCode: http.StatusMethodNotAllowed,
			Response:   []byte(""),
		},
		{
			Common: testutils.Common{
				Name:   "not found",
				Method: http.MethodGet,
				Path:   "/" + testutils.TestShortID,
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					ts.Repo.EXPECT().GetLongURL(gomock.Any(), testutils.TestShortID).
						Return("", false, db.ErrNotFound)
				},
			},
			StatusCode: http.StatusNotFound,
			Response:   []byte("Link not found\n"),
		},
		{
			Common: testutils.Common{
				Name:   "internal error",
				Method: http.MethodGet,
				Path:   path,
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					ts.Repo.EXPECT().GetLongURL(gomock.Any(), testutils.TestShortID).
						Return("", false, testutils.ErrSomeError)
				},
			},
			StatusCode: http.StatusInternalServerError,
			Response:   []byte("Internal Error\n"),
		},
		{
			Common: testutils.Common{
				Name:   "link deleted",
				Method: http.MethodGet,
				Path:   path,
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					ts.Repo.EXPECT().GetLongURL(gomock.Any(), testutils.TestShortID).
						Return(testutils.TestOriginalURL, true, nil)
				},
			},
			StatusCode: http.StatusGone,
			Response:   []byte("URL was deleted\n"),
		},
		{
			Common: testutils.Common{
				Name:   "success",
				Method: http.MethodGet,
				Path:   path,
				MockFunc: func(t testutils.TestReporter, ts *testutils.TestServer) {
					ts.Repo.EXPECT().GetLongURL(gomock.Any(), testutils.TestShortID).
						Return(testutils.TestOriginalURL, false, nil)
				},
			},
			StatusCode: http.StatusTemporaryRedirect,
			CheckReponse: func(t *testing.T, resp *resty.Response) {
				assert.Equal(t, testutils.TestOriginalURL, resp.Header().Get("Location"))
			},
		},
	}
	for tc := range slices.Values(testCases) {
		t.Run(tc.Name, server.RunTestCase(tc))
	}
}

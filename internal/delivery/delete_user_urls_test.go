package delivery_test

import (
	"net/http"
	"slices"
	"testing"

	"github.com/georgg2003/shortener/internal/pkg/testutils"
)

const deleteUserURLsPath = "/api/user/urls"

func TestDeleteUserURLs(t *testing.T) {
	server := testutils.NewTestServer(t)
	defer server.Close()

	testCases := []testutils.DeliveryTestCase{
		{
			Name:       "success",
			Method:     http.MethodDelete,
			Path:       deleteUserURLsPath,
			Body:       []string{"1"},
			StatusCode: http.StatusAccepted,
			Response:   []byte(""),
		},
		{
			Name:       "bad request",
			Method:     http.MethodDelete,
			Path:       deleteUserURLsPath,
			Body:       "12321",
			StatusCode: http.StatusBadRequest,
			Response:   []byte("failed to decode body\n"),
		},
	}
	for tc := range slices.Values(testCases) {
		t.Run(tc.Name, server.RunTestCase(tc))
	}
}

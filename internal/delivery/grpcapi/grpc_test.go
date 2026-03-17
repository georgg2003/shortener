package grpcapi_test

import (
	"context"
	"testing"

	"github.com/georgg2003/shortener/api"
	"github.com/georgg2003/shortener/internal/delivery/grpcapi"
	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/internal/pkg/testutils"
	"github.com/georgg2003/shortener/internal/usecase/mock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var testShortID = testutils.TestShortID
var testBaseURL = "localhost:3000"
var testOriginalURL = testutils.TestOriginalURL
var testShortURL = testBaseURL + "/" + testShortID

func TestShortenerService_ExpandURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc := mock.NewMockUseCase(ctrl)
	server := grpcapi.NewShortenerServer(uc, logrus.StandardLogger())

	makeCall := func() (*api.URLExpandResponse, error) {
		return server.ExpandURL(context.Background(), api.URLExpandRequest_builder{
			Id: &testShortID,
		}.Build())
	}

	makeMock := func() *gomock.Call {
		return uc.EXPECT().ProcessShortURL(gomock.Any(), testShortID)
	}

	t.Run("success", func(t *testing.T) {
		makeMock().Return(testutils.TestOriginalURL, false, nil)
		resp, err := makeCall()
		assert.NoError(t, err)
		assert.Equal(t, testutils.TestOriginalURL, resp.GetResult())
	})

	t.Run("internal", func(t *testing.T) {
		makeMock().Return("", false, testutils.ErrSomeError)
		resp, err := makeCall()
		assert.False(t, resp.HasResult())
		assert.Equal(t, codes.Internal, status.Code(err))
	})

	t.Run("deleted", func(t *testing.T) {
		makeMock().Return(testutils.TestOriginalURL, true, nil)
		resp, err := makeCall()
		assert.False(t, resp.HasResult())
		assert.Equal(t, codes.NotFound, status.Code(err))
	})
}

func TestShortenerService_ListUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc := mock.NewMockUseCase(ctrl)
	server := grpcapi.NewShortenerServer(uc, logrus.StandardLogger())

	makeCall := func() (*api.UserURLsResponse, error) {
		return server.ListUserURLs(context.Background(), api.ListUserURLsRequest_builder{}.Build())
	}

	makeMock := func() *gomock.Call {
		return uc.EXPECT().GetUserURLs(gomock.Any())
	}

	t.Run("success", func(t *testing.T) {
		makeMock().Return([]models.URLEntity{{
			ShortID:     testShortID,
			OriginalURL: testutils.TestOriginalURL,
			ShortURL:    testShortURL,
		}}, nil)
		resp, err := makeCall()
		assert.NoError(t, err)
		url := resp.GetUrl()
		require.Len(t, url, 1)
		assert.Equal(t, testShortURL, url[0].GetShortUrl())
		assert.Equal(t, testutils.TestOriginalURL, url[0].GetOriginalUrl())
	})

	t.Run("internal", func(t *testing.T) {
		makeMock().Return(nil, testutils.ErrSomeError)
		resp, err := makeCall()
		assert.Nil(t, resp)
		assert.Equal(t, codes.Internal, status.Code(err))
	})

	t.Run("empty", func(t *testing.T) {
		makeMock().Return([]models.URLEntity{}, nil)
		resp, err := makeCall()
		assert.Nil(t, resp)
		assert.NoError(t, err)
	})
}

func TestShortenerService_ShortenURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc := mock.NewMockUseCase(ctrl)
	server := grpcapi.NewShortenerServer(uc, logrus.StandardLogger())

	makeCall := func() (*api.URLShortenResponse, error) {
		return server.ShortenURL(context.Background(), api.URLShortenRequest_builder{
			Url: &testOriginalURL,
		}.Build())
	}

	makeMock := func() *gomock.Call {
		return uc.EXPECT().NewShortURL(gomock.Any(), testOriginalURL)
	}

	t.Run("success", func(t *testing.T) {
		makeMock().Return(testShortURL, nil)
		resp, err := makeCall()
		assert.NoError(t, err)
		assert.Equal(t, testShortURL, resp.GetResult())
	})

	t.Run("internal", func(t *testing.T) {
		makeMock().Return("", testutils.ErrSomeError)
		resp, err := makeCall()
		assert.False(t, resp.HasResult())
		assert.Equal(t, codes.Internal, status.Code(err))
	})
}

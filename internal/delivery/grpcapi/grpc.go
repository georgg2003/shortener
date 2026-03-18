package grpcapi

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/api"
	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShortenerServer struct {
	api.UnimplementedShortenerServiceServer

	uc     usecase.UseCase
	logger logrus.FieldLogger
}

func (s *ShortenerServer) ExpandURL(ctx context.Context, req *api.URLExpandRequest) (*api.URLExpandResponse, error) {
	longURL, isDeleted, err := s.uc.ProcessShortURL(ctx, req.GetId())
	if err != nil {
		msg := "failed to get original url"
		s.logger.WithError(err).Error(msg)
		return nil, status.Error(codes.Internal, msg)
	}
	if isDeleted {
		return nil, status.Error(codes.NotFound, "url was deleted")
	}

	return api.URLExpandResponse_builder{
		Result: &longURL,
	}.Build(), nil
}

func (s *ShortenerServer) ListUserURLs(ctx context.Context, req *api.ListUserURLsRequest) (*api.UserURLsResponse, error) {
	urls, err := s.uc.GetUserURLs(ctx)
	if err != nil {
		msg := "failed to get user urls"
		s.logger.WithError(err).Error(msg)
		return nil, status.Error(codes.Internal, msg)
	}
	if len(urls) == 0 {
		return nil, nil
	}

	urlData := make([]*api.URLData, 0, len(urls))
	for _, url := range urls {
		urlData = append(urlData, api.URLData_builder{
			ShortUrl:    &url.ShortURL,
			OriginalUrl: &url.OriginalURL,
		}.Build())
	}
	return api.UserURLsResponse_builder{
		Url: urlData,
	}.Build(), nil
}

func (s *ShortenerServer) ShortenURL(ctx context.Context, req *api.URLShortenRequest) (*api.URLShortenResponse, error) {
	shortURL, err := s.uc.NewShortURL(ctx, req.GetUrl())
	if errors.Is(err, usecase.ErrURLEntityAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, "url already exists")
	} else if err != nil {
		msg := "failed to create new short url"
		s.logger.WithError(err).Error(msg)
		return nil, status.Error(codes.Internal, msg)
	}
	return api.URLShortenResponse_builder{
		Result: &shortURL,
	}.Build(), nil
}

func NewShortenerServer(uc usecase.UseCase, logger logrus.FieldLogger) api.ShortenerServiceServer {
	return &ShortenerServer{
		uc:     uc,
		logger: logger,
	}
}

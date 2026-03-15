package grpcapi

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/api"
	"github.com/georgg2003/shortener/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ShortenerServer struct {
	api.UnimplementedShortenerServiceServer

	uc usecase.UseCase
}

func (s *ShortenerServer) ExpandURL(ctx context.Context, req *api.URLExpandRequest) (*api.URLExpandResponse, error) {
	longURL, isDeleted, err := s.uc.ProcessShortURL(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get original url: %w", err)
	}
	if isDeleted {
		return nil, status.Error(codes.NotFound, "url was deleted")
	}

	return api.URLExpandResponse_builder{
		Result: &longURL,
	}.Build(), nil
}

func (s *ShortenerServer) ListUserURLs(ctx context.Context, req *emptypb.Empty) (*api.UserURLsResponse, error) {
	urls, err := s.uc.GetUserURLs(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user urls: %w", err)
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
		return nil, status.Errorf(codes.Internal, "failed to create new short url: %w", err)
	}
	return api.URLShortenResponse_builder{
		Result: &shortURL,
	}.Build(), nil
}

func NewShortenerServer(uc usecase.UseCase) api.ShortenerServiceServer {
	return &ShortenerServer{
		uc: uc,
	}
}

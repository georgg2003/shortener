package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/georgg2003/shortener/internal/usecase/mock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type TestCase struct {
	name               string
	expectedRepoReturn error
	expectedReturn     error
}

var errConnectionFailed = errors.New("some connection error")

func TestPing(t *testing.T) {

	conf := config.New()

	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)

	uc := usecase.New(repo, conf, logrus.New())

	testCases := []TestCase{
		{
			name:               "success",
			expectedRepoReturn: nil,
			expectedReturn:     nil,
		},
		{
			name:               "error",
			expectedRepoReturn: errConnectionFailed,
			expectedReturn:     errConnectionFailed,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repo.EXPECT().Ping(gomock.Any()).Return(testCase.expectedRepoReturn)
			err := uc.Ping(context.Background())
			assert.Equal(t, err, testCase.expectedReturn)
		})
	}
}

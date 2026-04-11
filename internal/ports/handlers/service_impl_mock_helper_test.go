package handlers

import (
	"context"

	rpocts "github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/services/contracts"
	"github.com/golang/mock/gomock"
)

// Helpers para mockar operações do Service com suporte a transação e reutilização
func ExpectServiceCreate[C, R, U any, E rpocts.Entity](svc *MockService[C, R, U, E], fn func(ctx context.Context, dto C, response *R) error) *gomock.Call {
	return svc.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, dto C, response *R) error {
			if fn != nil {
				return fn(ctx, dto, response)
			}
			return nil
		},
	)
}

func ExpectServiceGetByID[C, R, U any, E rpocts.Entity](svc *MockService[C, R, U, E], fn func(ctx context.Context, id string, response *R) error) *gomock.Call {
	return svc.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, id string, response *R) error {
			if fn != nil {
				return fn(ctx, id, response)
			}
			return nil
		},
	)
}

func ExpectServiceList[C, R, U any, E rpocts.Entity](svc *MockService[C, R, U, E], fn func(ctx context.Context, params contracts.ListParams, response *[]R) (int64, error)) *gomock.Call {
	return svc.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, params contracts.ListParams, response *[]R) (int64, error) {
			if fn != nil {
				return fn(ctx, params, response)
			}
			return 0, nil
		},
	)
}

func ExpectServiceUpdate[C, R, U any, E rpocts.Entity](svc *MockService[C, R, U, E], fn func(ctx context.Context, id string, dto U, response *R) error) *gomock.Call {
	return svc.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, id string, dto U, response *R) error {
			if fn != nil {
				return fn(ctx, id, dto, response)
			}
			return nil
		},
	)
}

func ExpectServiceDelete[C, R, U any, E rpocts.Entity](svc *MockService[C, R, U, E], fn func(ctx context.Context, id string) error) *gomock.Call {
	return svc.EXPECT().Delete(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, id string) error {
			if fn != nil {
				return fn(ctx, id)
			}
			return nil
		},
	)
}

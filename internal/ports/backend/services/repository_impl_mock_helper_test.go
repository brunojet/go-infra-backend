package services

import (
	"context"

	rpocts "github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories/contracts"
	"github.com/golang/mock/gomock"
)

// Helper para mockar Create com suporte transparente a WithTx
func ExpectCreateWithTx[E rpocts.Entity](repo *MockRepository[E], fn func(*E)) *gomock.Call {
	repo.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, txFn func(context.Context) error) error {
			return txFn(ctx)
		},
	).AnyTimes()
	return repo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, inOut *E) error {
			if fn != nil {
				fn(inOut)
			}
			return nil
		},
	)
}

// Helper para mockar GetByID com suporte transparente a WithTx
func ExpectGetByIDWithTx[E rpocts.Entity](repo *MockRepository[E], fn func(*E)) *gomock.Call {
	repo.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, txFn func(context.Context) error) error {
			return txFn(ctx)
		},
	).AnyTimes()
	return repo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, _ map[string]any, out *E) error {
			if fn != nil {
				fn(out)
			}
			return nil
		},
	)
}

// Helper para mockar List com suporte transparente a WithTx
func ExpectListWithTx[E rpocts.Entity](repo *MockRepository[E], fn func(*[]E) error) *gomock.Call {
	repo.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, txFn func(context.Context) error) error {
			return txFn(ctx)
		},
	).AnyTimes()
	return repo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, _ interface{}, out *[]E) (int64, error) {
			if fn != nil {
				if err := fn(out); err != nil {
					return 0, err
				}
			}
			return int64(len(*out)), nil
		},
	)
}

// Helper para mockar Update com suporte transparente a WithTx
func ExpectUpdateWithTx[E rpocts.Entity](repo *MockRepository[E], fn func(*E)) *gomock.Call {
	repo.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, txFn func(context.Context) error) error {
			return txFn(ctx)
		},
	).AnyTimes()
	return repo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, _ map[string]any, inOut *E) error {
			if fn != nil {
				fn(inOut)
			}
			return nil
		},
	)
}

// Helper para mockar Delete com suporte transparente a WithTx
func ExpectDeleteWithTx[E rpocts.Entity](repo *MockRepository[E], fn func() error) *gomock.Call {
	repo.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, txFn func(context.Context) error) error {
			return txFn(ctx)
		},
	).AnyTimes()
	return repo.EXPECT().Delete(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, _ map[string]any) error {
			if fn != nil {
				return fn()
			}
			return nil
		},
	)
}

package sql_store

import (
	"context"

	v1 "github.com/superjcd/rivalservice/genproto/v1"
)

type SqlFactory interface {
	Rivals() RivalStore
	RivalChanges() RivalChangeStore
	ProductDetails() ProductDetailsStore
	Close() error
}

type RivalStore interface {
	Create(ctx context.Context, _ *v1.CreateRivalRequest) error
	List(ctx context.Context, _ *v1.ListRivalRequest) (*RivalList, error)
	Delete(ctx context.Context, _ *v1.DeleteRivalRequest) error
}

type RivalChangeStore interface {
	Append(ctx context.Context, _ *v1.AppendRivalChangesRequest) error
	List(ctx context.Context, _ *v1.ListRivalChangesRequest) (*UserRivalChangeList, error)
	Delete(ctx context.Context, _ *v1.DeleteRivalChangesRequest) error
}

type ProductDetailsStore interface {
	AppendActiveDetail(ctx context.Context, _ *v1.AppendRivalProductActiveDetailRequest) error
	DeleteActiveDetail(ctx context.Context, _ *v1.DeleteRivalActiveDetailRequest) error
	AppendInactiveDetail(ctx context.Context, _ *v1.AppendRivalProductInactiveDetailRequest) error
	DeleteInactiveDetail(ctx context.Context, _ *v1.DeleteRivalInactiveDetailRequest) error
}

package repository

import "context"

type Repository interface {
	Init(ctx context.Context) error
	CheckConnection(ctx context.Context) error
}

package repository

import (
	"context"

	"github.com/IFA-01/messenger/internal/repository/queries"
)

type Repository interface {
	CreateUser(ctx context.Context, arg queries.CreateUserParams) (queries.User, error)
	GetUserByEmail(ctx context.Context, email string) (queries.User, error)
	GetUserByUsername(ctx context.Context, username string) (queries.User, error)
	GetUserByID(ctx context.Context, id int64) (queries.User, error)
}

type repository struct {
	q *queries.Queries
}

func NewRepository(q *queries.Queries) Repository {
	return &repository{q: q}
}

func (r *repository) CreateUser(ctx context.Context, arg queries.CreateUserParams) (queries.User, error) {
	return r.q.CreateUser(ctx, arg)
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (queries.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (queries.User, error) {
	return r.q.GetUserByUsername(ctx, username)
}

func (r *repository) GetUserByID(ctx context.Context, id int64) (queries.User, error) {
	return r.q.GetUserByID(ctx, id)
}

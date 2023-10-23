package domain

import "context"

type User struct {
	Id          int64   `json:"id"`
	Name        *string `json:"name" validate:"required"`
	Surname     *string `json:"surname" validate:"required"`
	Patronymic  *string `json:"patronymic"`
	Age         *int64  `json:"age"`
	Gender      *string `json:"gender"`
	Nationality *string `json:"nationality"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, user *User) error
	GetById(ctx context.Context, id int64) (*User, error)
	GetFilteredUsers(ctx context.Context, count string, nationality string, minAge string, maxAge string, gender string) ([]*User, error)
}

type UserUsecase interface {
	Create(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, user *User) error
	GetById(ctx context.Context, id int64) (*User, error)
	GetFilteredUsers(ctx context.Context, count string, nationality string, minAge string, maxAge string, gender string) ([]*User, error)
}

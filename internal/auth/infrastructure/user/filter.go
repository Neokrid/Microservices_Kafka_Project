package user

import "github.com/google/uuid"

type UserFilter struct {
	Id    *uuid.UUID
	Email *string

	Limit uint64
}

type UserUpdateParams struct {
	Username       *string
	Email          *string
	Password       *string
	ImgUrl         *string
	ConfirmedEmail *bool
}

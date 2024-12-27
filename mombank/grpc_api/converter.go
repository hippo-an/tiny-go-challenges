package grpc_api

import (
	db "github.com/hippo-an/tiny-go-challenges/mombank/db/sqlc"
	"github.com/hippo-an/tiny-go-challenges/mombank/gpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(u *db.User) *gpb.User {
	return &gpb.User{
		ID:                u.ID,
		Username:          u.Username,
		FullName:          u.FullName,
		Email:             u.Email,
		PasswordChangedAt: timestamppb.New(u.PasswordChangedAt),
		CreatedAt:         timestamppb.New(u.CreatedAt),
	}
}

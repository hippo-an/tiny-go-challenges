package grpc_api

import (
	"context"
	"database/sql"
	"errors"
	db "github.com/hippo-an/tiny-go-challenges/mombank/db/sqlc"
	"github.com/hippo-an/tiny-go-challenges/mombank/gpb"
	"github.com/hippo-an/tiny-go-challenges/mombank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *gpb.LoginUserRequest) (*gpb.LoginUserResponse, error) {
	user, err := s.store.GetUserByUsername(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch user: %v", err)
	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid password: %v", err)
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(
		user.ID,
		s.config.AccessTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to make access token: %v", err)
	}

	refreshToken, refreshTokenPayload, err := s.tokenMaker.CreateToken(
		user.ID,
		s.config.RefreshTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to make refresh token: %v", err)
	}

	arg := db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		UserID:       user.ID,
		RefershToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	}

	session, err := s.store.CreateSession(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}

	rsp := &gpb.LoginUserResponse{
		User:                  convertUser(&user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
	}

	return rsp, nil
}

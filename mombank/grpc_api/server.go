package grpc_api

import (
	"fmt"
	db "github.com/hippo-an/tiny-go-challenges/mombank/db/sqlc"
	"github.com/hippo-an/tiny-go-challenges/mombank/gpb"
	"github.com/hippo-an/tiny-go-challenges/mombank/token"
	"github.com/hippo-an/tiny-go-challenges/mombank/util"
)

type Server struct {
	gpb.UnimplementedMomBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}

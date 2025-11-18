package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/hippo-an/tiny-go-challenges/netpro_99/encode_decode/housework/v1"
)

type Rosie struct {
	housework.UnimplementedRobotMaidServer
	mu     sync.Mutex
	chores []*housework.Chore
}

func (r *Rosie) Add(_ context.Context, chores *housework.Chores) (*housework.Response, error) {
	r.mu.Lock()
	r.chores = append(r.chores, chores.Chores...)
	r.mu.Unlock()

	return &housework.Response{Message: "ok"}, nil
}

func (r *Rosie) Complete(_ context.Context, req *housework.CompleteRequest) (*housework.Response, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.chores == nil || req.ChoreNumber < 1 || int(req.ChoreNumber) > len(r.chores) {
		return nil, fmt.Errorf("chore %d not found", req.ChoreNumber)
	}

	r.chores[req.ChoreNumber-1].Complete = true
	return &housework.Response{Message: "ok"}, nil
}

func (r *Rosie) List(_ context.Context, _ *housework.Empty) (*housework.Chores, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.chores == nil {
		r.chores = make([]*housework.Chore, 0)
	}

	return &housework.Chores{Chores: r.chores}, nil
}

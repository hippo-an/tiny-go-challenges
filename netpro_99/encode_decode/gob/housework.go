package gob

import (
	"encoding/gob"
	"io"

	"github.com/hippo-an/tiny-go-challenges/netpro_99/encode_decode/housework"
)

func Load(r io.Reader) ([]*housework.Chore, error) {
	var chores []*housework.Chore
	return chores, gob.NewDecoder(r).Decode(&chores)
}

func Flush(w io.Writer, chores []*housework.Chore) error {
	return gob.NewEncoder(w).Encode(chores)
}

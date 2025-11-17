package json

import (
	"encoding/json"
	"io"

	"github.com/hippo-an/tiny-go-challenges/netpro_99/encode_decode/housework"
)

func Load(r io.Reader) ([]*housework.Chore, error) {
	var chore []*housework.Chore

	return chore, json.NewDecoder(r).Decode(&chore)

}
func Flust(w io.Writer, chores []*housework.Chore) error {
	return json.NewEncoder(w).Encode(chores)
}

package protobuf

import (
	"io"

	"github.com/hippo-an/tiny-go-challenges/netpro_99/encode_decode/housework"
	pb "github.com/hippo-an/tiny-go-challenges/netpro_99/encode_decode/housework/v1"
	"google.golang.org/protobuf/proto"
)

func Load(r io.Reader) ([]*housework.Chore, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var pbChores pb.Chores
	if err := proto.Unmarshal(b, &pbChores); err != nil {
		return nil, err
	}

	chores := make([]*housework.Chore, len(pbChores.Chores))
	for i, pbChore := range pbChores.Chores {
		chores[i] = &housework.Chore{
			Complete:    pbChore.Complete,
			Description: pbChore.Description,
		}
	}

	return chores, nil
}

func Flush(w io.Writer, chores []*housework.Chore) error {
	pbChores := make([]*pb.Chore, len(chores))
	for i, chore := range chores {
		pbChores[i] = &pb.Chore{
			Complete:    chore.Complete,
			Description: chore.Description,
		}
	}

	b, err := proto.Marshal(&pb.Chores{Chores: pbChores})
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

package dungeon

import (
	"github.com/google/uuid"
	"math/rand/v2"
	"time"
)

type Dungeon struct {
	Groups []*Group
	Nodes  []*Node
}

type Group struct {
	ID    string
	Name  string
	Nodes []*Node
}
type Node struct {
	ID           string
	MapNumber    int
	DungeonGroup *Group
	Name         string
	Adjacent     []*Node
}

const (
	maxDungeonGroupNumber = 10
	maxDungeonNodeNumber  = 500
)

func NewDungeon() *Dungeon {
	// generate dungeon groups
	dungeonGroups := generateDungeonGroups()

	// generate dungeon nodes
	dungeonNodes := generateDungeonNodes(dungeonGroups)

	return &Dungeon{
		Groups: []*Group{},
		Nodes:  dungeonNodes,
	}
}

func generateDungeonGroups() []*Group {
	groups := make([]*Group, maxDungeonGroupNumber)
	for i := 0; i < maxDungeonGroupNumber; i++ {
		groups[i] = &Group{
			ID:   uuid.NewString(),
			Name: generateDungeonName(),
		}
	}
	return groups
}

func generateDungeonNodes(group []*Group) []*Node {
	nodes := make([]*Node, maxDungeonNodeNumber)

	for i := 0; i < maxDungeonNodeNumber; i++ {
		n := &Node{
			ID:        uuid.NewString(),
			MapNumber: i,
			Adjacent:  make([]*Node, 0),
			Name:      generateDungeonName(),
		}

		groupIndex := rand.IntN(maxDungeonGroupNumber)
		n.DungeonGroup = group[groupIndex]
		group[groupIndex].Nodes = append(group[groupIndex].Nodes, n)
		nodes[i] = n
	}

	// mapping adjacent dungeon node
	r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
	for i := 0; i < maxDungeonNodeNumber; i++ {
		numEdges := r.IntN(5) + 1
		for j := 0; j < numEdges; j++ {
			neighborIdx := GetValidNeighborIdx(i, maxDungeonNodeNumber, r)
			var neighbor *Node
			if r.IntN(10) < 8 {
				g := nodes[i].DungeonGroup
				if g != nil && len(g.Nodes) > 1 {
					neighbor = g.Nodes[r.IntN(len(g.Nodes))]
				} else {
					neighbor = nodes[neighborIdx]
				}
			} else {
				neighbor = nodes[neighborIdx]
			}

			nodes[i].Adjacent = append(nodes[i].Adjacent, neighbor)
		}
	}
	return nodes
}

func GetValidNeighborIdx(currentIdx, maxNum int, r *rand.Rand) int {
	var neighborIdx int
	for {
		neighborIdx = r.IntN(maxNum)
		if neighborIdx != currentIdx {
			break
		}
	}
	return neighborIdx
}

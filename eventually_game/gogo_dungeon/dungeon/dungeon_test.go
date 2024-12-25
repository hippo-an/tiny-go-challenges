package dungeon

import (
	"github.com/google/uuid"
	"math/rand/v2"
	"testing"
	"time"
)

func TestNewDungeon(t *testing.T) {
	dungeon := NewDungeon()

	if len(dungeon.Nodes) != maxDungeonNodeNumber {
		t.Errorf("Expected %d nodes, but got %d", maxDungeonNodeNumber, len(dungeon.Nodes))
	}

	if len(dungeon.Groups) != 0 {
		t.Errorf("Expected 0 groups, but got %d", len(dungeon.Groups))
	}
}

func TestGenerateDungeonGroups(t *testing.T) {
	groups := generateDungeonGroups()

	if len(groups) != maxDungeonGroupNumber {
		t.Errorf("Expected %d groups, but got %d", maxDungeonGroupNumber, len(groups))
	}

	for _, group := range groups {
		if _, err := uuid.Parse(group.ID); err != nil {
			t.Errorf("Invalid UUID for Group ID: %s", group.ID)
		}
		if group.Name == "" {
			t.Errorf("Group Name should not be empty")
		}
	}
}

func TestGenerateDungeonNodes(t *testing.T) {
	groups := generateDungeonGroups()
	nodes := generateDungeonNodes(groups)

	if len(nodes) != maxDungeonNodeNumber {
		t.Errorf("Expected %d nodes, but got %d", maxDungeonNodeNumber, len(nodes))
	}

	for _, node := range nodes {
		if _, err := uuid.Parse(node.ID); err != nil {
			t.Errorf("Invalid UUID for Node ID: %s", node.ID)
		}
		if node.Name == "" {
			t.Errorf("Node Name should not be empty")
		}
		if node.DungeonGroup == nil {
			t.Errorf("Node should belong to a Dungeon Group")
		}
	}
}

func TestNodeAdjacency(t *testing.T) {
	groups := generateDungeonGroups()
	nodes := generateDungeonNodes(groups)

	for _, node := range nodes {
		if len(node.Adjacent) == 0 {
			t.Errorf("Node %s should have at least one adjacent node", node.ID)
		}
	}
}

func TestGetValidNeighborIdx(t *testing.T) {
	r := rand.New(rand.NewPCG(uint64(time.Now().Unix()), uint64(time.Now().Unix())))

	for i := 0; i < 100; i++ {
		neighborIdx := GetValidNeighborIdx(i, maxDungeonNodeNumber, r)
		if neighborIdx == i {
			t.Errorf("Neighbor index should not be the same as the current index")
		}
	}
}

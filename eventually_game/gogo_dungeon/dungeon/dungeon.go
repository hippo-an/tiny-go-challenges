package dungeon

import (
	"fmt"
	"github.com/hippo-an/tiny-go-challenges/gogo-dungeon/utils"
	"log"
	"math/rand/v2"
)

var (
	dungeonAdjectives []string
	dungeonNouns      []string
)

func init() {
	file, err := utils.ReadFile("dungeon_adjectives")
	if err != nil {
		log.Fatal("error reading dungeon adjectives file:", err)
	}

	dungeonAdjectives = file

	file, err = utils.ReadFile("dungeon_nouns")
	if err != nil {
		log.Fatal("error reading dungeon nouns file:", err)
	}

	dungeonNouns = file
}

func generateDungeonName() string {
	adjective := dungeonAdjectives[rand.IntN(len(dungeonAdjectives))]
	noun := dungeonNouns[rand.IntN(len(dungeonNouns))]
	return fmt.Sprintf("%s %s", adjective, noun)
}

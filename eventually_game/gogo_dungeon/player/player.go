package player

import "github.com/hippo-an/tiny-go-challenges/gogo-dungeon/dungeon"

type Player struct {
	Name      string
	Health    int
	MaxHealth int
	Mana      int
	MaxMana   int
	Position  *dungeon.Node
}

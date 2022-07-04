// Copyright 2022. Motty Cohen
//
// Data model and testing data
//
package test

import (
	. "github.com/mottyc/yaaf-common/entity"
)

// region Heroes Test Model --------------------------------------------------------------------------------------------
type Hero struct {
	BaseEntity
	Key  int    `json:"key"`  // Key
	Name string `json:"name"` // Name
}

func (a Hero) TABLE() string { return "hero" }
func (a Hero) NAME() string  { return a.Name }

func NewHero() Entity {
	return &Hero{}
}

func NewHero1(id string, key int, name string) Entity {
	return &Hero{
		BaseEntity: BaseEntity{Id: id, CreatedOn: Now(), UpdatedOn: Now()},
		Key:        key,
		Name:       name,
	}
}

var list_of_heroes = []Entity{
	NewHero1("1", 1, "Ant man"),
	NewHero1("2", 2, "Aqua man"),
	NewHero1("3", 3, "Asterix"),
	NewHero1("4", 4, "Bat Girl"),
	NewHero1("5", 5, "Bat Man"),
	NewHero1("6", 6, "Bat Woman"),
	NewHero1("7", 7, "Black Canary"),
	NewHero1("8", 8, "Black Panther"),
	NewHero1("9", 9, "Captain America"),
	NewHero1("10", 10, "Captain Marvel"),
	NewHero1("11", 11, "Cat Woman"),
	NewHero1("12", 12, "Conan the Barbarian"),
	NewHero1("13", 13, "Daredevil"),
	NewHero1("14", 14, "Doctor Strange"),
	NewHero1("15", 15, "Elektra"),
	NewHero1("16", 16, "Ghost Rider"),
	NewHero1("17", 17, "Green Arrow"),
	NewHero1("18", 18, "Green Lantern"),
	NewHero1("19", 19, "Hawkeye"),
	NewHero1("20", 20, "Hellboy"),
	NewHero1("21", 21, "Iron Man"),
	NewHero1("22", 22, "Robin"),
	NewHero1("23", 23, "Spider Man"),
	NewHero1("24", 24, "Supergirl"),
	NewHero1("25", 25, "Superman"),
	NewHero1("26", 26, "Thor"),
	NewHero1("27", 27, "The Wasp"),
	NewHero1("28", 28, "Wolverine"),
	NewHero1("29", 29, "Wonder Woman"),
	NewHero1("30", 30, "X-Man"),
}

// endregion

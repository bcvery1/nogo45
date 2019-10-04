package main

var (
	acquiredUpgrades []upgrade
	allUpgrades      = []upgrade{
		movementControls,
		seeLevel,
		slowEnemies,
		mediumEnemies,
		fastEnemies,
	}
)

var (
	movementControls = upgrade{
		id:   uniqueID(),
		name: "Movement controls",
		desc: "Gives the players character the ability to move",
		cost: 1,
		next: []*upgrade{&seeLevel},
	}

	seeLevel = upgrade{
		id:   uniqueID(),
		name: "Level",
		desc: "Allows the player to see the level which includes coins",
		cost: 4,
		next: []*upgrade{&slowEnemies, &mediumEnemies, &fastEnemies},
	}

	slowEnemies = upgrade{
		id:   uniqueID(),
		name: "Basic enemies",
		desc: "Add slow moving enemies to the map.  Enemies drop coins on death",
		cost: 5,
		next: []*upgrade{},
	}

	mediumEnemies = upgrade{
		id:   uniqueID(),
		name: "Regular enemies",
		desc: "Add enemies to the map which can move a bit faster.  Enemies drop coins on death",
		cost: 25,
		next: nil,
	}

	fastEnemies = upgrade{
		id:   uniqueID(),
		name: "Advanced enemies",
		desc: "Add fast moving enemies to the map.  Enemies drop coins on death",
		cost: 35,
		next: nil,
	}
)

type upgrade struct {
	id   int
	name string
	desc string
	cost int
	next []*upgrade
}

func availableUpgrades() []upgrade {
	var availUps []upgrade

	for _, u := range acquiredUpgrades {
		for _, n := range u.next {
			availUps = append(availUps, *n)
		}
	}

	if len(availUps) == 0 {
		return []upgrade{movementControls}
	}

	return dedup(availUps)
}

func dedup(ups []upgrade) []upgrade {
	return ups
}

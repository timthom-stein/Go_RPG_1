package player

import (
	"fmt"
	"start/gear"
)

type Player struct {
	level int

	//creating experience system
	experience int
	exp_limit  int

	//creating stats
	hp          int
	att         int
	def         int
	rec         int
	hpPoints    int
	attPoints   int
	defPoints   int
	bonusPoints int

	stats [7]int

	//creating karma system
	karma int
}

// Gets value of karma to be added / subtracted
// Adds value to karma
// Returns Nothing
func (w *Player) UpdateKarma(value int) {
	w.karma += value
}

// Gets Nothing
// Does Nothing
// Returns karma
func (w *Player) GetKarma() int {
	return w.karma
}

// Gets Nothing
// Adjusts Player level + EXP
// Returns level
func (w *Player) levelUp() int {
	w.level = w.level + 1
	w.experience = 0
	w.exp_limit = (w.exp_limit * 6) / 5

	fmt.Println("Dein Level ist jetzt", w.level, "!")

	w.SetStats() //Updating Original Player Stats after levelUp

	return w.level
}

// Gets inventory, hp, att, def, rec
// Adjusts Stats according to current level
// Returns Stats (hp, att, def, rec)
func (w *Player) Level_Management(inventory [10]*gear.InventorySlot, hp, att, def, rec int) (int, int, int, int) {

	if w.level >= 1 && w.experience >= w.exp_limit {
		w.levelUp()
		w.UpdateSpStats()
		w.SetStats()

		hp, att, def, rec = w.CreateStats(inventory)

	} else if w.level >= 1 && w.experience < w.exp_limit {
		var neededExp int = w.exp_limit - w.experience
		_, att, def, rec = w.CreateStats(inventory)

		hp = hp + rec
		if hp > w.hp {
			hp = w.hp
		}

		fmt.Println("Du brauchst noch", neededExp, "Exp")
	} else {
		w.level = 1
		w.exp_limit = 5

		w.InitEXP()
		w.InitSpStats()
	}
	w.UpdateSpStats()

	return hp, att, def, rec
}

// Gets Nothing
// Allocates SpStats
// Returns Nothing
func (w *Player) UpdateSpStats() {
	var allocation int

	for w.bonusPoints >= 1 {
		fmt.Println("Du hast gerade", w.bonusPoints, "Bonuspunkt/e.")
		fmt.Println("Bitte lege deinen Bonuspunkt fest:")
		fmt.Println("1: HP")
		fmt.Println("2: Att")
		fmt.Println("3: Def")
		fmt.Scanln(&allocation)

		switch allocation {
		case 1:
			w.hpPoints = w.hpPoints + 1
		case 2:
			w.attPoints = w.attPoints + 1
		case 3:
			w.defPoints = w.defPoints + 1

		}
		w.bonusPoints = w.bonusPoints - 1
	}
}

// Gets Nothing
// initializing Special Stats
// Returns Nothing
func (w *Player) InitSpStats() {
	w.hpPoints = 0
	w.attPoints = 0
	w.defPoints = 0
	w.bonusPoints = 1
}

// Gets Nothing
// Setting Stats needed in Main with Stats[7] according to level + SpStats (primary)
// Returns Nothing
func (w *Player) SetStats() {

	//Kreiert stats Array
	var stats [7]int

	stats[0] = w.level
	stats[1] = w.exp_limit
	//erstellt Regeln für Stats
	w.hp = 10 + (3 * w.level) + (3 * w.hpPoints)
	w.att = 5 + w.level + w.attPoints
	w.def = 3 + w.level + w.defPoints

	w.bonusPoints = w.level - w.hpPoints - w.attPoints - w.defPoints
	//packt stats in array
	stats[2] = w.hp
	stats[3] = w.att
	stats[4] = w.def
	stats[5] = w.rec
	stats[6] = w.bonusPoints
	//gibt stats aus
	w.stats = stats
}

// Gets inventory
// Creates current Stats needed for Main (final)
// Returns Stats (hp, att, def, rec)
func (w *Player) CreateStats(inventory [10]*gear.InventorySlot) (int, int, int, int) {
	w.SetStats()
	w.SetStatsItems(inventory)
	var hp = w.GetStat(2)
	var att = w.GetStat(3)
	var def = w.GetStat(4)
	var rec = w.GetStat(5)

	return hp, att, def, rec
}

// Gets inventory, current stats, chosen Stat to upgrade
// Permanently Upgrades Stats by 1/10 of their current maximum -> better scaling
// Returns current stats
func (w *Player) UpgradeStat(inventory [10]*gear.InventorySlot, hp, att, def, rec int, chosenStat int) (int, int, int, int) {
	var maxHp, maxAtt, maxDef, maxRec = w.CreateStats(inventory)
	switch chosenStat {
	case 1:
		hp += maxHp / 10
		if hp > maxHp {
			hp = maxHp
		}
		w.hpPoints += maxHp / 10
	case 2:
		att += maxAtt / 10
		w.attPoints += maxAtt / 10
	case 3:
		def += maxDef / 10
		w.defPoints += maxDef / 10
	case 4:
		rec += maxRec / 10
		w.rec += maxRec / 10
	}
	return hp, att, def, rec
}

// Gets Nothing
// Gives out current Stats of Player (hp, att, def, rec)
// Returns Nothing
func (w *Player) SeePlayerStats(inventory [10]*gear.InventorySlot, hp, att, def, rec int) {
	var maxHp, maxAtt, maxDef, maxRec = w.CreateStats(inventory)
	fmt.Println("Deine Stats sind jetzt:")
	fmt.Println("HP:", hp, "/", maxHp)
	fmt.Println("Att:", att, "/", maxAtt)
	fmt.Println("Def:", def, "/", maxDef)
	fmt.Println("Rec:", rec, "/", maxRec)
}

// Gets Count of Stat in Array Stats[7]
// Returns chosen current Stat
func (w *Player) GetStat(i int) int {
	return w.stats[i]
}

func (w *Player) GetLevel() int {
	return w.level
}

// Gets inventory
// adjusts current Stats to include inventoryBonus
// Returns Nothing
func (w *Player) SetStatsItems(inventory [10]*gear.InventorySlot) {
	var att, def, rec int = gear.CreateStatsItems(inventory)
	w.stats[3], w.stats[4], w.stats[5] = w.stats[3]+att, w.stats[4]+def, w.stats[5]+rec
}

// Gets Nothing
// initiating EXP_Stat
// Returns Nothing
func (w *Player) InitEXP() {
	w.experience = 0
}

///////////////////////////////// Needs Update ///////////////////////////////////////

// Gets exp_value
// calculating Exp and adding it to player
// Returns Nothing
func (w *Player) Exp_Function(enemyStats [4]int) int {
	var value = (enemyStats[1] + enemyStats[2] + enemyStats[3] + enemyStats[0]) / 4
	w.experience = w.experience + value

	return w.experience
}

//////////////////////////////////////////////////////////////////////////

//////////////////////////// Player Creation /////////////////////////

// Gets Nothing
// Initializing Player ------------> Only Function used in Main during Player Creation
// Returns Player
func InitPlayer() *Player {
	var player = NewPlayer()
	player.InitSpStats()
	/////////Creating Empty Inventory To Load Level_Management////////
	var inventory = gear.NewInventory()
	player.Level_Management(inventory, 0, 0, 0, 0)
	//////////////////////////////////////////////////////////////////

	return player
}

// Gets Nothing
// Creates Player
// Returns Player
func NewPlayer() *Player {
	var player *Player = new(Player)
	return player
}

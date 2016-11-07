package elo

import (
	"sort"
	"strings"
	"time"
)

// Table holds the ratings of all registered players.
type Table struct {
	ConstantFactor int               `json:constantFactor`
	Players        map[string]Player `json:players`
	InitialRating  int               `json:initialRating`
	GameLog        GameLog           `json:gameLog`
}

// NewTable creates a new table.
func NewTable(constantFactor int, initialRating int) Table {
	return Table{
		ConstantFactor: constantFactor,
		InitialRating:  initialRating,
		Players:        make(map[string]Player),
		GameLog:        NewGameLog(),
	}
}

// Register adds a new player to the table.
func (table *Table) Register(playerName string) error {
	playerName = sanitiseName(playerName)
	if _, exists := table.Players[playerName]; exists {
		return PlayerAlreadyExists
	}
	table.Players[playerName] = Player{
		Rating: table.InitialRating,
		Name:   playerName,
	}
	return nil
}

func sanitiseName(name string) string {
	return strings.ToLower(name)
}

// AddResult adds the result of a match to the table.
func (table *Table) AddResult(winner, loser string) error {
	winningPlayer, exists := table.Players[winner]
	if !exists {
		return PlayerDoesNotExist
	}
	losingPlayer, exists := table.Players[loser]
	if !exists {
		return PlayerDoesNotExist
	}
	winningPlayer.Played++
	winningPlayer.Won++
	losingPlayer.Played++
	losingPlayer.Lost++
	wOld := winningPlayer.Rating
	lOld := losingPlayer.Rating
	wNew, lNew := CalculateRating(wOld, lOld, table.ConstantFactor)
	winningPlayer.Rating = wNew
	losingPlayer.Rating = lNew
	table.Players[winner] = winningPlayer
	table.Players[loser] = losingPlayer
	gle := GameLogEntry{
		Created: time.Now(),
		Winner:  winner,
		Loser:   loser,
		Notes:   "",
		WinnerChange: RatingChange{
			Before: wOld,
			After: wNew,
		},
		LoserChange: RatingChange{
			Before: lOld,
			After: lNew,
		},
	}
	table.GameLog.Entries = append(table.GameLog.Entries, gle)
	return nil
}

// GetPlayersSortedByRating returns a slice of players, sorted in desc order of rating.
func (table *Table) GetPlayersSortedByRating() []Player {
	count := len(table.Players)
	players := make([]Player, count)
	i := 0
	for _, player := range table.Players {
		players[i] = player
		i++
	}
	sort.Sort(Players(players))
	return players
}

// RecalculateRatingsFromLog will use the game log to recreate the ratings table.
// This is usually useful after having edited the game log manually.
func (table *Table) RecalculateRatingsFromLog() error {
	origGameLog := table.GameLog
	sort.Sort(origGameLog)
	table.Players = make(map[string]Player)
	table.GameLog.Entries = make([]GameLogEntry, 0)
	for _, entry := range origGameLog.Entries {
		if _, found := table.Players[entry.Winner]; !found {
			table.Register(entry.Winner)
		}
		if _, found := table.Players[entry.Loser]; !found {
			table.Register(entry.Loser)
		}
		table.AddResult(entry.Winner, entry.Loser)
		table.GameLog.Entries[len(table.GameLog.Entries)-1].Created = entry.Created
	}
	return nil
}


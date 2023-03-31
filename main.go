package main

import (
	"fmt"
	"math/rand"
	"time"
)

type ChangeSet struct {
	point      int
	removeDice int
	passDice   int
}

type Player struct {
	point     int
	diceCount int
	result    []int
}

type Game struct {
	players []*Player
}

func (g *Game) init(playerCount, dices int) {
	arrPlayer := make([]*Player, playerCount)
	for i := 0; i < len(arrPlayer); i++ {
		arrPlayer[i] = &Player{}
		arrPlayer[i].diceCount = dices
	}

	g.players = arrPlayer

}

func (g *Game) rollDice() {
	for i, p := range g.players {
		p.roll()
		fmt.Println(p.printPlayerDice(i + 1))
	}
}

func (g *Game) evaluate() bool {
	result := []ChangeSet{}

	for _, p := range g.players {
		result = append(result, p.evaluate())
	}

	for i, r := range result {
		index := i + 1
		if index == len(g.players) {
			index = 0
		}

		nextPlayer := g.players[index]
		nextPlayerResult := result[index]

		nextPlayer.diceCount += r.passDice - nextPlayerResult.removeDice - nextPlayerResult.passDice
		for i := 0; i < r.passDice; i++ {
			nextPlayer.addDice(1)
		}

		player := g.players[i]
		player.point += r.point

	}

	for i, r := range g.players {
		fmt.Println(r.printPlayerDice(i + 1))
	}

	activePlayerCount := 0
	for _, p := range g.players {
		if p.diceCount > 0 {
			activePlayerCount++
		}
	}

	return activePlayerCount <= 1
}

func (g *Game) getWinner() int {
	result := -1
	maxPoint := 0

	for i, p := range g.players {
		if p.point > maxPoint {
			result = i
			maxPoint = p.point
		}
	}

	return result
}

func (p *Player) evaluate() ChangeSet {
	result := ChangeSet{}
	for i := len(p.result) - 1; i >= 0; i-- {
		if p.result[i] == 6 {
			result.point++
			result.removeDice++
			p.removeDice(i)
		} else if p.result[i] == 1 {
			result.passDice++
			p.removeDice(i)
		}
	}
	return result
}

func (p *Player) addDice(value int) {
	p.result = append(p.result, value)
}

func (p *Player) removeDice(index int) {
	p.result = append(p.result[:index], p.result[index+1:]...)
}

func (p *Player) printPlayerDice(number int) string {
	if p.diceCount == 0 {
		return fmt.Sprintf("Pemain #%d (%d):_ (Berhenti Bermain karena tidak memiliki dadu)", number, p.point)
	}

	return fmt.Sprintf("Pemain #%d (%d):%v", number, p.point, p.result)
}

func (p *Player) roll() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	dices := make([]int, p.diceCount)
	for i := 0; i < p.diceCount; i++ {
		dices[i] = rand.Intn(6) + 1
	}
	p.result = dices
}

func main() {
	game := new(Game)

	playerCount := 50
	diceCount := 100
	game.init(playerCount, diceCount)

	var round int = 1

	fmt.Printf("Pemain = %d, Dadu = %d\n", playerCount, diceCount)
	fmt.Println("====================")
	for {
		fmt.Printf("Giliran %d lempar dadu:\n", round)
		game.rollDice()
		fmt.Println("Setelah Evaluasi: ")
		endgame := game.evaluate()
		fmt.Println("====================")

		round++

		if endgame {
			break
		}
	}

	winner := game.getWinner() + 1
	fmt.Printf("Game dimenangkan oleh pemain #%d karena memiliki poin lebih banyak dari pemain lainnya.\n", winner)

}

/*
Contributors:
Raman Gupta (300290648)
Aryan Gupta (300281987)
*/
package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type location struct {
	row int
	col int
}

type gridSize struct {
	numRows int
	numCols int
}

func makeMove(currLocation location, size gridSize) location {
	rowMoves := [4]int{1, -1, 0, 0}
	colMoves := [4]int{0, 0, 1, -1}

	for {
		move := rand.Intn(4)
		newLocation := location{currLocation.row + rowMoves[move], currLocation.col + colMoves[move]}

		isRowValid := newLocation.row < size.numRows && newLocation.row >= 0
		isColValid := newLocation.col < size.numCols && newLocation.col >= 0

		if isRowValid && isColValid {
			return newLocation
		}
	}
}

func police(position chan location, result chan int, size gridSize) {
	var currLocation location
	var moveResult int

	end := false
	for !end {
		moveResult = <-result

		if moveResult == -1 {
			end = true
		}

		if !end {
			currLocation = <-position
			position <- makeMove(currLocation, size)
		}
	}

	close(position)
	close(result)
}

func thief(position chan location, result chan int, size gridSize) {
	var currLocation location
	var moveResult int

	end := false
	for !end {
		moveResult = <-result

		if moveResult == -1 {
			end = true
		}

		if !end {
			currLocation = <-position
			position <- makeMove(currLocation, size)
		}
	}

	close(position)
	close(result)
}

func controller(posPolice chan location, resPolice chan int, posTheif chan location, resThief chan int, maxMovesPolice int) {
	policeLocation := location{0, 0}
	thiefLocation := location{4, 4}
	end := false

	fmt.Printf("The police is starting at: (%d, %d)\n", policeLocation.row, policeLocation.col)
	fmt.Printf("The thief is starting at: (%d, %d)\n\n", thiefLocation.row, thiefLocation.col)

	for !end {
		resPolice <- 0
		posPolice <- policeLocation
		policeLocation = <-posPolice
		maxMovesPolice -= 1

		fmt.Println("Police is making a move")
		fmt.Printf("Police: (%d, %d), ", policeLocation.row, policeLocation.col)
		fmt.Printf("Thief: (%d, %d)\n\n", thiefLocation.row, thiefLocation.col)

		hasCaughtThief := policeLocation.row == thiefLocation.row && policeLocation.col == thiefLocation.col

		if hasCaughtThief {
			fmt.Printf("The Police caught the thief at (%d, %d) and won the game.\n", policeLocation.row, policeLocation.col)
			resPolice <- -1
			resThief <- -1
			end = true
		} else if maxMovesPolice == 0 {
			fmt.Println("The Police ran out of moves and the Thief won the game.")
			resPolice <- -1
			resThief <- -1
			end = true
		}

		if end {
			return
		}

		resThief <- 0
		posTheif <- thiefLocation
		thiefLocation = <-posTheif

		fmt.Println("Thief is making a move")
		fmt.Printf("Police: (%d, %d), Thief (%d, %d)\n\n", policeLocation.row, policeLocation.col, thiefLocation.row, thiefLocation.col)

		hasRunIntoPolice := policeLocation.row == thiefLocation.row && policeLocation.col == thiefLocation.col
		hasWon := thiefLocation.row == 0 && thiefLocation.col == 0

		if hasRunIntoPolice && hasWon {
			fmt.Println("The game ends in a tie.")
			resThief <- -1
			resPolice <- -1
			end = true
		} else if hasRunIntoPolice {
			fmt.Printf("The Police caught the thief at (%d, %d) and won the game.\n", policeLocation.row, policeLocation.col)
			resThief <- -1
			resPolice <- -1
			end = true
		} else if hasWon {
			fmt.Println("The Thief escaped and won the game.")
			resThief <- -1
			resPolice <- -1
			end = true
		}
	}
}

func main() {
	resPolice := make(chan int)
	posPolice := make(chan location)

	resThief := make(chan int)
	posTheif := make(chan location)

	n := rand.Intn(191) + 10
	m := rand.Intn(191) + 10
	size := gridSize{n, m}

	maxMovesLowerBound := 2 * max(n, m)
	maxMovesUpperBound := 10 * max(n, m)
	maxMoves := rand.Intn(maxMovesUpperBound-maxMovesLowerBound+1) + maxMovesLowerBound

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		police(posPolice, resPolice, size)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		thief(posTheif, resThief, size)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		controller(posPolice, resPolice, posTheif, resThief, maxMoves)
		wg.Done()
	}()

	wg.Wait()

	fmt.Println()
	fmt.Printf("The Grid size was %d X %d\n", n, m)
	fmt.Printf("The Police had %d moves to catch the thief\n", maxMoves)
}

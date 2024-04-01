package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

//task 1
var i int

func makeCakeAndSend(cs chan string) {
	i += 1
	cakeName := "Strawberry Cake " + strconv.Itoa(i)
	fmt.Println("Making a cake and sending ...", cakeName)
	cs <- cakeName //send a strawberry cake 
}

func receiveCakeAndPack(cs chan string) {
	s := <- cs //get whaever cake is on the channel
	fmt.Println("Packing received cake: ", s)
}

//task 2

//creating a structure to hold number of students, professor,
//and class average

type C struct {numStudents int; professor string; avg float64}

//task 3

type treasureLocation struct {xCoord int; yCoord int}

func playerOne(location treasureLocation, boxsize int) {
	playerLocation := treasureLocation{0, 0}
	numMoves := 0

	for playerLocation != location {
		playerLocationAddress := &playerLocation
		*playerLocationAddress = treasureLocation{rand.Intn(boxsize), rand.Intn(boxsize)};
		numMoves++
	}

	response := fmt.Sprintf("Player 1 found the treasure after %d moves", numMoves)
	fmt.Println(response)
	
}

func playerTwo(location treasureLocation, boxsize int) {
	playerLocation := treasureLocation{boxsize, boxsize}
	numMoves := 0

	for playerLocation != location {
		playerLocationAddress := &playerLocation
		*playerLocationAddress = treasureLocation{rand.Intn(boxsize), rand.Intn(boxsize)};
		numMoves++
	}

	response := fmt.Sprintf("Player 2 found the treasure after %d moves", numMoves)
	fmt.Println(response)
	
}

func main() {

	//task 1
	cs := make(chan string)

	for i := 0; i < 3; i++ {
		go makeCakeAndSend(cs)
		go receiveCakeAndPack(cs)

		//sleep for a while so that the program doesn't exit
		//immediately and output is clear for illustration
		time.Sleep(1 * 1e9)
	}

	fmt.Println("")

	//task 2

	//creating a dynamic map with key datatype string
	//and value datatype the custom structure C
	m := make(map[string]C)

	//adding courses to map
	m["CSI2110"] = C{186, "Lang", 79.5}
	m["CSI2120"] = C{211, "Moura", 81}

	//iterating through map and printing 
	//key, value pairings
	for k, v := range m {
		fmt.Printf("Course Code: %s\n", k)
		fmt.Printf("Number of Students: %d\n", v.numStudents)
		fmt.Printf("Professor %s\n", v.professor)
		fmt.Printf("Average: %f\n\n", v.avg)
	}

	//task 3

	//generating random boxsizes and treasure locations
	boxSize := (rand.Intn(10 - 5) + 5)
	location := treasureLocation{rand.Intn(boxSize), rand.Intn(boxSize)}

	//calling go routines to initiate treasure hunt game
	go playerOne(location, boxSize)
	time.Sleep(1 * 1e6)
	go playerTwo(location, boxSize)
	time.Sleep(1 * 1e6)
	
}
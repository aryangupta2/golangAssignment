/*
Contributors:
Raman Gupta (300290648)
Aryan Gupta (300281987)
*/

//THIS QUESTION ONLY WORKS FOR macOS DEVICES, USES afplay

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

var (
	currSong  string
	isPlaying bool
	//change filepath for your own testing
	currDirectory string = "/Users/aryangupta/Desktop/uOttawa/Winter 2024/CSI 2120/golangAssignment"
	cmd *exec.Cmd
)

func musicPlayer(commands chan string, data chan string) {
	var command string
	end := false

	for !end {
		command = <-commands
		switch command {
		case "open":
			// Open file
			<-data
			commands <- "done"

		case "play":

			if (isPlaying) {
				err := cmd.Process.Signal(os.Interrupt)
				if err != nil {
					fmt.Println("Error pausing music:", err)
				}
				cmd.Wait()
			}

			isPlaying = true

			cmd = exec.Command("afplay", <-data)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Start()
			if err != nil {
				fmt.Println("Error playing music:", err)
			}

			commands <- "done"

		case "pause":
			err := cmd.Process.Signal(syscall.SIGSTOP)
			if err != nil {
				fmt.Println("Error pausing music:", err)
			}

			isPlaying = false
			commands <- "done"

		case "continue":
			isPlaying = true
			err := cmd.Process.Signal(syscall.SIGCONT)
				if err != nil {
					fmt.Println("Error resuming music:", err)
				}

			commands <- "done"

		case "exit":

			err := cmd.Process.Signal(os.Interrupt)
			if err != nil {
				fmt.Println("Error pausing music:", err)
			}
			cmd.Wait()

			end = true
		}
	}
	close(data)
}

func controller(commands chan string, data chan string) {
	var input1 string
	var input2 string
	end := false

	for !end {
		fmt.Println("Enter Command:")
		fmt.Scanf("%s %s", &input1, &input2)

		switch input1 {
		case "open":
			var newDir string

			if (strings.TrimSpace(input2) == "") {
				newDir = strings.TrimSpace(currDirectory)
			} else {
				newDir = strings.TrimSpace(currDirectory + "/" + input2)
			}

			if _, err := os.Stat(newDir); os.IsNotExist(err) {
				fmt.Println("Error, directory entered does not exist")
			} else {
				fmt.Println("The current working directory is changed to " + newDir)
				currDirectory = newDir
				commands <- "open"
				data <- input2
				<-commands
			}

		case "play": 
			if (strings.HasSuffix(input2, ".mp3")) {
				var fileCheck = currDirectory + "/" + input2
				if _, err := os.Stat(fileCheck); os.IsNotExist(err) {
					fmt.Println("File entered does not exist in current directory")
				} else {

					fmt.Printf("Playing %s...\n", input2)
					currSong = input2
					f, err := os.Open(currDirectory)
					if err != nil {
						fmt.Println("Error opening music file:", err)
						return
					}
					defer f.Close()
					commands <- "play"
					data <- fileCheck
					<-commands
				}
			} else {
				fmt.Println("Error, please select an mp3 file ")
			}

		case "pause":
			if !(currSong == "") && isPlaying {
				fmt.Printf("Pausing %s\n", currSong)
				commands <- "pause"
				<-commands
			} else {
				fmt.Println("Nothing playing right now")
			}

		case "continue":
			if !(currSong == "") {
				if !isPlaying {
					fmt.Printf("Continuing %s\n", currSong)
					commands <- "continue"
					<-commands
				} else {
					fmt.Println("Music already playing")
				}
			} else {
				fmt.Println("Select a song to play")
			}

		case "exit":
			fmt.Println("Thanks for listening to music with us, goodbye!")
			commands <- "exit"
			end = true
		default: 
			fmt.Println("Error, command not recognized")
		}
	}
	close(commands)
}

func main() {
	
	commandChannel := make(chan string)
	dataChannel := make(chan string)

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		controller(commandChannel, dataChannel)
		wg.Done()
	}()

	go func() {
		go musicPlayer(commandChannel, dataChannel)
		wg.Done()
	}()

	wg.Wait()
}

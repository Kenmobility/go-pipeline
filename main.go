package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

// GenerateNumberStage is a stage that generates a stream of random
// integers between 1 and 100, then send each generated random integer
// to a channel to be processed by another stage.
func GenerateNumberStage(count int, out chan<- int) {
	fmt.Println("*********INTEGER GENERATOR STAGE*********")
	defer close(out) // close the out channel once the stage is entirely done

	for i := 1; i <= count; i++ {
		intNumber := rand.Intn(100) + 1
		fmt.Printf("generated: %d\n", intNumber)
		out <- intNumber
	}
}

// SquareNumberStage takes in each integer, squares it and sends
// its result to another stage for further processing.
func SquareNumberStage(in <-chan int, out chan<- int) {
	fmt.Println("*********NUMBER SQUARE STAGE*********")
	defer close(out) // close the out channel once the stage is entirely done

	for num := range in {
		fmt.Printf("got: %d\n", num)
		sq := num * num // square each integer got from in chan
		out <- sq       // send sq to out chan
		fmt.Printf("sent: %d\n", sq)
	}
}

// SumAllStage sums all the squares and prints its result out
// to a console.
func SumAllStage(in <-chan int, done chan<- bool) {
	fmt.Println("*********SUM ALL STAGE*********") // info to know when this stage function is called
	sum := 0                                       // declare and initialise an integer variable to hold the sum of the streams of integers

	for num := range in { // use range to get each integer sent to in chan
		fmt.Printf("current Sum: %d\n", sum)
		fmt.Printf("got: %d\n", num)
		sum += num // add to sum
		fmt.Printf("Summed to: %d\n", sum)
	}

	fmt.Println("final sum:", sum)
	done <- true
}

func main() {
	count := flag.Int("count", 10000, "number of integers to generate") // flag to get the number of integers to generate or default to 10k if empty
	flag.Parse()

	start := time.Now() // records the current time

	// channels for pipeline stages
	genChan := make(chan int, *count)    // declare a buffered channel of {count} to hold the stream of generated integers
	squareChan := make(chan int, *count) // declare a buffered channel of {count} to hold the squares of each generated integers
	done := make(chan bool)              // declare an unbuffered channel to signal the end of the pipeline, after the last stage

	//run all stages in a goroutine
	go GenerateNumberStage(*count, genChan)
	go SquareNumberStage(genChan, squareChan)
	go SumAllStage(squareChan, done)

	<-done // listen on this channel to be signalled by the last stage (i.e SumAllStage) when it completes

	fmt.Printf("pipeline completed in %v\n", time.Since(start)) // print the pipeline duration
}

package main

import "testing"

func TestGeneratorStage(t *testing.T) {
	genChan := make(chan int, 10)
	go GenerateNumberStage(10, genChan)

	count := 0

	for range genChan {
		count++
	}

	if count != 10 {
		t.Errorf("Expected 10 numbers, got %d", count)
	}
}

func TestSquareStage(t *testing.T) {
	inChan := make(chan int, 5)
	outChan := make(chan int, 5)

	go func() {
		inChan <- 2
		inChan <- 3
		inChan <- 4

		close(inChan)
	}()

	go SquareNumberStage(inChan, outChan)

	squares := []int{4, 9, 16}
	for _, expected := range squares {
		if result := <-outChan; result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	}
}

func TestSumStage(t *testing.T) {
	inChan := make(chan int, 3)
	done := make(chan bool)

	go func() {
		inChan <- 4
		inChan <- 9
		inChan <- 16

		close(inChan)
	}()

	go SumAllStage(inChan, done)
	<-done
}

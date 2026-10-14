package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	philosophers = 5
	rounds       = 20
)

var forks [philosophers]*sync.Mutex

var (
	metricsMutex sync.Mutex

	currentEatingCount int
	maxEatingCount int
	
	completedCycles [philosophers]int
	
	totalWait time.Duration
	biggestWait time.Duration
)

func Philosopher(id, iteration int, waitGroup *sync.WaitGroup) {
	time_to_think := 2
	time_to_eat := 3

	for r := 1; r <= iteration; r++ {
		// Think
		for s := 1; s <= id*6; s++ {
			fmt.Printf(" ")
		}
		fmt.Printf(" T%d\n", r) // Thinking
		time.Sleep(time.Duration(time_to_think) * time.Second)

		//inicio timer de espera
		waitStart := time.Now()

		// take fork(i)
		forks[id].Lock()
		
		// take fork((i + 1) % N)
		forks[(id+1)%philosophers].Lock()
		
		//fim da espera
		waitTime := time.Since(waitStart)
		
		//atualiza contagem filosofos
		metricsMutex.Lock()
		currentEatingCount++
		if currentEatingCount > maxEatingCount {
			maxEatingCount = currentEatingCount
		}
		//adiciona ao tempo de espera total
		totalWait += waitTime
		
		//maior tempo de espera
		if waitTime > biggestWait{
			biggestWait = waitTime
		}
		metricsMutex.Unlock()
		
		// Eat
		for s := 1; s <= id*6; s++ {
			fmt.Printf(" ")
		}
		fmt.Printf(" E%d\n", r) // Eating
		time.Sleep(time.Duration(time_to_eat) * time.Second)

		//terminou de comer
		metricsMutex.Lock()
		currentEatingCount--
		completedCycles[id]++
		metricsMutex.Unlock()

		// put fork(i)
		forks[id].Unlock()

		// put fork((i + 1) % N)
		forks[(id+1)%philosophers].Unlock()
	}
	waitGroup.Done()

}

func main() {
	var waitGroup sync.WaitGroup

	// initialize forks
	for i := 0; i < philosophers; i++ {
		forks[i] = &sync.Mutex{}
	}

	fmt.Println("\n[P1] [P2] [P3] [P4] [P5]\n")

	start := time.Now()

	// run philosophers
	for j := 0; j < philosophers; j++ {
		waitGroup.Add(1)
		go Philosopher(j, rounds, &waitGroup)
	}

	// wait for all philosophers to finish
	waitGroup.Wait()

	elapsed := time.Since(start)

	averageWait := totalWait / time.Duration(philosophers * rounds)

	fmt.Printf("\nParallel Dinner took %s\n", elapsed)
	fmt.Printf("Max philosophers eating simultaneously: %d\n", maxEatingCount)
	fmt.Printf("Completed cycles per philosopher: %v\n", completedCycles)
	fmt.Printf("Average wait time: %v\n", averageWait)
	fmt.Printf("Biggest wait time: %v\n\n", biggestWait)
}

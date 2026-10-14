// Example code filo_seq.go - Sequential Dinning Philosophers
// Since is sequential, elapsed timne could be simply calculated by (roubnds*philosophers*(time_to_think+time_to_eat)// By Cesar De Rose

package main

import (
	"fmt"
	"time"
)

const (
	philosophers = 5
	rounds       = 20
	timeToThink  = 2
	timeToEat    = 3
)

func Philosopher(id, iteration int) {
	// think
	for s := 1; s <= id*6; s++ {
		fmt.Printf(" ")
	}
	fmt.Printf(" T%d\n", iteration) // Thinking
	time.Sleep(time.Duration(timeToThink) * time.Second)

	// eat

	for s := 1; s <= id*6; s++ {
		fmt.Printf(" ")
	}
	fmt.Printf(" E%d\n", iteration) // Eating
	time.Sleep(time.Duration(timeToEat) * time.Second)
}

func main() {
	completedCycles := [philosophers]int{}

	fmt.Println("\n[P1] [P2] [P3] [P4] [P5]\n")

	// run philosophers

	start := time.Now()

	for i := 1; i <= rounds; i++ {
		for j := 0; j < philosophers; j++ {
			Philosopher(j, i)
			completedCycles[j]++
		}
	}
	elapsed := time.Since(start)

	fmt.Printf("\nSequential Dinner took %s\n", elapsed)
	fmt.Printf("Max philosophers eating simultaneously: %d\n", 1)
	fmt.Printf("Completed cycles per philosopher: %v\n", completedCycles)
	fmt.Printf("Average wait time: %v\n", time.Duration(0))
	fmt.Printf("Biggest wait time: %v\n\n", time.Duration(0))
}

//Última atualização: terça-feira, 22 set. 2026, 14:58

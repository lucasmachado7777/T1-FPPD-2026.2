package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	thinkTime = 2 * time.Second
	eatTime   = 3 * time.Second

	totalWait   time.Duration
	biggestWait time.Duration
	waitCount   int
	eating      int
	maxEating   int
)

type Fork struct {
	ID    int
	Mutex sync.Mutex
}

func NewFork(id int) *Fork {
	return &Fork{ID: id}
}

type Philosopher struct {
	Name      string
	Left      *Fork
	Right     *Fork
	Conductor *sync.Mutex
	cycles    int
}

func NewPhilosopher(
	name string,
	left *Fork,
	right *Fork,
	conductor *sync.Mutex,
) *Philosopher {
	return &Philosopher{
		Name:      name,
		Left:      left,
		Right:     right,
		Conductor: conductor,
		cycles:    0,
	}
}

func (p *Philosopher) Run(wg *sync.WaitGroup) {
	for i := 0; i < 20; i++ {
		fmt.Println("Philosopher", p.Name, "is thinking.")
		time.Sleep(thinkTime)

		fmt.Println("Philosopher", p.Name, "is hungry.")
		start := time.Now()

		p.Conductor.Lock()

		p.Left.Mutex.Lock()
		fmt.Println("Philosopher", p.Name, "picked up left fork", p.Left.ID)

		p.Right.Mutex.Lock()
		fmt.Println("Philosopher", p.Name, "picked up right fork", p.Right.ID)

		waitTime := time.Since(start)
		totalWait += waitTime
		waitCount++
		if waitTime > biggestWait {
			biggestWait = waitTime
		}

		eating++
		if eating > maxEating {
			maxEating = eating
		}

		fmt.Println("Philosopher", p.Name, "is eating.")
		time.Sleep(eatTime)
		eating--

		fmt.Println("Philosopher", p.Name, "putting down right fork", p.Right.ID)
		p.Right.Mutex.Unlock()

		fmt.Println("Philosopher", p.Name, "putting down left fork", p.Left.ID)
		p.Left.Mutex.Unlock()

		p.cycles++
		p.Conductor.Unlock()
	}

	wg.Done()
}

func main() {
	conductor := &sync.Mutex{}
	start := time.Now()

	f1 := NewFork(1)
	f2 := NewFork(2)
	f3 := NewFork(3)
	f4 := NewFork(4)
	f5 := NewFork(5)

	p1 := NewPhilosopher("1", f1, f2, conductor)
	p2 := NewPhilosopher("2", f2, f3, conductor)
	p3 := NewPhilosopher("3", f3, f4, conductor)
	p4 := NewPhilosopher("4", f4, f5, conductor)
	p5 := NewPhilosopher("5", f5, f1, conductor)

	var wg sync.WaitGroup
	wg.Add(5)

	go p1.Run(&wg)
	go p2.Run(&wg)
	go p3.Run(&wg)
	go p4.Run(&wg)
	go p5.Run(&wg)

	wg.Wait()

	elapsed := time.Since(start)
	averageWait := totalWait / time.Duration(waitCount)
	completedCycles := []int{p1.cycles, p2.cycles, p3.cycles, p4.cycles, p5.cycles}

	fmt.Printf("\nParallel Dinner took %s\n", elapsed)
	fmt.Printf("Max philosophers eating simultaneously: %d\n", maxEating)
	fmt.Printf("Completed cycles per philosopher: %v\n", completedCycles)
	fmt.Printf("Average wait time: %v\n", averageWait)
	fmt.Printf("Biggest wait time: %v\n\n", biggestWait)
}

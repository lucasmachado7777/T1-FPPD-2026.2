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

// Barreira usada apenas para garantir a demonstração do deadlock.
// Todos precisam pegar o primeiro garfo antes de tentar o segundo.
var barrier sync.WaitGroup

var (
	metricsMutex sync.Mutex

	currentEatingCount int
	maxEatingCount     int

	completedCycles [philosophers]int

	totalWait   time.Duration
	biggestWait time.Duration
)

func Philosopher(id, iteration int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	time_to_think := 2
	time_to_eat := 3

	for r := 1; r <= iteration; r++ {

		// THINK
		for s := 1; s <= id*6; s++ {
			fmt.Printf(" ")
		}

		fmt.Printf(" T%d\n", r)

		time.Sleep(time.Duration(time_to_think) * time.Second)

		// Início da medição da espera
		waitStart := time.Now()

		// -----------------------------------------
		// PRIMEIRO GARFO
		// -----------------------------------------

		forks[id].Lock()

		fmt.Printf(
			"Filósofo %d pegou o garfo %d\n",
			id,
			id,
		)

		/*
			Cada filósofo informa que conseguiu
			pegar seu primeiro garfo.
		*/
		barrier.Done()

		/*
			Ninguém pode tentar pegar o segundo
			garfo até TODOS terem pegado o primeiro.
		*/
		barrier.Wait()

		// Todos chegaram até aqui segurando um garfo.

		fmt.Printf(
			"Filósofo %d possui o garfo %d e está esperando pelo garfo %d\n",
			id,
			id,
			(id+1)%philosophers,
		)

		// -----------------------------------------
		// SEGUNDO GARFO
		// -----------------------------------------

		/*
			Aqui ocorrerá o deadlock.

			P0 espera F1
			P1 espera F2
			P2 espera F3
			P3 espera F4
			P4 espera F0
		*/

		forks[(id+1)%philosophers].Lock()

		// O programa não deverá chegar abaixo
		// deste ponto durante a demonstração.

		waitTime := time.Since(waitStart)

		metricsMutex.Lock()

		currentEatingCount++

		if currentEatingCount > maxEatingCount {
			maxEatingCount = currentEatingCount
		}

		totalWait += waitTime

		if waitTime > biggestWait {
			biggestWait = waitTime
		}

		metricsMutex.Unlock()

		// EAT
		for s := 1; s <= id*6; s++ {
			fmt.Printf(" ")
		}

		fmt.Printf(" E%d\n", r)

		time.Sleep(time.Duration(time_to_eat) * time.Second)

		metricsMutex.Lock()

		currentEatingCount--
		completedCycles[id]++

		metricsMutex.Unlock()

		// Libera primeiro garfo
		forks[id].Unlock()

		// Libera segundo garfo
		forks[(id+1)%philosophers].Unlock()
	}
}

func main() {
	var waitGroup sync.WaitGroup

	// Inicializa os 5 garfos
	for i := 0; i < philosophers; i++ {
		forks[i] = &sync.Mutex{}
	}

	/*
		A barreira começa em 5.

		Cada filósofo chama barrier.Done()
		depois de pegar seu primeiro garfo.

		Só quando os cinco chegarem a zero
		o barrier.Wait() será liberado.
	*/
	barrier.Add(philosophers)

	fmt.Println()
	fmt.Println("======================================")
	fmt.Println("V1 - Demonstração forçada de deadlock")
	fmt.Println("======================================")
	fmt.Println()

	fmt.Println("[P0] [P1] [P2] [P3] [P4]")
	fmt.Println()

	start := time.Now()

	// Inicia os cinco filósofos concorrentemente
	for j := 0; j < philosophers; j++ {
		waitGroup.Add(1)

		go Philosopher(
			j,
			rounds,
			&waitGroup,
		)
	}

	// Espera os filósofos terminarem.
	// Como haverá deadlock, eles nunca terminarão.
	waitGroup.Wait()

	// Estas linhas não deverão ser alcançadas
	// durante a demonstração do deadlock.

	elapsed := time.Since(start)

	averageWait :=
		totalWait /
			time.Duration(philosophers*rounds)

	fmt.Printf(
		"\nParallel Dinner took %s\n",
		elapsed,
	)

	fmt.Printf(
		"Max philosophers eating simultaneously: %d\n",
		maxEatingCount,
	)

	fmt.Printf(
		"Completed cycles per philosopher: %v\n",
		completedCycles,
	)

	fmt.Printf(
		"Average wait time: %v\n",
		averageWait,
	)

	fmt.Printf(
		"Biggest wait time: %v\n\n",
		biggestWait,
	)
}
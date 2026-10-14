package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	philosophers = 5
	rounds       = 20
	timeToThink  = 2
	timeToEat    = 3
)

// Métricas individuais de cada filósofo
type PhilosopherMetrics struct {
	totalWait time.Duration
	maxWait   time.Duration
	cycles    int
}

// Mutex usado apenas para proteger as métricas
var metricsMutex sync.Mutex

var currentEating int
var maxEating int

func Philosopher(
	id int,
	left int,
	right int,
	chopsticks []sync.Mutex,
	metrics []PhilosopherMetrics,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	// Hierarquia de recursos:
	// sempre pegar primeiro o garfo de menor ID
	primeiro := left
	segundo := right

	if right < left {
		primeiro = right
		segundo = left
	}

	for iteration := 1; iteration <= rounds; iteration++ {

		// THINK
		for s := 1; s <= id*6; s++ {
			fmt.Print(" ")
		}

		fmt.Printf("T%d\n", iteration)

		time.Sleep(time.Duration(timeToThink) * time.Second)

		// Início da medição do tempo de espera
		waitStart := time.Now()

		// Pega os garfos seguindo a hierarquia
		chopsticks[primeiro].Lock()
		chopsticks[segundo].Lock()

		waitTime := time.Since(waitStart)

		// Atualiza métricas de espera
		metricsMutex.Lock()

		metrics[id].totalWait += waitTime

		if waitTime > metrics[id].maxWait {
			metrics[id].maxWait = waitTime
		}

		// Filósofo começou a comer
		currentEating++

		if currentEating > maxEating {
			maxEating = currentEating
		}

		metricsMutex.Unlock()

		// EAT
		for s := 1; s <= id*6; s++ {
			fmt.Print(" ")
		}

		fmt.Printf(
			"E%d  (P%d usa garfos %d e %d)\n",
			iteration,
			id,
			primeiro,
			segundo,
		)

		time.Sleep(time.Duration(timeToEat) * time.Second)

		// Atualiza métricas ao terminar o ciclo
		metricsMutex.Lock()

		currentEating--
		metrics[id].cycles++

		metricsMutex.Unlock()

		// Libera os garfos
		chopsticks[segundo].Unlock()
		chopsticks[primeiro].Unlock()
	}
}

func main() {

	fmt.Println()
	fmt.Println("===================================")
	fmt.Println("Dining Philosophers - V3")
	fmt.Println("Resource Hierarchy")
	fmt.Println("===================================")
	fmt.Println()

	fmt.Println("[P0]   [P1]   [P2]   [P3]   [P4]")
	fmt.Println()

	// Cria os 5 garfos.
	// Cada garfo é representado por um mutex.
	chopsticks := make([]sync.Mutex, philosophers)

	// Métricas de cada filósofo
	metrics := make([]PhilosopherMetrics, philosophers)

	var wg sync.WaitGroup

	start := time.Now()

	// Cria os filósofos concorrentemente
	for id := 0; id < philosophers; id++ {

		left := id
		right := (id + 1) % philosophers

		wg.Add(1)

		go Philosopher(
			id,
			left,
			right,
			chopsticks,
			metrics,
			&wg,
		)
	}

	// Espera todos os filósofos completarem seus ciclos
	wg.Wait()

	totalTime := time.Since(start)

	// =========================================
	// RESULTADOS
	// =========================================

	fmt.Println()
	fmt.Println("===================================")
	fmt.Println("RESULTADOS")
	fmt.Println("===================================")

	fmt.Printf("\nTempo total de execução: %s\n", totalTime)

	var totalWait time.Duration
	var largestWait time.Duration
	totalCycles := 0

	for id := 0; id < philosophers; id++ {

		totalWait += metrics[id].totalWait
		totalCycles += metrics[id].cycles

		if metrics[id].maxWait > largestWait {
			largestWait = metrics[id].maxWait
		}

		var averageWait time.Duration

		if metrics[id].cycles > 0 {
			averageWait =
				metrics[id].totalWait /
					time.Duration(metrics[id].cycles)
		}

		fmt.Printf("\nFilósofo %d:\n", id)

		fmt.Printf(
			"  Ciclos concluídos: %d\n",
			metrics[id].cycles,
		)

		fmt.Printf(
			"  Tempo total esperando: %s\n",
			metrics[id].totalWait,
		)

		fmt.Printf(
			"  Tempo médio de espera: %s\n",
			averageWait,
		)

		fmt.Printf(
			"  Maior tempo de espera: %s\n",
			metrics[id].maxWait,
		)
	}

	var overallAverage time.Duration

	if totalCycles > 0 {
		overallAverage =
			totalWait /
				time.Duration(totalCycles)
	}

	fmt.Println()
	fmt.Println("-----------------------------------")
	fmt.Println("MÉTRICAS GERAIS")
	fmt.Println("-----------------------------------")

	fmt.Printf(
		"Tempo total: %s\n",
		totalTime,
	)

	fmt.Printf(
		"Tempo médio de espera pelos garfos: %s\n",
		overallAverage,
	)

	fmt.Printf(
		"Maior tempo de espera observado: %s\n",
		largestWait,
	)

	fmt.Printf(
		"Máximo de filósofos comendo simultaneamente: %d\n",
		maxEating,
	)

	fmt.Printf(
		"Total de ciclos concluídos: %d\n",
		totalCycles,
	)

	fmt.Println("===================================")
}

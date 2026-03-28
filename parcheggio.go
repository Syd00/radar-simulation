package main

import (
	"fmt"
	"math/rand"
)

const Alpha = 3
const Pmax = 0.89
const O = 0.8

type Radar struct {
	Timestamp int
	TaskId    int
	Distance  float64
	X         float64
	Y         float64
	RCS       float64
	SNR       float64
	Speed     float64
}

// Genera il task id
// return 1 se classe dominante, 2 altrimenti
func GenerateTaskId(pMax float64) int {
	r := rand.Float64()

	if r < pMax {
		return 1
	}
	return 2
}

func generateRadarPacket() {
	task := GenerateTaskId(Pmax)
	fmt.Printf("%d", task)
}

func main() {
	for range 10 {
		generateRadarPacket()
	}
}

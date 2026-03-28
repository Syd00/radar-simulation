package main

import (
	"fmt"
	"math/rand"
	"time"
)

const Alpha = 3
const Pmax = 0.89
const Omitted = 0.8
const Offset = 0.1

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

func generateRadarPacket() Radar {
	task := GenerateTaskId(Pmax)
	var radar Radar

	if task == 1 {
		radar.Timestamp = int(time.Now().UnixMilli())
		radar.TaskId = task
		radar.Distance = 0
		radar.X = 0
		radar.Y = 0
		radar.RCS = 0
		radar.SNR = 0
		radar.Speed = 0
	} else {
		radar.Timestamp = int(time.Now().UnixMilli())
		radar.TaskId = task
		radar.Distance = 0
		radar.X = 0
		radar.Y = 0
		radar.RCS = 0
		radar.SNR = 0
		radar.Speed = 0
	}

	return radar
}

func main() {
	for range 10 {
		radar := generateRadarPacket()
		fmt.Printf("%+v\n", radar)
	}
}

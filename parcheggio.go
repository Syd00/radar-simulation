package main

import (
	"cmp"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const Alpha = 3
const Pmax = 0.89
const Omitted = 0.8
const Offset = 0.1
const radarMaxRange = 50

type Radar struct {
	Timestamp int
	Distance  float64
	Theta     float64
	X         float64
	Y         float64
	Speed     float64
	RCS       float64 // m^2
	SNR       float64 // dB
	TaskId    int
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

	if task == 1 { // libero: il radar non rileva oggetti
		radar.Timestamp = int(time.Now().UnixMilli())
		radar.Distance = -1.0
		radar.Theta = 0.0
		radar.X, radar.Y = calculatePosition(radar.Distance, radar.Theta) // no target                                                   // no target
		// radar.Z = 0.0 // no target
		radar.RCS = rand.Float64() * 0.01 // rumore di fondo [0 - 0.01]
		radar.SNR = rand.Float64() * 0.05 // rumore di fondo [0 - 0.05]
		radar.Speed = 0.0                 // no target
		radar.TaskId = classifyTask(radar.Distance, radar.X, radar.Y, radar.RCS, radar.SNR, radar.Speed)
	} else { // static: il radar rivela un oggetto statico
		radar.Timestamp = int(time.Now().UnixMilli())
		radar.TaskId = task
		radar.Distance = 0
		radar.Theta = (rand.Float64()*2 - 1) * 60
		radar.X = 0
		radar.Y = 0
		radar.RCS = 0
		radar.SNR = 0
		radar.Speed = 0
	}

	return radar
}

func classifyTask(distance, x, y, rcs, snr, speed float64) int {
	var task int
	switch {
	case cmp.Compare(distance, -1.0) == 0,
		cmp.Compare(x, 0.0) == 0,
		cmp.Compare(y, 0.0) == 0,
		rcs < 0.02,
		snr < 0.1,
		cmp.Compare(speed, 0.0) == 0:
		task = 1
	default:
		task = 2
	}
	return task
}

func calculatePosition(distance float64, theta float64) (float64, float64) {
	if cmp.Compare(distance, -1) != 0 {
		sin, cos := math.Sincos(theta)
		xPos := distance * sin
		yPos := distance * cos
		return xPos, yPos
	} else {
		return 0.0, 0.0
	}
}

func main() {
	for range 10 {
		radar := generateRadarPacket()
		fmt.Printf("%+v\n", radar)
	}
}

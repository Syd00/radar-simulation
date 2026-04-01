package main

import (
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
	radar.Timestamp = int(time.Now().UnixMilli())

	if task == 1 { // libero: il radar non rileva oggetti
		radar.Distance = -1.0
		radar.Theta = 0.0                 // no target
		radar.RCS = rand.Float64() * 0.01 // rumore di fondo [0 - 0.01]
		radar.SNR = rand.Float64() * 0.05 // rumore di fondo [0 - 0.05]
		radar.Speed = 0.0                 // no target
	} else { // static: il radar rivela un oggetto statico
		radar.Distance = rand.Float64() * radarMaxRange // da 0 a radarMaxRange
		radar.Theta = (rand.Float64()*2 - 1) * 60       // da -60 a +60 gradi
		radar.RCS = 0.5 + rand.Float64()*2
		radar.SNR = 15.0 + rand.Float64()*20
		radar.Speed = 0
	}

	radar.X, radar.Y = calculatePosition(radar.Distance, radar.Theta) // no target
	radar.TaskId = classifyTask(radar.Distance, radar.RCS, radar.SNR)

	return radar
}

func classifyTask(distance, rcs, snr float64) int {
	// il radar non trova niente
	if distance == -1.0 {
		return 1
	}

	// il radar trova qualcosa ma il segnale è troppo debole o il cross section troppo piccolo
	if snr < 0.1 || rcs < 0.02 {
		return 1
	}

	// è un oggetto reale
	return 2
}

func calculatePosition(distance float64, theta float64) (float64, float64) {
	if distance == -1.0 {
		return 0.0, 0.0
	}

	sin, cos := math.Sincos(theta * (math.Pi / 180))
	xPos := distance * sin
	yPos := distance * cos
	return xPos, yPos
}

func main() {
	for range 10 {
		radar := generateRadarPacket()
		fmt.Printf("%+v\n", radar)
	}
}

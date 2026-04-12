package main

import (
	"math"
	"math/rand"
	"time"
)

const RadarMaxRange = 50.0
const pMax = 0.5 // before 0.89 // probabilità classe dominante

// Radar rappresenta la struttura dati del sensore
type Radar struct {
	radarID   int
	Timestamp int64
	Range     float64
	Theta     float64
	X         float64
	Y         float64
	Speed     float64
	Rcs       float64
	Snr       float64
	TaskID    int
}

func newRadar(id int) Radar {
	return Radar{
		radarID: id,
	}
}

func (r *Radar) Equal(s Radar) bool {
	const Epsilon = 1e-9
	return math.Abs(r.Range-s.Range) < Epsilon &&
		math.Abs(r.Theta-s.Theta) < Epsilon &&
		math.Abs(r.X-s.X) < Epsilon &&
		math.Abs(r.Y-s.Y) < Epsilon &&
		math.Abs(r.Rcs-s.Rcs) < Epsilon &&
		math.Abs(r.Snr-s.Snr) < Epsilon &&
		r.TaskID == s.TaskID
}

// generateRadarScan simula il comportamento del processore mmWave
func GenerateRadarScan(r float64, radar Radar) Radar {

	var pOmitted = 0.5 // before 0.80

	radar.Timestamp = time.Now().Unix()

	if rand.Float64() > pOmitted { // radar will compute
		if r > pMax { // static object (task 2)
			theta := -60.0 + r*(60.0-(-60.0))
			rad := theta * (math.Pi / 180.0)
			rRange := 0.5 + r*(RadarMaxRange-0.5)
			X := math.Round(rRange*math.Sin(rad)*100) / 100
			Y := math.Round(rRange*math.Cos(rad)*100) / 100

			radar.Range = rRange
			radar.Theta = theta
			radar.X = X
			radar.Y = Y
			radar.Speed = 0.0
			radar.Rcs = 0.1 + r*(10.0-0.1)
			radar.Snr = 10.0 + r*(30.0-10.0)
		} else { // no object (task 1)
			radar.Range = -1.0
			radar.Theta = 0.0
			radar.X = 0.0
			radar.Y = 0.0
			radar.Speed = 0.0
			radar.Rcs = r * 0.09
			radar.Snr = r * 9.99
		}
		radar.TaskID = classifyTask(radar)
	} else { // radar will omit, returning task 1 class
		radar.Range = -1.0
		radar.Theta = 0.0
		radar.X = 0.0
		radar.Y = 0.0
		radar.Speed = 0.0
		radar.Rcs = rand.Float64() * 0.09
		radar.Snr = rand.Float64() * 9.99
		radar.TaskID = 1
	}

	return radar
}

func classifyTask(r Radar) int {
	if r.Range == -1.0 &&
		r.Theta == 0.0 &&
		r.X == 0.0 &&
		r.Y == 0.0 &&
		r.Rcs < 0.1 &&
		r.Snr < 10 {
		return 1
	}
	return 2
}

func verifySentinel(sentinel, scan Radar) bool {
	return sentinel.Equal(scan)
}

/* 	// Generatore per la logica di simulazione (non deterministico)
   	rLogic := rand.New(rand.NewSource(time.Now().UnixNano()))

   	willCompute := rLogic.Float64() > pOmitted
   	targetExists := rLogic.Float64() > pMax

   	currentTask := 1
   	if willCompute && targetExists {
   		currentTask = 2
   	}

   	// Se omette o non c'è target, restituisce default
   	if !willCompute || currentTask == 1 {
   		return Radar{
   			Timestamp: time.Now().Unix(),
   			Range:     -1.0,
   			Theta:     0.0,
   			X:         0.0,
   			Y:         0.0,
   			Speed:     0.0,
   			Rcs:       rLogic.Float64() * 0.09,
   			Snr:       rLogic.Float64() * 9.99,
   			TaskID:    1,
   		}
   	} */

/* // Se computa (Task 2), usiamo il seed deterministico per la "firma"
// Nota: rand.NewSource accetta int64, convertiamo il seed float64
rVerify := rand.New(rand.NewSource(int64(seedChallenge)))

//vRange := 0.5 + rVerify.Float64()*(RadarMaxRange-0.5)
//vTheta := -60.0 + rVerify.Float64()*(60.0-(-60.0))
//vSnr := 10.0 + rVerify.Float64()*(30.0-10.0)
//vRcs := 0.1 + rVerify.Float64()*(10.0-0.1)

// Calcolo coordinate cartesiane
rad := vTheta * (math.Pi / 180.0)
vX := math.Round(vRange*math.Sin(rad)*100) / 100
vY := math.Round(vRange*math.Cos(rad)*100) / 100

return Radar{
	Timestamp: time.Now().Unix(),
	Range:     vRange,
	Theta:     vTheta,
	X:         vX,
	Y:         vY,
	Speed:     0.0,
	Rcs:       vRcs,
	Snr:       vSnr,
	TaskID:    currentTask,
} */

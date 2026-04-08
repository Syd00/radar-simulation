package radar

import (
	"math"
	"math/rand"
	"time"
)

const RadarMaxRange = 50.0

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

// generateRadarScan simula il comportamento del processore mmWave
func GenerateRadarScan(seed float64) Radar {
	const pMax = 0.5     // before 0.89
	const pOmitted = 0.5 // before 0.80
	var r = rand.New(rand.NewSource(0))

	if rand.Float64() > pOmitted { // radar will compute
		if seed == 0.0 { // non sentinel, random seed
			r = rand.New(rand.NewSource(rand.Int63()))
		} else { // sentinel, predetermined seed
			r = rand.New(rand.NewSource(int64(seed)))
		}

		if rand.Float64() > pMax { // static object (task 2)
			theta := -60.0 + r.Float64()*(60.0-(-60.0))
			rad := theta * (math.Pi / 180.0)
			rRange := 0.5 + r.Float64()*(RadarMaxRange-0.5)
			X := math.Round(rRange*math.Sin(rad)*100) / 100
			Y := math.Round(rRange*math.Cos(rad)*100) / 100

			return Radar{
				Timestamp: time.Now().Unix(),
				Range:     rRange,
				Theta:     theta,
				X:         X,
				Y:         Y,
				Speed:     0.0,
				Rcs:       0.1 + r.Float64()*(10.0-0.1),
				Snr:       10.0 + r.Float64()*(30.0-10.0),
				TaskID:    2,
			}
		}
		// no object (task 1)
		return Radar{
			Timestamp: time.Now().Unix(),
			Range:     -1.0,
			Theta:     0.0,
			X:         0.0,
			Y:         0.0,
			Speed:     0.0,
			Rcs:       r.Float64() * 0.09,
			Snr:       r.Float64() * 9.99,
			TaskID:    1,
		}
	} else { // radar will omit, returning task 1 class
		return Radar{
			Timestamp: time.Now().Unix(),
			Range:     -1.0,
			Theta:     0.0,
			X:         0.0,
			Y:         0.0,
			Speed:     0.0,
			Rcs:       rand.Float64() * 0.09,
			Snr:       rand.Float64() * 9.99,
			TaskID:    1,
		}
	}
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

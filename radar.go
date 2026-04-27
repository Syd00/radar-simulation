package main

import (
	"math"
	"math/rand"
	"time"
)

const RadarMaxRange = 50.0
const pMax = 0.5   // before 0.89 // probabilità classe dominante
var pOmitted = 0.5 // before 0.80 // probabilità di saltare computazione

// Radar rappresenta la struttura dati del sensore
type Radar struct {
	RadarID   int
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
		RadarID: id,
	}
}

func (r *Radar) Equal(s Radar) bool {
	const Epsilon = 0.5
	return math.Abs(r.Range-s.Range) < Epsilon &&
		math.Abs(r.Theta-s.Theta) < Epsilon &&
		math.Abs(r.X-s.X) < Epsilon &&
		math.Abs(r.Y-s.Y) < Epsilon &&
		math.Abs(r.Rcs-s.Rcs) < Epsilon &&
		math.Abs(r.Snr-s.Snr) < Epsilon &&
		r.TaskID == s.TaskID
}

func addNoise(val float64, intensity float64) float64 {
	return val + (rand.NormFloat64() * intensity)
}

// generateRadarScan simula il comportamento del processore mmWave
func GenerateRadarScan(r float64, radar Radar) Radar {
	radar.Timestamp = time.Now().Unix()

	if r == 0.0 {
		r = rand.Float64()
	}

	if rand.Float64() > pOmitted { // radar will compute
		if r > pMax { // static object (task 2)
			theta := -60.0 + r*(60.0-(-60.0))
			rad := theta * (math.Pi / 180.0)
			rRange := 0.5 + r*(RadarMaxRange-0.5)
			X := math.Round(rRange*math.Sin(rad)*100) / 100
			Y := math.Round(rRange*math.Cos(rad)*100) / 100
			radar.Range = addNoise(rRange, 0.2)
			radar.Theta = theta
			radar.X = X
			radar.Y = Y
			radar.Speed = 0.0
			radar.Rcs = addNoise(0.1+r*(10.0-0.1), 0.05)
			radar.Snr = addNoise(10.0+r*(30.0-10.0), 0.5)
		} else { // no object (task 1)
			radar.Range = -1.0
			radar.Theta = 0.0
			radar.X = 0.0
			radar.Y = 0.0
			radar.Speed = 0.0
			radar.Rcs = addNoise(r*0.09, 0.05)
			radar.Snr = addNoise(r*9.99, 0.5)
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
	const Epsilon = 1e-9

	if math.Abs(r.Range-(-1.0)) < Epsilon &&
		math.Abs(r.Theta-0.0) < Epsilon &&
		math.Abs(r.X-0.0) < Epsilon &&
		math.Abs(r.Y-0.0) < Epsilon &&
		r.Rcs < 0.1 &&
		r.Snr < 10.0 {
		return 1
	}
	return 2
}

func verifySentinel(sentinel, scan Radar) bool {
	return sentinel.Equal(scan)
}

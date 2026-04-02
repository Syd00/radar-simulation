package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const RadarMaxRange = 50.0

// Radar rappresenta la struttura dati del sensore
type Radar struct {
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
func generateRadarScan(seedChallenge float64) Radar {
	const pMax = 0.89
	const pOmitted = 0.80

	// Generatore per la logica di simulazione (non deterministico)
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
	}

	// Se computa (Task 2), usiamo il seed deterministico per la "firma"
	// Nota: rand.NewSource accetta int64, convertiamo il seed float64
	rVerify := rand.New(rand.NewSource(int64(seedChallenge)))

	vRange := 0.5 + rVerify.Float64()*(RadarMaxRange-0.5)
	vTheta := -60.0 + rVerify.Float64()*(60.0-(-60.0))
	vSnr := 10.0 + rVerify.Float64()*(30.0-10.0)
	vRcs := 0.1 + rVerify.Float64()*(10.0-0.1)

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
	}
}

// verifyComputation controlla se i dati corrispondono alla firma del seed
func verifyComputation(data Radar, seed float64) bool {
	if data.Range == -1.0 {
		return false
	}

	// Ricostruiamo la sequenza deterministica
	rVerify := rand.New(rand.NewSource(int64(seed)))

	expectedRange := 0.5 + rVerify.Float64()*(RadarMaxRange-0.5)
	expectedTheta := -60.0 + rVerify.Float64()*(60.0)
	expectedSnr := 10.0 + rVerify.Float64()*(30.0-10.0)
	expectedRcs := 0.1 + rVerify.Float64()*(10.0-0.1)

	// Verifica con tolleranza (Epsilon)
	const eps = 0.1
	rangeOk := math.Abs(data.Range-expectedRange) < 0.01
	thetaOk := math.Abs(data.Theta-expectedTheta) < eps
	snrOk := math.Abs(data.Snr-expectedSnr) < eps
	rcsOk := math.Abs(data.Rcs-expectedRcs) < 0.01

	return rangeOk && thetaOk && snrOk && rcsOk
}

func main() {
	seed := 64.0

	fmt.Printf("%-3s | %-6s | %-8s | %-8s | %-20s\n", "N", "Task", "Range", "SNR", "Verifica")
	fmt.Println("-------------------------------------------------------------")

	for i := range 20 {
		res := generateRadarScan(seed)

		status := "OMISSIONE RILEVATA"
		if res.Range == -1.0 {
			status = "IDLE (T1)"
		} else if verifyComputation(res, seed) {
			status = "CALCOLO OK (FIRMA)"
		}

		fmt.Printf("%02d  | %-6d | %-8.2f | %-8.2f | %-20s\n",
			i, res.TaskID, res.Range, res.Snr, status)
	}
}

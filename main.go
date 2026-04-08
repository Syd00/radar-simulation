package main

import (
	"fmt"
	"math"
	"math/rand"
	"progetto/radar"
)

const RadarMaxRange = 50.0

// verifyComputation controlla se i dati corrispondono alla firma del seed
func verifyComputation(data radar.Radar, seed float64) bool {
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
	seed := 0.0

	for i := range 10 {
		res := radar.GenerateRadarScan(seed)

		fmt.Printf("%d: %#v", i, res)
		fmt.Print("\n--------------------------------\n")
	}
}

package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func generateSentinel(r float64) Radar {
	radar := newRadar(10)
	radar.Timestamp = time.Now().Unix()
	var pMax = 0.5
	if r > pMax { // task2 sentinel

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
		radar.TaskID = classifyTask(radar)
	} else {
		radar.Range = -1.0
		radar.Theta = 0.0
		radar.X = 0.0
		radar.Y = 0.0
		radar.Speed = 0.0
		radar.Rcs = r * 0.09
		radar.Snr = r * 9.99
		radar.TaskID = classifyTask(radar)
	}
	return radar
}

// verifyComputation controlla se i dati corrispondono alla firma del seed
/* func verifyComputation(data Radar, seed float64) bool {
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
} */

func main() {
	net := NewNetwork(5)
	t2Seed := 7 // TaskID 2 (oggetto presente)
	var r = rand.New(rand.NewSource(0))
	if float64(t2Seed) == 0.0 { // non sentinel, random seed
		r = rand.New(rand.NewSource(rand.Int63()))
	} else { // sentinel, predetermined seed
		r = rand.New(rand.NewSource(int64(t2Seed)))
	}
	random := r.Float64()
	fmt.Println(random)
	//t1Seed := 1 // TaskID 1 (no target)

	for i := 0; i < 1; i++ { // Simuliamo 10 cicli di scansione
		for _, rad := range net.Nodes {
			// 1. Il radar esegue lo scan (può essere pigro o meno)
			scan := GenerateRadarScan(random, rad)

			// 2. Il Monitor genera la "Verità" (un radar che NON può essere pigro)
			// trucco: per la verità chiamiamo GenerateRadarScan con pOmitted temporaneamente a 0
			truth := generateSentinel(random)

			// 3. Confronto
			if !verifySentinel(truth, scan) {
				fmt.Printf("[ALERT] Radar %d ha barato! Messaggio: %s\n", scan.radarID)
				fmt.Printf("%#v\n%#v\n---\n", scan, truth)
			} else {
				fmt.Printf("[OK] Radar %d ha computato correttamente\n", scan.radarID)
				fmt.Printf("%#v\n%#v\n---\n", scan, truth)
			}
		}
	}
}

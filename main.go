package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type RadarExport struct {
	Radar
	IsSentinel         bool
	PassedVerification bool
}

func generateSeed(targetTask int) float64 {
	// Assumiamo pMax = 0.5
	if targetTask == 2 {
		// Restituisce un valore tra 0.500001 e 1.0
		return 0.51 + (rand.Float64() * 0.48)
	}
	// Restituisce un valore tra 0.0 e 0.499999
	return rand.Float64() * 0.49
}

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

func saveToCSV(filename string, data []RadarExport) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("sep=,\n")
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Scrittura dell'Header (intestazione delle colonne)
	header := []string{"RadarID", "Timestamp", "Range", "Theta", "X", "Y", "Rcs", "Snr", "TaskID", "IsSentinel", "PassedVerification"}
	writer.Write(header)

	// Scrittura dei dati
	for _, r := range data {
		row := []string{
			strconv.Itoa(r.RadarID),
			strconv.FormatInt(r.Timestamp, 10),
			strconv.FormatFloat(r.Range, 'f', 4, 64),
			strconv.FormatFloat(r.Theta, 'f', 4, 64),
			strconv.FormatFloat(r.X, 'f', 4, 64),
			strconv.FormatFloat(r.Y, 'f', 4, 64),
			strconv.FormatFloat(r.Rcs, 'f', 4, 64),
			strconv.FormatFloat(r.Snr, 'f', 4, 64),
			strconv.Itoa(r.TaskID),
			strconv.FormatBool(r.IsSentinel),
			strconv.FormatBool(r.PassedVerification),
		}
		writer.Write(row)
	}
	return nil
}

func main() {
	numRadars := 5
	net := NewNetwork(numRadars)
	totalJobs := 200
	sentinelRate := 0.05                                   // percentuale di sentinels sui jobs totali
	numSentinels := int(float64(totalJobs) * sentinelRate) // numero sentinels

	var exportData []RadarExport

	fmt.Printf("Avvio simulazione: %d jobs, %d sentinelle\n", totalJobs, numSentinels)

	for i := range totalJobs { // Simuliamo 10 cicli di scansione

		radarId := i % numRadars
		radar := net.Nodes[radarId]

		var currentSeed float64
		isSentinel := i < numSentinels

		if isSentinel {
			targetTask := 1
			if i%2 == 0 {
				targetTask = 2
			}
			currentSeed = generateSeed(targetTask)
		} else {
			currentSeed = 0.0
		}

		scan := GenerateRadarScan(currentSeed, radar)
		entry := RadarExport{
			Radar:              scan,
			IsSentinel:         isSentinel,
			PassedVerification: true,
		}

		if isSentinel {
			truth := generateSentinel(currentSeed)
			passed := verifySentinel(truth, scan)
			entry.PassedVerification = passed
			if !passed {
				fmt.Printf("[ALERT] Sentinella FALLITA su Radar %d (Task Atteso: %d)\n",
					scan.RadarID, truth.TaskID)
			} else {
				fmt.Printf("[OK] Radar %d ha superato il test (Task: %d)\n",
					scan.RadarID, scan.TaskID)
			}
		}
		exportData = append(exportData, entry)
	}
	err := saveToCSV("radar_data.csv", exportData)
	if err != nil {
		fmt.Println("Errore nel salvataggio del file:", err)
	} else {
		fmt.Println("Dati salvati correttamente in radar_data.csv")
	}
}

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/fewstera/contact-tracing-crypto/pkg/tracing"
)

func main() {
	numberOfPeople := 50000
	numberOfKeysPerPerson := 10
	dailyKeys, err := generateDailyKeysForPeople(numberOfPeople, numberOfKeysPerPerson)
	if err != nil {
		panic(fmt.Errorf("Error generating daily keys: %w", err))
	}

	dividedDailyKeys := [][]tracing.DailyTracingKey{}
	numProcesses := runtime.NumCPU()
	chunkSize := (len(dailyKeys) + numProcesses - 1) / numProcesses
	for i := 0; i < len(dailyKeys); i += chunkSize {
		end := i + chunkSize

		if end > len(dailyKeys) {
			end = len(dailyKeys)
		}

		dividedDailyKeys = append(dividedDailyKeys, dailyKeys[i:end])
	}

	var mux sync.Mutex
	startTime := time.Now()
	progress := 0
	go func() {

		for {
			mux.Lock()
			progressPercent := float64(progress) / float64(len(dailyKeys)) * 100
			mux.Unlock()

			fmt.Printf("Processessing daily keys (%.f%%)\n", progressPercent)
			time.Sleep(time.Duration(1) * time.Second)
		}
	}()

	var wg sync.WaitGroup
	for _, dKeys := range dividedDailyKeys {
		wg.Add(1)

		go func(dKeys []tracing.DailyTracingKey, progress *int) {
			// Generate proximity tokens for each daily key
			for _, dailyKey := range dKeys {
				mux.Lock()
				*progress = *progress + 1
				mux.Unlock()

				for i := 0; i < 143; i++ {
					dailyKey.ProximityIdentifier(uint8(i))
				}
			}
			wg.Done()
		}(dKeys, &progress)
	}

	wg.Wait()

	duration := time.Now().Sub(startTime)
	fmt.Printf("Took %.2f seconds\n", float64(duration.Milliseconds())/1000)
}

func generateDailyKeysForPeople(numberOfPeople, numberOfKeysPerPerson int) ([]tracing.DailyTracingKey, error) {
	todayDayNumber := uint32(time.Now().Unix() / (60 * 60 * 24))
	dayNumbers := []uint32{}

	// Generate day numbers for past 10 days
	for x := 0; x < numberOfKeysPerPerson; x++ {
		dayNumbers = append(dayNumbers, todayDayNumber-uint32(x))
	}

	dailyKeys := []tracing.DailyTracingKey{}
	for i := 0; i < numberOfPeople; i++ {
		p, err := tracing.GeneratePerson()
		if err != nil {
			return nil, fmt.Errorf("generating new person %d: %w", i, err)
		}

		for _, dayNumber := range dayNumbers {
			dailyTracingKey, err := p.DailyTracingKey(dayNumber)
			if err != nil {
				return nil, fmt.Errorf("generating new person %d, daily key %d: %w", i, dayNumber, err)
			}
			dailyKeys = append(dailyKeys, dailyTracingKey)
		}
	}

	return dailyKeys, nil
}

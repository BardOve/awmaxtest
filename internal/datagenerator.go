package internal

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

func GenerateData() error {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	number := r.Float64() * 30

	data := make(InputData, 0, 50)

	for i := 0; i < 50; i++ {
		data = append(data, struct {
			Time   float64 `json:"time"`
			Weight float64 `json:"weight"`
		}{
			Time:   float64(i) / 10.0,
			Weight: number,
		})

		if number < 50 {
			number += r.Float64() * 5
		} else {
			number -= r.Float64() * 15
		}
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("./internal/data/input.json", jsonData, 0644)
}

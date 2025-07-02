package internal

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"net/url"
	"os"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

var sensorFilePath = "./internal/data/input.json"

type InputData []struct {
	Time   float64 `json:"time"`
	Weight float64 `json:"weight"`
}


func GetData(filepath string) (InputData, error) {
	
	var dataFromFile InputData

	inputBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil,	err
	}

	json.Unmarshal(inputBytes, &dataFromFile)
	return dataFromFile, nil
}

func GenerateChartHandler(c *gin.Context) {
	
	var req struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Exercise string `json:"exercise"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"success": false, "error": "Invalid request body"})
		return
	}

	sensorData, err := GetData(sensorFilePath)

	if err != nil {
		log.Printf("Error getting data: %v", err)
		c.JSON(500, gin.H{"success": false, "error": "Could not load sensor data"})
		return
	}
	
	if err := CreateChart(sensorData, req.Name, req.Exercise); err != nil {
		log.Printf("Error creating chart: %v", err)
		c.JSON(500, gin.H{"success": false, "error": "Could not create chart image"})
		return
	}

	chartDir := fmt.Sprintf("./charts/%s", req.Name)

	files, err := os.ReadDir(chartDir)
	if err != nil {
		log.Printf("Error reading chart directory: %v", err)
		c.JSON(500, gin.H{"success": false, "error": "Could not list charts"})
		return
	}

	var chartUrls []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".png") {
			url := fmt.Sprintf("/charts/%s/%s", url.PathEscape(req.Name), url.PathEscape(file.Name()))
			chartUrls = append(chartUrls, url)
		}
	}
	c.JSON(200, gin.H{"success": true, "chartUrls": chartUrls})
}

func CreateChart(data InputData, name string, excercise string) error {

	points := make(plotter.XYs, len(data))
	for i, record := range data {
		points[i].X = record.Time
		points[i].Y = record.Weight
	}

	
	var peakWeight float64
	var totalWeight float64
	var avgWeight float64

	if len(data) > 0 {
		peakWeight = data[0].Weight 
		for _, record := range data {
			if record.Weight > peakWeight {
				peakWeight = record.Weight
			}
			totalWeight += record.Weight
		}
		avgWeight = totalWeight / float64(len(data))
	}

	fmt.Printf("Creating chart with %d data points...\n", len(points))
	p := plot.New()
	p.Title.Text = fmt.Sprintf("%s Force Over Time for %s", excercise, name)
	p.X.Label.Text = "Time (seconds)"
	p.Y.Label.Text = "Force (kg)"
	p.Add(plotter.NewGrid())

	
	line, err := plotter.NewLine(points)
	if err != nil {
		return err 
	}
	line.Color = color.RGBA{R: 255, A: 255}
	line.Width = vg.Points(2)


	p.Add(line)
	p.Legend.Add(fmt.Sprintf("Peak Force: %.2f kg", peakWeight))
	p.Legend.Add(fmt.Sprintf("Avg Force: %.2f kg", avgWeight))

	
	if err := os.MkdirAll("./charts/"+name, 0755); err != nil {
		return fmt.Errorf("failed to create chart directory: %w", err)
	}
	date := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("./charts/%s/%s_%s_%s.png", name, name, excercise, date)

	if err := p.Save(10*vg.Inch, 4*vg.Inch, filename); err != nil {
		return fmt.Errorf("failed to save chart: %w", err)
	}
	os.Rename(sensorFilePath, fmt.Sprintf("./charts/%s/%s_%s_%s.json", name, name, excercise, date))

	log.Printf("Chart '%s' created successfully.", filename)
	return nil
}

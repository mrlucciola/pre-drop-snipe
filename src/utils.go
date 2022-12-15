package main

import (
	"log"
	"math"
	"os"
)

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

var logger *log.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

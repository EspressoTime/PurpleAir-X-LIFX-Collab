package main

import (
	// "fmt"
	"log"
	"os"
	"time"
)

func main() {
	f, err := os.OpenFile("/root/AQI.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(f)

	revertColor, err := listLights()
	if err != nil {
		log.Fatalln(err)
	}
	finalAQI := fetchPA()
	// finalAQI := float64(60)
	// fmt.Println(revertColor)

	var aqiColor string
	switch {
	case finalAQI < 10:
		aqiColor = "#0ed2c9" //light blue
		log.Printf("Lightblue AQI: %f", finalAQI)
	case finalAQI < 30:
		aqiColor = "#1a920b" //green
		log.Printf("Green AQI: %f", finalAQI)
	case finalAQI < 65:
		aqiColor = "hue:47.8125 saturation:1 brightness:1 kelvin:2500" //yellow
		log.Printf("Yellow AQI: %f", finalAQI)
	case finalAQI < 100:
		aqiColor = "hue:32.3383 saturation:1 brightness:1 kelvin:2500" //orange
		log.Printf("Orange AQI: %f", finalAQI)
	case finalAQI < 150:
		aqiColor = "hue:350 saturation:1 brightness:1 kelvin:2500" //red
		log.Printf("Red AQI: %f", finalAQI)
	default:
		aqiColor = "hue:350 saturation:1 brightness:0.15 kelvin:2500" //dark
		log.Printf("Awful ;_; AQI: %f", finalAQI)
	}
	// fmt.Println(aqiColor)

	setColor(aqiColor)

	time.Sleep(30 * time.Second)

	setColor(revertColor)
}

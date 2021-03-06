package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type KeyStruct struct {
	Key string `json:"key"`
}

type LightsStruct struct {
	Label  string      `json:"label"`
	Color  ColorStruct `json:"color"`
	Bright float64     `json:"brightness"`
	Power  string      `json:"power"`
}

type ColorStruct struct {
	Hue  float64 `json:"hue"`
	Sat  float64 `json:"saturation"`
	Kelv int     `json:"kelvin"`
}

func configKey() string {
	confContent, err := ioutil.ReadFile("/root/air/config.json")
	if err != nil {
		log.Fatal(err)
	}
	conf := &KeyStruct{}
	err = json.Unmarshal(confContent, conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf.Key
}

func listLights() (string, error) {
	key := configKey()
	bearer := fmt.Sprintf("Bearer %s", key)
	// ***List Lights***
	req, err := http.NewRequest("GET", "https://api.lifx.com/v1/lights/all", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", bearer)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	listBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	data := string(listBody)
	lightArr := make([]LightsStruct, 0)
	json.Unmarshal([]byte(data), &lightArr)

	if len(lightArr) < 1 {
		return "", errors.New("no lights found")
	}

	if lightArr[0].Power == "off" {
		return "", errors.New("light is off")
	}
	revertColor := fmt.Sprintf(
		"hue:%f saturation:%f brightness:%f kelvin:%d", 
		lightArr[0].Color.Hue, lightArr[0].Color.Sat, 
		lightArr[0].Bright, lightArr[0].Color.Kelv)
	// revertColor := `hue:` + fmt.Sprint(lightArr[0].Color.Hue) + ` saturation:` + fmt.Sprint(lightArr[0].Color.Sat) + ` brightness:` + fmt.Sprint(lightArr[0].Bright) + ` kelvin:` + fmt.Sprint(lightArr[0].Color.Kelv)
	// fmt.Println(revertColor)
	log.Printf("revertColor: %s", revertColor)

	return revertColor, nil
}

func setColor(paColor string) {
	key := configKey()
	bearer := fmt.Sprintf("Bearer %s", key)
	// ***Set Color***
	params := fmt.Sprintf("power=on&color=%s&duration=0", paColor)
	body := strings.NewReader(params)
	req, err := http.NewRequest("PUT", "https://api.lifx.com/v1/lights/all/state", body)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
}

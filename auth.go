package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Token string `json:"token"`
	Port  int    `json:"port"`
	Host  string `json:"host"`
}

var config Config

func configInit() {
	// READ FROM A .json FILE
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened config.json")
	byteValue, _ := io.ReadAll(jsonFile)
	//fmt.Println(byteValue)
	json.Unmarshal(byteValue, &config)
	fmt.Println(config.Token)
	defer jsonFile.Close()
}

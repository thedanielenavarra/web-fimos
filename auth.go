package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Config struct {
	Token string `json:"token"`
	Port  int    `json:"port"`
	Host  string `json:"host"`
}

var config Config

func configInit(filepath string) {
	// READ FROM A .json FILE
	jsonFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened config.json")
	byteValue, _ := io.ReadAll(jsonFile)
	//fmt.Println(byteValue)
	json.Unmarshal(byteValue, &config)
	defer jsonFile.Close()
}

func tokenRegeneration(hostname string, port int) {
	// WRITE TO A .json FILE
	jsonFile, err := os.Create("config.json")
	if err != nil {
		fmt.Println(err)
	}
	// generate a random token
	var seedRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	// a random number from 16 to 24
	tokenLength := seedRand.Intn(8) + 16
	b := make([]byte, tokenLength)
	for i := range b {
		b[i] = charset[seedRand.Intn(len(charset))]
	}
	log.Println("Generated token: ", string(b))
	hasher := sha512.New()
	hasher.Write(b)
	config.Token = hex.EncodeToString(hasher.Sum(nil))
	config.Port = port
	config.Host = hostname
	// write the token to the file
	jsonData, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err)
	}
	jsonFile.Write(jsonData)
	defer jsonFile.Close()
}

func verifyToken(token string) bool {
	// check if the token is the sha512 of the token in the config file
	hasher := sha512.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil)) == config.Token
}

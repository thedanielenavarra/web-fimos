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
	jsonFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)
	if config.Token == ""  && !*regen{
		tokenRegeneration()
	}
	defer jsonFile.Close()
}

func tokenRegeneration() {
	var seedRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	tokenLength := seedRand.Intn(8) + 16
	b := make([]byte, tokenLength)
	for i := range b {
		b[i] = charset[seedRand.Intn(len(charset))]
	}
	log.Println("Generated token: ", string(b))
	hasher := sha512.New()
	hasher.Write(b)
	config.Token = hex.EncodeToString(hasher.Sum(nil))
}

func writeConfig() {
	jsonFile, err := os.Create(*configFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully written ", *configFile)
	byteValue, _ := json.Marshal(config)
	jsonFile.Write(byteValue)
	defer jsonFile.Close()
}

func verifyToken(token string) bool {
	hasher := sha512.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil)) == config.Token
}

package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
)

func main() {
	log.Println("Starting KeyService")
	http.HandleFunc("/healthCheck", healthHandler)
	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const length = int64(len(characters))
const keySize = 32

// GenKey creates a new crypto secure key
func GenKey() string {
	bytes := make([]byte, keySize)
	for i := range bytes {
		randBigInt, err := rand.Int(rand.Reader, big.NewInt(length))
		if err != nil {
			log.Fatal("Failed to generate random big int", err)
		}
		bytes[i] = characters[randBigInt.Int64()]
	}
	return string(bytes)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Random handler enter")
	w.Write([]byte(GenKey()))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "")
}

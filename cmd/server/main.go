package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"pcmg/pkg/game"
)

var players = make(map[string]*game.Player)
var playerAddresses = make(map[string]string)
var mutex = &sync.Mutex{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Server is running")
	})
	http.HandleFunc("/sendPublicKey", sendPublicKeyHandler)
	http.HandleFunc("/sendSignature", sendSignatureHandler)
	http.HandleFunc("/sendNumber", sendNumberHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("test")
		log.Fatal(err.Error())
		return
	}

	fmt.Println("Server running on port 8080")
}

func sendPublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("id")
	publicKey := r.URL.Query().Get("publicKey")

	mutex.Lock()
	if _, exists := players[playerID]; !exists {
		players[playerID] = game.NewPlayer(playerID)
	}
	players[playerID].SetPublicKey(publicKey)
	mutex.Unlock()

	broadcastPublicKey(playerID, publicKey)

	_, err := fmt.Fprint(w, "Public key sent successfully")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

func sendSignatureHandler(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("id")
	signature := r.URL.Query().Get("signature")

	mutex.Lock()
	if _, exists := players[playerID]; !exists {
		players[playerID] = game.NewPlayer(playerID)
	}
	players[playerID].SetSignature([]byte(signature))
	mutex.Unlock()

	broadcastSignature(playerID, signature)

	_, err := fmt.Fprint(w, "Signature sent successfully")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

func sendNumberHandler(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("id")
	numberStr := r.URL.Query().Get("number")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}
	timestamp, err := time.Parse(time.RFC3339, r.URL.Query().Get("timestamp"))
	if err != nil {
		http.Error(w, "Invalid timestamp", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	if _, exists := players[playerID]; !exists {
		players[playerID] = game.NewPlayer(playerID)
	}
	players[playerID].SetNumber(number)
	players[playerID].SetTimestamp(timestamp)
	mutex.Unlock()

	broadcastNumber(playerID, numberStr, timestamp.String())

	_, err = fmt.Fprint(w, "Number and timestamp sent successfully")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

}

func broadcastPublicKey(id, publicKey string) {
	for playerID, address := range playerAddresses {
		if playerID != id {
			sendPublicKey(address, publicKey)
		}
	}
}

func broadcastSignature(id, signature string) {
	for playerID, address := range playerAddresses {
		if playerID != id {
			sendSignature(address, signature)
		}
	}
}

func broadcastNumber(id, number, timestamp string) {
	for playerID, address := range playerAddresses {
		if playerID != id {
			sendNumber(address, number, timestamp)
		}
	}
}

func sendPublicKey(address, publicKey string) {
	values := url.Values{}
	values.Set("publicKey", publicKey)

	resp, err := http.PostForm(fmt.Sprintf("http://%s/receivePublicKey", address), values)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer resp.Body.Close()
}

func sendSignature(address, signature string) {
	values := url.Values{}
	values.Set("signature", signature)

	resp, err := http.PostForm(fmt.Sprintf("http://%s/receiveSignature", address), values)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
}

func sendNumber(address, number, timestamp string) {
	values := url.Values{}
	values.Set("number", number)
	values.Set("timestamp", timestamp)

	resp, err := http.PostForm(fmt.Sprintf("http://%s/receiveNumber", address), values)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type ExchangeRate struct {
	Koken string             `json:"static_koken"`
	Time1 time.Time          `json:"time"`
	Base  string             `json:"base_code"`
	Rates map[string]float64 `json:"conversion_rates"`
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	interval := 2

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for range ticker.C {
		resp, err := http.Get(fmt.Sprintf("apiToken"))
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()

		var result ExchangeRate
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			log.Println(err)
		}
		staticKoken := "randomToken"
		result.Koken = staticKoken
		result.Time1 = time.Now()
		collection := client.Database("exchange").Collection("rub")
		_, err = collection.InsertOne(context.Background(), result)
		if err != nil {
			log.Println(err)
		}
		log.Println("обновлено")

	}
}

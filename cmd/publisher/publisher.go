package main

import (
	"WBTech0/internal/entity"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	sc, err := stan.Connect("test-cluster", "client-2", stan.NatsURL("192.168.31.158:4222"))
	if err != nil {
		log.Fatalf("Error connecting to NATS Streaming: %v", err)
	}
	defer sc.Close()

	order := getExampleOrder()
	subject := "subject"
	id := order.OrderUID
	for i := 1; ; i++ {
		order.OrderUID = id + strconv.Itoa(i)
		order.Payment.Transaction = order.OrderUID
		order.Items[0].ChrtID = i
		orderM, _ := json.Marshal(order)
		err := sc.Publish(subject, orderM)
		fmt.Printf("send %s \n", order.OrderUID)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(2 * time.Second)
	}

}

func getExampleOrder() entity.Order {
	var ord entity.Order
	exampleFile, err := os.Open("model.json")
	jsonParser := json.NewDecoder(exampleFile)
	if err = jsonParser.Decode(&ord); err != nil {
		fmt.Println("parsing config file", err.Error())
	}
	return ord
}

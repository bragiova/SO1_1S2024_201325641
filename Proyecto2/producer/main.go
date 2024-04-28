package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

type Data struct {
	Name  string `json:"name"`
	Album string `json:"album"`
	Year  string `json:"year"`
	Rank  string `json:"rank"`
}

func getData(w http.ResponseWriter, r *http.Request) {
	var serverData Data
	err := json.NewDecoder(r.Body).Decode(&serverData)
	if err != nil {
		http.Error(w, "Fallo decode json server grpc", http.StatusBadRequest)
		return
	}

	sendKafka(serverData)
}

func sendKafka(info Data) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	topic := os.Getenv("OPIC")
	kafkaBroker := os.Getenv("KAFKABROKER")

	// Configura el productor
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBroker, // Por ejemplo, "localhost:9092"
	})
	if err != nil {
		log.Fatalf("Error al crear el productor: %v\n", err)
	}

	defer producer.Close()

	// Manejar señales de interrupción
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Loop para enviar mensajes
	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			log.Printf("Terminando: %v\n", sig)
			run = false
		default:
			jsonSended := fmt.Sprintf(`{"name":"%s","album":"%s","year":"%s","rank":"%s"}`, info.Name, info.Album, info.Year, info.Rank)

			// Enviar mensaje
			err := producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte(jsonSended),
			}, nil)

			if err != nil {
				log.Printf("Error al enviar mensaje a Kafka: %v\n", err)
			} else {
				log.Println("Mensaje enviado correctamente a Kafka")
			}
		}
	}
}

func main() {
	http.HandleFunc("/sendProducer", getData)
}

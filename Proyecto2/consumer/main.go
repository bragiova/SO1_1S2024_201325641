package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

type Vote struct {
	Name  string `json:"name"`
	Album string `json:"album"`
	Year  string `json:"year"`
	Rank  string `json:"rank"`
}

func main() {
	// Cargar variables de entorno desde .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	topicKafka := os.Getenv("TOPIC")
	brokerKafka := os.Getenv("KAFKABROKER")

	// Configura el consumidor
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokerKafka, // Por ejemplo, "localhost:9092"
		"group.id":          "group-consumer",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Error al crear el consumidor: %s\n", err)
	}

	// Suscribirse a un tema (topic)
	err = consumer.SubscribeTopics([]string{topicKafka}, nil)
	if err != nil {
		log.Fatalf("Error al suscribirse al topic: %s\n", err)
	}

	// Manejar señales de interrupción
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			log.Printf("Terminando: %v", sig)
			run = false
		default:
			// Consumir mensajes
			msg, err := consumer.ReadMessage(-1)
			if err == nil {
				fmt.Printf("Mensaje recibido: %s\n", string(msg.Value))
			} else {
				log.Printf("Error al recibir mensaje: %v\n", err)
			}
		}
	}

	// Cerrar el consumidor
	consumer.Close()
}

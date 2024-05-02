package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Vote struct {
	Name  string `json:"name"`
	Album string `json:"album"`
	Year  string `json:"year"`
	Rank  string `json:"rank"`
}

func sendInfoRedis(voteBand string, rdb *redis.Client) {
	var voto Vote
	err := json.Unmarshal([]byte(voteBand), &voto)
	if err != nil {
		fmt.Printf("Error al decodificar el JSON: %v\n", err)
		return
	}

	// Rank
	saveRankeds(voto, rdb)

}

func saveRankeds(vote Vote, rdb *redis.Client) {
	fmt.Printf("SaveRankeds Redis:")
	// Convertir el rank a float64
	rank, err := strconv.ParseFloat(vote.Rank, 64)
	if err != nil {
		fmt.Println("Error al convertir string a float64:", err)
		return
	}

	// context, hash, field, value
	// Guarda incrementa el valor del rank a un álbum, es una sumatoria de ranks por álbum
	err = rdb.HIncrByFloat(context.Background(), "Rank", vote.Album, rank).Err()
	if err != nil {
		fmt.Printf("Error al enviar datos Rank Redis: %v\n", err)
		return
	}

	// Incrementa la cantidad de votos por Album
	err = rdb.HIncrBy(context.Background(), "Album", vote.Album, 1).Err()
	if err != nil {
		fmt.Printf("Error al enviar datos Album Redis: %v\n", err)
		return
	}

	// Incrementa cantidad de votos por banda
	err = rdb.HIncrBy(context.Background(), "Band", vote.Name, 1).Err()
	if err != nil {
		fmt.Printf("Error al enviar datos Band Redis: %v\n", err)
		return
	}

	fmt.Printf("Enviado a Redis: %+v\n", vote)
}

func calcAverageAlbum(rdb *redis.Client) {
	// Obtener todos los valores almacenados con el hash "Album"
	albumCount, err := rdb.HGetAll(context.Background(), "Album").Result()
	if err != nil {
		fmt.Printf("Error al obtener los valores de Redis: %v\n", err)
		return
	}

	// Obtener todos los valores almacenados con el hash "Rank"
	rankCount, err := rdb.HGetAll(context.Background(), "Rank").Result()
	if err != nil {
		fmt.Printf("Error al obtener los valores de Redis: %v\n", err)
		return
	}

	// Calcular el promedio para cada álbum
	for album, counter := range albumCount {
		counter64, err := strconv.ParseFloat(counter, 64)
		if err != nil {
			fmt.Printf("Error al convertir el valor a float64: %v\n", err)
			continue
		}

		// Se busca el valor en el map de Ranks, para obtener la sumatoria de ranks por el album
		rank64, err := strconv.ParseFloat(rankCount[album], 64)
		if err != nil {
			fmt.Printf("Error al convertir el valor a float64: %v\n", err)
			continue
		}

		// Promedio
		average := rank64 / counter64

		// Guardar el promedio de cada álbum
		err = rdb.HSet(context.Background(), "Average", album, average).Err()
		if err != nil {
			fmt.Printf("Error al enviar datos Average Redis: %v\n", err)
			continue
		}
	}
}

func saveDataMongo(vote string) {
	uriMongo := os.Getenv("MONGO_URI")

	// Crear cliente MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uriMongo))
	if err != nil {
		fmt.Printf("Error al conectar a MongoDB: %v\n", err)
		return
	}

	defer client.Disconnect(context.Background())

	// Obtenemos colección de logs
	collection := client.Database("dbp2").Collection("logs")

	currentDate := time.Now()

	stringLog := map[string]interface{}{
		"fecha": currentDate.Format("2006-01-02"),
		"hora":  currentDate.Format("15:04:05"),
		"voto":  vote,
	}

	// Insertar el documento en la colección
	_, err = collection.InsertOne(context.Background(), stringLog)
	if err != nil {
		fmt.Printf("Error al insertar log MongoDB: %v\n", err)
		return
	}

	fmt.Printf("Log enviado a MongoDB: %v\n", stringLog)
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

	// Cerrar el consumidor
	defer consumer.Close()

	// Suscribirse a un tema (topic)
	err = consumer.Subscribe(topicKafka, nil)
	if err != nil {
		log.Fatalf("Error al suscribirse al topic: %s\n", err)
	}

	// Conectar Redis
	hostRedis := os.Getenv("REDIS_HOST")
	passRedis := os.Getenv("REDIS_PASS")
	// addressRedis := hostRedis + ":6379"

	// Crear cliente Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     hostRedis + ":6379",
		Password: passRedis,
		DB:       0,
	})

	// Cerrar la conexión redis
	defer rdb.Close()

	// Consumir mensajes
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Mensaje recibido: %s\n", string(msg.Value))
			sendInfoRedis(string(msg.Value), rdb)
			calcAverageAlbum(rdb)
			saveDataMongo(string(msg.Value))
		} else {
			log.Printf("Error al recibir mensaje: %v\n", err)
			break
		}
	}
}

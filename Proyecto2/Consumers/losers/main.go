package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"time"
)

var ctx = context.Background()

// Estructura para deserializar el mensaje JSON
type Loser struct {
	Name       string `json:"name"`
	Age        int32  `json:"age"`
	Faculty    string `json:"faculty"`
	Discipline int32  `json:"discipline"`
}

// Configuración de Redis
var redisClient *redis.Client

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Verificar conexión a Redis
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error al conectar a Redis: %s", err)
	} else {
		log.Println("Conectado a Redis exitosamente:", pong)
	}
}

func consumeLosers() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_8_0_0

	consumer, err := sarama.NewConsumer([]string{os.Getenv("KAFKA_BROKER")}, config)
	if err != nil {
		log.Fatalf("Error al crear el consumidor: %s", err)
	}
	defer consumer.Close()

	topic := "losers"
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("Error al obtener particiones: %s", err)
	}

	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Error al consumir partición %d: %s", partition, err)
		}
		defer pc.AsyncClose()

		go func(pc sarama.PartitionConsumer) {
			for {
				select {
				case msg := <-pc.Messages():
					// Deserializar el mensaje JSON para obtener la facultad
					var loser Loser
					if err := json.Unmarshal(msg.Value, &loser); err != nil {
						log.Printf("Error al deserializar el mensaje: %s", err)
						continue
					}

					// Generar clave en Redis
					var key string
					if msg.Key != nil {
						key = fmt.Sprintf("losers:%s", string(msg.Key))
					} else {
						key = fmt.Sprintf("losers:%d", time.Now().UnixNano())
					}

					// Guardar el mensaje en Redis con la clave generada
					err := redisClient.Set(ctx, key, msg.Value, 0).Err()
					if err != nil {
						log.Printf("Error al guardar en Redis: %s", err)
					} else {
						fmt.Printf("Perdedor guardado con clave %s: %s\n", key, string(msg.Value))
					}

					// Incrementar el contador de facultad en Redis (CONTEO ALUMNO POR FACULTAD)
					if loser.Faculty == "Ingenieria" {
						err = redisClient.Incr(ctx, "ing-count").Err()
					} else if loser.Faculty == "Agronomia" {
						err = redisClient.Incr(ctx, "agro-count").Err()
					}
					if err != nil {
						log.Printf("Error al incrementar el contador de %s: %s", loser.Faculty, err)
					}

					// Incrementar el contador de la facultad que mas pedio (CONTEO ALUMNO POR FACULTAD)
					if loser.Faculty == "Ingenieria" {
						err = redisClient.Incr(ctx, "loser-inge-count").Err()
					} else if loser.Faculty == "Agronomia" {
						err = redisClient.Incr(ctx, "loser-agro-count").Err()
					}
					if err != nil {
						log.Printf("Error al incrementar el contador de %s: %s", loser.Discipline, err)
					}

				case err := <-pc.Errors():
					log.Printf("Error en el consumidor: %s", err)
				}
			}
		}(pc)
	}

	select {}
}

func main() {
	initRedis()
	consumeLosers()
}

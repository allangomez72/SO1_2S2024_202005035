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
type Winner struct {
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

func consumeWinners() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_8_0_0

	consumer, err := sarama.NewConsumer([]string{os.Getenv("KAFKA_BROKER")}, config)
	if err != nil {
		log.Fatalf("Error al crear el consumidor: %s", err)
	}
	defer consumer.Close()

	topic := "winners"
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
					var winner Winner
					if err := json.Unmarshal(msg.Value, &winner); err != nil {
						log.Printf("Error al deserializar el mensaje: %s", err)
						continue
					}

					// Generar clave en Redis
					var key string
					if msg.Key != nil {
						key = fmt.Sprintf("winners:%s", string(msg.Key))
					} else {
						key = fmt.Sprintf("winners:%d", time.Now().UnixNano())
					}

					// Guardar el mensaje en Redis con la clave generada
					err := redisClient.Set(ctx, key, msg.Value, 0).Err()
					if err != nil {
						log.Printf("Error al guardar en Redis: %s", err)
					} else {
						fmt.Printf("Ganador guardado con clave %s: %s\n", key, string(msg.Value))
					}

					// Incrementar el contador de facultad en Redis esto es para CONTEO ALUMNO POR FACULTAD
					if winner.Faculty == "Ingenieria" {
						err = redisClient.Incr(ctx, "ing-count").Err()
					} else if winner.Faculty == "Agronomia" {
						err = redisClient.Incr(ctx, "agro-count").Err()
					}
					if err != nil {
						log.Printf("Error al incrementar el contador de %s: %s", winner.Faculty, err)
					}

					// Incrementar el contador de disciplina en Redis esto es para CONTEO ALUMNO POR DISCIPLINA
					if winner.Discipline == 1 {
						err = redisClient.Incr(ctx, "natacion-count").Err()
					} else if winner.Discipline == 2 {
						err = redisClient.Incr(ctx, "atletismo-count").Err()
					} else if winner.Discipline == 3 {
						err = redisClient.Incr(ctx, "boxeo-count").Err()
					}
					if err != nil {
						log.Printf("Error al incrementar el contador de %s: %s", winner.Discipline, err)
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
	consumeWinners()
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

var ctx = context.Background()

// Configuración de Redis
var redisClient *redis.Client

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Verificar conexión a Redis
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error al conectar a Redis: %s", err)
	} else {
		log.Println("Conectado a Redis exitosamente")
	}
}

// Estructura para almacenar los datos de un "loser"
type Loser struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Faculty    string `json:"faculty"`
	Discipline int    `json:"discipline"`
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
					if msg == nil {
						log.Println("Mensaje recibido es nulo")
						continue
					}

					// Mensaje de depuración para confirmar recepción de datos
					fmt.Printf("Perdedor recibido: %s\n", string(msg.Value))

					// Verificar si el valor del mensaje no está vacío
					if len(msg.Value) == 0 {
						log.Println("Advertencia: el mensaje recibido está vacío")
						continue
					}

					var loser Loser
					// Deserializar el mensaje a la estructura de "loser"
					if err := json.Unmarshal(msg.Value, &loser); err != nil {
						log.Printf("Error al deserializar el mensaje: %s", err)
						continue
					}

					// Obtener la lista actual de losers desde Redis
					var losersList []Loser
					losersData, err := redisClient.Get(ctx, "losers:all").Result()
					if err == nil {
						// Si hay datos existentes, deserializarlos
						if err := json.Unmarshal([]byte(losersData), &losersList); err != nil {
							log.Printf("Error al deserializar la lista de losers: %s", err)
						}
					}

					// Agregar el nuevo loser a la lista
					losersList = append(losersList, loser)

					// Volver a serializar y guardar la lista actualizada en Redis
					updatedData, err := json.Marshal(losersList)
					if err != nil {
						log.Printf("Error al serializar la lista actualizada: %s", err)
						continue
					}

					err = redisClient.Set(ctx, "losers:all", updatedData, 0).Err()
					if err != nil {
						log.Printf("Error al guardar en Redis (clave: losers:all): %s", err)
					} else {
						log.Println("Lista de losers actualizada en Redis")
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

package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"

	"github.com/IBM/sarama"
)

var ctx = context.Background()

// Configuración de Redis
var redisClient = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_ADDR"),
	Password: os.Getenv("REDIS_PASSWORD"),
	DB:       0,
})

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
					fmt.Printf("Ganador recibido: %s\n", string(msg.Value))
					err := redisClient.Set(ctx, fmt.Sprintf("winners:%s", string(msg.Key)), msg.Value, 0).Err()
					if err != nil {
						log.Printf("Error al guardar en Redis: %s", err)
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
	consumeWinners()
}

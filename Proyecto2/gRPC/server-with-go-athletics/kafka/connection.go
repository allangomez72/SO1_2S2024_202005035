package kafka

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"server-with-go-athletics/estructura"
	"time"

	"github.com/segmentio/kafka-go"
)

// definir si gano o perdio
func isWinner() bool {
	// Se simula el lanzamiento de una moneda
	return rand.Intn(2) == 0
}

func Produce(value estructura.Data) {
	var topic string

	//verficar si gano o perdio el estudiante
	if isWinner() {
		topic = "winners"
		log.Println("El estudiante ha ganado: ", value.Name)
	} else {
		topic = "losers"
		log.Println("El estudiante ha perdido: ", value.Name)
	}

	partition := 0

	//Conectar a kafka en kubernentes usando el service que se le dio
	conn, err := kafka.DialLeader(context.Background(), "tcp", "mi-cluster-kafka-kafka-bootstrap.kafka:9092", topic, partition)
	if err != nil {
		log.Println("Error al conectar a kafka: ", err)
	}

	// Convertir el valor a bytes
	valueByte, err := json.Marshal(value)
	if err != nil {
		log.Fatalf("Error al convertir el valor a bytes: %s", err)
	}

	// definir un timeout de 10 segundos
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: valueByte},
	)
	if err != nil {
		log.Fatalf("Error al escribir el mensaje: %s", err)
	}

	//Cerrar la conexion
	if err := conn.Close(); err != nil {
		log.Fatalf("Error al cerrar la conexion: %s", err)
	}

	log.Println("Mensaje enviado al topico:", topic)

}

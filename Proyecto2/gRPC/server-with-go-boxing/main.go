package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"server-with-go-boxing/estructura"
	"server-with-go-boxing/kafka"
	pb "server-with-go-boxing/proto-go"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50053, "Puerto del servidor")
)

// Server es usado para poder implemetntar gRPC
type server struct {
	pb.UnimplementedFacultadServiceServer
}

// Metodo que recibe la solicutud del estudiante
func (ser *server) SendUserInfo(ctx context.Context, req *pb.Student) (*pb.StudentResponse, error) {
	//mostrar los datos recibidos
	log.Printf("Recieved: %v", req)

	//Crear  una estructura con los datos recibidos
	studentData := estructura.Data{
		Name:       req.Name,
		Age:        req.Age,
		Faculty:    req.Faculty,
		Discipline: req.Discipline,
	}

	//llamar a la funcion Produce para enviar los datos a kafka
	kafka.Produce(studentData)

	return &pb.StudentResponse{
		Message: "Hola cliente recibi al estudiante",
	}, nil
}

func main() {
	//Configuracion del servidor gRPC
	flag.Parse()
	port := fmt.Sprintf(":%d", *port)
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error al escuchar en el puerto %s: %v", port, err)
	}

	srv := grpc.NewServer()
	pb.RegisterFacultadServiceServer(srv, &server{})

	//iniciar el serviodor en el puerto 50051
	log.Printf("Servidor escuchando en el peurto %s", port)
	if err := srv.Serve(listen); err != nil {
		log.Fatalf("Error al iniciar el servidor gRPC: %v", err)
	}
}

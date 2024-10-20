package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "server-whit-go/proto-go"
)

var (
	port = flag.Int("port", 50051, "Puerto del servidor")
)

// Server es usado para poder implemetntar gRPC
type server struct {
	pb.UnimplementedFacultadServiceServer
}

// Metodo que recibe la solicutud del estudiante
func (ser *server) SendStudents(ctx context.Context, req *pb.Student) (*pb.StudentResponse, error) {
	//mostrar los datos recibidos
	log.Printf("Recieved: %v", req)
	log.Printf("Student name: %s", req.Name)
	log.Printf("Studen age: %d", req.Age)
	log.Printf("Student faculty: %s", req.Faculty)
	log.Printf("Student discipline: %d", req.Discipline)

	return &pb.StudentResponse{
		Message: true,
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

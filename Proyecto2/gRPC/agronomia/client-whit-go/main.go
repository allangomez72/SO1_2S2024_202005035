package main

import (
	pb "client-whit-go/proto-go"
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// se hace el struct de como van a llegar las solicitudes HTTP y por eso se va usar fiber para mas simplicidad
type Data struct {
	Name       string `json:"name"`
	Age        int32  `json:"age"`
	Faculty    string `json:"faculty"`
	Discipline int32  `json:"discipline"`
}

func getServerForDiscipline(discipline int32) string {
	switch discipline {
	case 1:
		return "localhost:50051" //servidor para natacion luego camhbiar el localhost
	case 2:
		return "localhost:50052" //servidor para el atletismo
	case 3:
		return "localhost:50053" //servidor para el boxeo
	default:
		return "localhost:50051" //por si da error
	}
}

// metodo para enviar los datos al server
func sendData_To_Server(fiberCtx *fiber.Ctx) error {
	var body Data
	if err := fiberCtx.BodyParser(&body); err != nil {
		return fiberCtx.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Agregar línea para depurar: mostrar los datos recibidos
	log.Printf("Received data: %+v\n", body)

	//Obtener el servidor correcto  dependiendo la disciplina que se ingrese
	serverAddress := getServerForDiscipline(body.Discipline) //-> aun no se usa pero es para ver a donde va cada cosa

	//Establecer la conexion gRPC con el servidor
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials())) //para cuando ya tenga las rutas de los deploymets cambiar
	if err != nil {
		log.Fatalf("No se puede conectar: %v", err)
	}
	defer conn.Close()

	//Crear un cliente gRPC
	client := pb.NewFacultadServiceClient(conn)

	// Crear un canal para recibir las respuestas y errores
	responseChann := make(chan *pb.StudentResponse) // para las respuestas de servidor
	errorChann := make(chan error)                  // para manejar cualquier error

	//ejecutar gorutine para enviar la solicitud y recivir la respuesta
	go func() {

		//Enviar solicitud de Student al servidor gRPC y escribir la respuesta
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		ret, err := client.SendUserInfo(ctx, &pb.Student{
			Name:       body.Name,
			Age:        int32(body.Age),
			Faculty:    body.Faculty,
			Discipline: int32(body.Discipline),
		})

		if err != nil {
			errorChann <- err
			return
		}

		responseChann <- ret
	}()

	select {
	case response := <-responseChann:
		return fiberCtx.JSON(fiber.Map{
			"message": response.GetMessage(),
		})
	case err := <-errorChann:
		return fiberCtx.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	case <-time.After(5 * time.Second):
		return fiberCtx.Status(500).JSON(fiber.Map{
			"errror": "timeout",
		})
	}

}

func main() {
	//Crear la app de FIber
	app := fiber.New()

	//definir cual es el endpoint que va a recibir los datos del estudiante
	app.Post("/send_student_agro", sendData_To_Server)

	//Iniciar el servidor en el puerto 3000
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Error al iniciar el servidor  %v", err)
	}
}

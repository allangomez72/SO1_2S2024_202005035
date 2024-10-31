use actix_web::{web, App, HttpServer, HttpResponse, Responder};
use facultad::{facultad_service_client::FacultadServiceClient, Student};
use serde::{Deserialize, Serialize};
use tokio::sync::mpsc;

pub mod facultad {
    tonic::include_proto!("facultad");
}

#[derive(Deserialize, Serialize, Clone)]
struct StudentData {
    name: String,
    age: i32,
    faculty: String,
    discipline: i32,
}

// Mapeo de los servidores que se van a usar para las disciplinas
fn get_server_for_discipline(discipline: i32) -> &'static str {
    match discipline {
        1 => "service-50051.grpc-app.svc.cluster.local:50051", // natación
        2 => "service-50052.grpc-app.svc.cluster.local:50052", // atletismo
        3 => "service-50053.grpc-app.svc.cluster.local:50053", // boxeo
        _ => "service-50051.grpc-app.svc.cluster.local:50051",   // valor por defecto
    }
}

async fn handle_student(student: web::Json<StudentData>) -> impl Responder {
    let student_data = student.clone();
    println!(
        "Received data: name = {}, age = {}, faculty = {}, discipline = {}",
        student_data.name, student_data.age, student_data.faculty, student_data.discipline
    );

    // Crear un canal para recibir la respuesta
    let (tx, mut rx) = mpsc::channel(1);

    // Obtener el servidor que corresponde a la disciplina del estudiante
    let server_address = get_server_for_discipline(student_data.discipline);
    //let server_address = "http://localhost:50051";
    // Hacer la solicitud gRPC en un task asíncrono

    //Hilo para manejar las solicit gRPC
    tokio::spawn(async move {
        // Crear el cliente gRPC y conectar al servidor específico
        let mut client = match FacultadServiceClient::connect(server_address).await {
            Ok(client) => client,
            Err(e) => {
                eprintln!("Failed to connect to gRPC server {}: {}", server_address, e);
                return;
            }
        };

        // Crear el request basado en los datos del estudiante
        let request = tonic::Request::new(Student {
            name: student_data.name.clone(),
            age: student_data.age,
            faculty: student_data.faculty.clone(),
            discipline: student_data.discipline,
        });

        // Hacer la llamada gRPC
        match client.send_user_info(request).await {
            Ok(response) => {
                println!("Respuesta recibida de la disciplina {}: {:?}", student_data.discipline, response);
                // Enviar la respuesta a través del canal
                tx.send(format!("Discipline {}: {:?}", student_data.discipline, response)).await.unwrap();
            }
            Err(e) => {
                eprintln!("Error al enviar la disciplina {}: {}", student_data.discipline, e);
            }
        }
    });

    // Recoger la respuesta
    while let Some(res) = rx.recv().await {
        return HttpResponse::Ok().json(res);
    }

    // Enviar la respuesta como JSON
    HttpResponse::Ok().json("No se han recibido respuestas")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("Iniciando el servidor en localhost:8080");
    HttpServer::new(|| {
        App::new()
            .route("/send_student_inge", web::post().to(handle_student))
    })
        .bind("127.0.0.1:8080")?
        .run()
        .await
}

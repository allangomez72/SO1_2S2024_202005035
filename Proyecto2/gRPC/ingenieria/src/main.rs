use std::future::Future;
use actix_web:: {web, App, HttpServer, HttpResponse, Responder};
use facultad::{facultad_service_client::FacultadServiceClient, Student};
use serde:: {Deserialize, Serialize};
use tokio::sync::mpsc;
use std::thread;

pub mod facultad {
    tonic::include_proto!("facultad");
}

#[derive(Deserialize, Serialize, Clone)]
struct StudentData{
    name:String,
    age: i32,
    faculty: String,
    discipline: i32,
}

//Mapeo de los servidores que se van a usar para las disciplinass
fn get_server_for_discipline(discipline: i32) -> &'static str {
    match discipline {
        1 => "localhost:50051", // para natacion luego cambiar localhost
        2 => "localhost:50052", // para atletistmo
        3 => "localhost:50053", // para boxeo
        _ => {"localhost:50051"}
    }
}

async fn handle_student(student: web::Json<StudentData>) -> impl Responder{
    let (tx, mut rx) = mpsc::channel(3); // Se debe de crear un canal para las 3 disciplinas

    // crear un thread por cada discicplina
    for discipline in 1..=3 {
        let tx = tx.clone();
        let student_data = student.clone();

        // Obtener el servidor que va a ir con cada disciplina
        let server_address = get_server_for_discipline(discipline);

        thread::spawn(move ||{
            let runtime = tokio::runtime::Runtime::new().unwrap();
            runtime.block_on(async {
                // crear el cliente gRPC y conectar al servidor especifico
                let mut client = match FacultadServiceClient::connect(server_address).await {
                    Ok(client) => client,
                    Err(e) => {
                        eprint!("Failed to connect to gRPC serve {}: {}", server_address, e);
                        return;
                    }
                };

                //Crear el request basado en los datos del estudiante y la disciiplina
                let request = tonic::Request::new(Student {
                    name: student_data.name.clone(),
                    age: student_data.age,
                    faculty: student_data.faculty.clone(),
                    discipline

                });

                //Hacer la llamada gRPC
                match client.send_user_info(request).await {
                    Ok(response) => {
                        println!("Received responde from discipline {}: {:?}", discipline,response);
                        tx.send(format!("Discipline {}: {:?}", discipline, response)).await.unwrap();
                    }
                    Err(e) => {
                        eprint!("Failed to send to discipline {}: {}", discipline,e);
                    }
                }
            });
        });
    }
    // Recoger las respuestas de todos los threads
    let mut responses = Vec::new();
    while let Some(res) = rx.recv().await {
        responses.push(res);

    }
    HttpResponse::Ok().json(responses)
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("Iniciando el servidor en localhost:8080");
    HttpServer::new(|| {
        App::new()
            .route("/faculty", web::post().to(handle_student))
    })
        .bind("127.0.0.1:8080")?
        .run()
        .await
}
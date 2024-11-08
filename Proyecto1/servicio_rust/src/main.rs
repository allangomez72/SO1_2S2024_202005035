use std::fs::File;
use std::io::{self, Read};
use std::path::Path;
use serde::{Deserialize, Serialize};
use std::error::Error;
use reqwest::blocking::Client;
use chrono::{DateTime, Utc};
use reqwest::StatusCode;
use serde_json::{Value, json};

use std::process::{Command, Stdio};
use std::sync::{Arc, atomic::{AtomicBool, Ordering}};
use std::thread;
use std::time::Duration;
use ctrlc;

use std::os::unix::process::ExitStatusExt;

#[derive(Debug, Serialize, Deserialize)]
struct SystemInfo {
    #[serde(rename = "Total RAM")]
    total_ram: u64,
    #[serde(rename = "Free RAM")]
    free_ram: u64,
    #[serde(rename = "Used Ram")]
    used_ram: u64,
    #[serde(rename = "Processes")]
    processes: Vec<Process>
}

/*
    Además de esto, estamos implementando los traits Eq, Ord y PartialOrd para poder comparar
    los procesos en base a su uso de CPU y memoria.

    La estructura de datos representa un proceso en el sistema operativo, con los siguientes campos:
    - pid: El identificador del proceso.
    - name: El nombre del proceso.
    - cmd_line: La línea de comandos que se utilizó para ejecutar el proceso.
    - memory_usage: La cantidad de memoria que está utilizando el proceso.
    - cpu_usage: El porcentaje de uso de CPU que está utilizando el proceso.

    Serde nos deja implementar macros a cada campo de la estructura de datos para poder renombrar
    los campos en el JSON que se genere.
*/
//Siguiendo la nueva estructura json es esta:
#[derive(Debug, Serialize, Deserialize, PartialEq)]
struct Process {
    #[serde(rename = "PID")]
    pid: u32,
    #[serde(rename = "Name")]
    name: String,
    #[serde(rename = "Cmdline")]
    cmd_line: String,
    #[serde(rename = "VSZ")]
    vsz: u64,
    #[serde(rename = "RSS")]
    rss: u64,
    #[serde(rename = "MemoryUsage")]
    memory_usage: f64,
    #[serde(rename = "CPUUsage")]
    cpu_usage: f64,
}

//no lo cambio por que es indiferente solo es para capturar los procesos que se matan
#[derive(Debug, Serialize, Clone)]
struct LogProcess {
    pid: u32,
    container_id: String,
    name: String,
    vsz:u64,
    rss:u64,
    memory_usage: f64,
    cpu_usage: f64,
}

#[derive(Debug, Serialize, Clone)]
struct MemInfo {
    total_ram: u64,
    free_ram: u64,
    used_ram: u64,
    timestamp:String,
}

// IMPLEMENTACIÓN DE MÉTODOS

/*
    Función para sobreescribir el campo cmd_line de cada proceso por el id del contenedor.
*/
impl Process {
    fn get_container_id(&self) -> &str {
        let parts: Vec<&str> = self.cmd_line.split_whitespace().collect();
        for (i, part) in parts.iter().enumerate() {
            if *part == "-id" {
                if let Some(id) = parts.get(i + 1) {
                    return id;
                }
            }
        }
        "N/A"
    }
}

// IMPLEMENTACIÓN DE TRAITS

/*
    Este trait no lleva ninguna implementación, pero es necesario para poder comparar ya que debe satisfacer
    la propiedad de reflexividad, es decir, que un proceso es igual a sí mismo.
*/
impl Eq for Process {}


impl Ord for Process {
    fn cmp(&self, other: &Self) -> std::cmp::Ordering {
        self.cpu_usage.partial_cmp(&other.cpu_usage).unwrap_or(std::cmp::Ordering::Equal)
            .then_with(|| self.memory_usage.partial_cmp(&other.memory_usage).unwrap_or(std::cmp::Ordering::Equal))
    }
}

impl PartialOrd for Process {
    fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
        Some(self.cmp(other))
    }
}


// FUNCIONES

/*
    Función para matar un contenedor de Docker.
    - id: El identificador del contenedor que se quiere matar.
    - Regresa un std::process::Output que contiene la salida del comando que se ejecutó.
*/
fn kill_container(id: &str) -> std::process::Output {
    //Obtener el ID del contenedor FastAPI
    let fastapi_container_id = std::process::Command::new("docker")
    .arg("ps")
    .arg("--filter")
    .arg("name=log_container")
    .arg("--format")
    .arg("{{.ID}}")
    .output()
    .expect("fallo al obtener el ID del contenedor FastAPI");

    // Convertir el resultado a String
    let fastapi_id_str = String::from_utf8_lossy(&fastapi_container_id.stdout).trim().to_string();

    // Imprimir el ID del contenedor FastAPI
    println!("ID del contenedor FastAPI: {}", fastapi_id_str);

    // Continuar con la validación o cualquier otro proceso
    if id.starts_with(&fastapi_id_str) {
        println!("El contenedor FastAPI no será detenido.");
        return std::process::Output {
            status: std::process::ExitStatus::from_raw(0),
            stdout: Vec::new(),
            stderr: Vec::new(),
        };
    }

    let  output = std::process::Command::new("sudo")
        .arg("docker")
        .arg("stop")
        .arg(id)
        .output()
        .expect("failed to execute process");

    println!("Matando contenedor con id: {}", id);

    output
}

fn analyzer( system_info:  SystemInfo) {


    // Creamos un vector vacío para guardar los logs de los procesos.
    let mut log_proc_list: Vec<LogProcess> = Vec::new();


    /*
        Creamos un vector vacío para guardar los logs del sistema.
        En este caso, no se guardará nada, pero se puede modificar para guardar
        información del sistema.
    */
    let mut processes_list: Vec<Process> = system_info.processes;


    processes_list.sort();


    // Dividimos la lista de procesos en dos partes iguales.
    let (lowest_list, highest_list) = processes_list.split_at(processes_list.len() / 2);


    // Hacemos un print de los contenedores de bajo consumo en las listas.
    println!("\n*************** Bajo consumo ***************");
    for process in lowest_list {
        println!("PID: {}, Name: {}, Container ID: {}, Memory Usage: {}, CPU Usage: {}, VSZ: {}, RSS: {}",
                 process.pid,
                 process.name,
                 process.get_container_id(),
                 process.memory_usage,
                 process.cpu_usage,
                 process.vsz,
                 process.rss);
    }

    println!("------------------------------");

    println!("\n*************** Alto consumo ***************");
    for process in highest_list {
        println!("PID: {}, Name: {}, Container ID: {}, Memory Usage: {}, CPU Usage: {}, VSZ: {}, RSS: {}",
                 process.pid,
                 process.name,
                 process.get_container_id(),
                 process.memory_usage,
                 process.cpu_usage,
                 process.vsz,
                 process.rss);
    }

    println!("------------------------------");

    /*
        En la lista de bajo consumo, matamos todos los contenedores excepto los 3 primeros.
        antes
        | 1 | 2 | 3 | 4 | 5 |

        después
        | 1 | 2 | 3 |
    */

    if lowest_list.len() > 3 {
        // Iteramos sobre los procesos en la lista de bajo consumo.
        for process in lowest_list.iter().skip(3) {
            let log_process = LogProcess {
                pid: process.pid,
                container_id: process.get_container_id().to_string(),
                name: process.name.clone(),
                vsz: process.vsz,
                rss: process.rss,
                memory_usage: process.memory_usage,
                cpu_usage: process.cpu_usage,
            };

            log_proc_list.push(log_process.clone());

            // Matamos el contenedor.
            let _output = kill_container(&process.get_container_id());

        }
    }

    /*
        En la lista de alto consumo, matamos todos los contenedores excepto los 2 últimos.
        antes
        | 1 | 2 | 3 | 4 | 5 |

        después
                    | 4 | 5 |
    */
    if highest_list.len() > 2 {
        // Iteramos sobre los procesos en la lista de alto consumo.
        for process in highest_list.iter().take(highest_list.len() - 2) {
            let log_process = LogProcess {
                pid: process.pid,
                container_id: process.get_container_id().to_string(),
                name: process.name.clone(),
                vsz: process.vsz,
                rss: process.rss,
                memory_usage: process.memory_usage,
                cpu_usage: process.cpu_usage
            };

            log_proc_list.push(log_process.clone());

            // Matamos el contenedor.
            let _output = kill_container(&process.get_container_id());
        }
    }

    // TODO: ENVIAR LOGS AL CONTENEDOR REGISTRO
    send_logs(&log_proc_list);

    // Hacemos un print de los contenedores que matamos.
    println!("\n=============== Contenedores matados ===============");

    for process in log_proc_list {
        println!("PID: {}, Name: {}, Container ID: {}, Memory Usage: {}, CPU Usage: {} ",
                 process.pid,
                 process.name,
                 process.container_id,
                 process.memory_usage,
                 process.cpu_usage);
    }

    println!("------------------------------");

    //Enviar la informacion de la memoria y tiempo
    let now: DateTime<Utc> = Utc::now();
    let formatted_datetime = now.format("%d/%m/%Y %H:%M:%S").to_string();
    let mem_info = MemInfo{
        total_ram: system_info.total_ram,
        free_ram: system_info.free_ram,
        used_ram: system_info.used_ram,
        timestamp: formatted_datetime
    };
    if let Err(e) = send_meminfo(&mem_info) {
        eprintln!("Failed to send memory info: {}", e);
    }

}
//Eseta es mi fucnion para enviar los logs al servicio de python
fn send_logs(processes: &Vec<LogProcess>) -> Result<(), Box<dyn Error>> {
    let client = Client::new();
    //let timestamp = Utc::now().to_rfc3339();

    let log_message = json!(processes);

    let response = client
        .post("http://localhost:8000/logs")
        .json(&log_message)
        .send()?;

    if !response.status().is_success() {
        println!("La respuesta del servidor no fue exitosa: {}", response.status());
    }

    Ok(())
}

//Funcion para enviar la info de la memoria:
fn send_meminfo(mem_info: &MemInfo) -> Result<(), Box<dyn Error>>{
    let client = Client::new();
    let response = client
        .post("http://localhost:8000/meminfo")
        .json(mem_info)
        .send()?;

    if !response.status().is_success(){
        println!("La respuesta del servidor no fue exitosa: {}",response.status());
    }

    Ok(())
}

/*
    Función para leer el archivo proc
    - file_name: El nombre del archivo que se quiere leer.
    - Regresa un Result<String> que puede ser un error o el contenido del archivo.
*/
fn read_proc_file(file_name: &str) -> io::Result<String> {
    // Se crea un Path con el nombre del archivo que se quiere leer.
    let path  = Path::new("/proc").join(file_name);

    /*
        Se abre el archivo en modo lectura y se guarda en la variable file.
        En caso de que haya un error al abrir el archivo, se regresa un error.
        El signo de interrogación es un atajo para regresar un error en caso de que haya uno.
    */
    let mut file = File::open(path)?;

    // Se crea una variable mutable content que se inicializa con un String vacío.
    let mut content = String::new();

    // Se lee el contenido del archivo y se guarda en la variable content.
    file.read_to_string(&mut content)?;


    // Se regresa el contenido del archivo.
    Ok(content)
}

/*
    Función para deserializar el contenido del archivo proc a un vector de procesos.
    - json_str: El contenido del archivo proc en formato JSON.
    - Regresa un Result<> que puede ser un error o un SystemInfo.
*/
fn parse_proc_to_struct(json_str: &str) -> Result<SystemInfo, serde_json::Error> {
    // Se deserializa el contenido del archivo proc a un SystemInfo.
    let system_info: SystemInfo = serde_json::from_str(json_str)?;

    // Se regresa el SystemInfo.
    Ok(system_info)
}

fn get_img(){

}


fn main() {
    // Creamos una bandera para indicar cuándo se debe detener el proceso.
    let running = Arc::new(AtomicBool::new(true));
    let r = running.clone();

    // Capturamos la señal de Ctrl+C.
    ctrlc::set_handler(move || {
        println!("\n\n\nRecibido Ctrl+C! Deteniendo procesos...");
        // Cambiamos la bandera para detener el loop.
        r.store(false, Ordering::SeqCst);
        // Desinstalamos el cronjob
        uninstall_cronjob();

        println!("\n\nGenerando imagenes :D ");
        // Llamamos a las funciones bloqueantes
        get_img_process().unwrap();
        get_img_memory().unwrap();
        
        // Detenemos el contenedor del servicio de Python.
        let _ = Command::new("docker-compose")
            .arg("-f")
            .arg("/home/gomez/Documentos/SO1_2S2024_202005035/Proyecto1/server_python/docker-compose.yaml")  // Ruta relativa al archivo
            .arg("down")
            .stdout(Stdio::inherit())
            .stderr(Stdio::inherit())
            .output()
            .expect("fallo al detener docker-compose");
    }).expect("Error al configurar el manejador de Ctrl+C");

    //antes del docker-compose y del loop instalar el modulo
    install_sysinfo();
    install_cronjob();

    // Iniciamos el servicio de Docker Compose.
    println!("Iniciando Docker Compose...");
    let _docker_up = Command::new("docker-compose")
        .arg("-f")
        .arg("/home/gomez/Documentos/SO1_2S2024_202005035/Proyecto1/server_python/docker-compose.yaml")  // Ruta relativa al archivo
        .arg("up")
        .arg("-d")  // Modo detacheado
        .stdout(Stdio::inherit())
        .stderr(Stdio::inherit())
        .output()
        .expect("fallo al iniciar docker-compose");

    // Esperamos un tiempo fijo para que Docker Compose levante los contenedores.
    println!("Esperando 5 segundos para asegurar que Docker Compose haya levantado los contenedores...");
    thread::sleep(Duration::from_secs(5));

    // Loop principal.
    while running.load(Ordering::SeqCst) {
        // Creamos una estructura de datos SystemInfo con un vector de procesos vacío.
        let system_info: Result<SystemInfo, _>;

        // Leemos el contenido del archivo proc y lo guardamos en la variable json_str.
        let json_str = read_proc_file("sysinfo_202005035").unwrap();

        // Deserializamos el contenido del archivo proc a un SystemInfo.
        system_info = parse_proc_to_struct(&json_str);

        // Dependiendo de si se pudo deserializar el contenido del archivo proc o no, se ejecuta una u otra rama.
        match system_info {
            Ok(info) => {
                println!("\n============= GENERAL INFORMATION =============");
                println!("Total RAM: {} KB", info.total_ram);
                println!("Free RAM: {} KB", info.free_ram);
                println!("Used RAM: {} KB", info.used_ram);

                analyzer(info);
            }
            Err(e) => println!("Failed to parse JSON: {}", e),
        }

        // Dormimos el hilo principal por 10 segundos.
        thread::sleep(Duration::from_secs(10));
    }

    println!("\n\n================ Proceso finalizado ================");
    uninstall_sysinfo(); //Va desinstalar el modulo cuando termine la ejecucion
}

// Función para instalar el módulo sysinfo.ko
fn install_sysinfo() {
    println!("Instalando el módulo sysinfo.ko...");
    let _ = Command::new("sudo")
        .arg("insmod")
        .arg("/home/gomez/Documentos/SO1_2S2024_202005035/Proyecto1/ModuleKernel/sysinfo.ko")
        .stdout(Stdio::inherit())
        .stderr(Stdio::inherit())
        .output()
        .expect("Fallo al instalar sysinfo.ko");
}

// Función para desinstalar el módulo sysinfo.ko
fn uninstall_sysinfo() {
    println!("Desinstalando el módulo sysinfo.ko...");
    let _ = Command::new("sudo")
        .arg("rmmod")
        .arg("sysinfo")
        .stdout(Stdio::inherit())
        .stderr(Stdio::inherit())
        .output()
        .expect("Fallo al desinstalar sysinfo.ko");
}

// Función para obtener información de memoria 
fn get_img_process() -> Result<(), Box<dyn Error>> {
    let client = Client::new();
    let response = client
        .get("http://localhost:8000/generate_memory_graph")
        .send()?;

    let status = response.status();
    let body = response.text()?;

    if !status.is_success() {
        println!("La respuesta del servidor no fue exitosa: {}", status);
    } else {
        let json: Value = serde_json::from_str(&body)?;
        if json.get("error").is_some() {
            println!("Error: {}", json["error"].as_str().unwrap_or("Desconocido"));
        } else if let Some(plot_path) = json.get("plot_path") {
            println!("Plot path: {}", plot_path.as_str().unwrap_or("No disponible"));
        }
    }

    Ok(())
}

fn get_img_memory() -> Result<(), Box<dyn Error>> {
    let client = Client::new();
    let response = client
        .get("http://localhost:8000/generate_meminfo_graph")
        .send()?;

    let status = response.status();
    let body = response.text()?;

    if !status.is_success() {
        println!("La respuesta del servidor no fue exitosa: {}", status);
    } else {
        let json: Value = serde_json::from_str(&body)?;
        if json.get("error").is_some() {
            println!("Error: {}", json["error"].as_str().unwrap_or("Desconocido"));
        } else if let Some(plot_path) = json.get("plot_path") {
            println!("Plot path: {}", plot_path.as_str().unwrap_or("No disponible"));
        }
    }

    Ok(())
}

// Función para instalar el cronjob
fn install_cronjob() {
    let cron_job = "* * * * * /bin/bash /home/gomez/Documentos/SO1_2S2024_202005035/Proyecto1/Containers/generate_containers.sh\n";
    
    // Agregamos el cronjob directamente al crontab
    let status = Command::new("sh")
        .arg("-c")
        .arg(format!("(crontab -l 2>/dev/null; echo \"{}\") | crontab -", cron_job))
        .status()
        .expect("Failed to install cronjob");

    if status.success() {
        println!("Cronjob instalado.");
    } else {
        println!("Error al instalar el cronjob.");
    }
}

// Función para desinstalar el cronjob
fn uninstall_cronjob() {
    // Elimina el cronjob que contiene el nombre del script
    let status = Command::new("sh")
        .arg("-c")
        .arg("crontab -l | grep -v 'generate_containers.sh' | crontab -")
        .status()
        .expect("Failed to uninstall cronjob");

    if status.success() {
        println!("Cronjob desinstalado.");
    } else {
        println!("Error al desinstalar el cronjob.");
    }
}

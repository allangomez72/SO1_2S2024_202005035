
# Manual Técnico: Instalación y Uso del Proyecto
## Introducción

Este manual describe los pasos necesarios para la instalación, configuración y ejecución del proyecto que incluye contenedores Docker, un módulo de kernel en Linux, y servicios escritos en Python y Rust. El proyecto se compone de varios elementos, como la creación y administración de contenedores, la implementación de un módulo de kernel, y la ejecución de servicios API en contenedores.

---

## Requisitos Previos

- **Docker** instalado y funcionando en tu sistema. Verifica la instalación usando `docker --version`.
- **Python 3** con las dependencias necesarias.
- **Rust** instalado.
- **Permisos de administrador (sudo)** para ejecutar algunos comandos.

---

## 1. Instalación de Imágenes Docker

### 1.1 Creación de Imágenes Docker

El proyecto contiene varias imágenes Docker que deben ser generadas para crear contenedores. Las imágenes se generan a partir de los archivos `Dockerfile` ubicados en diferentes directorios.

#### Instrucciones para generar las imágenes:

1. Dirígete a la ruta donde se encuentra el archivo `Dockerfile` correspondiente.

   Para fines de ejemplo se usa esta ruta imágenes de alto consumo (HIGH):
   ```bash
   cd Proyecto1/Containers/high/highCPU
   ```

   Para imágenes de bajo consumo (LOW):
   ```bash
   cd Proyecto1/Containers/low/lowCPU
   ```

2. Ejecuta el siguiente comando para construir la imagen:
   ```bash
   sudo docker build -t <nombre_imagen> .
   ```

   Ejemplos:
    - Para crear la imagen `high_cpu_image`:
      ```bash
      sudo docker build -t high_cpu_image .
      ```
    - Para la imagen `low_cpu_image`:
      ```bash
      sudo docker build -t low_cpu_image .
      ```

Ojo tu archivo dockerfile, debera de verse similar a este ejemmplo:
```dockerfile
FROM python:3.9-slim

RUN pip install flask

COPY log_cpu.py /log_cpu.py

CMD ["python", "log_cpu.py"]
```
Pasa muchas veces que docker se encuentra desativado, por lo que se tiene que activar, por lo cual de bebe de hacer lo siguiente:
### 1.2 Verificación de Docker Activo

Si Docker no está activo, usa el siguiente comando para iniciarlo:
```bash
sudo systemctl start docker
```

---

## 2. Generación de Contenedores

Una vez generadas las imágenes, puedes usar un script de bash para crear los contenedores de manera aleatoria. A continuación, se explica cómo ejecutar dicho script.

Ejemplo del script de bash para Generar contenedores, tambien se puede consultar [generate_container](https://github.com/allangomez72/SO1_2S2024_202005035/blob/main/Proyecto1/Containers/generate_containers.sh)
```bash
#!/bin/bash

# Definir las imágenes para contenedores de bajo y alto consumo
LOW_CONSUMPTION_IMAGE="low_cpu_image"
HIGH_CONSUMPTION_IMAGE="high_cpu_image"

# Función para generar un nombre aleatorio
generate_name() {
  cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 8 | head -n 1
}

# Crear 3 contenedores de bajo consumo
for i in {1..3}; do
  CONTAINER_NAME=$(generate_name)
  echo "Creando contenedor de bajo consumo: $CONTAINER_NAME"
  sudo docker run -d --name "$CONTAINER_NAME" "$LOW_CONSUMPTION_IMAGE"
done

# Crear 2 contenedores de alto consumo
for i in {1..2}; do
  CONTAINER_NAME=$(generate_name)
  echo "Creando contenedor de alto consumo: $CONTAINER_NAME"
  sudo docker run -d --name "$CONTAINER_NAME" "$HIGH_CONSUMPTION_IMAGE"
done

# Verificar los contenedores creados
echo "Contenedores creados:"
sudo docker ps

```
### 2.1 Ejecución Manual del Script

Ejecuta el script de bash para generar los contenedores:

```bash
sudo ./generate_containers.sh
```

> [!NOTE]
> Se recomienda **no agregar el script a cronjob inmediatamente** hasta que el proyecto esté más estructurado, ya que podría sobrecargar el equipo.

---

## 3. Módulo de Kernel

El proyecto incluye un módulo de kernel que proporciona información sobre el sistema, como la memoria RAM. A continuación, se explican los comandos necesarios para compilar e instalar el módulo.

### 3.1 Compilación e Instalación del Módulo

Se debe crear el archivo sysinfo.c, a continuacion se muestra un ejemplo muy simple del módulo de kernle en C
```c++
#include <linux/module.h>   
#include <linux/kernel.h>   
#include <linux/init.h>  

// Función que se ejecuta cuando se carga el módulo
static int __init sysinfo_basic_init(void) {
    printk(KERN_INFO "Módulo de kernel cargado: Hola, sistema!\n");
    return 0;  // 0 indica una carga exitosa
}

// Función que se ejecuta cuando se remueve el módulo
static void __exit sysinfo_basic_exit(void) {
    printk(KERN_INFO "Módulo de kernel removido: Adiós, sistema!\n");
}

// Definir las funciones de entrada y salida del módulo
module_init(sysinfo_basic_init);
module_exit(sysinfo_basic_exit);

MODULE_LICENSE("GPL");              // Licencia del módulo
MODULE_AUTHOR("Tu Nombre");         // Autor del módulo
MODULE_DESCRIPTION("Módulo básico de ejemplo para obtener información del sistema");  // Descripción

```
Tambien se puede verificar directamente del proyecto el modulo construido:
[sysinfo.ko](https://github.com/allangomez72/SO1_2S2024_202005035/blob/main/Proyecto1/ModuleKernel/sysinfo.c)


1. Compilar el módulo con:
   ```bash
   make
   ```

2. Instalar el módulo en el sistema:
   ```bash
   sudo insmod sysinfo.ko
   ```

3. Para desinstalar el módulo:
   ```bash
   sudo rmmod sysinfo
   ```

4. Verifica la información del sistema con:
   ```bash
   cat /proc/sysinfo_202005035
   ```

---

## 4. Ejecución del Proyecto en Rust

El componente en Rust se encarga de la lógica principal del proyecto. A continuación se describe cómo configurarlo y ejecutarlo.

### 4.1 Dependencias

El archivo `Cargo.toml` debe incluir las siguientes dependencias:
```toml
[dependencies]
reqwest = { version = "0.11", features = ["blocking", "json"] }
chrono = "0.4"
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"
ctrlc = "3.2"
```

### 4.2 Compilación y Ejecución

1. Para compilar el proyecto, usa:
   ```bash
   cargo build
   ```

2. Para ejecutar el proyecto:
   ```bash
   cargo run
   ```
   
### 4.3 Ejemplo basico

Como ejemplo, puedes crear un pequeño programa llamado hello_rust.rs dentro de tu proyecto.

```rust
fn main() {
    println!("Hello, Rust!");
}
```

---

## 5. Servicio en Python con FastAPI

Este proyecto también incluye un servicio en Python que utiliza **FastAPI**. A continuación se detallan los pasos para configurar y ejecutar el entorno virtual y las dependencias.

### 5.1 Instalación de Dependencias

1. Activa el entorno virtual:
   ```bash
   source env/bin/activate
   ```

2. Instala las dependencias desde el archivo `requirements.txt`:
   ```bash
   pip install -r requirements.txt
   ```

   Las principales dependencias son **FastAPI** y **Matplotlib**.

### 5.2 Dockerización del Servicio Python

1. Construye la imagen Docker para el servicio:
   ```bash
   sudo docker compose build
   ```

2. Levanta el servicio:
   ```bash
   sudo docker compose up -d
   ```

3. Para detener el servicio:
   ```bash
   sudo docker compose down
   ```

Ejemplo corto del servicio en Python con FastAPI
```python
from fastapi import FastAPI
import json
import os
from typing import List
from pydantic import BaseModel

app = FastAPI()

# Definimos el modelo para los logs
class LogProcess(BaseModel):
    pid: int
    name: str
    memory_usage: float

# Ruta para recibir y guardar los logs en un archivo JSON
@app.post("/logs")
def add_logs(logs: List[LogProcess]):
    logs_file = 'logs.json'

    # Si el archivo existe, cargamos los logs previos
    if os.path.exists(logs_file):
        with open(logs_file, 'r') as f:
            existing_logs = json.load(f)
    else:
        existing_logs = []

    # Añadimos los nuevos logs
    new_logs = [log.dict() for log in logs]
    existing_logs.extend(new_logs)

    # Guardamos los logs actualizados en el archivo
    with open(logs_file, 'w') as f:
        json.dump(existing_logs, f, indent=4)

    return {"status": "Logs recibidos"}

```


---

## 6. Automatización con Rust

El sistema puede automatizar la ejecución del módulo de kernel, levantar el servicio Docker, y generar contenedores automáticamente desde el código Rust. Se recomienda iniciar estos procesos manualmente al principio y luego agregar la automatización al cronjob cuando el proyecto esté completamente estructurado.

---

## Ejemplos de Uso

- Generar contenedores:
  ```bash
  sudo ./generate_containers.sh
  ```

- Verificar la información del sistema:
  ```bash
  cat /proc/sysinfo_202005035
  ```

- Levantar el servicio de FastAPI:
  ```bash
  sudo docker compose up -d
  ```

---

## Conclusión

- Este manual proporciona los pasos esenciales para instalar y ejecutar los componentes clave del proyecto, que abarca contenedores Docker, un módulo de kernel, y servicios en Python y Rust. En un inicio se vuelve tendioso el tener tanto comandos, pero se recomienda tambien tener un archivo extern, donde tener los comandos esenciales apuntados para no perder tanto tiempo buscando


- A lo largo del proyecto se ha evidenciado la importancia de contar con una estructura detallada, que permita el manejo eficiente de servicios, contenedores y procesos. La documentacion que se presenta en el manual detalla ademas de algunos pasos tecnicos, presenta un enfoque de organizacion para abordar la coplejidad del sistema. 


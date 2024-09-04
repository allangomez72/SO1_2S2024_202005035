#!/bin/bash

# Lista de imágenes construidas
IMAGES=("high_cpu_image" "high_ram_image" "low_cpu_image", "low_ram_image")

# Archivo de salida para registrar los contenedores creados
output_file="container_log.txt"
> $output_file # Limpiar el archivo de salida si existe

# Número de contenedores a crear
CONTAINERS_TOTAL=10

for i in $(seq 1 $CONTAINERS_TOTAL); do
    # Selección aleatoria de la imagen
    IMAGE=${IMAGES[$RANDOM % ${#IMAGES[@]}]}  #IMAGE va ser para elegir el nombre done Random es el aleatorio y IMAGES[@] cantidad del arreglo

    # Generación de un nombre único para el contenedor
    CONTAINER_NAME=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 8)

    # Creación y ejecución del contenedor
    CONTAINER_ID=$(docker run -d --name "$CONTAINER_NAME" "$IMAGE")

    # Registro de la creación del contenedor
    echo "Creando Contenedor  $i: $CONTAINER_NAME | Imagen: $IMAGE | ID: $CONTAINER_ID" | tee -a $output_file

    # Obtención del PID del contenedor y registro (opcional)
    CONTAINER_PID=$(docker inspect --format '{{.State.Pid}}' "$CONTAINER_ID")
    echo "PID: $CONTAINER_PID" >> $output_file
    echo "-----------------------------" >> $output_file
done

echo "Registro completo de los contenedores creados en $output_file."

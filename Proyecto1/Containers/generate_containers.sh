#!/bin/bash

# Lista de imágenes construidas
IMAGES=("high_cpu_image" "high_ram_image" "low_cpu_image" "low_ram_image")

# Archivos de salida para registrar los contenedores y PIDs creados
output_file="container_log.txt"
pid_file="/home/gomez/Documentos/SO1_2S2024_202005035/Proyecto1/ModuleKernel/pid_log.txt"

# Limpiar el archivo de salida de contenedores si existe
> $output_file

# Número de contenedores a crear
CONTAINERS_TOTAL=10

for i in $(seq 1 $CONTAINERS_TOTAL); do
    # Selección aleatoria de la imagen
    IMAGE=${IMAGES[$RANDOM % ${#IMAGES[@]}]}

    # Generación de un nombre único para el contenedor
    CONTAINER_NAME=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 8)

    # Creación y ejecución del contenedor
    CONTAINER_ID=$(docker run -d --name "$CONTAINER_NAME" "$IMAGE")
    if [ $? -ne 0 ]; then
        echo "Error al crear el contenedor con la imagen $IMAGE" | tee -a $output_file
        continue
    fi

    # Registro de la creación del contenedor
    echo "Creando Contenedor $i: $CONTAINER_NAME | Imagen: $IMAGE | ID: $CONTAINER_ID" | tee -a $output_file

    # Obtención del PID del contenedor y registro en ambos archivos
    CONTAINER_PID=$(docker inspect --format '{{.State.Pid}}' "$CONTAINER_ID")
    echo "PID: $CONTAINER_PID" >> $output_file
    echo $CONTAINER_PID >> $pid_file  # Añadir el PID al archivo pid_log.txt sin sobrescribir
    echo "-----------------------------" >> $output_file
done

echo "Registro completo de los contenedores creados en $output_file."
echo "PIDs de los contenedores creados en $pid_file."

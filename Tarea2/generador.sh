#!/bin/bash

# Función para generar un nombre aleatorio
generate_random_name() {
    echo "container_$(tr -dc A-Za-z0-9 </dev/urandom | head -c 8)"
}

# Crear 10 contenedores y como se pide usar alpine, se agrega la imagen
for i in {1..10}; do
    name=$(generate_random_name)
    docker run -d --name $name alpine sleep 3600 #sleep para que se ejecuten los contenedores y se puedan listar
done

# Mostrar los contenedores creados
docker ps

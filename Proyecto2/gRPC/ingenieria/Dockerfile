# Definir la imagen base
FROM rust:1.72 as builder

# Instalar las dependencias necesarias
RUN apt-get update && apt-get install -y protobuf-compiler && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app

# Copiar todo el contenido 
COPY . .

# Constuir la aplicacion
RUN cargo build --release

# Definir la imagen base
FROM debian:buster-slim

# Instalar las librerias necesarias
RUN apt-get update && apt-get install -y libssl1.1 ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Se copia el binario de la aplicacion final
COPY --from=builder /usr/src/app/target/release/ingenieria .
# es el nombre que se definio en el toml
CMD ["./ingenieria"]
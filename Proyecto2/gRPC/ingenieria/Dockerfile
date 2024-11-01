# Etapa 1: Construcción
FROM rust:1.72-slim as builder

# Instala protobuf y las dependencias necesarias
RUN apt-get update && apt-get install -y protobuf-compiler && rm -rf /var/lib/apt/lists/*

# Configura el directorio de trabajo
WORKDIR /app

# Copia los archivos Cargo para gestionar dependencias
COPY Cargo.toml Cargo.lock ./

# Construye las dependencias de Rust para que se almacenen en caché
RUN mkdir src && echo 'fn main() {}' > src/main.rs
RUN cargo build --release && rm -rf src

# Copia el código fuente
COPY . .

# Ejecuta el build.rs para compilar los protos
RUN cargo build --release

# Etapa 2: Imagen ligera para ejecución
FROM debian:bookworm-slim

# Instala las dependencias mínimas necesarias para ejecutar la aplicación
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copia el binario desde la etapa de construcción
COPY --from=builder /app/target/release/ingenieria /usr/local/bin/ingenieria

# Comando de inicio
CMD ["ingenieria"]

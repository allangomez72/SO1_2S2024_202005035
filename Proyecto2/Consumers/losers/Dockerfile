# Se define una imagen de base con la versión de Go 1.23 para compilar el código fuente.
FROM golang:1.23-bullseye AS build

# Se define el directorio de trabajo
WORKDIR /app

# Se copian los archivos go.mod y go.sum al directorio de trabajo
COPY go.mod go.sum ./

# Se descargan las dependencias del proyecto que están en el archivo go.mod
RUN go mod download

# Se copia todo el código fuente al directorio de trabajo
COPY . .

# Compilar la aplicación con las banderas para enlazar estáticamente el binario
RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o kafka-redis-consumer

# Nueva imagen a partir de scratch
FROM scratch
# Se copian los certificados de la imagen base que se definió en build
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Se copia el binario de la aplicación final
COPY --from=build /app/kafka-redis-consumer kafka-redis-consumer

# Comando para ejecutar el binario
CMD ["/kafka-redis-consumer"]

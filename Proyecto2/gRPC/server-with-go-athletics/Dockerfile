# Se define una imagen de base con la versión de Go 1.23 para compilar el código fuente.
FROM golang:1.23-bullseye AS build

# Se define el directorio de trabajo
WORKDIR /app

# Se copian los archivos go.mod y go.sum al directorio de trabajo
COPY go.mod go.sum ./

# Se descargan las dependencias del proyecto que estan en el archivo go.mod
RUN go mod download

# Se copia todo el codigo fuente al directorio de trabajo
COPY . .

# Compliar la aplicacion con las banderas, ldflags se utiliza para enlazar estaticamente el de tag debe usar su red y el main para defini el nombre del binario
RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o main

# Nueva imagen a partir de la imagen de alpine
FROM alpine:latest

# Es necesario para conexiones HTTP seguras
RUN apk --no-cache add ca-certificates

# Se copia el binario de la aplicacion final
COPY --from=build /app/main /main

CMD ["/main"]
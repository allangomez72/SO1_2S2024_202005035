Para el cliente servidor que se usa en Golang se debe instalar lo siguiente para el funcionamiento de gRPC:

### Prerequisitos
```bash
#Esto es aplicado a linux (administrador de paquetes apt)
apt install -y protobuf-compiler
# Verificar la correcta instalacion
protoc --version
```
Para mas informacion verificar [documentación](https://grpc.io/docs/protoc-installation/)

Una vez instalado protoc, necesitarás instalar el plugin de gRPC para Go:

```bash
go get google.golang.org/grpc

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# Actualizar el PATH para que encuentre los plugins
export PATH="$PATH:$(go env GOPATH)/bin"
```
## Instaslar fiber para crear las APIS

```bash
go get github.com/gofiber/fiber/v2 # fiber
```
## El comando para ejecutar el proto
```bash
# compilar el agro.proto
protoc --go_out=./proto-go/ --go_opt=paths=source_relative --go-grpc_out=./proto-go/ --go-grpc_opt=paths=source_relative agro.proto
# otra opcion
protoc --go_out=./proto-go --go-grpc_out=./proto-go agro.proto
```

Comando usados extra
- Crear go mod
- Crear archivo vacio
- Crear estrutura Rust
- Verificar se si instalaron los plugins
```bash
go mod init <carpeta_general_archivo>
touch <archivo.extension>
cargo new ingenieria

#plugins
protoc-gen-go --version
protoc-gen-go-grpc --version
 ```
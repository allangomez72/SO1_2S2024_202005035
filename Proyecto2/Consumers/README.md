Iniciamos con la creacion de un archivo go, en el cual se creara un consumidor de kafka y un consumidor de redis, para ello se debe de crear un archivo go y se debe de inicializar el modulo de go con el siguiente comando

```bash
go mod init <carpeta_general_archivo>
```

Para ambos debemos de instalar las biblitcas necesarias para redis y para la conexion con kafka
    
```bash
go get github.com/IBM/sarama
go get github.com/go-redis/redis/v8
```

Para subir al docker hub se debe hacer lo siguiente:
```bash
docker build -t <nombre_usuario>/<nombreimagen>:<version> .
# Para subir la imagen al docker hub
docker push <nombre_usuario>/<nombreimagen>:<version>
```
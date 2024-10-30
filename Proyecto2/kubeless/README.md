Verificar que se esta apuntando al cluster
```bash
gcloud container clusters get-credentials sopes-ht-1 --region us-central1
```
Subir los archivos yaml
Una vez creados todos los archivos yaml, se deben subir al cluster de kubernetes con el siguiente comando:
```bash
kubectl apply -f namespace.yaml
kubectl apply -f grpc-client.yaml
kubectl apply -f grpc-server.yaml
kubectl apply -f service-grpc.yaml
kubectl apply -f rote-ingress.yaml
```
Verificar los namespaces
```bash
kubectl get namespaces
```

Verificar los deployments ojo dentro del namespace grpc-app
```bash
kubectl get deployments -n grpc-app
```

Para ver los pods
```bash
kubectl get pods -n grpc-app
```
Para ver los servicios
```bash
kubectl get services -n grpc-app
```
Para ver los ingress
```bash
kubectl get ingress -n grpc-app
```
Como se ocupo la IP publica del load balancer, para la comunicacion con locust:
```bash
kubectl get services -n grpc-app
```
Buscar la columna `External-IP` y copiar la IP publica

## Locust
Para correr locust en local y apuntar al cluster de kubernetes, se debe correr el siguiente comando:
```bash
locust -f <locustfile.py> --host=http://YOUR_LOAD_BALANCER_IP:3000
```

Verificar los pods
```bash
kubectl get pods -n grpc-app
```

Para poder ver los logs que se definieron en el servidor:
```bash
kubectl logs <nombre-del-pod-grpc-server> -n grpc-app
```
kubectl logs grpc-server-7d95dbd98b-kqhtt -n grpc-app

Para verificar que el cliente este corriendo:
```bash
kubectl logs <nombre-del-pod-grpc-client> -n grpc-app
```
Y ya ezz

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


---
## PROYECTO 2 FLUJO DE TRABAJO

Primero debemos subir el namespace, y el deployment de los clientes, para ello se debe correr el siguiente comando:
```bash
kubectl apply -f namespace.yaml
kubectl apply -f client-agronomia-grpc.yaml
kubectl apply -f client-ingenieria-grpc.yaml
```
Creamos el namespace para lo que va ser kafka
```bash
kubectl create namespace kafka
#instalamos el operador de strimzi
kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka
```
Creamos el cluster de kafka
```bash
kubectl apply -f kafka-cluster.yaml
# los topicos de kafka
kubectl apply -f topics-kafka.yaml
```
Y ahora si ya subimos los sevidores para que consuman los clientes
```bash
kubectl apply -f server-athletics-grpc.yaml
kubectl apply -f server-boxing-grpc.yaml
kubectl apply -f server-swimming-grpc.yaml
```
Se instala redis con helm
```bash
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis-release bitnami/redis
```

Subimos los consumers de kafka
```bash
kubectl apply -f consumer-losers.yaml
kubectl apply -f consumer-winners.yaml
```
Subimos el HPA
```bash
kubectl apply -f HPA-client.yaml
kubectl apply -f HPA-server.yaml
```

Posterior a eso subir el ingress
```bash
kubectl apply -f ingress.yaml
```
Para la cli de redis y verificar los datos:
```bash
kubectl exec -it redis-release-master-0 -- redis-cli
```
Para ver lo de la CLI de redis
```bash
> keys winners:*
> keys losers:*
```
Para ver las ips de los load balancer
```bash
kubectl get services -n grpc-app
```
para ver los pods
```bash
kubectl get pods -n grpc-app
```
Para ver los logs de los pods
```bash
kubectl logs <nombre-del-pod> -n grpc-app
```

```bash
#para la contraseña de redis
kubectl get secret --namespace default redis-release -o jsonpath="{.data.redis-password}" | base64 -d
```

```bash
#para la contraseña de grafana
echo "Password: $(kubectl get secret grafana-admin --namespace default -o jsonpath="{.data.GF_SECURITY_ADMIN_PASSWORD}" | base64 -d)"
```

Instalar grafana con la configuracion del loadbalancer
```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install grafana bitnami/grafana --set plugins="redis-datasource" --set security.allowEmbedding=true --set service.type=LoadBalancer --set http.cors.allowOrigins=http://localhost:3000 --namespace default
```

Recuperar contraseña de grafana
```bash
echo "Password: $(kubectl get secret grafana-admin --namespace default -o jsonpath="{.data.GF_SECURITY_ADMIN_PASSWORD}" | base64 -d)"
```

Instalar prometheus con los exportadores de las metricas:
```bash
helm install prometheus bitnami/prometheus --namespace default
helm install kube-state-metrics bitnami/kube-state-metrics
helm install node-exporter bitnami/node-exporter
```

kubectl delete -f ingress.yaml -n grpc-app
kubectl get ingress -n grpc-app

kubectl get pods -n kafka

kubectl logs sever-athletics-grpc-6d8d89f977-8zlb2 -n grpc-app
kubectl logs sever-boxing-grpc-565dbfbbb-r9q9l -n grpc-app
kubectl logs sever-swimming-grpc-764fbd48c4-j845f -n grpc-app

kubectl logs client-agronomia-grpc-747dd95988-85ckx -n grpc-app
kubectl logs client-ingenieria-grpc-5968996f46-765kq -n grpc-app

kubectl logs <nombre_del_pod> (-n <nombre_del_namespace)

kubectl logs consumer-losers-56cfc6b676-2rpk4
kubectl logs consumer-winners-567f585c9b-km5p9

consumer-losers-56cfc6b676-2rpk4
consumer-winners-567f585c9b-km5p9

eliminar todo lo de client
kubectl delete -f client-ingenieria-grpc.yaml

eliminar todo el loser
kubectl delete -f consumer-losers.yaml

eliminar todo el winner
kubectl delete -f consumer-winners.yaml

ver errores con logs
kubectl logs client-ingenieria-grpc-75c7b48674-29dj5 -n grpc-app


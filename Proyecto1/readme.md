## Para los contenedores de docker
Entrar a la ruta de cada uno donde esta el Dockerfile

Para entrar a la ruta donde se encuentra los contenedores de HIGH
```bash
cd Containers/high/highCPU
```
Pra los de LOW

```bash
cd Containers/low/lowCPU
```

Para el de high_CPU
```bash
sudo docker build -t high_cpu_image .
```
Para el de high_RAM
```bash
sudo docker build -t high_ram_image .
```
Para el de low_CPU
```bash
sudo docker build -t low_cpu_image .
```
Para el de low_RAM
```bash
sudo docker build -t low_ram_image .
```
Para encender docker que a veces esta apagado
```bash
sudo systemctl start docker
```


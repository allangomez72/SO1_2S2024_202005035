#Script para consumir CPU
import time
import math

while True:
    for _ in range(150):  # Menos iteraciones
        _ = [math.sqrt(i) for i in range(25)]  # Rango más pequeño
    #Reducir tiempo de descanso
    time.sleep(0.01)



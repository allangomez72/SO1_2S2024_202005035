#Script para el consumo de RAM
import time
import numpy as np

#Lista grande para consumir RAM
large_data = []

while True:
    #Generar una matriz grande y añadir a la lista
    large_data.append(np.random.rand(100,100))
    #Esperar
    time.sleep(1)
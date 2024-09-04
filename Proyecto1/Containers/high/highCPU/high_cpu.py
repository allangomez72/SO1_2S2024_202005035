#Script para consumir CPU
import time
import math

while True:
    for _ in range(10000):
        _ = [math.sqrt(i) for i in range(100)]
    print("Iteration completed")
    #reducir tiempo de descanso
    time.sleep(0.01)


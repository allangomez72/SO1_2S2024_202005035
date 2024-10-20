from locust import HttpUser, task, between
from faker import  Faker #libreria para datos aleatorios
import random

fake = Faker()  # Esto es para inicializar la generacion de datos aleatorios
faculties = ["Ingenieria", "Agronomia"]
disiciplines = [1,2,3]

class WebsiteUser(HttpUser):
    wait_time = between(1,5)

    @task
    def send_data_students(self):
        #Generamos los datos random
        student_data = {
            "name": fake.name(), #esto es para el nombre y apellido
            "age": random.randint(18,30), #edad entre 18 y 30
            "faculty": random.choice(faculties), #eleccion aleatoria de la facultad
            "discipline": random.choice(disiciplines) #eleccion aleatoria de las disicplinas
        }
        #para enviar los daatos al endpoint
        self.client.post("/sendstudent", json=student_data )
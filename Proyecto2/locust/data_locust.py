from locust import HttpUser, task, between
from faker import  Faker #libreria para datos aleatorios
import  random


fake = Faker()  # Esto es para inicializar la generacion de datos aleatorios
faculties = ["Ingenieria", "Agronomia"]
disiciplines = [1,2,3]

class WebsiteUser(HttpUser):
    wait_time = between(1,5)
    @task
    def send_data_students_ingenieria(self):
        #Generamos los datos random
        student_data = {
            "name": f"{fake.name()} {fake.last_name()}", #esto es para el nombre y apellido
            "age": random.randint(18,30), #edad entre 18 y 30
            "faculty": faculties[0], #eleccion aleatoria de la facultad
            "discipline": random.choice(disiciplines) #eleccion aleatoria de las disicplinas
        }
        #para enviar los daatos al endpoint
        self.client.post("/send_student_inge", json = student_data )

    @task
    def send_data_students_agronomia(self):
        #Generamos los datos random
        student_data = {
            "name": f"{fake.name()} {fake.last_name()}", #esto es para el nombre y apellido
            "age": random.randint(18,30), #edad entre 18 y 30
            "faculty": faculties[1], #eleccion aleatoria de la facultad
            "discipline": random.choice(disiciplines) #eleccion aleatoria de las disicplinas
        }
        #para enviar los daatos al endpoint
        self.client.post("http://34.83.224.70:3000/send_student_agro", json = student_data )
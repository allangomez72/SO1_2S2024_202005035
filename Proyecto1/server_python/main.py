from fastapi import FastAPI  # type: ignore
import os
import json
import matplotlib.pyplot as plt
from typing import List,Dict
from models.models import LogProcess
from models.models import LogSystem

app = FastAPI()


@app.get("/")
def read_root():
    return {"Hello": "World"}


@app.post("/logs")
def get_logs(logs_proc: List[LogProcess]):
    logs_file = 'logs/logs.json'

    # Checamos si existe el archivo logs.json
    if os.path.exists(logs_file):
        # Leemos el archivo logs.json
        with open(logs_file, 'r') as file:
            existing_logs = json.load(file)
    else:
        # Sino existe, creamos una lista vacía
        existing_logs = []

    # Agregamos los nuevos logs a la lista existente
    new_logs = [log.dict() for log in logs_proc]
    existing_logs.extend(new_logs)

    # Escribimos la lista de logs en el archivo logs.json
    with open(logs_file, 'w') as file:
        json.dump(existing_logs, file, indent=4)

    return {"received": True}

@app.post("/meminfo")
def get_meminfo(mem_info: LogSystem):
    logs_file = 'logs/meminfo.json'

    # Checamos si el archivo existe y tiene contenido
    if os.path.exists(logs_file) and os.stat(logs_file).st_size > 0:
        try:
            # Leemos el archivo logs.json
            with open(logs_file, 'r') as file:
                existing_logs = json.load(file)
        except json.JSONDecodeError:
            existing_logs = []
    else:
        existing_logs = []

    # Agregamos los nuevos datos de memoria a la lista existente
    new_logs = [mem_info.dict()] #Convertimos el objeto LogSystem en un diccionario y lo ponemos en una lista
    existing_logs.extend(new_logs)

    # Escribimos la lista de datos en el archivo logs.json
    with open(logs_file, 'w') as file:
        json.dump(existing_logs, file, indent=4)

    return {"received": True}

@app.get("/generate_meminfo_graph")
def generate_meminfo_plot():
    logs_file = 'logs/meminfo.json'
    img_folder = 'img'

    if not os.path.exists(logs_file):
        return {"error": "No meminfo logs file found."}

    if not os.path.exists(img_folder):
        os.makedirs(img_folder)

    try:
        with open(logs_file, 'r') as file:
            logs = json.load(file)

        timestamps = [log['timestamp'] for log in logs]
        total_ram = [log['total_ram'] for log in logs]
        free_ram = [log['free_ram'] for log in logs]
        used_ram = [log['used_ram'] for log in logs]

        plt.figure(figsize=(12, 8))
        plt.fill_between(timestamps, [0]*len(free_ram), free_ram, color='lightblue', label='Free RAM')
        plt.fill_between(timestamps, free_ram, [f + u for f, u in zip(free_ram, used_ram)], color='salmon', label='Used RAM')
        plt.plot(timestamps, total_ram, color='black', linestyle='--', label='Total RAM')
        plt.xlabel('Timestamp')
        plt.ylabel('RAM (MB)')
        plt.title('RAM Usage Over Time')
        plt.xticks(rotation=45)
        plt.legend()
        plt.tight_layout()

        plot_path = os.path.join(img_folder, 'meminfo_plot.png')
        plt.savefig(plot_path)
        plt.close()

    except Exception as e:
        return {"error": f"Error generating meminfo plot: {str(e)}"}

    return {"plot_path": plot_path}

@app.get("/generate_memory_graph")
def generate_memory_plot():
    logs_file = 'logs/logs.json'
    img_folder = 'img'

    if not os.path.exists(logs_file):
        return {"error": "No logs file found."}

    if not os.path.exists(img_folder):
        os.makedirs(img_folder)

    try:
        with open(logs_file, 'r') as file:
            logs = json.load(file)

        process_names = [log['name'] for log in logs]
        memory_usages = [log['memory_usage'] for log in logs]

        plt.figure(figsize=(10, 6))
        plt.bar(process_names, memory_usages, color='skyblue')
        plt.xlabel('Process Name')
        plt.ylabel('Memory Usage (%)')
        plt.title('Memory Usage by Process')
        plt.xticks(rotation=45)
        plt.tight_layout()

        plot_path = os.path.join(img_folder, 'memory_usage_plot.png')
        plt.savefig(plot_path)
        plt.close()

    except Exception as e:
        return {"error": f"Error generating memory plot: {str(e)}"}

    return {"plot_path": plot_path}
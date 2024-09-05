# Dockerfile para un contenedor de bajo consumo
FROM python:3.9-slim

RUN pip install flask

COPY log_cpu.py /log_cpu.py

CMD ["python", "log_cpu.py"]
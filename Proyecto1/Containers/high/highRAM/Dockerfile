#Dockerfile para alto consumo de RAM
FROM python:3.8-slim

RUN pip install numpy
COPY high_ram.py /high_ram.py

CMD ["python","high_ram.py"]
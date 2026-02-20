FROM python:3.10-alpine
WORKDIR /code
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY app.py .
EXPOSE 8000
ENTRYPOINT ["python3"]
CMD ["app.py"]

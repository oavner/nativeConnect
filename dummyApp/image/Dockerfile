FROM python:3-slim
ARG SERVER_BASE_DIR=/usr/local/bin
WORKDIR $SERVER_BASE_DIR
COPY app.py ./app.py
EXPOSE 9999/tcp
RUN pip3 install flask
ENTRYPOINT ["python3","app.py"]

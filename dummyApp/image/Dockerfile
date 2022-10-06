FROM python:3-slim
COPY ./dummyApp/image/app.py ./
EXPOSE 9999/tcp
RUN pip3 install flask
ENTRYPOINT ["python3","app.py"]

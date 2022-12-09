FROM python:3.9-alpine

RUN  pip3 install flask requests

RUN apk update
RUN apk upgrade
RUN apk add curl

WORKDIR /
COPY . .

EXPOSE	5000
CMD	python3 frontend.py
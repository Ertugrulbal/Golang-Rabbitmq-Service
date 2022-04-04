# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
RUN mkdir programA
COPY ./programA/*.go ./programA

# RUN go run ./programA/*.go 
CMD ["go", "run", "./programA/programA.go"]



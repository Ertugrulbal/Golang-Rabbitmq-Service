# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
RUN mkdir programB
COPY ./programB/*.go ./programB

# RUN go run ./programB/*.go 
CMD ["go", "run", "./programB/programB.go"]



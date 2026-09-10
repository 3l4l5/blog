FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY template ./template/
COPY src ./src/
COPY articles ./articles/

RUN go build -o main ./src/
RUN ./main build


FROM nginx:1.31.5

COPY nginx.conf .

COPY --from=builder /app/dist /usr/share/nginx/html/





FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY template ./template/
COPY src ./src/
COPY articles ./articles/

RUN go build -o main ./src/
RUN ./main build

FROM node:26-alpine3.23 AS tailwind

WORKDIR /app
COPY package.json package-lock.json ./
COPY assets/ ./assets
COPY --from=builder /app/dist ./dist
RUN npm ci
RUN npx @tailwindcss/cli \
  -i ./assets/input.css \
  -o ./dist/assets/style.css

FROM nginx:1.31.5

COPY nginx.conf .

COPY --from=tailwind /app/dist /usr/share/nginx/html/





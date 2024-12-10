# Etapa 1: Construcción de la aplicación Node.js
FROM node:lts-alpine AS build-node

WORKDIR /app/staticv2

COPY ./staticv2/package*.json ./
RUN npm install
RUN npm audit fix --force

COPY ./staticv2/ ./
RUN npm run build

# Etapa 2: Construcción de la aplicación Go
FROM golang:1.21-alpine AS build-go

WORKDIR /app

COPY . .
RUN go get ./src
RUN go build -o ./gosfv2 ./src

# Etapa 3: Imagen final
FROM alpine

WORKDIR /app

# Copiar los archivos necesarios de las etapas anteriores
COPY --from=build-node /app/staticv2/dist ./static
COPY --from=build-go /app/gosfv2 .
RUN chmod +x ./gosfv2

LABEL Name=gosfv2 Version=1.0.0
CMD [ "/app/gosfv2", "-config", "/app/config/config.json" ]

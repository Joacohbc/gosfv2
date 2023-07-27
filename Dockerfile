# La imagen base es la oficial de Go (ultima versión estable) en Alpine 
FROM golang:alpine 
RUN apk add --no-cache git gcc libc-dev

WORKDIR /app

# Copio los cosas que necesito para compilar
COPY . .

# Compilo y creo el ejecutable
RUN go get ./src
RUN go build -o ./gosfv2 ./src
RUN apk --no-cache add ca-certificates

# Agrego permisos de ejecución
RUN chmod +x ./docker-entrypoint.sh
RUN chmod +x ./gosfv2

LABEL Name=gosfv2 Version=1.0.0
EXPOSE 80
ENTRYPOINT [ "/app/docker-entrypoint.sh" ]
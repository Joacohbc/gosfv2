FROM golang:1.21 

WORKDIR /app

COPY . .

RUN apt update
RUN apt install nodejs npm -y
WORKDIR /app/staticv2
RUN npm install
RUN npm audit fix
RUN npm run build
WORKDIR /app

# Compilo y creo el ejecutable
RUN go get ./src
RUN go build -o ./gosfv2 ./src
RUN apt install ca-certificates

# Agrego permisos de ejecución
RUN chmod +x ./docker-entrypoint.sh
RUN chmod +x ./gosfv2

LABEL Name=gosfv2 Version=1.0.0
EXPOSE 80
ENTRYPOINT [ "/app/docker-entrypoint.sh" ]
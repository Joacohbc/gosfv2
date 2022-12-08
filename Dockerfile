FROM golang:alpine
RUN apk add --no-cache git gcc libc-dev

WORKDIR /app

# Copio los cosas que necesito para compilar
COPY . .

# Compilo y creo el ejecutable
RUN go get ./src
RUN go build -o ./gosfv2 ./src
RUN apk --no-cache add ca-certificates

LABEL Name=gosfv2 Version=1.0.0
CMD ["/app/gosfv2", "-config", "/app/config.json"]

# # build stage
# FROM golang:alpine
# RUN apk add --no-cache git gcc libc-dev
# WORKDIR /go/src/app
# COPY . .
# RUN go get -d -v ./...
# RUN go build -o /go/src/app/src/gosfv2 ./src 
# LABEL Name=gosfv2 Version=1.0.0
# RUN apk --no-cache add ca-certificates
# EXPOSE 3000
# CMD ["/go/src/app/src/gosfv2", "-config", "/go/src/app/src/config.json"]
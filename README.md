# GOSF v2

GOSF es un servidor de archivos que permite compartir archivos entre usuario de forma sencilla y segura, de forma pública o privada.

## Características

1. Autenticación de usuarios con **JWT Token**
    - Cantidad de **tiempo de expiración** configurable
    - **Cantidad de Tokens** por usuario configurable (cantidad sesiones activas)
2. **Subir**, **descargar**, **eliminar** y **modificar** archivos de forma sencilla
3. **Compartir archivos** de forma pública o privada, siempre y cuando el usuario esté autenticado  
4. **Control de acceso** a los archivos (solo el usuario que subió el archivo puede modificarlo o eliminarlo)
5. Altamente **configurable**, se puede configurar:
    - El puerto en el que se ejecuta el servidor
    - La ruta en la que se guardan los archivos
    - La ruta en la que se guardan los logs
    - Se pueden modificar la características de las conexiones a la base de datos de Redis y MySQL (host, puerto, usuario, contraseña, etc)

## Tecnologías

- Go & [Echo Framework](https://echo.labstack.com/guide/)
- HTML/CSS y JavaScript
- React / React Router
- Redis & MySQL
- Docker y Docker Compose
- REST API

## Tipos de Instalación

La instalación de tiene 2 opciones para su despliegue:

1. Utilizando Docker Compose, que despliega un contenedor de GOSF y otro de MySQL y Redis (Ya configurados)
2. Como binario que se ejecuta en el sistema operativo y se conecta a una base de datos MySQL y Redis (Previamente instaladas y configuradas por el usuario)

## Pre-requisitos

- Docker & Docker Compose (Solo instalación 1)
- Go (Solo instalación 2)
- MySQL (Solo instalación 2)
- Redis (Solo instalación 2)

Como configurar los pre-requisitos [aquí](./readme/Prerequisites.md).

## Instalación

### Opción 1 (Recomendada)

Para iniciar basta con "ejecutar" el [docker-compose.yml](./docker-compose.yml).

Por defecto el servicio corre en el puerto 80, para cambiar esto basta con modificar el puerto que se expone en el [config.env](./config.env). Ademas se puede cambiar otros parámetros de configuración del servicio (como la Volume path, y versiones de Tags, etc) en el mismo archivo.

- Clonar el repositorio

```bash
git clone https://github.com/Joacohbc/gosfv2; cd ./gosfv2;
```

- Crear los contenedores

```bash
docker compose build --no-cache
```

- Iniciar los contenedores

```bash
docker compose --env-file config.env up
```

### Opción 2

Para iniciar basta con ejecutar el binario (con las 2 base de datos corriendo)

Utilizando el binario pre-compilado:

- Clonar el repositorio

```bash
git clone https://github.com/Joacohbc/gosfv2; cd ./gosfv2;
```

- Ejecutar el binario pre-compilado (./gosfv2) y le indicio que utilice el archivo de configuración ./config.json

```bash

# Linux (64 bits)
./bin/gosfV2-64.bin -config ./config.json;

# Linux (32 bits)
./bin/gosfV2-32.bin -config ./config.json;

# Windows
./bin/gosfV2.exe -config ./config.json;
```

Si se quiere compilar el binario, se debe ejecutar el siguiente comando:

- Clonar el repositorio

```bash
git clone https://github.com/Joacohbc/gosfv2; cd ./gosfv2;
```

- Obtener todas las dependencias indicadas en el go.mod

```bash
go get ./src;
```

Compilar el código fuente

```bash
go build -o ./gosfv2 ./src;
```

- Ejecutar el binario anteriormente compilado (./gosfv2) y le indicio que utilice el archivo de configuración ./config.json

```bash
./gosfv2 -config ./config.json;
```

## Frontend Implementado

![Login and Register Page](/readme/Login%20and%20Register.png)

![User, Home, Notes](/readme//Main%20Pages.png)

![Share Overlay](/readme/Share%20Overlay.png)

![Delete Overlay PC](/readme/Delete%20Recording%20PC.gif)

![Delete Overlay Mobile](/readme/Delete%20Recording%20Mobile.gif)

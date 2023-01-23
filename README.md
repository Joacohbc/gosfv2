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

- Go
- [Echo Framework](https://echo.labstack.com/guide/)
- Redis
- MySQL
- Docker y Docker Compose
- HTML/CSS/JavaScript ES6 y JSON

## Instalación

La instalación de tiene 2 opciones para su despliegue:

1. Como binario que se ejecuta en el sistema operativo y se conecta a una base de datos MySQL y Redis (Previamente instaladas y configuradas por el usuario)
2. Utilizando Docker Compose, que despliega un contenedor de GOSF y otro de MySQL y Redis (Ya configurados)

## Requisitos

### Go

**_NOTA_**: Necesario solo si se utiliza la primera opción de despliegue

Necesita tener instalado [Go](https://golang.org/doc/install). Solo en el caso que se desea compilar el binario, ya que se puede utilizar el binario ya compilado que se encuentra en la carpeta [bin](./bin).

```bash
# Clono el repositorio
git clone https://github.com/Joacohbc/gosfv2;

# Obtener todas las dependencias indicadas en el go.mod
go get ./src

# Compilar el código fuente
go build -o ./gosfv2 ./src
```

### Docker

Necesita tener instalado [Docker](https://docs.docker.com/get-docker/). Solo en el caso que se desea utilizar la segunda opción de despliegue.

### MySQL

Necesita tener instalado [MySQL en el sistema operativo](https://dev.mysql.com/downloads/) o desde un contenedor [contenedor Docker](https://hub.docker.com/_/mysql). Se recomienda utilizar la versión 8.0, pero es compatible con la version 5.7 también.

Ademas se debe crear (si no existe) un usuario y una base de datos para que el servicio pueda crear las tablas necesarias para su correcto funcionamiento. Las variables que indican las propiedades de la conexión con el servicio de MySQL son:

Para el contenedor de Docker de GOSF son (dentro del archivo [config-gosf.env](./config-gosf.env)):

```bash
// ... 
DB_HOST_SQL="mysql-gosf"
DB_USER_SQL="gosf"
DB_PASSWORD_SQL="gosf"
DB_NAME_SQL="gosf"
DB_PORT_SQL=3306
DB_CHARSET_SQL="utf8"
// ... 
```

Para el binario de GOSF son (dentro del archivo [config.json](./config.json)):

```json
"db_host_sql": "localhost",
"db_user_sql": "gosf",
"db_password_sql": "gosf",
"db_name_sql": "gosf",
"db_port_sql": 3306,
"db_charset_sql": "utf8",
```

_Si se quiere utilizar MySQL como contenedor de Docker (con port-forwarding):_  

_NOTA_: Las configuraciones de MySQL indicadas son las predeterminadas, indicadas en los archivos de configuración. Si se desea cambiarlas, se debe modificar los archivos de configuración también.

Ejemplo de ejecución de un contenedor de MySQL 8.0:

```bash
docker run -d -p 3306:3306 --name mysql-docker \
    -e MYSQL_DATABASE="gosf" \
    -e MYSQL_USER="gosf" \
    -e MYSQL_PASSWORD="gosf" \
    -e MYSQL_ROOT_PASSWORD=1234 \
    mysql:8.0-debian
```

Ejemplo de ejecución de un contenedor de MySQL 5.7:

```bash
docker run -d -p 3306:3306 --name mysql-docker \
    -e MYSQL_DATABASE="gosf" \
    -e MYSQL_USER="gosf" \
    -e MYSQL_PASSWORD="gosf" \
    -e MYSQL_ROOT_PASSWORD=1234 \
    mysql:5.7-debian
```

### Redis

Necesita tener instalado [Redis en el sistema operativo](https://redis.io/topics/quickstart) o desde un [contenedor de Docker](https://hub.docker.com/_/redis). Se recomienda utilizar la version 7.0 o superior, pero es compatible con las version 6.0 y 6.2.7.

Ademas se debe configurar el servicio adecuadamente (contraseña de la base de datos, puerto, base de datos utilizables, etc) para que el servicio pueda crear las estructuras de datos para su correcto funcionamiento. Las variables que indican las propiedades de la conexión con el servicio de Redis son:

Para el contenedor de Docker de GOSF son (dentro del archivo [config-gosf.env](./config-gosf.env):

```bash
// ...
REDIS_HOST="localhost"
REDIS_PORT=6379
REDIS_PASSWORD=""
REDIS_DB=0
// ...
```

Para el binario de GOSF son (dentro del archivo [config.json](./config.json)):

```json
"redis_host": "localhost",
"redis_port": 6379,
"redis_password": "",
"redis_db": 0,
```

Ejemplo de ejecución de un contenedor de Redis 7.0.5:

```bash
docker run -d -p 6379:6379 --name redis-docker redis:7.0.5
```

Ejemplo de ejecución de un contenedor de Redis 6.2.7 (Debian Bullseye):

```bash
docker run -d -p 6379:6379 --name redis-docker redis:redis:6.2.7-bullseye
```

Ejemplo de ejecución de un contenedor de Redis 6.0 (Debian Bullseye):

```bash
docker run -d -p 6379:6379 --name redis-docker redis:6.0-bullseye
```

### Docker Compose

Necesita tener instalado [Docker Compose](https://docs.docker.com/compose/install/). Unicamente para la opción de despliegue numero 2.

## Iniciar (Opción 1)

Para iniciar basta con ejecutar el binario (con las 2 base de datos corriendo)

```bash
# Clono el repositorio
git clone https://github.com/Joacohbc/gosfv2; cd ./gosfv2;

# Obtener todas las dependencias indicadas en el g.mod
go get ./src;

# Compilar el código fuente
go build -o ./gosfv2 ./src;

# Ejecuto el binario anteriormente compilado (./gosfv2) y 
# le indicio que utilice el archivo de configuración ./config.json
./gosfv2 -config ./config.json;
```

## Iniciar (Opción 2)

Para iniciar basta con "ejecutar" el [docker-compose.yml](./docker-compose.yml).

Por defecto el servicio corre en el puerto 80 (hace port forwarding del puerto 80 donde ser escucha el servidor), para cambiar esto basta con modificar el puerto que se expone en el [config.env](./config.env). Ademas se puede cambiar otros parámetros de configuración del servicio (como la Volume path, y versiones de Tags, etc) en el mismo archivo.

```bash
# Clono el repositorio
git clone https://github.com/Joacohbc/gosfv2; cd ./gosfv2;

# Creo los contenedores y inicio los contenedores
docker compose build --no-cache; 
docker compose --env-file config.env up
```

# GOSF v2

GOSF es un servidor de archivos que permite compartir archivos de entre usuario de forma sencilla y segura. Se permite compartir archivos de forma pública o privada.

## Características

1. Autenticación de usuarios con JWT Token
    - Cantidad de tiempo de expiración configurable
    - Cantidad de Tokens por usuario configurable (cantidad sesiones activas)
2. Permite subir, descargar, eliminar y modificar archivos de forma sencilla
3. Permite compartir archivos de forma pública o privada, siempre y cuando el usuario esté autenticado  
4. Control de acceso a los archivos (solo el usuario que subió el archivo puede modificarlo o eliminarlo)
5. Altamente configurable, se puede configurar:
    - El puerto en el que se ejecuta el servidor
    - La ruta en la que se guardan los archivos
    - La ruta en la que se guardan los logs
    - Se pueden modificar la características de las conexiones a la base de datos de Redis y MySQL (host, puerto, usuario, contraseña, etc)

## Instalación

La instalación de GOSF es muy sencilla, y tiene 2 opciones para su despliegue

1. Como binario que se ejecuta en el sistema operativo y se conecta a una base de datos MySQL y Redis (Previamente instaladas y configuradas por el usuario)
2. Utilizando Docker Compose, que despliega un contenedor de GOSF y otro de MySQL y Redis (Ya configurados)

## Requisitos

### Go / Golang

*NOTA*: Necesario solo si se utiliza la primera opción de despliegue

Necesita tener instalado [Go](https://golang.org/doc/install) (Solo en el caso que se desea compilar el binario). Se recomienda utilizar la ultima versión.

```bash
# Clono el repositorio
git clone https://github.com/Joacohbc/gosfv2;

# Obtener todas las dependencias indicadas en el g.mod
go get ./src

# Compilar el código fuente
go build -o ./gosfv2 ./src
```

### Docker

Necesita tener instalado [Docker](https://docs.docker.com/get-docker/).

### MySQL

Necesita tener instalado [MySQL en el sistema operativo](https://dev.mysql.com/downloads/) o desde un contenedor [contenedor Docker](https://hub.docker.com/_/mysql). Se recomienda utilizar la versión 8.0, pero es compatible con la version 5.7 también.

Ademas se debe crear (si no existe) un usuario y una base de datos para que el servicio pueda crear las tablas necesarias para su correcto funcionamiento. Las variables que indican las propiedades de la conexión con el servicio de MySQL son:

Para el contenedor de Docker de GOSF son (dentro del archivo [config-docker-compose.env](./config-docker-compose.env)):

```bash
DB_HOST_SQL="localhost"
DB_USER_SQL="gosf"
DB_PASSWORD_SQL="gosf"
DB_NAME_SQL="gosf"
DB_PORT_SQL=3306
DB_CHARSET_SQL="utf8"
```

Para el binario de GOSF son (dentro del archivo [config.json](./config.json)):

```json
"db_host_sql": "localhost",
"db_user_sql":"gosf",
"db_password_sql":"gosf",
"db_name_sql":"gosf",
"db_port_sql": 3306,
"db_charset_sql": "utf8",
```

*Si se quiere utilizar MySQL como contenedor de Docker:*

*NOTA*: Las configuraciones de MySQL indicadas son las predeterminadas, indicadas en los archivos de configuración. Si se desea cambiarlas, se debe modificar los archivos de configuración también.

```bash
# Ejemplo de ejecución de un contenedor de MySQL 8.0
docker run -d -p 3306:3306 --name mysql-docker \
    -e MYSQL_DATABASE="gosf" \
    -e MYSQL_USER="gosf" \
    -e MYSQL_PASSWORD="gosf" \
    -e MYSQL_ROOT_PASSWORD=1234 \
    mysql:8.0-debian

# Ejemplo de ejecución de un contenedor de MySQL 5.7
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

Para el contenedor de Docker de GOSF son (dentro del archivo [config-docker-compose.env](./config-docker-compose.env)):

```bash
REDIS_HOST="localhost"
REDIS_PORT=6379
REDIS_PASSWORD=""
REDIS_DB=0
```

Para el binario de GOSF son (dentro del archivo [config.json](./config.json)):

```json
"redis_host": "localhost",
"redis_port": 6379,
"redis_password": "",
"redis_db": 0,
```

```bash
# Ejemplo de ejecución de un contenedor de Redis 7.0.5
docker run -d -p 6379:6379 --name redis-docker redis:7.0.5

# Ejemplo de ejecución de un contenedor de Redis 6.2.7 (Debian Bullseye)
docker run -d -p 6379:6379 --name redis-docker redis:redis:6.2.7-bullseye

# Ejemplo de ejecución de un contenedor de Redis 6.0 (Debian Bullseye)
docker run -d -p 6379:6379 --name redis-docker redis:6.0-bullseye
```

### Docker Compose

Necesita tener instalado [Docker Compose](https://docs.docker.com/compose/install/). Unicamente para la opción de despliegue numero 2.

## Iniciar (Opción 1)

Para iniciar basta con ejecutar el binario (con las 2 base de datos corriendo)

```bash
# Clono el repositorio
git clone https://github.com/Joacohbc/gosfv2;

# Obtener todas las dependencias indicadas en el g.mod
go get ./src;

# Compilar el código fuente
go build -o ./gosfv2 ./src;

# Ejecuto el binario anteriormente compilado (./gosfv2) y 
# le indicio que utilice el archivo de configuración ./config.json
./gosfv2 -config ./config.json;
```

## Iniciar (Opción 2)

Para iniciar basta con "ejecutar" el [docker-compose.yml](./docker-compose.yml). Dentro de la [carpeta docker](./docker/), hay más docker-compose.yml para las versiones de MySQL y Redis anteriores.

Por defecto el servicio corre en el puerto 3000, para cambiar esto basta con modificar el la variable de ambiente `PORT` en el [config-docker-compose.env](./config-docker-compose.env) y el puerto que se expone en el [docker-compose.yml](./docker-compose.yml)

```bash
# Clono el repositorio
git clone https://github.com/Joacohbc/gosfv2;

# Compilo los binarios y inicio los contenedores
docker-compose build --no-cache; 
docker-compose up;
```

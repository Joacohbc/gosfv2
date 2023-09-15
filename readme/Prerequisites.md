# Configuración de Pre-requisitos

- Go (Solo instalación 1)
- Docker & Docker Compose (Solo instalación 2)
- MySQL (Solo instalación 1)
- Redis (Solo instalación 1)

## Go

**_NOTA_**: Solo si se utiliza la primera opción de despliegue de manera opcional

Necesita tener instalado [Go](https://golang.org/doc/install). Solo en el caso que se desea compilar el binario, ya que se puede utilizar el binario ya compilado que se encuentra en la carpeta [bin](./bin).

## Docker

**_NOTA_**: Solo si se utiliza la segunda opción de despliegue o para ejecutar el contenedor de MySQL y Redis

Necesita tener instalado [Docker](https://docs.docker.com/get-docker/).

## Docker Compose

**_NOTA_**: Necesario solo si se utiliza la segunda opción de despliegue

Necesita tener instalado [Docker Compose](https://docs.docker.com/compose/install/). Unicamente para la opción de despliegue numero 2.

## MySQL

**_NOTA_**: Necesario solo si se utiliza la primera opción de despliegue

Necesita tener instalado [MySQL en el sistema operativo](https://dev.mysql.com/downloads/) o desde un contenedor [contenedor Docker](https://hub.docker.com/_/mysql). Se recomienda utilizar la versión 8.0, pero es compatible con la version 5.7 también.

Ademas se debe crear (si no existe) un usuario y una base de datos para que el servicio pueda crear las tablas necesarias para su correcto funcionamiento.

Las variables que indican las propiedades de la conexión con el servicio de MySQL con el **binario de GOSF** son: ([config.json](./config_example.json))

```json
"db_host_sql": "localhost",
"db_user_sql": "gosf",
"db_password_sql": "gosf",
"db_name_sql": "gosf",
"db_port_sql": 3306,
"db_charset_sql": "utf8",
```

Estas variables se pueden modificar en el archivo de configuración ([config.json](./config_example.json)).

### MySQL con Docker

Si se quiere utilizar MySQL como contenedor de Docker (con port-forwarding) se debe ejecutar el siguiente comando:

**_NOTA_**: Las configuraciones de MySQL indicadas son las predeterminadas, indicadas en los archivos de configuración. Si se desea cambiarlas, se debe modificar los archivos de configuración también.

Ejemplo de ejecución de un contenedor de MySQL 8.0:

```bash
docker run -d -p 3306:3306 --name mysql-docker \
    -e MYSQL_DATABASE="gosf" \
    -e MYSQL_USER="gosf" \
    -e MYSQL_PASSWORD="gosf" \
    -e MYSQL_ROOT_PASSWORD="gosf" \
    mysql:8.0-debian
```

Ejemplo de ejecución de un contenedor de MySQL 5.7:

```bash
docker run -d -p 3306:3306 --name mysql-docker \
    -e MYSQL_DATABASE="gosf" \
    -e MYSQL_USER="gosf" \
    -e MYSQL_PASSWORD="gosf" \
    -e MYSQL_ROOT_PASSWORD="gosf" \
    mysql:5.7-debian
```

## Redis

**_NOTA_**: Necesario solo si se utiliza la primera opción de despliegue

Necesita tener instalado [Redis en el sistema operativo](https://redis.io/topics/quickstart) o desde un [contenedor de Docker](https://hub.docker.com/_/redis). Se recomienda utilizar la version 7.0 o superior, pero es compatible con las version 6.0 y 6.2.7.

Ademas se debe configurar el servicio adecuadamente (contraseña de la base de datos, puerto, base de datos utilizables, etc) para que el servicio pueda crear las estructuras de datos para su correcto funcionamiento.

Las variables que indican las propiedades de la conexión con el servicio de Redis con el **binario de GOSF** son:

```json
"redis_host": "localhost",
"redis_port": 6379,
"redis_password": "",
"redis_db": 0,
```

Estas variables se pueden modificar en el archivo de configuración ([config.json](./config_example.json)).

### Redis con Docker

Si se quiere utilizar Redis como contenedor de Docker (con port-forwarding):  

**_NOTA_**: Las configuraciones de Redis indicadas son las predeterminadas, indicadas en los archivos de configuración. Si se desea cambiarlas, se debe modificar los archivos de configuración también.

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

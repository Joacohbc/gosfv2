<div align="center">

# GOSF v2

**La solución definitiva para la gestión y compartición de archivos: potencia, seguridad y simplicidad en un solo servidor auto-alojado.**

[![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?&logo=go&logoColor=white)](#)
[![React](https://img.shields.io/badge/React-%2320232a.svg?logo=react&logoColor=%2361DAFB)](#)
[![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=fff)](#)
[![MySQL](https://img.shields.io/badge/MySQL-4479A1?logo=mysql&logoColor=fff)](#)
[![Redis](https://img.shields.io/badge/Redis-%23DD0031.svg?logo=redis&logoColor=white)](#)

</div>

---

GOSF v2 no es solo un servidor de archivos; es tu propia nube privada, rápida y altamente configurable. Diseñado para quienes buscan simplicidad sin sacrificar la seguridad, permite gestionar tus datos con total control, ofreciendo autenticación robusta y un despliegue ágil mediante contenedores. Ya sea para uso público o privado, GOSF v2 es el puente seguro entre tus archivos y tus usuarios.

## 🛠️ Tech Stack

### Backend
[![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?&logo=go&logoColor=white)](#)
[![Echo](https://img.shields.io/badge/Echo-00ADD8?logo=echo&logoColor=white)](#)
[![ent](https://img.shields.io/badge/ent-5487A6?logo=ent&logoColor=white)](#)
[![Postman](https://img.shields.io/badge/Postman-FF6C37?logo=postman&logoColor=white)](#)

### Frontend
[![React](https://img.shields.io/badge/React-%2320232a.svg?logo=react&logoColor=%2361DAFB)](#)
[![React Router](https://img.shields.io/badge/React_Router-CA4245?logo=react-router&logoColor=white)](#)
[![Bootstrap](https://img.shields.io/badge/Bootstrap-7952B3?logo=bootstrap&logoColor=fff)](#)
[![JavaScript](https://img.shields.io/badge/JavaScript-F7DF1E?logo=javascript&logoColor=000)](#)

### Infraestructura & Base de Datos
[![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=fff)](#)
[![MySQL](https://img.shields.io/badge/MySQL-4479A1?logo=mysql&logoColor=fff)](#)
[![Redis](https://img.shields.io/badge/Redis-%23DD0031.svg?logo=redis&logoColor=white)](#)

---

## Características

1. **Autenticación de usuarios con JWT Token**
    - Cantidad de **tiempo de expiración** configurable.
    - **Cantidad de Tokens** por usuario configurable (sesiones activas).
2. **Gestión completa**: Subir, descargar, eliminar y modificar archivos de forma sencilla.
3. **Compartición Flexible**: Archivos públicos o privados para usuarios autenticados.
4. **Control de Acceso**: Seguridad garantizada; solo el propietario puede modificar o eliminar sus archivos.
5. **Altamente Configurable**: Personaliza puertos, rutas de almacenamiento, logs y conexiones a bases de datos.

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
docker compose up -d
```

- Dentro del volumen del contenedor, se encuentra el archivo de configuración `config.json` que debe ser modificado para que se ajuste a las necesidades de su sistema

```bash
# gosf-data es el volumen del contenedor
# Revisar el config_example.json como base para la configuración
nano gosf-data/config/config.json
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

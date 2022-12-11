#!/bin/sh

if [ -z ${JWT_KEY+x} ]; then
    echo "JWT_KEY variable must be set"
    exit 1
fi

if [ -z ${JWT_MINUTES+x} ]; then
    JWT_MINUTE=60
fi

if [ -z ${SERVER_PORT+x} ]; then
    SERVER_PORT=80
fi

if [ -z ${DB_HOST_SQL+x} ]; then
    MYSQL_HOST="localhost"
fi

if [ -z ${DB_USER_SQL+x} ]; then
    MYSQL_USER="gosf"
fi

if [ -z ${DB_PASSWORD_SQL+x} ]; then
    DB_PASSWORD_SQL="gosf"
fi

if [ -z ${DB_NAME_SQL+x} ]; then
    DB_NAME_SQL="gosf"
fi

if [ -z ${DB_PORT_SQL+x} ]; then
    DB_PORT_SQL=3306
fi

if [ -z ${DB_CHARSET_SQL+x} ]; then
    DB_CHARSET_SQL="utf8"
fi

if [ -z ${REDIS_HOST+x} ]; then
    REDIS_HOST="localhost"
fi

if [ -z ${REDIS_PORT+x} ]; then
    REDIS_PORT=6379
fi

if [ -z ${REDIS_PASSWORD+x} ]; then
    REDIS_PASSWORD=""
fi

if [ -z ${REDIS_DB+x} ]; then
    REDIS_DB=0
fi

if [ -z ${MAX_TOKEN_PER_USER+x} ]; then
    MAX_TOKEN_PER_USER=5
fi

echo "
{
    \"jwt_key\": \"$JWT_KEY\",
    \"jwt_minutes\": $JWT_MINUTES,
    \"port\": $PORT,
    \"log_dir_path\": \"./logs\",
    \"db_host_sql\": \"$DB_HOST_SQL\",
    \"db_user_sql\":\"$DB_USER_SQL\",
    \"db_password_sql\":\"$DB_PASSWORD_SQL\",
    \"db_name_sql\":\"$DB_NAME_SQL\",
    \"db_port_sql\": $DB_PORT_SQL,
    \"db_charset_sql\": \"$DB_CHARSET_SQL\",
    \"redis_host\": \"$REDIS_HOST\",
    \"redis_port\": $REDIS_PORT,
    \"redis_password\": \"$REDIS_PASSWORD\",
    \"redis_db\": $REDIS_DB,
    \"files_directory\": \"./files\",
    \"static_files\": \"./static\",
    \"max_token_per_user\": $MAX_TOKEN_PER_USER
}
" > /app/config.json

sh -c /app/gosfv2 -config /app/config.json
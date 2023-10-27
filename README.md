# env
The blocks below are the env settings for the modules 

## app
```bash
APP_TITLE=app_title
APP_HOST=app.example.com
APP_PATH=/path
APP_PORT=80
GIN_MODE=debug
TIMEZONE=
TRUST_PROXY=localhost
```

## mysql
```bash
AUTO_CREATE_DB_SCHEMA=false
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_NAME=dbname
DB_MAX_OPEN=10
DB_MAX_IDLE=10
DB_LIFE_TIME=120
DB_IDLE_TIME=90

DB_INIT_USER=user
DB_INIT_PASSWORD=password
DB_INIT_PARAMS=parseTime\=true&&multiStatements\=true

DB_USER=user
DB_PASSWORD=password
DB_PARAMS=parseTime\=true
```

## smtp
```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=user@example.com
SMTP_PASSWORD=password
SMTP_DISPLAY_NAME=name
SMTP_DISPLAY_EMAIL=user@example.com
```

## jwt
```bash
JWT_SIGNING_KEY=my_jwt_signing_key
JWT_SIGNING_METHOD=HS512 # Values: HS256, HS384, HS512
```

## face
```bash
DIR_FACE_RECOGNIZATION_MODELS=/path/to/the/models
```

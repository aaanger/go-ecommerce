# REST API онлайн магазина
- Используется фреймворк [gin-gonic/gin](https://github.com/gin-gonic/gin)
- Работа с БД PostgreSQL с использованием драйвера [jackc/pgx](github.com/jackc/pgx/v5), запуск в Docker, миграции осуществляются с помощью [pressly/goose](https://github.com/pressly/goose)
- Авторизация с JWT токенами
- Graceful Shutdown
- Структура приложения построена с подходом чистой архитектуры
- Конфигурация приложения с помощью библиотеки [spf13/viper](https://github.com/spf13/viper)
- Загрузка .env файла с [joho/godotenv](https://github.com/joho/godotenv)
# Как запустить
- ```make build``` сборка приложения
- ```make migrate``` миграции БД, если приложение запускается впервые
- ```make run``` запуск приложения

  ecommerce-app:
  build: ./
  ports:
  - "3000:3000"
  depends_on:
  - db
  environment:
  PSQL_HOST: db
  PSQL_USER: ${PSQL_USER}
  PSQL_PASSWORD: ${PSQL_PASSWORD}
  PSQL_DBNAME: ${PSQL_DBNAME}
  REDIS_ADDR: redis:6379
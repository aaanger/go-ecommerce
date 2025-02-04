# REST API онлайн магазина
- Используется фреймворк [gin-gonic/gin](https://github.com/gin-gonic/gin)
- Работа с БД PostgreSQL с использованием драйвера pgx, запуск в Docker, генерация файлов миграций с помощью [pressly/goose](https://github.com/pressly/goose)
- Авторизация с JWT токенами
- Graceful Shutdown
- Структура приложения построена с подходом чистой архитектуры
- Конфигурация приложения с помощью библиотеки [spf13/viper](https://github.com/spf13/viper)
- Загрузка .env файла с [joho/godotenv](https://github.com/joho/godotenv)
# Как запустить
- ```docker compose up``` поднимает БД
- ```make migrate``` миграция БД
- ```make run``` запускает приложение

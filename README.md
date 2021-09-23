# Short-link
## Использованы следующие технологии / концепции 

- Разработка Веб-Приложения на Go по REST API.  
- Работа с фреймворком gin-gonic/gin.  
- Hexagonal architecture
- Работа с БД Postgres используя библиотеку sqlx. 
- Регистрация и аутентификация. Работа с JWT. Middleware.
- Запуск интеграционных тестов в отдельном окружении
- Swagger
- Graceful Shutdown


### Команды makefile

- Запуск приложения первый раз с миграциями

```make start initial```

- Повторный запуск без миграций 

```make start```

- Запуск интеграционных тестов в отдельном окружении

```make test```

- Остановка приложения

```make down```

### Документация

[Swagger](http://localhost:8080/swagger/index.html)



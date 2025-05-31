# Мини-сервис "Цитатник"

REST API сервис для хранения и управления цитатами, написанный на Go.

## Функциональность

- Добавление новой цитаты (POST /quotes)
- Получение всех цитат (GET /quotes)
- Получение случайной цитаты (GET /quotes/random)
- Фильтрация по автору (GET /quotes?author={author})
- Удаление цитаты по ID (DELETE /quotes/{id})

## Требования

- Go 1.16+
- Для тестов: `github.com/stretchr/testify` (устанавливается автоматически)

## Установка и запуск

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/yourusername/quotes-service.git
   cd quotes-service
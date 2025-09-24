### User Reward Server 🎯
HTTP сервер для управления пользователями и системой вознаграждений на Go.

🚀 Возможности
✅ JWT аутентификация

✅ Управление пользователями

✅ Система выполнения заданий

✅ Реферальная система

✅ Leaderboard пользователей

✅ PostgreSQL база данных

✅ Docker контейнеризация

📋 API Endpoints
Публичные эндпоинты
* POST	/login	Аутентификация пользователя

Защищенные эндпоинты (требуют JWT)
* GET	/users/{id}/status	Информация о пользователе
* GET	/users/leaderboard	Топ пользователей по балансу
* POST	/users/{id}/task/complete	Выполнение задания
* POST	/users/{id}/referrer	Добавление реферера

## Запуск 
* docker-compose up --build

## Примеры запросов

/login POST | 
```
{
    "username": "Alice Johnson",
    "password": "password1"
}
```
Пример ответа:
```
{
    "message": "Login successful",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTgzODY3NTgsInVzZXJJRCI6MSwidXNlcm5hbWUiOiJBbGljZSBKb2huc29uIn0.qdHjwx3TQGtgqrULXba3ROSF3Pudv7utEYAlgoce8IA"
}
```
/users/{id}/status GET | Bearer Token {отправляем тот же токен, что пользователь получил при авторизации}

Пример ответа (id=3):
```
{
    "id": 3,
    "name": "Charlie Brown",
    "balance": 75,
    "created_at": "2025-09-19T16:45:03.127492Z"
}
```
C эндпоинтом /users/leaderboard принцип проверки схож.

/users/{id}/task/complete POST | Bearer Token {отправляем тот же токен, что пользователь получил при авторизации}

Тело запроса: 
```
{
    "task_type": "follow_twitter"
}
```
Пример ответа (id=2):
```
{
    "message": "Task completed successfully",
    "new_balance": 330,
    "reward": 30,
    "success": true,
    "task_type": "follow_twitter"
}
```
/users/{id}/referrer POST | Bearer Token {отправляем тот же токен, что пользователь получил при авторизации}

Тело запроса:
```
{
    "referrer_id": 1
}
```
Пример ответа (id=2):
```
{
    "success": true,
    "message": "Referrer added successfully",
    "referrer_id": 1,
    "reward": 100,
    "new_balance": 250
}
```


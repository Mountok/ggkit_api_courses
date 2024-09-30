# API for Courses Site

---
## Стек:
1. Go
2. Docker
3. Postgres
---
## Инструкция по запуску:

- Установить Docker и утилиту migrate

- После установки докер:

Установка образа postgres
```bash
docker pull postgres
```
Разворачиваем контейрен postgres
``` bash
docker run --name=postgtres -e POSTGRES_PASSWORD='admin' -p 5436:5432 -d --rm postgres 
```
Установка образа Go: из докер файла
```bash
docker build -t go-app .
```
Разворачиваем контейнер с Go
```bash
docker  run --name=go-praxis-back -p 80:8080 go-praxis
```
С помошью migrate надо выполнить миграцию
Находясь в главном каталоге выполните команду (постфик ```up``` - для поднятия, ```down``` - для отката):
```bash
migrate -path ./schema -database 'postgres://postgres:admin@localhost:5436/postgres?sslmode=disable' up  
```
 
# !!! Все что идет дальше не важно !!!

---
Запрос для рейтинга пользователей по счету:
```sql
SELECT p.user_id, p.full_name, p.score 
FROM users u 
JOIN profiles p 
ON u.id = p.user_id 
ORDER BY p.score DESC 
LIMIT 10;
```

postgresql://postgres:JMBBpmeyasyiQWhdpLxjESwTwsocyehv@junction.proxy.rlwy.net:38705/railway


Для создания файлом миграции:
```bash
migrate create -ext sql -dir ./schema -seq init
```
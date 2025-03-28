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
```bash
   docker run --name my-redis -d -p 6379:6379 redis
   
```

### Команды для выгрузки бд с докера 

docker exec -t postgtres pg_dump -U postgres -d postgres -f /tmp/dump.sql

docker cp postgtres:/tmp/dump.sql ./dump.sql

### Команды для загруски бд в дрокер с локальной выгрузки

cat dump.sql | docker exec -i postgres_container psql -U myuser -d mydatabase



# API DOC

### Subjects
    --GET--
    /api/subject
    --Request--
    none
    --Response--
    {
    "data": [{
      "id": int,
      "title": string,
      "image": string,
      "description": string }],
    "result": "ok"
    }
    --info--
    Возращает массив обьектов (курсов)
---
    --GET--
    /api/subject/{id}
    --Request--
    none
    --Response--
    {
    "data": [{
      "id": int,
      "title": string,
      "image": string,
      "description": string }],
    "result": "ok"
    }
    --info--
    Возращает массив с обьектом (курс)
---
    --POST--
    /api/subject
    --Request--
    type FormValue
    {
      title: string
      description: string
      image: string
    }
    --Response--
    {
      result: string,
      subject_id: int
    }
    --info--
    Возращает id созданного курса
---
 ### Theme
---
<br>
<br>
<br>

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
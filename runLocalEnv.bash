docker run --name=postgtres -e POSTGRES_PASSWORD='admin' -p 5436:5432 -d --rm postgres 
migrate -path ./schema -database 'postgres://postgres:admin@localhost:5436/postgres?sslmode=disable' up  
go run cmd/mian.go
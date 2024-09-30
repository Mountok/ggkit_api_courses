FROM golang:latest

COPY ./ ./
RUN go build -o main cmd/mian.go
CMD [ "./main" ]

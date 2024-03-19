FROM golang:1.21-alpine
WORKDIR /app
COPY go.mod go.sum ./
COPY ../generated ./generated
RUN go mod download
COPY . .
RUN go build -o user_service .

CMD ["./user_service"]
FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./app

EXPOSE 8080

# Command to run the application
CMD ["./main"]

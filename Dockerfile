FROM golang:1.20
WORKDIR /app
RUN go version

COPY ./ ./
RUN go mod download
COPY . .
RUN go build -o main ./main.go

CMD ["./main"]
FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o proxy proxy.go

CMD ["./proxy"]

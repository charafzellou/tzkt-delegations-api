FROM --platform=linux/amd64 golang:1.20.5

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY ./src/ ./

CMD ["air", "-c", ".air.toml"]
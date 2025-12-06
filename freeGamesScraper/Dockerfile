FROM golang:1.23

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN rm -rf /app/games_info.json

COPY . .
RUN go build -o myapp
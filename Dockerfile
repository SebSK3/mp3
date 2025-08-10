FROM golang:tip-alpine

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

RUN apk add --no-cache yt-dlp

RUN apk add --no-cache python3 mutagen


COPY assets ./assets
COPY cmd ./cmd
COPY internal ./internal

ENV DOWNLOAD_DIR="/downloads"
ENV DB_DSN="/app/db.sqlite"

RUN CGO_ENABLED=0 GOOS=linux go build -o /mp3 ./cmd/web

RUN rm -rf ./*

RUN mkdir /downloads

EXPOSE 3435

CMD ["/mp3"]

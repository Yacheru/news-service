FROM golang:alpine AS builder

WORKDIR /news-service

ADD ../go.mod .

COPY . .

RUN go build -o ./build/news.exe ./cmd/main.go

FROM alpine

WORKDIR /news-service

COPY ./configs/.env /news-service/configs/.env
COPY --from=builder /news-service/build /news-service

CMD ["/news-service/news.exe"]
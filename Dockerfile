FROM golang:latest AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app .

#Second stage of build
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=build /app .

EXPOSE 1381

CMD ["./app", "api"]

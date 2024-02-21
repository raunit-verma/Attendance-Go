FROM golang

COPY . /app

WORKDIR /app

RUN go mod download

RUN go mod tidy

RUN go run -mod=mod github.com/google/wire/cmd/wire

RUN go build -o server .

EXPOSE 1025

CMD [ "./server" ]
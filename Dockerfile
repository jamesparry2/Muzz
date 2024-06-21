FROM golang:1.22.4 as builder

ENV CGO_ENABLED=0

COPY go.mod go.sum  ./
RUN go mod download

COPY . .

FROM builder as compile
RUN go build -o /muzz-api main.go

EXPOSE 5001

CMD ["/muzz-api"]
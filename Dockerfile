FROM golang:1.21

ADD . .
WORKDIR /
RUN go mod tidy
RUN go build -o /server ./cmd/server/
RUN go build -o /cli ./cmd/cli

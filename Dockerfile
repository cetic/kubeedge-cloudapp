# build stage
FROM golang as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build /app/cmd/main.go
## final stage

FROM scratch
COPY --from=builder /app/main /app/
COPY --from=builder /app/configs/config.yaml /app/configs/
CMD ["/app/main","-c","/app/configs/config.yaml"]
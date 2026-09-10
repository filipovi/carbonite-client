FROM golang:1-alpine AS builder
WORKDIR /go/src/admin
COPY . /go/src/admin
RUN go mod vendor
RUN ./bin/templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin ./cmd/client/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/client/bin/admin .
COPY --from=builder /go/src/client/.env.docker .env 
COPY --from=builder /go/src/client/images images/
CMD ["./admin"]

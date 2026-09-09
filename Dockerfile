FROM golang:1-alpine AS builder
WORKDIR /go/src/admin
COPY . /go/src/admin
RUN go mod vendor
RUN ./bin/templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin ./cmd/admin/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/admin/bin/admin .
COPY --from=builder /go/src/admin/.env.docker .env 
COPY --from=builder /go/src/admin/images images/
CMD ["./admin"]

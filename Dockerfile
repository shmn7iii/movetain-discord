FROM golang:1.17.8 as builder
WORKDIR /movetain-discord
COPY . /movetain-discord
RUN CGO_ENABLED=0 GOOS=linux go build -o main && chmod +x ./main

FROM alpine:3.15
WORKDIR /movetain-discord
RUN apk --no-cache add ca-certificates
COPY --from=builder /movetain-discord/main ./
CMD ["./main"]

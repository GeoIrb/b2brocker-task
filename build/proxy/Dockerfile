FROM golang:1.15-alpine AS build

WORKDIR /go/src/proxy
COPY . ./

RUN go build \
    -o /out/service \
    ./cmd/proxy/main.go

FROM alpine
WORKDIR /app
COPY --from=build /out/service /app/service
CMD ["/app/service"]

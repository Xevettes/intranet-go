FROM docker.io/golang:1.24-alpine AS builder

RUN apk update && apk upgrade

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLE=0 GOOS=linux go build -o /api ./cmd

FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=builder /api .

EXPOSE 5555

USER nonroot:nonroot

ENTRYPOINT ["/api"]
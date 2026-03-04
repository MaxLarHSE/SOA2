FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache make

RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.0

ENV PATH="/go/bin:${PATH}"


COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

###
FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/api.yaml .

EXPOSE 8080

CMD ["./app"]
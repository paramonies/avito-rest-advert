FROM golang:1.15-alpine

RUN mkdir /app
WORKDIR /app

COPY ./ ./
RUN go mod download

RUN apk update && apk upgrade && apk add --no-cache postgresql-client make gcc musl-dev

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

RUN go build -o apiserver ./cmd/apiserver/main.go
CMD [ "./apiserver" ]
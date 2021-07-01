FROM golang:1.15-alpine

RUN mkdir /app
WORKDIR /app

COPY ./ ./
RUN go mod download

# install psql
RUN apk update && apk add --no-cache postgresql-client
#RUN apt-get update
#RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

RUN go build -o apiserver ./cmd/apiserver/main.go
CMD [ "./apiserver" ]
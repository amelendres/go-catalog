FROM golang:1.17-alpine3.14

ARG APP_NAME
WORKDIR /go/src/${APP_NAME}
RUN apk add --no-cache make gcc musl-dev

COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["catalog"]
FROM golang:1.18-alpine3.16

ENV APP_HOME /go/src/hasmail

RUN mkdir -p "$APP_HOME"

WORKDIR "$APP_HOME"

COPY . "$APP_HOME"

EXPOSE 8000

CMD [ "go", "run", "./cmd/mailservice/main.go" ]

FROM golang:1.18-bullseye

ENV APP_HOME /go/src/hasmail
RUN mkdir -p "$APP_HOME"

WORKDIR "$APP_HOME"
COPY . .
EXPOSE 8000

CMD [ "go", "run", "./cmd/mailservice/main.go" ]

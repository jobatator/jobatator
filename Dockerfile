FROM golang:latest
LABEL maintainer="spamfree@matthieubessat.fr"

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o jobatator main.go
RUN go test

EXPOSE 8962
EXPOSE 8952

CMD ["jobatator"]
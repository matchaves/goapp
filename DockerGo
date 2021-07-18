FROM golang:1.15-alpine

WORKDIR /go/src/app
ADD main.go .

RUN apk update && apk add git
RUN go get -d -v ./...
RUN go install -v ./...

#ENV URLAPP1=http://app1/post
#ENV URLAPP2=http://app2/post
#ENV AUTH=http://auth/auth

EXPOSE 80
EXPOSE 8080

CMD ["app"]
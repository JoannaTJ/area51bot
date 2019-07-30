FROM golang:1.8
LABEL maintainer="Robert D"

WORKDIR /go/src/app
COPY . .

# Download all dependencies
RUN go get -d -v ./...
RUN go install -v ./...

# expose port 8080 to outside
EXPOSE 8080

CMD ["app"]

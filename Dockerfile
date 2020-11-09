FROM golang:alpine
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV HTTPPORT 19093
RUN go install
# RUN ls /go/bin
# RUN pwd
ENTRYPOINT sensibull-test
# CMD ['/go/bin/sensibull-test']
EXPOSE 19093

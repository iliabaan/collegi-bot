FROM golang:1.20

RUN go version
ENV GOPATH=/
ENV GOGC=400

COPY ./ ./

# build go app
RUN go mod download
RUN go build -o collegi-bot ./cmd/main.go

CMD ["./collegi-bot"]

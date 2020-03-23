FROM golang:1.13 AS build_img
ENV APP_DIR=/app
RUN mkdir -p $APP_DIR
COPY *.go $APP_DIR
WORKDIR $APP_DIR
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -o /kamino

ENTRYPOINT /kamino

FROM scratch

COPY --from=build_img /villip /usr/bin/kamino

ENTRYPOINT ["/usr/bin/kamino" ]

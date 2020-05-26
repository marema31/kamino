FROM golang:1.14 AS build_img
ENV APP_DIR=/app
RUN mkdir -p $APP_DIR
COPY *.go $APP_DIR
WORKDIR $APP_DIR
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -gcflags "all=-N -l" -o /kamino

ENTRYPOINT /kamino

FROM alpine

ENV GOSU_VERSION=1.12

RUN apk --update add git less openssh-client ca-certificates wget bash \
    && apk add --no-cache --virtual .gosu-deps \
        dpkg \
        gnupg \
        openssl \
        tzdata \
    && dpkgArch="$(dpkg --print-architecture | awk -F- '{ print $NF }')" \
    && wget -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch" \
    && wget -O /usr/local/bin/gosu.asc "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch.asc" \
    && export GNUPGHOME="$(mktemp -d)" \
    && gpg --batch --keyserver hkps://keys.openpgp.org --recv-keys B42F6819007F00F88E364FD4036A9C25BF357DD4 \
    && gpg --batch --verify /usr/local/bin/gosu.asc /usr/local/bin/gosu \
    && rm -r /usr/local/bin/gosu.asc \
    && chmod +x /usr/local/bin/gosu \
    && gosu nobody true \
    && apk del .gosu-deps \
    && rm -rf /var/lib/apt/lists/* \
    && rm /var/cache/apk/* \
    && mkdir /.ssh \
    && chmod 700 /.ssh \
    && touch /.ssh/id_rsa \
    && chmod u=r,g=,o= /.ssh/id_rsa \
    && echo -e "Host *" > /.ssh/config \
    && echo -e "  StrictHostKeyChecking no" >> /.ssh/config \
    && echo -e "  UserKnownHostsFile=/dev/null" >> /.ssh/config \
    && echo -e "  IdentityFile ~/.ssh/id_rsa" >> /.ssh/config

COPY --from=build_img /kamino /usr/bin/kamino
COPY ./docker-entrypoint.sh /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh" ]

FROM alpine:3.3

ADD *.go /nativerw/

RUN apk add --update bash \
  && apk --update add git bzr \
  && apk --update add go \
  && export GOPATH=/gopath \
  && REPO_PATH="github.com/Financial-Times/nativerw" \
  && mkdir -p $GOPATH/src/${REPO_PATH} \
  && mv nativerw/* $GOPATH/src/${REPO_PATH} \
  && cd $GOPATH/src/${REPO_PATH} \
  && go get -t ./... \
  && go build \
  && mv nativerw /nativerw-app \
  && apk del go git bzr \
  && rm -rf $GOPATH /var/cache/apk/*

CMD [ "/nativerw-app" ]

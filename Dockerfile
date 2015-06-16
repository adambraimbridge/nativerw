FROM fedora

ENV GOPATH /usr/local

RUN dnf -y install git bzr mercurial golang
RUN go get git.svc.ft.com/cp/nativerw

CMD $GOPATH/bin/nativerw -mongos=$MONGODB_URIS $GOPATH/src/git.svc.ft.com/cp/nativerw/config.json


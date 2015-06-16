# Native Store Reader Writer (nativerw)

__Writes any raw content/data from native CMS in mongoDB without transformation.
The same data can then be read from here just like from the original CMS.__

## Installation

You need [Go to be installed](https://golang.org/doc/install). Please read about Go and about [How to Write Go Code](https://golang.org/doc/code.html) before jumping right in. For example you will need Git, Mercurial, Bazaar installed and working, so that Go can use them to retrieve dependencies. For this additionally you will also need a computer etc. Hope this helps.

for the first time: `go get git.svc.ft.com/scm/cp/nativerw.git` or update: `go get -u git.svc.ft.com/scm/cp/nativerw.git`
	
`go install git.svc.ft.com/scm/cp/nativerw.git`


## Running

`$GOPATH/bin/nativerw.git $GOPATH/src/git.svc.ft.com/scm/cp/nativerw.git/config.json`

You can override the mongos with -mongos flag, e.g.

`$GOPATH/bin/nativerw.git -mongos=mongo1:port,mongo2:port $GOPATH/src/git.svc.ft.com/scm/cp/nativerw.git/config.json`

## Try it!

`curl -XPUT -H "X-Request-Id: 123" -H "Content-Type: application/json" localhost:8080/methode/221da02e-c853-48ed-8753-3d1540fa190f --data '{"uuid":"221da02e-c853-48ed-8753-3d1540fa190f", "test": "test" }'`

`curl -H "X-Request-Id: 123" localhost:8080/methode/221da02e-c853-48ed-8753-3d1540fa190f`

Look in your mongoDB for database _native-store_ and collection _methode_ and notice the things you've written.

Healthchecks: [http://localhost:8080/__health](http://localhost:8080/__health)

Good-to-go: [http://localhost:8080/__gtg](http://localhost:8080/__gtg)


## Managing the app on the FT remote hosts

You can easily start or stop the app and see the logs on this page: [http://remote-hostname:9001/](http://ftapp08074-lvpr-uk-int:9001/)

The following commands are also useful:

Check the app's status: `sudo supervisorctl status nativerw`

Starting or stopping the app: `sudo supervisorctl stop/start/restart nativerw`

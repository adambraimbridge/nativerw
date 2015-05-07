# Native Store Reader Writer (nativerw)

__Writes any raw content/data from native CMS in mongoDB without transformation.
The same data can then be read from here just like from the original CMS.__

## Installation

for first time: `go get git.svc.ft.com/scm/cp/nativerw.git` or update: `go get -u git.svc.ft.com/scm/cp/nativerw.git`
	
`go install git.svc.ft.com/scm/cp/nativerw.git`


## Running

`$GOPATH/bin/nativerw.git $GOPATH/src/git.svc.ft.com/scm/cp/nativerw.git/config.json`


## Try it!

`curl -XPUT -H "X-Request-Id: 123" -H "Content-Type: application/json" localhost:8080/methode/221da02e-c853-48ed-8753-3d1540fa190f --data '{"uuid":"221da02e-c853-48ed-8753-3d1540fa190f", "test": "test" }'`

`curl -H "X-Request-Id: 123" localhost:8080/methode/221da02e-c853-48ed-8753-3d1540fa190f`

Look in your mongoDB for database _native-store_ and collection _methode_ and notice the things you've written.

Healthchecks: [http://localhost:8080/__health](http://localhost:8080/__health)

Good-to-go: [http://localhost:8080/__gtg](http://localhost:8080/__gtg)


## Managing the app

You can easily start or stop the app and see the logs on this page: [http://remote-hostname:9001/](http://ftapp08074-lvpr-uk-int:9001/)

The following commands are useful to manage the application on the FT hosts:

Check the app's status: `sudo supervisorctl status nativerw`

Starting or stopping the app: `sudo supervisorctl stop/start/restart nativerw`

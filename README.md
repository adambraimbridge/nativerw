# Native Store Reader Writer (nativerw)

- backed by mongoDB
- mapping from idiomatic json representations into bson/mongodb types (e.g., converting RFC3339 dates to bson date types)

## Installation

	go get git.svc.ft.com/scm/cp/nativerw.git
	
	go install git.svc.ft.com/scm/cp/nativerw.git

## Running

    $GOPATH/bin/nativerw.git $GOPATH/src/git.svc.ft.com/scm/cp/nativerw.git/config.json

## Try it!

    curl localhost:8080/__health

	curl -XPUT -H "X-Request-Id: 123" localhost:8080/content/221da02e-c853-48ed-8753-3d1540fa190f --data '{"uuid":"221da02e-c853-48ed-8753-3d1540fa190f","publishedDate":"2014-11-12T20:05:47.000Z", "foo":"bar","child":{"child2":{"uuid":"bfa8e9c9-1b53-46ac-a786-7cd296d5cbd4"}}, "num":135}'

	curl -H "X-Request-Id: 123" localhost:8080/content/221da02e-c853-48ed-8753-3d1540fa190f

Look in your mongodb for database "testdb" and collection "content" and notice things with nice bson types.

---
## Manage the app with supervisord

The following commands are useful to manage the application on the FT hosts (deployed with puppet):

1. check status:

        a. supervisorctl status nativerw
        b. service supervisord status nativerw


2. stop / start / restart (without making conf changes available)

        supervisorctl stop/start/restart nativerw


3. restart applying conf changes

        supervisorctl update

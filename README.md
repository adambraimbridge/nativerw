#jsonapi

An experimental example of a generic-ish json read/write api in golang, backed by mongodb, demonstrating a proof of concept mapping from idiomatic json representations into bson/mongodb types (e.g., converting RFC3339 dates to bson date types)

This is an example or how to write mediocre go code with no tests - don't try to use this for anything serious.

##Installation

You will need a mongodb database called "testdb", then:

go get git.svc.ft.com/scm/cp/jsonapi-go.git/jsonapi

##Running

jsonapi

##Try it!

You now have an endpoint listening on 0.0.0.0:8082


#jsonapi

An experimental example of a generic-ish json read/write api in golang, backed by mongodb, demonstrating a proof of concept mapping from idiomatic json representations into bson/mongodb types (e.g., converting RFC3339 dates to bson date types)

This is an example or how to write mediocre go code with no tests - don't try to use this for anything serious.

##Installation

	go get git.svc.ft.com/scm/cp/jsonapi-go.git/jsonapi

##Running

	jsonapi

##Try it!

You now have an endpoint listening on 0.0.0.0:8082

Try:

	curl -XPUT localhost:8082/content/221da02e-c853-48ed-8753-3d1540fa190f --data '{"uuid":"221da02e-c853-48ed-8753-3d1540fa190f","publishedDate":"2014-11-12T20:05:47.000Z", "foo":"bar","child":{"child2":{"uuid":"bfa8e9c9-1b53-46ac-a786-7cd296d5cbd4"}}, "num":135}'

then:

	curl localhost:8082/content/221da02e-c853-48ed-8753-3d1540fa190f

Look in your mongodb for database "testdb" and collection "content" and notice things with nice bson types.


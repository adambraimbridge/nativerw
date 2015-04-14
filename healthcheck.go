package jsonapi

import "git.svc.ft.com/scm/gl/fthealth.git"

var HealthCheck = fthealth.Check{
    BusinessImpact:   "blah",
    Name:             "My check",
    PanicGuide:       "Don't panic",
    Severity:         1,
    TechnicalSummary: "Something technical",
    Checker:          func() error { return nil }, //TODO: create the real check
}

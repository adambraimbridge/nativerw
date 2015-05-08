package main

import (
	"git.svc.ft.com/scm/gl/fthealth.git"
)

const healthcheckColl = "healthcheck"

var sampleResource = Resource{
	UUID:        "cda5d6a9-cd25-4d76-8bad-9eaa35e85f4a",
	ContentType: "application/json",
	Content:     "{\"foo\": [\"a\",\"b\"], \"bar\": 10.4}",
}

func (m *MgoApi) writeHealthCheck() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Publishing won't work. Writing content to native store is broken.",
		Name:             "CanWriteToMongoDB",
		PanicGuide:       "https://sites.google.com/a/ft.com/technology/systems/dynamic-semantic-publishing/extra-publishing/nativerw-runbook",
		Severity:         1,
		TechnicalSummary: "Writing to mongoDB is broken. Check mongoDB is up, its disk space, ports, network.",
		Checker:          m.checkWritable,
	}
}

func (m *MgoApi) checkWritable() error {
	return m.Write(healthcheckColl, sampleResource)
}

var sampleUUID = "cda5d6a9-cd25-4d76-8bad-9eaa35e85f4a"

func (m *MgoApi) readHealthCheck() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Reading content from native store is broken.",
		Name:             "CanReadFromMongoDB",
		PanicGuide:       "https://sites.google.com/a/ft.com/technology/systems/dynamic-semantic-publishing/extra-publishing/nativerw-runbook",
		Severity:         1,
		TechnicalSummary: "Reading from mongoDB is broken. Check mongoDB is up, its disk space, ports, network.",
		Checker:          m.checkReadable,
	}
}

func (m *MgoApi) checkReadable() error {
	_, _, err := m.Read(healthcheckColl, sampleUUID)
	return err
}

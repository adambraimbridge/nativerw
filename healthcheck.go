package main

import "git.svc.ft.com/scm/gl/fthealth.git"

const healthcheckColl = "healthcheck"

var sampleResource = Resource{
	UUID:        "cda5d6a9-cd25-4d76-8bad-9eaa35e85f4a",
	ContentType: "application/json",
	Content:     "{\"foo\": [\"a\",\"b\"], \"bar\": 10.4}",
}

func (m *MgoApi) buildHealthCheck() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Publishing won't work. Writing content to native store is broken.",
		Name:             "Writing to mongoDB.",
		PanicGuide:       "https://sites.google.com/a/ft.com/technology/systems/dynamic-semantic-publishing/extra-publishing/nativerw-runbook",
		Severity:         1,
		TechnicalSummary: "Writing to mongoDB is broken. Check mongoDB is up, its disk space, ports, network between.",
		Checker:          m.checkWritable,
	}
}

func (m *MgoApi) checkWritable() error {
	if err := m.Write(healthcheckColl, sampleResource); err != nil {
		return err
	}
	return nil
}

package main

import fthealth "github.com/Financial-Times/go-fthealth"

const healthcheckColl = "healthcheck"

var sampleResource = resource{
	UUID:        "cda5d6a9-cd25-4d76-8bad-9eaa35e85f4a",
	ContentType: "application/json",
	Content:     "{\"foo\": [\"a\",\"b\"], \"bar\": 10.4}",
}

func (m *mgoAPI) writeHealthCheck() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Publishing won't work. Writing content to native store is broken.",
		Name:             "Write to mongoDB",
		PanicGuide:       "https://sites.google.com/a/ft.com/technology/systems/dynamic-semantic-publishing/extra-publishing/native-store-reader-writer-run-book",
		Severity:         2,
		TechnicalSummary: "Writing to mongoDB is broken. Check mongoDB is up, its disk space, ports, network.",
		Checker:          m.checkWritable,
	}
}

func (m *mgoAPI) checkWritable() error {
	return m.Write(healthcheckColl, sampleResource)
}

var sampleUUID = "cda5d6a9-cd25-4d76-8bad-9eaa35e85f4a"

func (m *mgoAPI) readHealthCheck() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Reading content from native store is broken.",
		Name:             "Read from mongoDB",
		PanicGuide:       "https://sites.google.com/a/ft.com/technology/systems/dynamic-semantic-publishing/extra-publishing/native-store-reader-writer-run-book",
		Severity:         2,
		TechnicalSummary: "Reading from mongoDB is broken. Check mongoDB is up, its disk space, ports, network.",
		Checker:          m.checkReadable,
	}
}

func (m *mgoAPI) checkReadable() error {
	_, _, err := m.Read(healthcheckColl, sampleUUID)
	return err
}

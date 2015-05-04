class nativerw::monitoring {

  nagios::nrpe_checks::check_tcp {
    "${::certname}/1":
      host          => "localhost",
      port          => 8082,
      notes         => "check ${::certname} [$hostname] listening on HTTP port 8082 ";
  }

  nagios::nrpe_checks::check_http {
    "${::certname}/1":
      url           => "http://${::hostname}/__health",
      port          => 8082,
      expect        => "OK",
      size          => 1,
      action_url    => 'https://sites.google.com/a/ft.com/technology/systems/dynamic-semantic-publishing/extra-publishing/nativerw-runbook',
      notes         => "Severity 1 \\n Service unavailable \\n Native Reader Writer healthchecks are failing. Please check http://${::hostname}:8082/__health \\n\\n"
  }
}

class binary_writer::monitoring {

nagios::nrpe_checks::check_tcp {
  "${::certname}/1":
  host          => "localhost",
  port          => 8080,
  notes         => "check ${::certname} [$hostname] listening on api HTTP port 8080 ";

  "${::certname}/2":
  host          => "localhost",
  port          => 8081,
  notes         => "check ${::certname} [$hostname] listening on admin port 8081 ";
}

nagios::nrpe_checks::check_http {
"${::certname}/1":
url        => "http://localhost/build-info",
port       => "8080",
expect     => 'binary-writer',
size       => 1,
action_url => 'https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/image-binary-writer',
notes      => "check if build-info page shows ${::certname} [$hostname] as running";

"${::certname}/2":
url        => "http://localhost/ping",
port       => "8081",
expect     => 'pong',
size       => 1,
action_url => 'https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/image-binary-writer',
notes      => "check dropwizard ping healthcheck on ${::certname} [$hostname]";

"${::certname}/3":
url           => "http://${::hostname}/healthcheck",
port          => "8081",
expect        => 'OK',
size          => 1,
action_url    => 'https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/image-binary-writer',
notes         => "Severity 1 \\n Service unavailable \\n Binary Writer healthchecks are failing. Please check http://${::hostname}:8081/healthcheck \\n\\n"
}

}

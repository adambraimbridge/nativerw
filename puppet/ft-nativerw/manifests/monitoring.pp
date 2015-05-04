class nativerw::monitoring {

  $port = "8082"
  $cmd_check_http_json = "/usr/lib64/nagios/plugins/check_http_json.py --host http://$hostname:$port --path /__health --key_equals \"\$expression\$\""
  $nrpe_cmd_check_http_json = '/usr/lib64/nagios/plugins/check_nrpe -H $HOSTNAME$ -c check_http_json -a "$expression$"'
  $action_url = 'https://sites.google.com/a/ft.com/technology/systems/dynamic-semantic-publishing/extra-publishing/native-store-reader-writer-run-book'

  satellitesubscribe{"gateway-epel": channel_name => 'epel'}

  file { 'removeTemporaryInstallDir':
    ensure => absent,
    path => '/tmp/pip-build-root/pymongo',
    recurse => true,
    purge => true,
    force => true,
  }

  package {
    'python-pip':
      ensure  => 'installed',
      require => Satellitesubscribe["gateway-epel"];
    'argparse':
      ensure  => 'installed',
      provider => pip,
      require  => [Package['python-pip'], File['removeTemporaryInstallDir']];
  }

  # https://github.com/kovacshuni/nagios-http-json ; hash: 3b048f66c54e48607d195eef84e1746589492f39
  file { '/usr/lib64/nagios/plugins/check_http_json.py':
    ensure          => 'present',
    mode            => 0755,
    source          => 'puppet:///modules/nativerw/check_http_json.py',
  }

  file { '/etc/nrpe.d/check_http_json.cfg':
    ensure          => 'present',
    mode            => 0644,
    content         => "command[check_http_json]=${$cmd_check_http_json}\n"
  }

  @@nagios_command { "${hostname}_check_http_json":
    command_line => $nrpe_cmd_check_http_json,
    tag => $content_platform_nagios::client::tags_to_apply
  }

  @@nagios_service { "${hostname}_check_http_json_health_1":
    use                 => "generic-service",
    host_name           =>  "${::certname}",
    check_command       => "${hostname}_check_http_json!checks[0].ok,True",
    check_interval      => 1,
    action_url          => $action_url,
    notes_url           => $action_url,
    notes               => "See https://github.com/mzupan/nagios-plugin-mongodb#check-connection",
    service_description => "Check each host that is listed in the Mongo Servers group",
    display_name        => "${hostname}_check_mongodb",
    tag                 => $content_platform_nagios::client::tags_to_apply,
  }

  @@nagios_service { "${hostname}_check_http_json_health_2":
    use                 => "generic-service",
    host_name           =>  "${::certname}",
    check_command       => "${hostname}_check_http_json!checks[1].ok,True",
    check_interval      => 1,
    action_url          => $action_url,
    notes_url           => $action_url,
    notes               => "See https://github.com/mzupan/nagios-plugin-mongodb#check-connection",
    service_description => "Check each host that is listed in the Mongo Servers group",
    display_name        => "${hostname}_check_mongodb",
    tag                 => $content_platform_nagios::client::tags_to_apply,
  }

  nagios::nrpe_checks::check_tcp {
    "${::certname}/1":
      host          => "localhost",
      port          => 8082,
      notes         => "check ${::certname} [$hostname] listening on HTTP port 8082 ";
  }

  nagios::nrpe_checks::check_http {
    "${::certname}/1":
      url           => "http://${::hostname}/healthcheck",
      port          => 8082,
      expect        => "OK",
      size          => 1,
      action_url    => 'https://sites.google.com/a/ft.com/technology/systems/dynamic-semantic-publishing/extra-publishing/nativerw-runbook',
      notes         => "Severity 1 \\n Service unavailable \\n Native Reader Writer healthchecks are failing. Please check http://${::hostname}:8082/healthcheck \\n\\n"
  }
}

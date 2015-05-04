class nativerw::monitoring {

  $health_path = "http://$hostname:8082/__health"
  $cmd_check_http_json = "/usr/lib64/nagios/plugins/check_http_json.rb -u \"$health_path\" --element \"\$element\$\" --result \"\$result\$\" --warn \"\$warn\$\" --crit \"\$crit\$\""
  $nrpe_cmd_check_http_json = '/usr/lib64/nagios/plugins/check_nrpe -H $HOSTNAME$ -c check_http_json -a "$element$" "$result$" "$warn$" "$crit$"'
  $action_url = 'https://sites.google.com/a/ft.com/technology/systems/dynamic-semantic-publishing/extra-publishing/native-store-reader-writer-run-book'

  # check_http_json v1.3.1 https://github.com/phrawzty/check_http_json
  file { '/usr/lib64/nagios/plugins/check_http_json.rb':
    ensure          => 'present',
    mode            => 0755,
    source          => 'puppet:///modules/content_store_api_mongo/check_mongodb.py',
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

  @@nagios_service { "${hostname}_check_http_json":
    use                 => "generic-service",
    host_name           =>  "${::certname}",
    check_command       => "${hostname}_check_http_json!checks[0].ok!true!2!4",
    check_interval      => 1,
    action_url          => $action_url,
    notes_url           => $action_url,
    notes               => "See https://github.com/mzupan/nagios-plugin-mongodb#check-connection",
    service_description => "Check each host that is listed in the Mongo Servers group",
    display_name        => "${hostname}_check_mongodb",
    tag                 => $content_platform_nagios::client::tags_to_apply,
  }

  @@nagios_service { "${hostname}_check_http_json":
    use                 => "generic-service",
    host_name           =>  "${::certname}",
    check_command       => "${hostname}_check_http_json!checks[1].ok!true!2!4",
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

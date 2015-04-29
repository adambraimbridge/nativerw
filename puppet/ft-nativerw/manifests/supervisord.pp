class nativerw::supervisord {

  $supervisord_user = "supervisord"
  $supervisors_group = "supervisors"
  $supervisord_init_file = "/etc/init.d/supervisord"
  $supervisord_config_file = "/etc/supervisord.conf"
  $supervisord_log_dir = "/var/log/supervisor/"
  $binary_name = "nativerw"

  satellitesubscribe { 'gateway-epel':
      channel_name  => 'epel'
  }

  package { 'python-pip':
      ensure    => 'installed',
      require   => Satellitesubscribe["gateway-epel"];
  }

  package { 'supervisor':
      provider  => pip,
      ensure    => present,
      require   => [ Package['python-pip'] ]
  }

  user { $supervisord_user:
    ensure    => present
  }

  group { $supervisors_group:
    ensure    => present,
    members   => [ $supervisord_user, $binary_name ],
    require   => [ User[ $supervisord_user ] ]
  }

  file {
    $supervisord_init_file:
      mode      => "0755",
      content    => template("$module_name/supervisord.init.erb"),
      owner     => $supervisord_user,
      group     => $supervisors_group,
      require   => [ Package['supervisor'], User[$supervisord_user], Group[$supervisors_group] ];

    $supervisord_config_file:
      mode      => "0664",
      content    => template("$module_name/supervisord.conf.erb"),
      owner     => $supervisord_user,
      group     => $supervisors_group,
      require   => [ Package['supervisor'], User[$supervisord_user], Group[$supervisors_group] ];

    $supervisord_log_dir:
      ensure    => directory,
      owner     => $supervisord_user,
      group     => $supervisors_group,
      mode      => "0664",
      require   => [ Package['supervisor'], User[$supervisord_user], Group[$supervisors_group] ];
  }

  service { 'supervisord':
    ensure      => running,
    restart     => true,
    require     => [ File[$supervisord_init_file], File[$supervisord_config_file], File[$supervisord_log_dir] ];
  }
}

class nativerw::supervisord {

  $supervisord_init_file = "/etc/init.d/supervisord"
  $supervisord_config_file = "/etc/supervisord-custom.conf"

  satellitesubscribe {
    'gateway-epel':
      channel_name  => 'epel'
  }

  package {
    'python-pip':
      ensure    => 'installed',
      require   => Satellitesubscribe["gateway-epel"];
  }

  package {
    'supervisor':
      provider  => pip,
      ensure    => present,
      require   => [ Package['python-pip'] ]
  }

  service { 'supervisord':
    ensure      => running,
    restart     =>
  }

  file {
    $supervisord_init_file:
      mode      => "0755",
      content    => template("$module_name/supervisord.init.erb"),
      owner     => 'root',
      group     => 'root',
      require   => [ Package['supervisor'] ];

    $supervisord_config_file:
      mode      => "0664",
      source    => "puppet:///modules/$module_name/supervisord.conf",
      owner     => 'root',
      group     => 'root',
      require   => [ Package['supervisor'] ];
  }
}

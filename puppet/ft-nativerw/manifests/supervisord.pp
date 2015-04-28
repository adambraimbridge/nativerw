class nativerw::supervisord {

  $supevisord_init_file = "/etc/init.d/supervisord"

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

  file {
    $supevisord_init_file:
      mode      => "0755",
      source    => "puppet:///modules/$module_name/supervisord.init",
      owner     => 'root',
      group     => 'root',
      require   => [ Package['supervisor'] ]
  }
}

class nativerw {

  $binary_name = "nativerw"
  $install_dir = "/usr/local/$binary_name"
  $binary_file = "$install_dir/$binary_name"
  $log_dir = "/var/log/apps/$binary_name"
  $config_file = "/etc/$binary_name.json"
  $supervisord_user = "supervisord"

  class { 'common_pp_up': }
  class { "${module_name}::supervisord": }

  group { $binary_name:
    ensure    => present,
    members   => [ $binary_name, $supervisord_user ]
  }

  user { $binary_name:
    ensure    => present,
  }

  file {
    $install_dir:
      mode    => "0664",
      ensure  => directory,
      owner   => $binary_name,
      group   => $binary_name;

    $binary_file:
      ensure  => present,
      source  => "puppet:///modules/$module_name/$binary_name",
      owner   => $binary_name,
      group   => $binary_name,
      mode    => "0755",
      require => File[$install_dir];

    $config_file:
      content => template("$module_name/config.json.erb"),
      owner   => $binary_name,
      group   => $binary_name,
      mode    => "0664";

    $log_dir:
      ensure  => directory,
      owner   => $binary_name,
      group   => $binary_name,
      mode    => "0664",
  }
}

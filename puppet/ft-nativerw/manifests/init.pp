class nativerw {

  $binary_name = "nativerw"
  $install_dir = "/usr/local/${binary_name}"
  $binary_file = "${install_dir}/${binary_name}"
  $log_dir = "/var/log/apps"
  $config_file = "/etc/${binary_name}.json"

  class { 'common_pp_up': }

  supervisor::service { 'nativerw':
      ensure      => present,
      command     => "$binary_file $config_file",
      user        => 'root',
      group       => 'root',
      require     => [ File[$config_file] ];
  }

  file {
    $install_dir:
      mode    => "0755",
      ensure  => directory;

    $binary_file:
      path   => "/usr/local/$binary_name/$binary_name",
      ensure => present,
      source => "puppet:///modules/$module_name/$binary_name",
      mode   => "0755";

    $log_dir:
     mode   => "0755",
     ensure => directory;

    $config_file:
      mode    => "0755",
      content => template("$module_name/config.json.erb");
  }

  File[$install_dir]
  -> File[$binary_file]
  -> File[$config_file]
}

class nativerw {

  $binary_name = "nativerw"
  $install_dir = "/usr/local/${binary_name}"
  $binary_file = "${install_dir}/${binary_name}"
  $log_dir = "/var/log/apps"
  $config_file = "/etc/${binary_name}.json"

  class { 'common_pp_up': }
  class { "${module_name}::supervisord": }

  file {
    $install_dir:
      mode    => "0755",
      ensure  => directory;

    $binary_file:
      ensure  => present,
      source  => "puppet:///modules/$module_name/$binary_name",
      mode    => "0755",
      require => File[$install_dir];

    $log_dir:
     mode   => "0755",
     ensure => directory;

    $config_file:
      mode    => "0755",
      content => template("$module_name/config.json.erb");
  }
}

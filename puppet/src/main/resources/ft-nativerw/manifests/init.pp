class binary_writer {

  $dir_heap_dumps = "/var/log/apps/binary-writer-heap-dumps"
  $jar_name = "binary-writer-service.jar"

  class { "${module_name}::monitoring": }
  class { 'common_pp_up': }

  content_runnablejar { "${module_name}_runnablejar":
    service_name        => "${module_name}",
    service_description => 'Binary Writer',
    jar_name            => "${jar_name}",
    config_file_content => template("${module_name}/config.yml.erb"),
    artifact_location   => "${module_name}/${jar_name}",
    status_check_url    => "http://localhost:8081/ping";
  }

  file { "sysconfig":
    path    => "/etc/sysconfig/${module_name}",
    ensure  => 'present',
    content => template("${module_name}/sysconfig.erb"),
    mode    => '0644';
  }

  file { "heap-dumps-dir":
    path    => "${dir_heap_dumps}",
    owner   => "${module_name}",
    group   => "${module_name}",
    ensure  => 'directory',
    mode    => '0755';
  }

File['sysconfig']
  -> Content_runnablejar["${module_name}_runnablejar"]
  -> Class["${module_name}::monitoring"]
}

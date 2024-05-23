# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.

# More doc at https://app.vagrantup.com/mbr/boxes/postgres

## install vagrant and virtualbox
# $ vagrant up
Vagrant.configure("2") do |config|
  config.vm.box = "mbr/postgres"
  config.vm.network "forwarded_port", guest: 5432, host: 5432, host_ip: "127.0.0.1"
end
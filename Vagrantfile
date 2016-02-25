# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure(2) do |config|

  config.vm.box = "dduportal/boot2docker"
  config.vm.network :forwarded_port, guest: 9999, host: 9999
  config.vm.provision :docker do |d|

  config.vm.provider "virtualbox" do |vb|
    vb.customize ["modifyvm", :id, "--natdnsproxy1", "off"]
    vb.customize ["modifyvm", :id, "--natdnshostresolver1", "off"]
  end

    # BUILD APP
    d.build_image "/vagrant",
      args: "-t digibib/build -f /vagrant/Dockerfile.build"

    # COPY COMPILED APP
    d.run "builder",
      image: "digibib/build",
      daemonize: false,
      restart: false,
      args: "--rm -v /vagrant/build:/app",
      cmd: "cp app /app"
  end

  config.vm.provision :docker do |d|
    # BUILD MINIMAL DOCKER IMAGE WITH COMPILED APP
    d.build_image "/vagrant",
      args: "-t digibib/docker-compose-dot -f /vagrant/Dockerfile"

  end

end

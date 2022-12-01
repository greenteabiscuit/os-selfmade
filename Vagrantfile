# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.define "haribote-build" do |v|
  end

  config.vm.provider "virtualbox" do |vb|
    vb.customize ["modifyvm", :id, "--usb", "on"]
    vb.customize ["modifyvm", :id, "--usbehci", "off"]
    vb.customize ["modifyvm", :id, "--cableconnected1", "on"]
  end

  config.vm.box = "minimal/xenial64"

  config.vm.synced_folder "./", "/home/vagrant/workspace"

  config.vm.provision "shell", inline: <<-SHELL
    apt-get update
    apt-get install -y nasm make xorriso binutils gcc mtools
  SHELL
end

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/groovy64"
  config.vm.network "private_network", ip: "192.168.0.10"
  config.vm.provider "virtualbox" do |box|
    box.memory = "4096"
  end

  config.vm.provision "shell", inline: <<-SHELL
    apt-get update

    export DEBIAN_FRONTEND=noninteractive
    apt-get install -y mariadb-server

    printf "[mysqld]\nbind-address = 0.0.0.0" >> /etc/mysql/mariadb.cnf

    sudo service mysql restart
    mysql -uroot <<-EOL
      CREATE DATABASE queue;
      CREATE USER scyther IDENTIFIED BY '123';
      GRANT ALL ON queue.* TO scyther;

      CREATE TABLE queue.queues
      (
        id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
        name CHAR(64) NOT NULL,
        capacity BIGINT
      )
EOL
  SHELL
end

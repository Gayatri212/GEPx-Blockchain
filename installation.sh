#! /bin/bash
sudo yum install docker
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -a -G docker ec2-user
sudo chmod 666 /var/run/docker.sock
sudo curl -L "https://github.com/docker/compose/releases/download/1.27.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo yum install git
curl https://dl.google.com/go/go1.15.2.linux-amd64.tar.gz --output go1.15.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.15.2.linux-amd64.tar.gz
sudo chmod 777 $GOPATH/src
sudo amazon-linux-extras install java-openjdk11
sudo yum install java-1.8.0-openjdk
alternatives --config java
sudo yum install -y gcc-c++ make
curl -sL https://rpm.nodesource.com/setup_12.x | sudo -E bash -
sudo yum install -y nodejs

# Set envs permenantly
# open .bash_profile file using vi editor
# vi .bash_profile
# Copy paste following to that file
# export PATH=$PATH:/usr/local/go/bin
# export GOPATH=/usr/local/go


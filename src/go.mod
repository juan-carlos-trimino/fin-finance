//The default name for the generated executable would be finance.
module finance

//To find the system architecture type, execute:
//$ dpkg --print-architecture
//To upgrade golang
//$ sudo apt-get update
//$ sudo apt-get -y upgrade
//Notice the version and system architecture type: 1.22.0.linux-xxxxx
//$ wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
//$ sudo tar -xvf go1.22.0.linux-amd64.tar.gz -C /usr/local
//Set the Go path.
//$ echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
//$ source ~/.profile
//Verify the installation.
//$ go version
//Set up the Go workspace, if one is not set up.
//$ mkdir ~/go
//$ echo "export GOPATH=$HOME/go" >> ~/.profile
//$ source ~/.profile
go 1.22.0

require (
	github.com/google/uuid v1.6.0
	golang.org/x/crypto v0.18.0
)

require (
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

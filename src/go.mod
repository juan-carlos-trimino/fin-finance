//Uninstalling a library (or module) from a Go project primarily involves managing the go.mod file
//and the module cache.
//(1) Remove the Dependency from go.mod
//    The first step is to remove the line corresponding to the library you want to uninstall from
//    your project's go.mod file. This file lists all the direct dependencies of your module.
//(2) Run go mod tidy
//    After removing the dependency from go.mod, execute the following command in your terminal
//    within your project's root directory:
//    $ go mod tidy
//
//The default name for the generated executable would be finance.
module finance

go 1.24.3

require (
	github.com/juan-carlos-trimino/go-middlewares v1.1.3
	github.com/juan-carlos-trimino/gplogger v1.0.7
	github.com/juan-carlos-trimino/gposu v1.0.1
	github.com/juan-carlos-trimino/gps3storage v1.0.1
	github.com/juan-carlos-trimino/gpsessions v1.0.1
	golang.org/x/crypto v0.42.0
)

require (
	github.com/aws/aws-sdk-go v1.55.7 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/text v0.29.0 // indirect
)

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
	github.com/juan-carlos-trimino/go-middlewares v1.1.4
	github.com/juan-carlos-trimino/go-os v1.1.1
	github.com/juan-carlos-trimino/gplogger v1.0.7
	github.com/juan-carlos-trimino/gposu v1.0.1
	github.com/juan-carlos-trimino/gps3storage v1.0.1
	github.com/juan-carlos-trimino/gpsessions v1.0.1
	golang.org/x/crypto v0.46.0
)

require (
	github.com/aws/aws-sdk-go v1.55.7 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.6 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)

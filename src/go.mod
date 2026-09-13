//Uninstalling a library (or module) from a Go project primarily involves managing the go.mod file and the module cache.
//(1) Remove the Dependency from go.mod
//    The first step is to remove the line corresponding to the library you want to uninstall from your project's go.mod file.
//    This file lists all the direct dependencies of your module.
//(2) Run go mod tidy
//    After removing the dependency from go.mod, execute the following command in your terminal within your project's root directory:
//    $ go mod tidy
//
//The default name for the generated executable would be finance.
module finance

//How to update the Go version
//1. Download the new Go archive
//   $ wget -c https://golang.org/dl/go1.26.4.linux-amd64.tar.gz
//2. Remove the Old Version and Extract the New One
//   $ ls -al /usr/local/go
//   $ sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.26.4.linux-amd64.tar.gz
//3. Verify the Installation Path
//   $ echo $PATH | grep "/usr/local/go/bin"
//4. Confirm the Upgrade
//   $ go version
//5. Update the Project Configuration
//   The go.mod file should be located at the root directory of the Go project/repository.
//   $ go mod edit -go=1.26.4
//6. Update All Dependencies (Optional)
//   To update all of the project's current dependencies to their latest compatible minor or patch versions, run from the root directory of the Go
//   project/repository:
//   $ go get -u ./...
//   Add the -t flag to update dependencies used in the test files.
//   $ go get -u -t ./...
//7. Clean Up the go.mod File
//   Every time you add, update, or remove code, you must sync your go.mod and go.sum files:
//   $ go mod tidy
go 1.26.4

require (
	github.com/jackc/pgx/v5 v5.10.0
	github.com/juan-carlos-trimino/go-logger v1.0.9
	github.com/juan-carlos-trimino/go-middlewares v1.1.15
	github.com/juan-carlos-trimino/go-os v1.1.2
	github.com/juan-carlos-trimino/go-s3storage v1.0.4
	github.com/juan-carlos-trimino/go-sessions v1.0.8
	golang.org/x/crypto v0.57.0
)

require (
	github.com/aws/aws-sdk-go v1.55.8 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/juan-carlos-trimino/go-os v1.1.2 // indirect
	github.com/juan-carlos-trimino/go-sessions v1.0.8 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/text v0.42.0 // indirect
)

//To debug code inside a dependency module (moduleA) while running your main module, you need to point the Go  compiler to a local copy
//of moduleA instead of its cached, read-only version in $GOPATH/pkg/mod.
//
//Use a replace Directive
// You can explicitly force your main module to look at your local file system for moduleA.
// 1. Edit the Main go.mod
// 2. Add a replace directive at the bottom of your main module's go.mod file pointing to the relative or absolute path of your local dependency:
//
// module main-module
//
// go 1.26
//
// require github.com/xxx/moduleA v1.1.0
//
// Force Go to use your local copy for compiling and debugging.
// replace github.com/xxx/moduleA => ../moduleA
//
// 3. Run Delve or your IDE Debugger
// Run your debugger from the main module. You can now open files from ../moduleA, insert breakpoints, and step into them. The folder path
// specified in your replace directive must exactly match the folder path you have open in your IDE.
//
// (Note: Remember to remove this line before pushing your go.mod file to production.)
// replace github.com/juan-carlos-trimino/go-middlewares => ../../../go-middlewares/

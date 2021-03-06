# go-UCSActiveServers
This application was constructed to query one or more UCS Domains and report on the currently active UCS Servers.  Each UCS Domain is queried for all the discovered servers, then each server is interrogated and finally a CSV report is generated for the year, month, day and time.

An active server for this application is defined by;
1. A Service Profile is associated to the server
2. The server is powered up

If both of these factors are true at the time of the application running, then the blade will be considered to be in use.


## Setting up your GO environment
Depending on your particular environment, there are a number of ways to setup and install GO.  This repo was developed on a MAC and was installed using Brew.  For instructions on installing HomeBrew, please check [here](https://brew.sh/); and then entering;
```fish
> brew install go
```

If you do not want to use HomeBrew or you are running on a different platform, you can install the GO language using a binary from here;
https://golang.org/dl/

Once this has completed, open a cmd or terminal window and check GO has been installed and configured correctly;

Enter <b>echo $GOPATH</b>, hopefully you will be presented with a path and should be ready to go.

```fish
> echo $GOPATH
/path/to/go/bin/src/pkg folders
```

## Testing your GO environment
Once you have completed the above, its time to create a very simple test script to ensure everything is ready.

Go to a path where you are happy to store the source code for your application, this could be anywhere, including your desktop, documents, root folder, etc.

Create a folder and enter the directory.  Create a new file called "main.go" and enter the following code into it;

```go
package main

import "fmt"

func main() {
    fmt.Println("GO is working!")
}
```

At the command line, change directory using cd to the directory where your main.go file is and execute the following;
```fish
> go run main.go
```

You should see as output, something similar to;

"GO is working!"

If you reached this point, everything is working and you are ready to run the included code!

## Getting the code
There are a couple of ways you can get the code, depending on how comfortable you are with the command line and development environments;

You could download the zip file, [here](https://github.com/robjporter/go-UCSActiveServers/archive/master.zip).

You could use the command line git command to clone the repository to your local machine;
1. At the command line, change directory using cd to the directory where the repository will be stored.
2. Enter, git clone https://github.com/robjporter/go-UCSActiveServers.git
3. You will see output similar to the following while it is copied.
```fish
Cloning into `go-UCSActiveServers`...
remote: Counting objects: 10, done.
remote: Compressing objects 100% (8/8), done.
remove: Total 10 (delta 1), reused 10 (delta 1)
unpacking objects: 100% (10/10), done.
```
4. Change into the new directory, cd go-UCSActiveServers.
5. Move onto setting up the application.

## Application dependencies
For the application to work correctly, we need to get one dependency and we can achieve that with the following, via the cmd line.
```fish
> go get -u github.com/robjporter/go-functions
```

## Setting up the application
You need to add the UCS and UCS Performance Manager systems to the application.  Your password will be encrypted before it is stored, however usernames will remain in plain text.  This should be a read only account on both systems, so should not cause too much of a security risk.

### Add UCS Domain
Repeat this process as many times as needed.
```go
> go run main.go add ucs --ip=<IP> --username=<USERNAME> --password=<PASSWORD>
```

### Update UCS Domain
The update process will only succeed if the IP of the UCS Domain is already in the config file.
```go
> go run main.go update ucs --ip=<IP> --username=<USERNAME> --password=<PASSWORD>
```

### Delete UCS Domain
The delete process will only succeed if the IP of the UCS Domain is already in the config file.
```go
> go run main.go delete ucs --ip=<IP>
```
### Show UCS Domains
To show the current configuration details for a UCS System;
```go
> go run main.go show ucs --ip=<IP>
```

### Show All discoverable systems
To show all the currently entered system information;
```go
> go run main.go show all
```

## Running the application
Once the UCS Domains have been added, the application is now ready to run.
```go
> go run main.go run
```

## Building to a Binary
One of the great advantages of GO is the ability to compile the code and all dependencies into a single binary file.  This is enhanced by building for multiple platforms.  I have included a short script to compile to most of the common formats and place them in the ./bin folder.  You may need to add the execute ability onto the script, as this maybe removed during the download process, on a Mac, you can complete this by doing;
```fish
> chmod +x buildall.sh
```
To build the application run this;
```fish
> ./buildall.sh
```

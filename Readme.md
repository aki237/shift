# shift

Proxy Configuration manager for [proGY](https://github.com/aki237/proGY).

## Download

+ Go to the [releases](https://github.com/aki237/shift/releases) page
+ Grab the latest release for your platform
+ Extract the archive and simple run `./shift`

## Options

+ `-help`
  
  Prints the help information
  
  ```
  $ ./shift -help
  shift : proxy manager for proGY, v0.01 Linux x86_64 30-04-2017 05:01
  Usage : /usr/local/bin/shift <options>
	  -config string
    	  Path pointing to the .progy config file. (default "/Users/stinson/.progy")
	  -proxylist string
    	  Path pointing to the proxy list file. (TOML only accepted) (default "/Users/stinson/.config/proxy")
	  -version
    	  Prints the version of the application
  ```

+ `-config`
  
  Option for specifying the location of the proGY config file.
  
  ```
  $ ./shift -config /etc/progy.conf
  ...
  ```
  
  If not specified, by default shift tried opening .progy from the corresponding `$HOME` directory.

+ `-proxylist`
  
  Option for specifying the location of the proxy list file. A Toml like configuration is expected for
  the list file. The format is specified as follows :
  
  ```toml
  [remoteProxyAddress1]
  username1 = password1
  username2=password2
  [ remoteProxyAddress2 ]
  username3 = password3
  ```
  
  Example :
  
  ```toml
  [172.2.0.1:8080]
  dudeNoOne = someAwkwardPassword
  dudeNoTwo = someMoreAwkwardPassword
  divaNoThree = brassKnuckles
  [172.2.0.12:8080]
  somePoorDude1 = dontStealMyPassword
  somePoorDude2 = password
  someDudeThinksHimselfAsAwesome3 = ee2e9&Chh8d7h*
  ```

+ `-version`
  
  Prints the version and build information of the program.
  
  ```
  $ ./shift -version
  shift : proxy manager for proGY, v0.01 Linux x86_64 01-01-1 00:00
  ```

## systemd service
A systemd service file has been added in the repo. Place that in your systemd system service folder.
Make sure that you have changed the User to a valid user and the binary location in `Exec` line in the
service file before copying.

Enable the service :

```
$ sudo systemctl enable shift.service
```

Start the service :

```
$ sudo systemctl start shift.service
```

## Build

Simple `go install` will do the trick. But for all the build information to be included
special linker flags have to be passed to the go compiler.

```
$ cd $GOPATH/src/github.com/aki237/shift/cmd/shift/
$ CGO_ENABLED=0 go build -ldflags "-w -s -X 'main.appBuildDate=$(date '+%d-%m-%Y %H:%M')' -X 'main.appVersion=$(cat ../../VERSION)' -X 'main.appBuildPlatform=$(echo $(uname) $(uname --machine))'"
```

Also `-w` and `-s` options will strip any debug symbols, thus reducing the executable size.

Any Issues or PRs for bugs or improvements are welcome.

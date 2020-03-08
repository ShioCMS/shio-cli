![shio_banner.png](https://shiocms.github.io/shio/img/shio_banner.png) 
------

**Viglet Shio Client** - Shio CMS command line client to connect to remote Shio CMS and execute actions.

**If you'd like to contribute to Viglet Shio Client, be sure to review the [contribution
guidelines](CONTRIBUTING.md).**

**We use [GitHub issues](https://github.com/ShioCMS/shio-client/issues) for tracking requests and bugs.**

# Installation

## Download

```shell
$ git clone https://github.com/ShioCMS/shio-client.git
$ cd shio-client
```
## Deploy

### Build
Use Go lang to build Shio CMS Client.

```shell
$ go build shio-cli.go
```

## Usage

### Connection

**shio-cli.ini** file is present in same directory of shio-client executable.
Edit this file and add login, password and server URL of your Shio CMS.

```ini
login="admin"
password="admin"
server="http://localhost:2710"
```

### Command line
`./shio-cli new name_of_site`: Create a new site

For example
```shell
$ ./shio-cli new SampleSite
```
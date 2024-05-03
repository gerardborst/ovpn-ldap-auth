# VPN authentication app

## Development

### Install dependencies

```sh
$ go get -u ./...
```

### Build

```sh
$ ./build.sh
```

### Release

Create a production release by adding a tag with a [Semantic Version](https://semver.org/) and run:

```
$ ./release.sh
```

This will result in 2 tar files in the `target` directory.

```
ldap-auth-beta.darwin-<version>.tar.gz
ldap-auth-beta.linux-<version>.tar.gz
```

The released linux tar file should be placed in the `openvpn-server-config` Git project in the `binaries` directory. Also change
the `ldap_auth_version` variable in `group_vars/all.yaml`.

Without a tag on the current commit the version will be beta.

## Usage

### Configuration

Create a file `/etc/openvpn/auth/ldap-auth.conf` and copy-paste the following

```
LDAP_BASE="OU=Accounts,dc=eu,dc=example,dc=com"
LDAP_BINDDN="CN=OpenVPN,DC=example,DC=com"
LDAP_BINDPASSWORD="SECRET"
LDAP_GROUPFILTER="(memberUid=%s)"
LDAP_HOST="ldaps.example.com"
LDAP_PORT=636
LDAP_USERFILTER="(CN=%s)"
LDAP_SERVERNAME="ldaps.example.com"
```

#### Optional parameters

```
LDAP_USESSL="true"
LDAP_SKIPTLS="false"
LDAP_LOGLEVEL="info"
LDAP_LOGTOFILE="false"
LDAP_LOGFILE="/var/log/openvpn/auth/ldap.log"
LDAP_CHECKCN="true"
LDAP_CHECKCNFAIL="true"
```

If you don't want to use tls set `LDAP_USESSL` to false and `LDAP_SKIPTLS` to true. And change the port to the non-tls port.

Logging to file can be enabled by setting `LDAP_LOGTOFILE` to true. Optionally the log file path can be changed
with a valid path in `LDAP_LOGFILE`. **This file should be writable by de user to which the openvpn server drops privileges.**
Set `LDAP_CHECKCN` to true to check if the `CN` is equal to the `username`.

If `LDAP_CHECKCNFAIL` is set to true and the `CN` is not equal to the `username` the authentication is rejected.
With set to false only a message is logged.

### Run

The `ldap-auth` process extracts the `username`, `password` and `common_name` from the environment.

```sh
$ export username="user@example.com"
$ export password="SECRET123"
$ export common_name="user@example.com"
```

Now start the authentication flow

```
$ go run cmd/ovpn-ldap-auth

```


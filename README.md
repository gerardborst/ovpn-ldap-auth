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

Create a file `/etc/openvpn/auth/ovpn-auth-config.yaml` and copy-paste the following

```yaml
ldapClient:
  base: "dc=example,dc=org"
  groupFilter: "(memberOf=%s)"
  host: "ldaps.example.org"
  port: 636
  useSsl: true
  useStartTls: true
  vpnGroupFilter: "(&(uid=%s)(memberOf=cn=openvpn,ou=users,dc=example,dc=org))"
  serverName: "ldaps.example.org"
  binddn: "cn=admin,dc=example,dc=org"
  bindPassword: "123456"
log:  
  level: "info"
  logToFile: true
  logFile: "/var/log/openvpn/auth/ldap.log"
cn:
  check: true
  fail: true
```

If you don't want to use tls set `ldapClient.useSsl` to false and if the ldap server doesn´t support it `ldapClient.useStartTls` also to false. And change the port to the non-tls port.

Logging to file can be enabled by setting `log.logToFile` to true. Optionally the log file path can be changed
with a valid path in `log.logFile`. **This file should be writable by de user to which the openvpn server drops privileges.**
Set `cn.check` to true to check if the `CN` is equal to the `username`.

If `cn.fail` is set to true and the `CN` is not equal to the `username` the authentication is rejected.
With set to false only a message is logged.

### Run

The `ovpn-ldap-auth` process extracts the `username`, `password` and `common_name` from the environment.

```sh
export username="user01"
export password="password1"
export common_name="user01"
export auth_control_file="./.vscode/auth_control_file.txt"
```

Now start the authentication flow

```
docker-compose -f ./tests/openldap/docker-compose.yaml up -d
go run ./cmd/ovpn-ldap-auth

```


# gobsq

**gobsq** is a collection of tools to ease and automate the querying of the open build service. 

## Tools

 - **lsrepos** - list update repos of currently open release requests
 - **lsrr** - list open/unassigned release request
 - **lsmu** - list open/unassigned maintenance updates
 
## Examples
```
# get open requests update repositories for the qam-test group
$ lsrepos -group qam-test
  2021-07-29T08:33:47 http://download.suse.de/ibs/SUSE:/Maintenance:/10398/SUSE_SLE-15-SP1_Update/SUSE:Maintenance:10398.repo
  2021-07-30T10:22:02 http://download.suse.de/ibs/SUSE:/Maintenance:/20287/SUSE_SLE-15_Update/SUSE:Maintenance:20287.repo
  2021-07-30T11:30:06 http://download.suse.de/ibs/SUSE:/Maintenance:/30176/SUSE_SLE-12-SP2_Update/SUSE:Maintenance:30176.repo

# get open release requests with a short summary
# requests marked with an `!` have higher priority
$ lsrr -group qam-test
  https://maintenance.suse.de/request/167250 (Recommended update for farbtest)
  https://maintenance.suse.de/request/258145 (Recommended update for yast2)
! https://maintenance.suse.de/request/349031 (Security update for bluez)

```

## Installation
```
# for lsrr
$ go get -u github.com/fgerling/gobsq/cmd/lsrr
# 
# for lsrepos
$ go get -u github.com/fgerling/gobsq/cmd/lsrepos
# for lsmu
$ go get -u github.com/fgerling/gobsq/cmd/lsmu
```

## Configuration

The tools are using ~/.gobs.toml as config file.
```
cat ~/.gobs.toml
Username = "username"
Password = "password"
Group = "qam-test"
Server = "https://api.suse.de/"
```

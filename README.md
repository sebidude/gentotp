# gentotp
Commandline tool to generate TOTPs 

Install
-------

Install go

```
export GOPATH=<path_to_go_dev>
go get github.com/r3ek0/gentotp...
go build github.com/r3ek0/gentotp

go install github.com/r3ek0/gentotp
```

Run
---
```
export TOTP_SECRET=<your-totp-secret>
gentotp
```


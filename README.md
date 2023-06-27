# zig-version-downloader

### <sub>This is a small program which downloads the latest (master) version of the Zig programming language as an archive from here:</sub>
<sub>[Zig Download Index JSON](https://ziglang.org/download/index.json)</sub>
------------------------------------------------------------------


### How to use it:
- You would need to have Go installed.
- Clone the repo.
- Follow the below steps to run the program and the tests.

*P.S. Since I'm not a millionaire and don't own a Mac, I haven't been able to test whether the current implementation works for iOS users, sorry.*.

##### 1. To download the latest archive:
```
go run main.go
```

##### 2. To run all tests in the subpackages from the current package:
```
go test ./...

Because tests are being cached and most probably you would like to see that the unzip functionality works, you can use a flag to uncache them:

go test -count=1 ./...
```

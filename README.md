## go-verifier

> go-verifier is small utility that allows you to calculate the MD5 and SHA1 hashes of one or more files in your system.

```plain
Usage of go-verifier:
        go-verifier [flags]                Runs on all files in current directory
        go-verifier [flags] [file|dir]     Compute for given files or directories
Flags:
  -hash string
        Specify hashtype, values: md5, sha1, sha256 (default "md5")
  -help
        Show this help information
  -nopath
        Without full path name
  -upper
        Get hash in uppercase
  -verify string
        Read and verify the checksum file
```

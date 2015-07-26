# gost
Gist client in go

# Installation
```
go get github.com/zbindenren/gost
```

# Usage

## List gists
```
# gost ls
2f310e76d950a71706d0 file3.txt, file1.txt, file2.txt - new description
0fa48a250b4a6292ed52               cli.go - codegangsta cli example
32d777d529a50e5950e5               zip.go - append zip archive to go binary
057d2c763ca61afb36a5             log15.go - log15 exit on critical
```
or
```
# gost ls 2f310e76d950a71706d0
           file3.txt          6 https://gist.githubusercontent.com/zbindenren/2f310e76d950a71706d0/raw/7c8ac2f8d82a1eb5f6aaece6629ff11015f91eb4/file3.txt
           file1.txt          6 https://gist.githubusercontent.com/zbindenren/2f310e76d950a71706d0/raw/e2129701f1a4d54dc44f03c93bca0a2aec7c5449/file1.txt
           file2.txt         12 https://gist.githubusercontent.com/zbindenren/2f310e76d950a71706d0/raw/340890024e71054982dcda2036e07fb3a020eb4c/file2.txt
```

## Create new gist
```
# gost -d "my description for new gist" main.go
```

or multiple files
```
# gost -d "other gist" file1.txt file2.txt file3.txt
```

## View Gist
```
#gost cat e89fddd4bbe0b405960d
main.go:
package main

import (
        "fmt"
)

func main() {
        fmt.Println("Hello Gist")
}
```
or
```
# gost cat 2f310e76d950a71706d0 -f file1.txt
file1.txt:
file1
```

or view it in the browser

```
# gost cat -b e89fddd4bbe0b405960d
```

## Save Gist to local disk
```
# gost get e89fddd4bbe0b405960d
```
or just one file
```
# gost get e89fddd4bbe0b405960d -f filename.txt
```

## Delete/Remove Gist
```
# gost rm e89fddd4bbe0b405960d
```
or just one file
```
# gost rm 2f310e76d950a71706d0 -f file3.txt
```

## Update Gist
```
# gost update 2f310e76d950a71706d0 -f file1.txt
```
or
```
# gost update 2f310e76d950a71706d0 -f file1.txt -f file2.txt
```
or
```
# gost -d "new description" update 2f310e76d950a71706d0
```

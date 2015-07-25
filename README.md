# gost
Gist client in go

# Installation
```shell
go get github.com/zbindenren/gost
```

# Usage

## Add Gist
```
gost -a -d "my cool gist" main.go
```
or
```
gost -a -d "other gist" file1.txt file2.txt file3.txt
```
## List Gists
```
gost -l
8aeaae66d4b6bde78371              main.go - my cool gist
9e530353f3fe19026e84         CHANGELOG.md -
32e777d529a50e5950e5               zip.go - zip example
058d2c763ca61afb36a5             log15.go - log15 exit on critical
```
## View Gist
```
gost -v 8aeaae66d4b6bde78371
main.go:
package main

import (
        "fmt"
)

func main() {
        fmt.Println("Hello Gist")
}
```

## Save Gist to local disk
```
gost -s 8aeaae66d4b6bde78371
```

## Open Gist in browser
```
gost -b 8aeaae66d4b6bde78371
```

## Delete/Remove Gist
```
gost -rm 8aeaae66d4b6bde78371
```

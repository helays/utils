package main

import (
	"fmt"
	"helay.net/go/utils/v3/close/vclose"
	"os"
)

func main() {
	file, err := os.Open("clsoe.go")
	defer vclose.Close(file)
	fmt.Println(err, "文件")
}

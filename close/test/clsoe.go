package main

import (
	"fmt"
	"github.com/helays/utils/v2/close/vclose"
	"os"
)

func main() {
	file, err := os.Open("clsoe.go")
	defer vclose.Close(file)
	fmt.Println(err, "文件")
}

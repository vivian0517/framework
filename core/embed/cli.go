package embed

import (
	_ "embed"
	"fmt"
	"os"
	"path"
)

//go:embed first_test.txt
var first []byte

//go:embed second_test.txt
var second []byte

func Book() {
	fmt.Println("示例程序将会在当前目录创建case文件夹,并包含两份测试代码示例,如果你当前存在case文件夹可能会报错,忽略即可")
	err := os.Mkdir("case", 0755)
	if err != nil {
		fmt.Println("创建case目录失败:", err)
	}

	err = os.MkdirAll(path.Join("case", "one"), 0755)
	if err != nil {
		fmt.Println("创建one目录失败:", err)
	}
	err = os.MkdirAll(path.Join("case", "two"), 0755)
	if err != nil {
		fmt.Println("创建two目录失败:", err)
	}

	firstpath := path.Join("case", "one", "first_test.go")
	file, _ := os.Create(firstpath)
	file.Write(first)

	secondparh := path.Join("case", "two", "second_test.go")
	file2, _ := os.Create(secondparh)
	file2.Write(second)
}

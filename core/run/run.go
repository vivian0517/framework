package run

import (
	"bufio"
	"fmt"
	logf "framework/core/log"
	"framework/core/manager"
	fileutil "framework/util/file"
	"framework/util/ini"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var importPath []string
var modulePath = "framework"

func Start() {

	fmt.Printf("%+v\n", manager.GetAllCase())

	cmd := runAllTestCases("testFileName")

	if err := cmd.Run(); err != nil {
		fmt.Println("测试失败:", err)
	} else {
		fmt.Println("测试完成，结果存入 test_output.log")
	}
}

// 执行用例
func RunCase() {
	logs, err := logf.RedirectOutputToFile("test.log")
	if err != nil {
		fmt.Println("无法创建日志文件:", err)
		return
	}
	defer logs.Close()

	var packageMap = make(map[string][]testing.InternalTest)

	// 按包名归类测试函数
	for k, v := range manager.Registry {
		pkg := v.PkgName
		test := testing.InternalTest{
			Name: k,
			F:    v.Callback,
		}
		packageMap[pkg] = append(packageMap[pkg], test)
	}

	// 分组运行
	for pkgName, tests := range packageMap {
		fmt.Printf("+++ 运行包 [%s] 的测试，共 %d 个\n", pkgName, len(tests))
		manager.TestManager(tests)
	}

}

func ReadFileWriteFile() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("当前路径", path)

	folderName := filepath.Join(path, "case")
	// 尝试打开目录
	if dir, err := os.Open(folderName); err == nil {
		// 目录存在
		fmt.Println("文件夹存在")
		defer dir.Close()
		cache := filepath.Join(path, "cache")
		err := os.Mkdir(cache, 0755)
		if err != nil {
			fmt.Println("创建cache目录失败:", err)
		}
		casefilename, err := dir.Readdirnames(-1)
		if err != nil {
			fmt.Println("读取目录内容时出错:", err)
		}

		for _, filename := range casefilename {
			//fmt.Println("模块文件名", filename)
			importmodel := fmt.Sprintf("_ \"%s/cache/%s\"", modulePath, filename)
			importPath = append(importPath, importmodel)
			model := filepath.Join(cache, filename)
			err := os.Mkdir(model, 0755)
			if err != nil {
				fmt.Println("创建cache下子目录失败:", err)
			}
			// 读取目录中的_test.go
			modulepath := filepath.Join(folderName, filename)
			files, err := os.ReadDir(modulepath)
			if err != nil {
				fmt.Println("Error reading directory:", err)
			}

			// 遍历文件列表并检查文件名
			for _, file := range files {
				if !file.IsDir() && strings.HasSuffix(file.Name(), "_test.go") {
					fmt.Println(file.Name()) // 打印文件名

					// 打开文件
					testpath := filepath.Join(modulepath, file.Name())
					dst := filepath.Join(model, strings.TrimSuffix(file.Name(), ".go")+"_."+"go")
					err := fileutil.CopyFile(testpath, dst)
					if err != nil {
						log.Fatal(err)
					}
					file, err := os.Open(testpath)
					if err != nil {
						log.Fatal(err)
					}
					defer file.Close()

					// 创建正则表达式，用于匹配 manager.RegisterCase("描述", "p0", Test_ceshi05)
					// 正则表达式匹配格式：manager.RegisterCase("描述", "p0", Test_ceshi05)
					// 捕获内容：描述(p0)、Test_ceshi05（不带引号的部分）
					pattern := `manager\.RegisterCase\("([^"]+)",\s*"([^"]+)",\s*([A-Za-z0-9_]+)\)`
					re := regexp.MustCompile(pattern)

					// 使用scanner逐行读取文件内容
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						line := scanner.Text()
						// 如果行包含manager.RegisterCase
						if re.MatchString(line) {
							// 提取匹配的内容
							matches := re.FindStringSubmatch(line)
							if len(matches) > 3 {
								//description := matches[1]
								//p0 := matches[2]
								//testCase := matches[3]
								//fmt.Printf("描述: %s, 级别: %s, 函数名: %s\n", description, p0, testCase)
							}
						}
					}

					// 检查扫描时是否出现错误
					if err := scanner.Err(); err != nil {
						log.Fatal(err)
					}

				}
			}
			gen_imports()
		}

	} else if os.IsNotExist(err) {
		// 目录不存在
		fmt.Println("测试目录不存在")
	} else {
		// 出现其他错误
		fmt.Println("读取测试目录文件夹时出错:", err)
	}
}

func choiceRuntype() {
	if ini.Iniconfig["runtype"] == "all" {

	}
}

func gen_imports() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Print(err)
	}
	cache := filepath.Join(path, "gen", "main")
	err = os.MkdirAll(cache, 0755)
	if err != nil {
		fmt.Println("创建目录失败:", err)
	}
	filename := filepath.Join(cache, "main.go")
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintf(f, "package %s\n\n", "main")
	fmt.Fprintf(f, "import (\n")
	for _, imp := range importPath {
		fmt.Fprintf(f, "    %s\n", imp)
	}
	fmt.Fprintf(f, "	\"log\"\n")
	//添加	"framework/core/run"
	runpath := filepath.Base(path)
	fmt.Fprintf(f, "	\"%s/core/run\"\n", runpath)
	fmt.Fprintf(f, ")\n")

	fmt.Fprintf(f, "func main() {\n")
	fmt.Fprintf(f, "	log.Println(\"code generate\")\n")
	fmt.Fprintf(f, "	run.RunCase()\n")
	fmt.Fprintf(f, "}\n")

}

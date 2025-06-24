package run

import (
	"os/exec"
)

// 运行单一测试用例
func runSingleTestCase(testCaseName string) *exec.Cmd {
	//go test -run "^(Test_ceshi05|Test_ceshi02|Test_xxxx)$"
	cmd := exec.Command("go", "test", "-timeout", "30s", "-run", testCaseName)
	return cmd
}

// 运行文件中的测试用例
func runFileCases(testFileName string) *exec.Cmd {
	
	cmd := exec.Command("go", "test", "-timeout", "30s", testFileName)
	return cmd
}

// 运行所有包中的测试用例
func runAllTestCases(testFileName string) *exec.Cmd {
	cmd := exec.Command("go", "test", "-timeout", "30s", testFileName)
	return cmd
}

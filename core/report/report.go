package report

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

//go:embed templete.html
var temple string

type TestCase struct {
	Name     string
	Status   string
	Output   string
	Duration float64
}

type PackageResult struct {
	Name      string
	TestCases []TestCase
	Passed    int
	Failed    int
	Duration  float64 // 单位秒
}

type ReportData struct {
	Packages []PackageResult
	Total    int
	Passed   int
	Failed   int
	Duration float64
}

func Report() {
	logFile := "test.log"
	outputFile := "report.html"

	pkgs, err := parseLogFile(logFile)
	if err != nil {
		fmt.Println("解析日志失败:", err)
		return
	}

	err = generateHTMLReport(pkgs, outputFile)
	if err != nil {
		fmt.Println("生成 HTML 失败:", err)
		return
	}

	fmt.Println("测试报告已生成:", outputFile)
}
func parseLogFile(filename string) (*ReportData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	packageStartRe := regexp.MustCompile(`^\+\+\+ 运行包 \[(.+)] 的测试，共 (\d+) 个`)
	runRe := regexp.MustCompile(`^=== RUN\s+(.+)$`)
	resultRe := regexp.MustCompile(`^--- (PASS|FAIL): (.+) \(([\d.]+)s\)$`)

	var (
		report        ReportData
		currentPkg    PackageResult
		currentTest   TestCase
		testOutputBuf []string
		inTestOutput  bool
	)

	for scanner.Scan() {
		line := scanner.Text()

		// 新包开始
		if match := packageStartRe.FindStringSubmatch(line); len(match) == 3 {
			if currentPkg.Name != "" {
				report.Packages = append(report.Packages, currentPkg)
				report.Passed += currentPkg.Passed
				report.Failed += currentPkg.Failed
				report.Total += len(currentPkg.TestCases)
				report.Duration += currentPkg.Duration
			}
			currentPkg = PackageResult{
				Name: match[1],
			}
			continue
		}

		// 开始测试
		if match := runRe.FindStringSubmatch(line); len(match) == 2 {
			if inTestOutput {
				currentTest.Output = strings.Join(testOutputBuf, "\n")
				currentPkg.TestCases = append(currentPkg.TestCases, currentTest)
				testOutputBuf = nil
				currentTest = TestCase{}
			}
			currentTest.Name = match[1]
			inTestOutput = true
			continue
		}

		// 测试结果 + 时间
		if match := resultRe.FindStringSubmatch(line); len(match) == 4 {
			status := match[1]
			testName := match[2]
			duration, _ := strconv.ParseFloat(match[3], 64)

			currentTest.Status = status
			currentTest.Name = testName
			currentTest.Duration = duration
			currentTest.Output = strings.Join(testOutputBuf, "\n")

			currentPkg.TestCases = append(currentPkg.TestCases, currentTest)
			currentPkg.Duration += duration
			if status == "PASS" {
				currentPkg.Passed++
			} else {
				currentPkg.Failed++
			}

			testOutputBuf = nil
			currentTest = TestCase{}
			inTestOutput = false
			continue
		}

		// 收集输出
		if inTestOutput {
			testOutputBuf = append(testOutputBuf, line)
		}
	}

	// 最后一个包
	if currentPkg.Name != "" {
		report.Packages = append(report.Packages, currentPkg)
		report.Passed += currentPkg.Passed
		report.Failed += currentPkg.Failed
		report.Total += len(currentPkg.TestCases)
		report.Duration += currentPkg.Duration
	}

	return &report, nil
}

func generateHTMLReport(data *ReportData, output string) error {

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	// 编译并渲染模板
	t := template.Must(template.New(output).Parse(temple))
	return t.Execute(f, data)
}

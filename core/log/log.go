package log

import (
	"os"
)

func RedirectOutputToFile(logfile string) (*os.File, error) {
	logf, err := os.Create(logfile)
	if err != nil {
		return nil, err
	}
	// 重定向 stdout 和 stderr 到日志文件
	os.Stdout = logf
	os.Stderr = logf
	return logf, nil
}

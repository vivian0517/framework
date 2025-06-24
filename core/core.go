package core

import (
	"os"
	"path/filepath"

	"framework/core/report"
	"framework/util/file"
	"framework/util/ini"
)

func Core() {

	ini.CatchINI()

	path, _ := os.Getwd()
	run := filepath.Join(path, "core", "run", "main", "main.go")

	file.RunGoFile(run)

	run2 := filepath.Join(path, "gen", "main", "main.go")

	file.RunGoFile(run2)
	os.RemoveAll(path + "\\cache")
	report.Report()

}

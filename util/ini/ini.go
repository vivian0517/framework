package ini

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

var Iniconfig = make(map[string]string)

// 读取INI文件并返回配置值
func ReadINIFile(filePath string) (*ini.File, error) {
	// 加载INI文件
	cfg, err := ini.Load(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load ini file: %v", err)
	}
	return cfg, nil
}

// 读取指定路径的INI文件
func CatchINI() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("当前路径", path)
	ini := filepath.Join(path, "env", "basic.ini")
	cfg, err := ReadINIFile(ini)
	if err != nil {
		log.Fatalf("Error reading INI file: %v", err)
	}

	// 读取数据库相关配置
	runtype := cfg.Section("basic").Key("runtype").String()
	timeout := cfg.Section("basic").Key("timeout").String()
	env := cfg.Section("basic").Key("env").String()
	fmt.Println("runtype", runtype)
	fmt.Println("timeout", timeout)
	fmt.Println("env", env)
	Iniconfig["runtype"] = runtype
	Iniconfig["timeout"] = timeout
	Iniconfig["env"] = env
	GetConfig()
}

// 获取设置的环境配置
func GetConfig() {

	path, _ := os.Getwd()
	filename := filepath.Join(path, "env", Iniconfig["env"], "config.ini")
	cfg, err := ini.Load(filename)
	if err != nil {
		log.Fatal("读取配置文件失败:", err)
	}

	for _, section := range cfg.Sections() {
		sectionName := section.Name()
		if sectionName == ini.DefaultSection {
			sectionName = ""
		}
		for _, key := range section.Keys() {
			fullKey := key.Name()
			if sectionName != "" {
				fullKey = sectionName + "." + key.Name()
			}
			Iniconfig[fullKey] = key.Value()
		}
	}
	for k, v := range Iniconfig {
		fmt.Printf("%s = %s\n", k, v)
	}
}

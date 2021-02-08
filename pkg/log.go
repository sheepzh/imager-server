package pkg

import (
	"fmt"
	"log"
	"os"
	"os/user"
)

var level string

var logPath string

var allLevels = [...]string{"debug", "info", "warn", "error"}

func InitLog(appName string) {
	dir := GetArg("-lf", "--log-file-dir")
	if dir == "" {
		user, err := user.Current()
		if nil == err {
			dir = user.HomeDir + "/imager-log"
		} else {
			dir = "./logs"
		}
	}
	res2, err := CreateDir(dir) //创建文件夹
	if res2 == false {
		panic(err)
	}

	logPath = dir + "/" + appName

	levelArg := GetArg("-ll", "--log-level")

	if levelArg == "" {
		levelArg = "info"
	} else if !validLevel(levelArg) {
		log0("[ERROR]", "invalid log level: "+levelArg+". 'info' will be used.", "error")
		levelArg = "info"
	}

	level = levelArg

	Logi("initialized log, level=" + level)
}

func log0(prefix string, c string, level string) {
	fmt.Println(prefix + c)
	var filePath string
	if level == "error" {
		filePath = logPath + "-" + level + ".log"

	} else {
		filePath = logPath + ".log"
	}
	file, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()

	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime)
	// prefix
	log.SetPrefix(prefix)
	log.Println(c)
}

func validLevel(level string) bool {
	for _, allowed := range allLevels {
		if allowed == level {
			return true
		}
	}
	return false
}

func Logd(c string) {
	if level == "debug" {
		log0("[DEBUG] ", c, "debug")
	}
}

func Logi(c string) {
	if level == "debug" || level == "info" {
		log0("[INFO ] ", c, "info")
	}
}

func Logw(c string) {
	if level == "debug" || level == "info" || level == "warn" {
		log0("[WARN ] ", c, "warn")
	}
}

func Loge(c string) {
	log0("[ERROR] ", c, "error")
}

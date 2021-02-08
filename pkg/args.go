package pkg

import (
	"os"
	"strings"
)

// parse the argments
// @param shortName      the short name of argument
//                       e.g.
//                          -s=/home/ubuntu, input -s, output /home/ubuntu
// @param fullName       the full name of argument
//                       e.g.
//                          --seged-file-dir=/home/ubuntu, input --seged-file-dir, output /home/ubuntu
func getArg(shortName string, fullName string, required bool) string {
	for _, v := range os.Args {
		if strings.HasPrefix(v, shortName+"=") {
			return v[len(shortName)+1:]
		} else if strings.HasPrefix(v, fullName+"=") {
			return v[len(fullName)+1:]
		}
	}
	if required {
		Loge(fullName + " or " + shortName + " is required")
	}
	return ""
}

func GetArg(shortName string, fullName string) string {
	return getArg(shortName, fullName, false)
}

func GetArgRequired(shortName string, fullName string) string {
	return getArg(shortName, fullName, true)
}

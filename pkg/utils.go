package pkg

import (
	"net/http"
	"os"
	"strconv"
)

func CreateDir(dir string) (bool, error) {
	_, err := os.Stat(dir)

	if err == nil {
		//directory exists
		return true, nil
	}

	err2 := os.MkdirAll(dir, 0755)
	if err2 == nil {
		return true, err2
	}
	return false, err2
}

func GetPageInfo(req *http.Request) (int, int) {
	query := req.URL.Query()
	page := query.Get("pn")
	if page == "" {
		page = "10"
	}
	size := query.Get("size")
	if size == "" {
		size = "10"
	}

	pn, e1 := strconv.Atoi(page)
	ps, e2 := strconv.Atoi(size)

	if e1 != nil || e2 != nil {
		return 1, 0
	}

	return pn, ps
}

// get the inclusive startIndex and exclusive endIndex for page query
// @param pageNum   the number of page, starts with 1
// @param pageSize  the size of page
// @param total     the total of
// @return si int, ei int
func GetStartAndEnd(pageNum int, pageSize int, total int) (int, int) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 1
	}
	startIndex := (pageNum - 1) * pageSize
	if startIndex >= total {
		return total - 1, total - 1
	}
	endIndex := startIndex + pageSize
	if endIndex > total {
		endIndex = total
	}
	return startIndex, endIndex
}

func Int2Str(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

// get integer param of request, need defaultValue
// @param req            request
// @param paramName      name of parameter
// @param defaultValue   default value
func GetIntParam(req *http.Request, paramName string, defaultValue int) int {
	param := req.FormValue(paramName)
	if param == "" {
		return defaultValue
	} else {
		i, err := strconv.Atoi(param)
		if err != nil {
			Logd("failed to parse integer parameter: " + err.Error())
			return defaultValue
		} else {
			return i
		}
	}
}

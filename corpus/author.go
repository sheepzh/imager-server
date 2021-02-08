// author service
package corpus

import (
	"io/ioutil"
	"net/http"
	"pkg"
	"strings"
	"sync"
)

type Author struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var __onceAuthor sync.Once

var __allAuthors []Author

/*  */
type __AuthorPage struct {
	List  []Author `json:"list"`
	Total int      `json:"total"`
	Page  int      `json:"page"`
	Size  int      `json:"size"`
}

// list authors by page
func ListAuthors(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.WriteHeader(404)
		return
	}
	page, size := pkg.GetPageInfo(req)
	response := __ListAuthors(page, size)

	pkg.WriteJsonOfResponce(w, response)
}

func __ListAuthors(pageNum int, pageSize int) __AuthorPage {
	__onceAuthor.Do(__ScanAllAuthors)

	if __allAuthors == nil {
		return __AuthorPage{List: []Author{}, Total: 0, Page: pageNum, Size: 0}
	}

	total := len(__allAuthors)

	startIndex, endIndex := pkg.GetStartAndEnd(pageNum, pageSize, total)

	var list []Author

	pkg.Logi("query aothor:" +
		" start=" + pkg.Int2Str(startIndex) +
		" end=" + pkg.Int2Str(endIndex) +
		" count=" + pkg.Int2Str(total))
	list = __allAuthors[startIndex:endIndex]
	return __AuthorPage{List: list, Total: total, Page: pageNum, Size: len(list)}
}

// parse the author info with the directory of seged files
func __ParseAuthor(dirName string) Author {
	strArr := strings.Split(dirName, "_")
	return Author{ID: strArr[1], Name: strArr[0]}
}

func __ScanAllAuthors() {
	segedPoemDir := pkg.GetArgRequired("-s", "--seged-poem-dir")
	dir, err := ioutil.ReadDir(segedPoemDir)
	if err != nil {
		__allAuthors = nil
		pkg.Logw("directory of seged poem is not exist: " + err.Error())
		pkg.Logw(segedPoemDir)
	} else {
		var authors []Author
		for _, pd := range dir {
			if pd.IsDir() {
				author := __ParseAuthor(pd.Name())
				authors = append(authors, author)
			}
		}
		__allAuthors = authors
	}
}

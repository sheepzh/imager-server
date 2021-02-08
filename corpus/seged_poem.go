// poem service
package corpus

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"pkg"
	"strconv"
	"strings"
	"sync"
)

type SegedPoem struct {
	Name        string `json:"name"`
	Author      Author `json:"author"`
	WritenYear  int    `json:"writenYear"`
	WritenMonth int    `json:"writenMonth"`
	WritenDay   int    `json:"writenDay"`
	WritenDate  string `json:"writeDate"`
	Content     string `json:"content"`
	path        string
}

var __allPoems []SegedPoem

var __onceSegedPoem sync.Once

// GET /poems/seged
func ListSegedPoems(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.WriteHeader(404)
		return
	}

	page, size := pkg.GetPageInfo(req)
	authorKey := req.FormValue("ak")
	authorName := req.FormValue("an")
	year := pkg.GetIntParam(req, "year", 0)
	response := __ListPoems(page, size, authorKey, authorName, year)

	pkg.WriteJsonOfResponce(w, response)
}

type __SegedPoemPage struct {
	List  []SegedPoem `json:"list"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// query poems
// @param page         page page num
// @param size         page size
// @param authorKey    the key of author, full-match if not blank
// @param authorName   the name of author, fuzzy-match if not blank
// @param year         the writen year, if not zero
func __ListPoems(page int, size int, authorKey string, authorName string, year int) __SegedPoemPage {
	__onceSegedPoem.Do(__ScanAllPoems)
	pkg.Logd("count of poems: " + pkg.Int2Str(len(__allPoems)))
	list := []SegedPoem{}
	pkg.Logi("query seged poems: page=" + pkg.Int2Str(page) +
		" size=" + pkg.Int2Str(size) +
		" authorKey=" + authorKey +
		" authorName=" + authorName +
		" year=" + pkg.Int2Str(year))
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	index := 0
	for _, p := range __allPoems {
		ak := authorKey == "" || authorKey == p.Author.ID
		an := authorName == "" || strings.Contains(p.Author.Name, authorName)
		y := year == 0 || year == p.WritenYear
		if ak && an && y {
			if index >= startIndex && index < endIndex {
				__ReadPoemContent(&p)
				list = append(list, p)
			}
			index++
		}
	}
	return __SegedPoemPage{List: list, Total: index, Page: page, Size: len(list)}
}

func __ReadPoemContent(p *SegedPoem) {
	exist := (*p).Content
	if exist == "" {
		path := (*p).path

		fi, err := os.Open(path)
		if err != nil {
			pkg.Loge("Failed to open seged poem file: " + (*p).path)
			return
		}
		defer fi.Close()
		br := bufio.NewReader(fi)
		index := 0

		lines := []string{}
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			if index >= 3 {
				lines = append(lines, string(a))
			}
			index = index + 1
		}
		(*p).Content = strings.Join(lines, "\n")
	}
}

func __ScanAllPoems() {
	segedPoemDir := pkg.GetArgRequired("-s", "--seged-poem-dir")
	dir, err := ioutil.ReadDir(segedPoemDir)
	if err != nil {
		__allPoems = nil
		pkg.Logw("directory of seged poem is not exist: " + err.Error())
		pkg.Logw(segedPoemDir)
	} else {
		allPoems := []SegedPoem{}
		for _, pd := range dir {
			if pd.IsDir() {
				pdn := pd.Name()
				// @see author.go
				author := __ParseAuthor(pdn)

				authorPath := segedPoemDir + string(os.PathSeparator) + pdn
				authorDir, _ := ioutil.ReadDir(authorPath)
				for _, pf := range authorDir {
					if pf.IsDir() {
						continue
					}
					pfn := pf.Name()
					if strings.HasSuffix(pfn, "pt") {
						// is poem file
						poem, err := __ParsePoem(authorPath, pfn, author)
						if err != nil {
							pkg.Loge("read seged poem error: " + err.Error())
						}
						allPoems = append(allPoems, poem)
					}
				}
			}
		}
		__allPoems = allPoems
	}
}

func __ParsePoem(authorPath string, pfn string, author Author) (SegedPoem, error) {
	path := authorPath + string(os.PathSeparator) + pfn
	fi, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return SegedPoem{}, err
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	index := 0
	title := ""
	year := 0
	month := 0
	day := 0
	date := ""
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if index == 0 {
			// title line
			// title:
			title = string(a[6:])
		} else if index == 1 {
			// date line
			// date:
			date = string(a[5:])
			date = strings.Trim(date, " ")
			l := len(date)
			if l >= 4 {
				year, _ = strconv.Atoi(date[0:4])
			}
			if l >= 6 {
				month, _ = strconv.Atoi(date[4:6])
			}
			if l >= 8 {
				day, _ = strconv.Atoi(date[6:8])
			}
			break
		}
		index = index + 1
	}
	return SegedPoem{Author: author, Name: title, WritenYear: year, WritenMonth: month, WritenDay: day, WritenDate: date, path: path}, nil
}

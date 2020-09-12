package nginx

import (
	"bufio"
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
	"strings"
	"time"
)

type Request struct {
	RemoteIP string
	Date   time.Time
	Method string
	URL    string
	Code   int
	length int
	Header string
}

type Access struct {
	Name string
	List []*Request
}

func catSub(line string, begin, end string) (string, int) {
	idx1 := strings.Index(line, begin)
	idx2 := strings.Index(line, end)
	if idx1 == -1 || idx2 == -1 {
		return "",0
	}
	if idx1 > idx2 {
		return "",0
	}
	return line[idx1+len(begin): idx2], idx2+len(end)
}

func catSub2(line string, flag rune) (string, int) {
	begin, end := -1, -1
	for i, v:= range line {
		if v != flag {
			continue
		}
		if begin == -1 {
			begin = i
			continue
		}
		if end == -1 {
			end = i
		}
	}
	if begin == -1 || end == -1 {
		return "",0
	}
	return line[begin+1 : end], end+1
}

func Split(line string) []string {
	return strings.Split(line, " ")
}

func MonthGet(line string) time.Month {
	var month = []string {
		"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
	}
	for idx, v:= range month {
		if v == line {
			return time.Month(idx + 1)
		}
	}
	return 0
}

func TimeGet(line string) (time.Time , error) {
	var day, year, hour, minute, second int
	var month string

	month, _ = catSub2(line, '/')
	if month == "" {
		return time.Time{}, fmt.Errorf("parse month fail")
	}

	// "12/Sep/2020:07:33:35 +0000"
	cnt , err := fmt.Sscanf(line, "%d/" + month + "/%d:%d:%d:%d", &day, &year, &hour, &minute, &second)
	if err != nil  {
		return time.Time{}, err
	}

	if cnt != 5 {
		return time.Time{}, fmt.Errorf("time get fail, %d", cnt)
	}

	local, err := time.LoadLocation("")
	if err != nil {
		return time.Time{}, err
	}

	return time.Date( year, MonthGet(month), day, hour, minute, second, 0, local), nil
}

func parseLine(line string) *Request {
	var req Request
	var err error

	idx := strings.Index(line, " ")
	if idx != -1 {
		req.RemoteIP = line[:idx]
	}

	body, idx2 := catSub(line, "[", "]")
	if idx2 == 0 {
		return &req
	}
	line = line[idx2:]

	req.Date, err = TimeGet(body)
	if err != nil {
		fmt.Println(err.Error())
		return &req
	}

	body, idx2 = catSub2(line, '"')
	if idx2 == 0 {
		return &req
	}
	line = line[idx2:]

	list := Split(body)
	req.Method = list[0]
	req.URL = list[1]

	cnt, _ := fmt.Sscanf(line, "%d %d", &req.Code, &req.length)
	if cnt != 2 {
		return &req
	}

	_, idx2 = catSub2(line, '"')
	if idx2 == 0 {
		return &req
	}
	line = line[idx2:]

	body, idx2 = catSub2(line, '"')
	if idx2 == 0 {
		return &req
	}
	req.Header = body

	return &req
}

func ParseAccessFile(file string) *Access {
	fileHandler, err := os.Open(file)
	if err != nil {
		logs.Error("parse access fail", err.Error())
	}
	defer fileHandler.Close()
	var access Access
	reader := bufio.NewReader(fileHandler)
	for  {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		req := parseLine(string(line))
		if req != nil {
			access.List = append(access.List, req)
		}
	}
	return &access
}
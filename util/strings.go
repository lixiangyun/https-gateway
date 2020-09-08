package util

import (
	"fmt"
	"strconv"
)

func StringList(list []string) string {
	var body string
	for idx,v := range list {
		if idx == len(list) - 1 {
			body += fmt.Sprintf("%s",v)
		}else {
			body += fmt.Sprintf("%s;",v)
		}
	}
	return body
}

func StringToInt(value string) int {
	number, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return number
}

func ByteView(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	} else if size < (1024 * 1024) {
		return fmt.Sprintf("%.1fKB", float64(size)/float64(1024))
	} else if size < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.1fMB", float64(size)/float64(1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.1fGB", float64(size)/float64(1024*1024*1024))
	} else {
		return fmt.Sprintf("%.1fTB", float64(size)/float64(1024*1024*1024*1024))
	}
}

func IntDiff(oldlist []int, newlist []int) ([]int, []int) {
	del := make([]int, 0)
	add := make([]int, 0)
	for _,v1 := range oldlist {
		flag := false
		for _,v2 := range newlist {
			if v1 == v2 {
				flag = true
				break
			}
		}
		if flag == false {
			del = append(del, v1)
		}
	}
	for _,v1 := range newlist {
		flag := false
		for _,v2 := range oldlist {
			if v1 == v2 {
				flag = true
				break
			}
		}
		if flag == false {
			add = append(add, v1)
		}
	}
	return del, add
}

func StringDiff(oldlist []string, newlist []string) ([]string, []string) {
	del := make([]string, 0)
	add := make([]string, 0)
	for _,v1 := range oldlist {
		flag := false
		for _,v2 := range newlist {
			if v1 == v2 {
				flag = true
				break
			}
		}
		if flag == false {
			del = append(del, v1)
		}
	}
	for _,v1 := range newlist {
		flag := false
		for _,v2 := range oldlist {
			if v1 == v2 {
				flag = true
				break
			}
		}
		if flag == false {
			add = append(add, v1)
		}
	}
	return del, add
}

func BoolToString(flag bool) string {
	if flag {
		return "true"
	}else {
		return "false"
	}
}
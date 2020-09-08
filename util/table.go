package util

func TablePage(page, limit int, cnt int) (int, int) {
	if page == 0 {
		return 0, cnt
	}
	if (page - 1) * limit >= cnt {
		return 0, cnt
	}
	begin := (page - 1) * limit
	if cnt - begin > limit {
		cnt = limit
	}else {
		cnt = cnt - begin
	}
	return begin, begin + cnt
}

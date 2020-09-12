package nginx

import (
	"testing"
)

var TEST_LINE string = "92.222.92.237 - - [12/Sep/2020:07:33:35 +0000] \"GET /wp-login.php HTTP/1.1\" 200 2373 \"-\" \"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:62.0) Gecko/20100101 Firefox/62.0\""

func TestParseAccess(t *testing.T){
	access := parseLine(TEST_LINE)
	if access.Header == "" {
		t.Error("parse line fail")
	}
	t.Errorf("%v", access)
}
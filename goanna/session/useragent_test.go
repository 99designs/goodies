package session

import (
	"testing"
)

type uaTest struct {
	ua              string
	expectedVersion float64
}

var ieUAs = []uaTest{
	uaTest{`Mozilla/5.0 (compatible; MSIE 10.6; Windows NT 6.1; Trident/5.0; InfoPath.2; SLCC1; .NET CLR 3.0.4506.2152; .NET CLR 3.5.30729; .NET CLR 2.0.50727) 3gpp-gba UNTRUSTED/1.0`, 10.6},
	uaTest{`Mozilla/5.0 (Windows; U; MSIE 9.0; WIndows NT 9.0; en-US))`, 9.0},
	uaTest{`Mozilla/5.0 (compatible; MSIE 8.0; Windows NT 6.1; Trident/4.0; GTB7.4; InfoPath.2; SV1; .NET CLR 3.3.69573; WOW64; en-US)`, 8.0},
	uaTest{`Mozilla/5.0 (Windows; U; MSIE 7.0; Windows NT 6.0; en-US)`, 7.0},
	uaTest{`Mozilla/5.0 (Windows; U; MSIE ; Windows NT 6.0; en-US)`, 0.0},
	uaTest{`foobar`, 0.0},
}

func TestIEDetection(t *testing.T) {

	for _, uatest := range ieUAs {
		result, _ := getInternetExplorerVersion(uatest.ua)
		if result != uatest.expectedVersion {
			t.Errorf("Expected %f, got %f", uatest.expectedVersion, result)
		}
	}
}

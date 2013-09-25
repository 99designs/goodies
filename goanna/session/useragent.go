package session

import (
	"errors"
	"regexp"
	"strconv"
)

var uare = regexp.MustCompile(`MSIE ([0-9]{1,}[\.0-9]{0,})`)

// Returns the version of Internet Explorer or a -1
// (indicating the use of another browser).
// http://msdn.microsoft.com/en-us/library/ms537509(v=vs.85).aspx
func getInternetExplorerVersion(ua string) (float64, error) {
	matches := uare.FindStringSubmatch(ua)

	if len(matches) == 0 {
		return 0.0, errors.New("Useragent is not IE")
	}

	return strconv.ParseFloat(matches[1], 64)
}

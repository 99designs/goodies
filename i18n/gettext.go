// Package i18n implements a gettext based translation library with XLiff support
package i18n

import (
	"io/ioutil"
	"os"
	"strings"
)

type Catalog map[string]string
type LanguagesCatalog map[string]Catalog
type Parser func(fp *os.File) (Catalog, error)

func LoadCatalogs(translationsPath string, parser Parser) (LanguagesCatalog, error) {
	dirs, _ := ioutil.ReadDir(translationsPath)
	catalogs := make(map[string]Catalog)
	for _, fileInfo := range dirs {
		parts := strings.Split(fileInfo.Name(), ".")
		lang := parts[0]

		fp, err := os.Open(translationsPath + fileInfo.Name())
		if err != nil {
			return nil, err
		}
		defer fp.Close()

		catalog, err := parser(fp)
		if err != nil {
			return nil, err
		}

		catalogs[lang] = catalog
	}

	return catalogs, nil
}

func (catalog Catalog) String(message string) string {
	tmsg := catalog[message]
	if tmsg == "" {
		return message
	}
	return tmsg
}

func (catalog Catalog) StringN(msgid1, msgid2 string, n int) (tmsg string) {
	//
	// if n != 1 then it is plural
	//
	if n == 1 {
		tmsg = catalog[msgid1]
	} else {
		tmsg = catalog[msgid2]
	}
	if tmsg == "" {
		if n == 1 {
			tmsg = msgid1
		} else {
			tmsg = msgid2
		}
	}
	return tmsg
}

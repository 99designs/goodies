// Package i18n implements a gettext based translation library with XLiff support
package i18n

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Catalog map[string]string
type LanguagesCatalog map[string]Catalog
type Parser func(fp *os.File) (Catalog, error)

func LoadCatalogs(translationsPath string, parser Parser) LanguagesCatalog {
	dirs, _ := ioutil.ReadDir(translationsPath)
	catalogs := make(map[string]Catalog)
	for _, fileInfo := range dirs {
		filename := fileInfo.Name()

		// ignore dot files and directories
		if filename[0] == '.' || fileInfo.IsDir() {
			continue
		}

		parts := strings.Split(filename, ".")
		lang := parts[0]

		fp, err := os.Open(filepath.Join(translationsPath, filename))
		if err != nil {
			log.Printf("Couldn't open %s: %s\n", filename, err.Error())
			continue
		}
		defer fp.Close()

		catalog, err := parser(fp)
		if err != nil {
			log.Printf("Parser couldn't read %s: %s\n", filename, err.Error())
			continue
		}

		catalogs[lang] = catalog
	}

	return catalogs
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

package i18n

import (
	"crypto/md5"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

type Xliff struct {
	XMLName xml.Name `xml:"xliff"`
	File    File     `xml:"file"`
	Version string   `xml:"version,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
}
type File struct {
	Body           Body   `xml:"body"`
	Original       string `xml:"original,attr"`
	Datatype       string `xml:"datatype,attr"`
	SourceLanguage string `xml:"source-language,attr"`
	TargetLanguage string `xml:"target-language,attr,omitempty"`
}
type Body struct {
	TransList []TransUnit `xml:"trans-unit"`
}
type TransUnit struct {
	Id     string   `xml:"id,attr,omitempty"`
	Source string   `xml:"source"`
	Target string   `xml:"target"`
	Note   []string `xml:"note"`
}

func XliffParser(fp *os.File) (Catalog, error) {
	defer fp.Close()
	xmldata, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, errors.New("Couldn't read file")
	}

	q := Xliff{}
	err = xml.Unmarshal(xmldata, &q)
	if err != nil {
		return nil, err
	}

	catalog := make(Catalog)

	for _, transunit := range q.File.Body.TransList {
		catalog[transunit.Source] = transunit.Target
	}

	return catalog, nil
}
func (b *Body) Add(source, target string) {
	transunit := TransUnit{
		Id:     fmt.Sprintf("%x", md5.Sum([]byte(source))),
		Source: source,
		Target: target,
	}
	b.TransList = append(b.TransList, transunit)
}

type Alphabetically []TransUnit

func (a Alphabetically) Len() int           { return len(a) }
func (a Alphabetically) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Alphabetically) Less(i, j int) bool { return a[i].Source < a[j].Source }

func (b *Body) SortAlphabetically() {
	sort.Sort(Alphabetically(b.TransList))
}

func CreateXliff(catalog Catalog, sourceLanguage, targetLanguage string) []byte {

	xliff := Xliff{
		File: File{
			Body:           Body{},
			Original:       "file.ext",
			Datatype:       "plaintext",
			SourceLanguage: sourceLanguage,
			TargetLanguage: targetLanguage,
		},
		Version: "1.2",
		Xmlns:   "urn:oasis:names:tc:xliff:document:1.2",
	}
	for source, target := range catalog {
		xliff.File.Body.Add(source, target)
	}

	xliff.File.Body.SortAlphabetically()

	xml, err := xml.Marshal(xliff)
	if err != nil {
		panic(err)
	}

	return xml
}

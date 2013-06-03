package i18n

import (
	"io/ioutil"
	"os"
	"testing"
)

const fixture = `<?xml version="1.0" encoding="UTF-8"?>
<xliff xmlns="urn:oasis:names:tc:xliff:document:1.2" version="1.2">
  <file source-language="en" datatype="plaintext">
    <body>
      <trans-unit id="1">
        <source>Source test</source>
        <target>Target test</target>
      </trans-unit>
    </body>
  </file>
</xliff>
`

var fixtureFile string

func init() {
	f, _ := ioutil.TempFile(os.TempDir(), "")
	f.WriteString(fixture)
	f.Close()

	fixtureFile = f.Name()
}

func TestXliffParser(t *testing.T) {
	f, _ := os.Open(fixtureFile)
	defer f.Close()

	result, _ := XliffParser(f)

	source := "Source test"

	if result["Source test"] != "Target test" {
		t.Errorf("Lookup of '%s' failed", source)
	}
}

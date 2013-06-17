// App generator sets up a new golang webapp with all the good stuff.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/99designs/goodies/app_generator/templates"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

type Config struct {
	Name          string
	TargetPath    string
	TemplateDir   string
	AuthorName    string
	DevAssetPort  int
	DevTempPort   int
	DevServerPort int
}

var config Config
var mode os.FileMode
var t *template.Template

func init() {
	var err error
	flag.StringVar(&config.Name, "name", "New Application", "Application Name")
	flag.StringVar(&config.TargetPath, "path", "", "Path to generate application (default == name)")
	flag.StringVar(&config.TemplateDir, "templateDir", "./src/github.com/99designs/goodies/app_generator/templates", "Path to application templates")
	flag.StringVar(&config.AuthorName, "authorName", "Daniel Heath", "Name of author")
	flag.IntVar(&config.DevAssetPort, "assetPort", 8003, "Port to serve development assets on")
	flag.IntVar(&config.DevTempPort, "devTempPort", 8002, "Port to run per-request server")
	flag.IntVar(&config.DevServerPort, "devServerPort", 8001, "Port to serve development requests on")
	flag.Parse()

	if config.TargetPath == "" {
		config.TargetPath = config.Name
	}
	if config.Name == "" {
		panic("Name must be specified (use `-name MyApp`)")
	}
	config.TargetPath, err = filepath.Abs(config.TargetPath)
	panicIf(err)

	t = template.Must(
		template.New("all").
			Delims("{%%", "%%}").
			ParseGlob(config.TemplateDir + "/*"),
	)

	pwd, err := os.Getwd()
	panicIf(err)
	dir, err := os.Open(pwd)
	panicIf(err)
	stat, err := dir.Stat()
	panicIf(err)
	mode = stat.Mode()
}

func main() {
	// Empty directories with ignore files
	md("bin")
	putFile("bin/.gitignore", []byte(templates.GitIgnore))
	md("pkg")
	putFile("pkg/.gitignore", []byte(templates.GitIgnore))
	md("public")
	putFile("public/.gitignore", []byte(templates.GitIgnore))

	// Toolchain files
	putTemplate("Gemfile", "Gemfile")
	putTemplate("Makefile", "Makefile")
	putTemplate("Procfile", "Procfile")
	putTemplate("Rakefile", "Rakefile")

	// Human-friendly files
	putTemplate("Licence", "Licence")
	putTemplate("README.md", "README.md")

	// Sprockets-based asset pipeline
	md("assets")
	putFile("assets/application.css", []byte(templates.CssFile))
	putFile("assets/application.js", []byte(templates.JsFile))

	// Application config
	md("config")
	putTemplate("config/config.json", "config.json")

	// Config parser
	md("src/" + config.Name + "/application")
	putTemplate("src/"+config.Name+"/application/config.go", "config.go.templ")

	// Html templates
	md("templates")
	putTemplate("templates/sample.html", "sample.html")

	// Html template parser
	md("src/" + config.Name + "/templates")
	putTemplate("src/"+config.Name+"/templates/templates.go", "templates.go.templ")

	// Controllers
	md("src/" + config.Name + "/controllers/sample")
	putTemplate("src/"+config.Name+"/controllers/sample/sample.go", "controller.go.templ")
	putTemplate("src/"+config.Name+"/controllers/sample/routes.go", "routes.go.templ")

	// Main app file
	putTemplate("src/"+config.Name+"/"+config.Name+".go", "sample.go.templ")

	// Subprocesses need gopath set.
	os.Setenv("GOPATH", config.TargetPath)
	os.Setenv("PWD", config.TargetPath)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		subproc("bundle")
		wg.Done()
	}()

	subproc("git", "init", ".")
	subproc("git", "add", ".")
	github("gorilla/context")
	github("gorilla/mux")
	github("99designs/goodies")
	github("DanielHeath/shotgun-go")
	github("daaku/go.httpgzip")
	subproc("make")
	wg.Wait()
	subproc("git", "add", ".")
	subproc("git", "commit", "-m", "Add Submodules")
	fmt.Println("Now run 'GOPATH=`pwd` bundle exec foreman start' from your app directory.")
	fmt.Printf("Then, navigate to http://localhost:%d\n", config.DevServerPort)
}

func github(name string) {
	subproc("git", "submodule", "add", "git://github.com/"+name+".git", "src/github.com/"+name)
}

func subproc(name string, args ...string) {
	fmt.Println("Running subshell: ", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Dir = config.TargetPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func md(p string) {
	p = config.TargetPath + "/" + p
	fmt.Println("Creating directory: ", p)
	os.MkdirAll(p, mode)
}

func putTemplate(p string, name string) {
	p = config.TargetPath + "/" + p
	var wr bytes.Buffer
	err := t.ExecuteTemplate(&wr, name, config)
	panicIf(err)
	byt, err := ioutil.ReadAll(&wr)
	panicIf(err)
	fmt.Println("Writing file: ", p)
	panicIf(ioutil.WriteFile(p, byt, mode))

}

func putFile(p string, content []byte) {
	p = config.TargetPath + "/" + p
	fmt.Println("Writing file: ", p)
	ioutil.WriteFile(p, content, mode)
}

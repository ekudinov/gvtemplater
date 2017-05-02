package main

import (
	"log"

	"github.com/ekudinov/gvtemplater/cmd"
	flag "github.com/spf13/pflag"
)

func main() {
	app, err := cmd.Create()
	if err != nil {
		log.Fatal(err)
	}
	//get arguments from command line if not defined set as defaults
	var workingDir *string = flag.StringP("dir", "d", app.WorkDir, "Directory for scan")
	var packageName *string = flag.StringP("package", "p", app.Package, "Set package name for generated file")
	var templateName *string = flag.StringP("template", "n", app.TemplateName, "Name for generated file")
	flag.Parse()
	//redefine arguments
	app.WorkDir = *workingDir
	app.Package = *packageName
	app.TemplateName = *templateName
	app.Run()
}

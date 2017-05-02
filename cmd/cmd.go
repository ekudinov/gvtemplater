package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	ERR_EXISTS    = "%s already exists as %s!You must have different template names"
	ERR_EMPTY     = "%s empty file!You must check it content"
	ERR_NOT_FOUND = "No any templates found in:%s!"
	HEAD_MSG      = `package %s
/*
This file was generated from vuejs templates!Do not edit it here!
Author:Evgeniy Kudinov 2017 https://github.com/ekudinov
*/

List of files:
%s
const (

`
	VUE_EXT               = ".vue"
	DEFAULT_TEMPLATE_NAME = "templates.go"
)

//application var
var app *App

//Data keeps path to and content of file
type Data struct {
	//Absolute path to file
	Path string
	//Content of file
	Content string
}

//Application struct
type App struct {
	//map [name of file used as constant name] Data
	//with path to file and content of file
	//to use as value of constant name
	Names map[string]Data
	//working directory for scan
	WorkDir string
	//name of package for generated template file
	Package string
	//name of generated template go file(without go extension)
	TemplateName string
}

//function create app and fills it default values
//you may after creation fill yours values
func Create() (*App, error) {
	app = &App{Names: make(map[string]Data, 0)}
	//get working normal
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	app.WorkDir = dir
	//set package name as current normal name
	app.Package = filepath.Base(app.WorkDir)
	app.TemplateName = DEFAULT_TEMPLATE_NAME
	return app, nil
}

//collect data:scan all dirs to find files
//with vue extensions and get it content
func (app *App) Scan() error {
	if err := filepath.Walk(app.WorkDir, visit); err != nil {
		return err
	}
	//is found any template?
	if len(app.Names) == 0 {
		return fmt.Errorf(ERR_NOT_FOUND, app.WorkDir)
	}
	return nil
}

//function find vue files and add data to map
func visit(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !f.IsDir() {
		fullName := f.Name()
		if filepath.Ext(fullName) == VUE_EXT {
			name := strings.TrimSuffix(fullName, filepath.Ext(fullName))
			//check if name exists return collision error
			if exists, ok := app.Names[name]; ok {
				return fmt.Errorf(ERR_EXISTS, exists.Path, path)
			}
			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			//check if file is empty
			content := string(buf)
			if content == "" {
				return fmt.Errorf(ERR_EMPTY, path)
			}
			app.Names[name] = Data{Path: path, Content: content}
		}
	}
	return nil
}

//create go file with constants as templates files names
func (app *App) MakeFile() error {
	//create ordered list of keys
	var keys []string
	for k := range app.Names {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	//create list of files
	var buffer bytes.Buffer
	for i, k := range keys {
		buffer.WriteString(
			fmt.Sprintf("%d. %s - %s\n", i+1, k, app.Names[k].Path))
	}

	header := fmt.Sprintf(HEAD_MSG, app.Package, buffer.String())

	file, err := os.Create(app.TemplateName)
	defer file.Close()
	if err != nil {
		return err
	}
	//write header
	if _, err := file.WriteString(header); err != nil {
		return err
	}

	//write body
	for name, data := range app.Names {
		if _, err := file.WriteString(name + " = `\n"); err != nil {
			return err
		}
		if _, err := file.WriteString(data.Content + "`\n\n"); err != nil {
			return err
		}
	}
	end := ")"
	if _, err := file.WriteString(end); err != nil {
		return err
	}
	return nil
}

//function to run
//must be invoke after Create()
func (app *App) Run() {
	check(app.Scan())
	check(app.MakeFile())
}

//check is error exists and if true exit 1
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

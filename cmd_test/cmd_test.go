package cmd_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ekudinov/gvtemplater/cmd"
)

func TestCreateApp(t *testing.T) {
	app, _ := cmd.Create()
	dir, _ := os.Getwd()
	check(dir, app.WorkDir, t)
	check("cmd_test", app.Package, t)
	check(cmd.DEFAULT_TEMPLATE_NAME, app.TemplateName, t)
}

func TestNormalScan(t *testing.T) {
	//test normal scan
	app, _ := cmd.Create()
	app.WorkDir = "normal"
	app.Scan()
	check(len(app.Names), 2, t)
}

func TestCollisionFileNameScan(t *testing.T) {
	workDir := "existstemplate2"
	app, _ := cmd.Create()
	app.WorkDir = workDir
	err := app.Scan()
	expect := fmt.Errorf(cmd.ERR_EXISTS, workDir +"/exists/template2.vue",
		workDir +"/template2.vue")
	check(err.Error(), expect.Error(), t)
}

func TestEmptyFileScan(t *testing.T) {
	workDir := "empty"
	app, _ := cmd.Create()
	app.WorkDir = workDir
	err := app.Scan()
	expect := fmt.Errorf(cmd.ERR_EMPTY, workDir + "/empty.vue")
	checkErr(err, expect, t)
}

func TestNoTemplatesFoundScan(t *testing.T) {
	workDir := "notemplates"
	app, _ := cmd.Create()
	app.WorkDir = workDir
	err := app.Scan()
	expect := fmt.Errorf(cmd.ERR_NOT_FOUND, workDir)
	checkErr(err, expect, t)
}

func TestMakeFile(t *testing.T) {
	const checking = "forcheck_test"
	const generated = "templates_test"
	app, _ := cmd.Create()
	app.WorkDir = "normal"
	app.TemplateName = generated
	app.Run()
	che := loadFile(checking, t)
	gen := loadFile(generated, t)
	check(gen, che, t)
}

func loadFile(filename string, t *testing.T) string {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error(err)
		t.SkipNow()
	}
	return string(buf)
}
func check(result, expected interface{}, t *testing.T) {
	if result != expected {
		t.Error("Must be:", expected, " but got:", result)
	}
}

func checkErr(result, expected error, t *testing.T) {
	if result != nil && expected != nil {
		res := result.Error()
		exp := expected.Error()
		check(res, exp, t)
	} else {
		t.Errorf("Errors must be not nil - result:%s, expected:%s", result, expected)
	}
}

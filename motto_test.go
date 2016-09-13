// Copyright 2014 dong<ddliuhb@gmail.com>.
// Licensed under the MIT license.
//
// Motto - Modular Javascript environment.
package motto

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/robertkrimen/otto"

	//	_ "./underscore"

	//	"reflect"

	//	"encoding/json"

	"bytes"
	//	"errors"
	//	"reflect"
	"strings"

	//	"github.com/go-errors/errors"
	//	"log"
	"os"
	//	"sync"
	//	"sync/atomic"
	//	"time"

	"github.com/fatih/color"
)

func TestModule(t *testing.T) {
	_, v, err := Run("tests/index.js")
	if err != nil {
		t.Error(err)
	}

	i, _ := v.ToString()
	//	v.IsObject()
	//	v.
	if i != "rat" {
		t.Error("testing result: ", i, "!=", "rat")
	}
}

func SkipTestModuleBench(t *testing.T) {
	for x := 0; x < 1000; x++ {
		_, v, err := Run("tests/index.js")
		if err != nil {
			t.Error(err)
		}

		i, _ := v.ToString()
		//	v.IsObject()
		//	v.
		if i != "rat" {
			t.Error("testing result: ", i, "!=", "rat")
		}
	}
}

func TestCompileThenRun(t *testing.T) {
	vm := New()
	jsRequire := func(call otto.FunctionCall) otto.Value {
		name, _ := call.Argument(0).ToString()
		var pwd = "tests"
		if len(call.ArgumentList) >= 2 {
			pwd, _ = call.Argument(1).ToString()
		} else {

		}
		result, err := vm.Require(name, pwd)
		if err != nil {
			t.Logf("Require error %s", err.Error())
		}
		return result
	}

	vm.Set("require", jsRequire)

	jsData, err := ioutil.ReadFile("tests/index_compile.js")
	if err != nil {
		t.Error(err)
	}
	s, err := vm.Compile("tindex.js", jsData)
	if err != nil {
		t.Error(err)
	}

	for x := 0; x < 1000; x++ {

		v, err := vm.Eval(s)

		//_, v, err := Run("tests/index.js")
		if err != nil {
			t.Error(err)
		}

		i, _ := v.ToString()
		if i != "rat" {
			t.Error("testing result: ", i, "!=", "rat")
		}
	}
}

func TestCompileThenRun2(t *testing.T) {
	vm := New()

	jsData, err := ioutil.ReadFile("tests/index_compile.js")
	if err != nil {
		t.Error(err)
	}
	s, err := vm.CompilePrepare("tindex.js", jsData, "tests")
	if err != nil {
		t.Error(err)
	}

	for x := 0; x < 1000; x++ {

		v, err := vm.Run(s)

		//_, v, err := Run("tests/index.js")
		if err != nil {
			t.Error(err)
		}

		i, _ := v.ToString()
		if i != "rat" {
			t.Error("testing result: ", i, "!=", "rat")
		}
	}
}

func TestCompileExport(t *testing.T) {
	vm := New()

	//	jsData, err := ioutil.ReadFile("tests/index_export.js")
	//	if err != nil {
	//		t.Error(err)
	//	}

	exports, err := vm.RunDirect("tests/index_export.js")
	if err != nil {
		t.Error(err)
	}

	search, err := exports.Object().Get("search")

	if err != nil {
		t.Error(err)
	}

	for x := 0; x < 1000/5; x++ {

		for y := 0; y < 5; y++ {
			v, err := search.Call(search, y)

			//_, v, err := Run("tests/index.js")
			if err != nil {
				t.Error(err)
			}

			i, _ := v.ToString()
			//fmt.Printf("data is, %s\n", i)
			//t.Logf("data is, %s", i)
			//			if i != "rat" {
			//				t.Error("testing result: ", i, "!=", "rat")
			//			}
			_ = i
		}
	}
	fmt.Printf("data is, %d\n", 0)
}

type Message struct {
	Payload  *string
	HostName *string
}

type PipelinePack struct {
	Message Message
}
type PipelinePackProxy struct {
	*otto.Object
	Output *PipelinePack
}

func (proxy *PipelinePackProxy) GetTarget(call otto.FunctionCall) (target *PipelinePack) {
	target = nil
	if name_, err := call.This.Object().Get("name"); err == nil {
		if name, err := name_.ToString(); err == nil {
			switch name {
			case "output":
				target = proxy.Output
				break
			default:
				target = proxy.Output
				break
			}
		}
	}
	return
}

func NewPipelinePackProxy(vm *Motto, proxy *PipelinePackProxy) {
	var (
		buffer bytes.Buffer
		err    error
	)

	proto_tmpl := `
	set{{.Name}}:function(value){
		console.log("set{{.Name}} called:" + String(value));
	},
	get{{.Name}}:function(){
		return "fake {{.Name}}";
	},`
	_ = proto_tmpl

	property_tmpl := `
	{{.Name}}:{
		get:function(){this.get{{.Name}}2();},
		set:function(v){this.set{{.Name}}2(v);}
	},`
	buffer.WriteString("(Object.create(")
	buffer.WriteString("{")
	//	buffer.WriteString(strings.Replace(proto_tmpl, "{{.Name}}", "Playload", -1))
	//	buffer.WriteString(strings.Replace(proto_tmpl, "{{.Name}}", "Type", -1))
	//	buffer.WriteString(strings.Replace(proto_tmpl, "{{.Name}}", "Logger", -1))
	//	buffer.WriteString(strings.Replace(proto_tmpl, "{{.Name}}", "Hostname", -1))
	//	buffer.WriteString(strings.Replace(proto_tmpl, "{{.Name}}", "EnvVersion", -1))
	buffer.WriteString("\n name:\"in/out\"")

	buffer.WriteString("},{")
	buffer.WriteString(strings.Replace(property_tmpl, "{{.Name}}", "Playload", -1))
	buffer.WriteString(strings.Replace(property_tmpl, "{{.Name}}", "Type", -1))
	buffer.WriteString(strings.Replace(property_tmpl, "{{.Name}}", "Logger", -1))
	buffer.WriteString(strings.Replace(property_tmpl, "{{.Name}}", "Hostname", -1))
	buffer.WriteString(strings.Replace(property_tmpl, "{{.Name}}", "EnvVersion", -1))
	buffer.WriteString("\n}")
	buffer.WriteString("\n))")

	packobjtmpl := buffer.String()
	buffer.Truncate(0)

	proxytmpl := strings.Replace(fmt.Sprintf(`(pipe={
 input:%s,
 output:%s
	}, pipe)`, packobjtmpl, packobjtmpl), "	", "    ", -1)

	color.Set(color.FgGreen)
	fmt.Fprintln(os.Stderr, proxytmpl)
	color.Unset()

	proxy.Object, err = vm.Object(proxytmpl)
	if err != nil {
		fmt.Printf("js vm init error: %s", err)
		panic(err)
	}

	return
}

func checkNative(proxy *PipelinePackProxy, name string) {

	getPlayload := func(call otto.FunctionCall) otto.Value {
		//		data := call.Argument(0).String()
		if target := proxy.GetTarget(call); target != nil {
			val, err := otto.ToValue(target.Message.Payload)
			_ = err
			return val
		} else {
			return otto.NullValue()
		}
	}

	setPlayload := func(call otto.FunctionCall) otto.Value {
		data := call.Argument(0).String()
		if target := proxy.GetTarget(call); target != nil {
			target.Message.Payload = &data
		}
		return otto.Value{}
	}

	input_, err := proxy.Get(name)
	if err != nil {
		panic(err)
	}
	_ = err
	input := input_.Object()
	input.Set("getPlayload2", getPlayload)
	input.Set("setPlayload2", setPlayload)
	input.Call("getPlayload2", "")
	input.Call("setPlayload2", "playdddd")

}

func TestCompileExportWithNative(t *testing.T) {
	vm := New()

	//	jsData, err := ioutil.ReadFile("tests/index_export.js")
	//	if err != nil {
	//		t.Error(err)
	//	}

	proxy := &PipelinePackProxy{}
	s1 := "Payload"
	s2 := "HostName"
	Msg := Message{
		Payload:  &s1,
		HostName: &s2,
	}
	proxy.Output = &PipelinePack{
		Message: Msg,
	}

	NewPipelinePackProxy(vm, proxy)

	checkNative(proxy, "input")

	proxy.Object.Get("input")
	exports, err := vm.RunDirect("tests/index_export.js")
	if err != nil {
		t.Error(err)
	}

	search, err := exports.Object().Get("search")

	if err != nil {
		t.Error(err)
	}

	for x := 0; x < 1000/5; x++ {

		for y := 0; y < 5; y++ {
			v, err := search.Call(search, y)

			//_, v, err := Run("tests/index.js")
			if err != nil {
				t.Error(err)
			}

			i, _ := v.ToString()
			//fmt.Printf("data is, %s\n", i)
			//t.Logf("data is, %s", i)
			//			if i != "rat" {
			//				t.Error("testing result: ", i, "!=", "rat")
			//			}
			_ = i
		}
	}
	fmt.Printf("data is, %d\n", 0)
}

func __TestNpmModule(t *testing.T) {
	_, v, err := Run("tests/npm/index.js")

	if err != nil {
		t.Error(err)
	}

	i, _ := v.ToInteger()

	if i != 1 {
		t.Error("npm test failed: ", i, "!=", 1)
	}
}

func TestCoreModule(t *testing.T) {
	vm := New()
	vm.AddModule("fs", fsModuleLoader)

	v, err := vm.RunDirect("tests/core_module_test.js")
	if err != nil {
		t.Error(err)
	}

	s, _ := v.ToString()
	if s != "cat" {
		t.Error("core module test failed: ", s, "!=", "cat")
	}
}

func fsModuleLoader(vm *Motto) (otto.Value, error) {
	fs, _ := vm.Object(`({})`)
	fs.Set("readFileSync", func(call otto.FunctionCall) otto.Value {
		filename, _ := call.Argument(0).ToString()
		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			return otto.UndefinedValue()
		}

		v, _ := call.Otto.ToValue(string(bytes))
		return v
	})

	fs.Call("readFileSync", "tests/data.json")

	return vm.ToValue(fs)
}

package components

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/dop251/goja"
	"github.com/rotisserie/eris"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/core"
)

// CodejsComponent 代码执行组件
type CodeComponent struct {
	engine *GojaJsEngine
}

var codeComponentPool = sync.Pool{
	New: func() interface{} {
		return &CodeComponent{}
	},
}

type CodeConfig struct {
	Code              string `json:"code"`
	Language          string `json:"language"`
	ErrorHandlingMode string `json:"errorHandlingMode"`
	Retry             int64  `json:"retry"`
	Timeout           int64  `json:"timeout"`
}

const (
	languagesJs     = "javascript"
	languagesGolang = "golang"
	languagesPython = "python"

	errorHandlingModes_abort = "abort" // 中断
	errorHandlingModes_retry = "retry" // 重试
)

func (c *CodeComponent) Name() string {
	return Code
}

func NewCodeComponent(config any) (*CodeComponent, error) {
	jsonConfig, err := sonic.Marshal(config)
	if err != nil {
		return nil, eris.Wrap(err, "failed to parse code execution component config")
	}
	component := codeComponentPool.Get().(*CodeComponent)
	var codeConfig CodeConfig
	if err := sonic.Unmarshal(jsonConfig, &codeConfig); err != nil {
		return nil, eris.Wrap(err, "failed to parse code execution component config")
	}
	logx.Debugf("[Code Execution] Config content: %s", codeConfig.Code)
	engine, err := NewGojaJsEngine(codeConfig.Code, nil, codeConfig.Timeout)
	if err != nil {
		logx.Errorf("[Code Execution] Failed to create engine [error:%v]", err)
		return nil, eris.Wrap(err, "failed to create code execution component")
	}
	component.engine = engine
	return component, nil
}

func (c *CodeComponent) Clear() {
	codeComponentPool.Put(c)
}
func (c *CodeComponent) Execute(ctx context.Context, input any) (*core.Result, error) {
	result, err := c.engine.Execute("main", input)
	if err != nil {
		return nil, eris.Wrap(err, "failed to execute code component")
	}

	return &core.Result{
		Output: result,
		Route:  []string{Success},
	}, nil
}

func (c *CodeComponent) Validate() []core.ValidationError {
	return nil
}

func (c *CodeComponent) AnalyzeInputs(ctx context.Context) (any, error) {
	return nil, nil
}

// GojaJsEngine goja js engine
type GojaJsEngine struct {
	vmPool            sync.Pool
	config            Config
	jsScript          *goja.Program
	jsUdfProgramCache map[string]*goja.Program
}

type Config struct {
	Udf     map[string]interface{}
	Timeout int64
}

// NewGojaJsEngine Create a new instance of the JavaScript engine
func NewGojaJsEngine(jsScript string, fromVars map[string]interface{}, timeout int64) (*GojaJsEngine, error) {
	config := Config{
		Udf:     make(map[string]interface{}),
		Timeout: timeout,
	}
	if config.Udf == nil {
		config.Udf = make(map[string]interface{})
	}
	name := "gmd5"
	config.Udf[name] = func(v string) string {
		return fmt.Sprintf("%x", md5.Sum([]byte(v)))
	}

	program, err := goja.Compile("", jsScript, true)
	if err != nil {
		logx.Errorf("[Code Execution] Failed to compile JS script [error:%v]", err)
		return nil, eris.Wrap(err, "failed to compile JS script")
	}
	jsEngine := &GojaJsEngine{
		config:   config,
		jsScript: program,
	}
	if err = jsEngine.PreCompileJs(config); err != nil {
		logx.Errorf("[Code Execution] Failed to pre-compile JS script [error:%v]", err)
		return nil, eris.Wrap(err, "failed to pre-compile JS script")
	}
	jsEngine.vmPool = sync.Pool{
		New: func() interface{} {
			return jsEngine.NewVm(config, fromVars)
		},
	}
	return jsEngine, nil
}

// PreCompileJs Precompiled UDF JavaScript file
func (g *GojaJsEngine) PreCompileJs(config Config) error {
	var jsUdfProgramCache = make(map[string]*goja.Program)
	for k, v := range config.Udf {
		if jsFuncStr, ok := v.(string); ok {
			if p, err := goja.Compile(k, jsFuncStr, true); err != nil {
				logx.Errorf("[Code Execution] Failed to compile UDF script [function:%s] [error:%v]", k, err)
				return eris.Wrap(err, "failed to compile UDF script")
			} else {
				jsUdfProgramCache[k] = p
			}
		} else if script, scriptOk := v.(Script); scriptOk {
			if script.Type == Js || script.Type == "" {
				if c, ok := script.Content.(string); ok {
					if p, err := goja.Compile(k, c, true); err != nil {
						logx.Errorf("[Code Execution] Failed to compile script content [function:%s] [error:%v]", k, err)
						return eris.Wrap(err, "failed to compile script content")
					} else {
						jsUdfProgramCache[k] = p
					}
				} else if p, ok := script.Content.(*goja.Program); ok {
					jsUdfProgramCache[k] = p
				}
			}
		}
	}
	g.jsUdfProgramCache = jsUdfProgramCache

	return nil
}

// NewVm new a js VM
func (g *GojaJsEngine) NewVm(config Config, fromVars map[string]any) *goja.Runtime {
	vm := goja.New()
	vars := make(map[string]interface{})
	if fromVars != nil {
		for k, v := range fromVars {
			vars[k] = v
		}
	}

	// Add global custom functions to the JavaScript runtime
	for k, v := range config.Udf {
		var err error
		if _, ok := v.(string); ok {
			if p, ok := g.jsUdfProgramCache[k]; ok {
				_, err = vm.RunProgram(p)
			}
		} else if script, scriptOk := v.(Script); scriptOk {
			if script.Type == Js || script.Type == "" {
				// parse  JS script
				if _, ok := script.Content.(string); ok {
					if p, ok := g.jsUdfProgramCache[k]; ok {
						_, err = vm.RunProgram(p)
					}
				} else if _, ok := script.Content.(*goja.Program); ok {
					if p, ok := g.jsUdfProgramCache[k]; ok {
						_, err = vm.RunProgram(p)
					}
				} else {
					funcName := strings.Replace(k, Js+"#", "", 1)
					vars[funcName] = vm.ToValue(script.Content)
				}
			}
		} else {
			// parse go func
			vars[k] = vm.ToValue(v)
		}
		if err != nil {
			logx.Errorf("[Code Execution] Failed to parse JS script [script:%s] [error:%v]", k, err)
		}
	}
	for k, v := range vars {
		if err := vm.Set(k, v); err != nil {
			logx.Errorf("[Code Execution] Failed to set variable [variable:%s] [error:%v]", k, err)
		}
	}

	state := g.setTimeout(vm)

	_, err := vm.RunProgram(g.jsScript)
	// If there is no timeout, state=0; otherwise, state=-2
	closeStateChan(state)

	if err != nil {
		logx.Errorf("[Code Execution] JS VM execution failed [error:%v]", err)
	}
	return vm
}

// Execute Execute JavaScript script
func (g *GojaJsEngine) Execute(functionName string, argumentList ...any) (out interface{}, err error) {
	defer func() {
		if caught := recover(); caught != nil {
			logx.Errorf("[Code Execution] Panic occurred during execution [function:%s] [error:%v]", functionName, caught)
			err = eris.New(fmt.Sprintf("%s", caught))
		}
	}()

	vm := g.vmPool.Get().(*goja.Runtime)

	// vm.Set(CtxKey, ctx)

	state := g.setTimeout(vm)

	f, ok := goja.AssertFunction(vm.Get(functionName))
	if !ok {
		logx.Errorf("[Code Execution] Function does not exist [function:%s]", functionName)
		return nil, eris.New(functionName + " is not a function")
	}
	var params []goja.Value
	for _, v := range argumentList {
		params = append(params, vm.ToValue(v))
	}
	res, err := f(goja.Undefined(), params...)
	// If there is no timeout, state=0; otherwise, state=-2
	closeStateChan(state)
	// Put back to the pool
	g.vmPool.Put(vm)
	if err != nil {
		params, _ := sonic.Marshal(argumentList)
		logx.Errorf("[Code Execution] Failed to execute function [function:%s] [params:%s] [error:%v]", functionName, string(params), err)
		return nil, eris.Wrap(err, "failed to execute function")
	}
	return res.Export(), nil
}

func (g *GojaJsEngine) Stop() {
}

// setTimeout if timeout interrupt the js script execution
func (g *GojaJsEngine) setTimeout(vm *goja.Runtime) chan int {
	state := make(chan int, 1)
	state <- 0
	time.AfterFunc(time.Duration(g.config.Timeout)*time.Second, func() {
		if <-state == 0 {
			state <- 2
			vm.Interrupt("execution timeout")
		}
	})
	return state
}

func closeStateChan(state chan int) {
	if <-state == 0 {
		state <- 1
	}
	close(state)
}

// Script is used to register native functions or custom functions defined in Go.
type Script struct {
	// Type is the script type, default is Js.
	Type string
	// Content is the script content or custom function.
	Content interface{}
}

const (
	Js = "Js" // Represents JavaScript scripting language.
)

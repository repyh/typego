package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dop251/goja"
	"github.com/repyh3/typego/bridge"
	"github.com/repyh3/typego/compiler"
	"github.com/repyh3/typego/engine"
	"github.com/spf13/cobra"
)

type NativeTools struct {
	StartTime string
}

func (n *NativeTools) GetRuntimeInfo() string {
	return "TypeGo Supercharged Runtime v1.0"
}

var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Run a TypeScript file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		absPath, _ := filepath.Abs(filename)

		res, err := compiler.Compile(absPath, nil)
		if err != nil {
			fmt.Printf("Build Error: %v\n", err)
			os.Exit(1)
		}

		eng := engine.NewEngine(128*1024*1024, nil)

		cliBuffer := make([]byte, 1024)
		bridge.MapSharedBuffer(eng.VM, "cliBuffer", cliBuffer)

		tools := &NativeTools{StartTime: "2026-01-16"}
		_ = bridge.BindStruct(eng.VM, "native", tools)

		eng.EventLoop.RunOnLoop(func() {
			val, err := eng.Run(res.JS)
			if err != nil {
				fmt.Printf("Runtime Error: %v\n", err)
				os.Exit(1)
			}

			if val != nil && !goja.IsUndefined(val) && !goja.IsNull(val) {
				if obj := val.ToObject(eng.VM); obj != nil {
					then := obj.Get("then")
					if then != nil && !goja.IsUndefined(then) {
						if _, ok := goja.AssertFunction(then); ok {
							eng.EventLoop.WGAdd(1)
							done := eng.VM.ToValue(func(goja.FunctionCall) goja.Value {
								eng.EventLoop.WGDone()
								return goja.Undefined()
							})
							thenFn, _ := goja.AssertFunction(then)
							_, _ = thenFn(val, done, done)
						}
					}
				}
			}
		})

		eng.EventLoop.Start()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

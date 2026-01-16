package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/repyh3/typego/compiler"
	"github.com/spf13/cobra"
)

var buildOut string
var minify bool

var buildCmd = &cobra.Command{
	Use:   "build [file]",
	Short: "Build and bundle a TypeScript file for production",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		absPath, _ := filepath.Abs(filename)

		fmt.Printf("📦 Building %s...\n", absPath)
		res, err := compiler.Compile(absPath)
		if err != nil {
			fmt.Printf("Build Error: %v\n", err)
			os.Exit(1)
		}

		// 2. Prepare Temp Directory
		tmpDir := ".typego_build_tmp"
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			fmt.Printf("Error creating temp dir: %v\n", err)
			os.Exit(1)
		}
		defer os.RemoveAll(tmpDir) // Cleanup

		// 3. Generate Shim (main.go)
		// We embed the JS code directly into the Go binary
		shimContent := fmt.Sprintf(shimTemplate, fmt.Sprintf("%q", res.JS))

		shimPath := filepath.Join(tmpDir, "main.go")
		if err := os.WriteFile(shimPath, []byte(shimContent), 0644); err != nil {
			fmt.Printf("Error writing shim: %v\n", err)
			os.Exit(1)
		}

		// 4. Generate go.mod for the shim
		// We assume we are in the TypeGo repo, so we can point to local source
		// In a real release, this would point to the github version
		currWd, _ := os.Getwd()
		// Calculate relative path from tmpDir back to root (usually just "..")
		// But let's use absolute path for safety in the replace directive

		goModContent := fmt.Sprintf(`module typego_app

go 1.21

require github.com/repyh3/typego v0.0.0

replace github.com/repyh3/typego => %s
`, filepath.ToSlash(currWd))

		if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
			fmt.Printf("Error writing go.mod: %v\n", err)
			os.Exit(1)
		}

		// 5. Build Binary
		outputName := buildOut
		if outputName == "" {
			outputName = "app.exe"
		}
		// Make output absolute so go build puts it in the right place
		absOut, _ := filepath.Abs(outputName)

		fmt.Println("🧹 Resolving dependencies...")
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = tmpDir
		tidyCmd.Stdout = os.Stdout
		tidyCmd.Stderr = os.Stderr
		if err := tidyCmd.Run(); err != nil {
			fmt.Printf("go mod tidy failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("🔨 Compiling binary...")
		buildCmd := exec.Command("go", "build", "-o", absOut, ".")
		buildCmd.Dir = tmpDir
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr

		if err := buildCmd.Run(); err != nil {
			fmt.Printf("Compilation failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✨ Binary created: %s\n", outputName)
	},
}

const shimTemplate = `package main

import (
	"fmt"
	"os"

	"github.com/dop251/goja"
	"github.com/repyh3/typego/bridge"
	"github.com/repyh3/typego/engine"
)

const jsBundle = %s

type NativeTools struct {
	StartTime string
}

func (n *NativeTools) GetRuntimeInfo() string {
	return "TypeGo Standalone v1.0"
}

func main() {
	eng := engine.NewEngine(128*1024*1024, nil)

	// Initialize Shared Buffer
	cliBuffer := make([]byte, 1024)
	bridge.MapSharedBuffer(eng.VM, "cliBuffer", cliBuffer)

	// Initialize Native Tools
	tools := &NativeTools{StartTime: "2026-01-16"}
	_ = bridge.BindStruct(eng.VM, "native", tools)

	// Run on EventLoop
	eng.EventLoop.RunOnLoop(func() {
		val, err := eng.Run(jsBundle)
		if err != nil {
			fmt.Printf("Runtime Error: %v\n", err)
			os.Exit(1)
		}

		// Handle Top-Level Async (Promises)
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
}
`

func init() {
	buildCmd.Flags().StringVarP(&buildOut, "out", "o", "dist/index.js", "Output bundle path")
	buildCmd.Flags().BoolVarP(&minify, "minify", "m", false, "Minify output")
	rootCmd.AddCommand(buildCmd)
}

package compiler

import (
	"fmt"

	"github.com/evanw/esbuild/pkg/api"
)

type Result struct {
	JS        string
	SourceMap string
}

func Compile(entryPoint string) (*Result, error) {
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{entryPoint},
		Bundle:      true,
		Write:       false,
		LogLevel:    api.LogLevelSilent,
		Target:      api.ESNext,
		Format:      api.FormatIIFE,
		Sourcemap:   api.SourceMapInline,
		Plugins: []api.Plugin{
			{
				Name: "typego-virtual",
				Setup: func(build api.PluginBuild) {
					build.OnResolve(api.OnResolveOptions{Filter: `^go/.*`}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
						return api.OnResolveResult{Path: args.Path, Namespace: "typego-internal"}, nil
					})
					build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "typego-internal"}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
						var content string
						switch args.Path {
						case "go/memory":
							content = "export const Ptr = (globalThis as any).Ptr; const mf = (globalThis as any).__go_memory_factory__; export const makeShared = mf.makeShared;"
						case "go/fmt":
							content = "const f = (globalThis as any).__go_fmt__; export const Println = f.Println;"
						case "go/os":
							content = "const o = (globalThis as any).__go_os__; export const WriteFile = o.WriteFile;"
						case "go/net/http":
							content = "const h = (globalThis as any).__go_http__; export const Get = h.Get; export const Fetch = h.Fetch;"
						case "go/sync":
							content = "const s = (globalThis as any).__go_sync__; export const Spawn = s.Spawn; export const Sleep = s.Sleep; export const Chan = (globalThis as any).Chan;"
						default:
							return api.OnLoadResult{Errors: []api.Message{{Text: "Unknown virtual module: " + args.Path}}}, nil
						}
						return api.OnLoadResult{Contents: &content, Loader: api.LoaderTS}, nil
					})
				},
			},
		},
	})

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("compilation failed: %v", result.Errors[0].Text)
	}

	res := &Result{}
	for _, file := range result.OutputFiles {
		if file.Path == "<stdout>.js" || len(result.OutputFiles) == 1 {
			res.JS = string(file.Contents)
		} else if file.Path == "<stdout>.js.map" {
			res.SourceMap = string(file.Contents)
		}
	}

	if res.JS == "" && len(result.OutputFiles) > 0 {
		res.JS = string(result.OutputFiles[0].Contents)
	}

	return res, nil
}

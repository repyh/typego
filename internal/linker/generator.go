package linker

import (
	"fmt"
	"strings"
)

// GenerateShim creates the Go code to bind the package to the VM
func GenerateShim(info *PackageInfo, variableName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\n\t// Bind %s\n", info.ImportPath))
	sb.WriteString(fmt.Sprintf("\t%s := eng.VM.NewObject()\n", variableName))

	for _, fn := range info.Exports {
		// Direct binding lets Goja handle reflection and argument mapping
		// This supports variadic, primitives, and even structs (to some extent)
		sb.WriteString(fmt.Sprintf("\t%s.Set(%q, %s.%s)\n", variableName, fn.Name, info.Name, fn.Name))
	}

	sb.WriteString(fmt.Sprintf("\teng.VM.Set(%q, %s)\n", "_go_hyper_"+info.Name, variableName))
	return sb.String()
}

// GenerateTypes creates the TypeScript definition with JSDoc
func GenerateTypes(info *PackageInfo) string {
	var sb strings.Builder

	// Add Marker for additive linking
	sb.WriteString(fmt.Sprintf("// MODULE: go:%s\n", info.Name))

	sb.WriteString(fmt.Sprintf("declare module \"go:%s\" {\n", info.ImportPath))

	for _, fn := range info.Exports {
		if fn.Doc != "" {
			sb.WriteString("\t/**\n")
			lines := strings.Split(fn.Doc, "\n")
			for _, line := range lines {
				sb.WriteString(fmt.Sprintf("\t * %s\n", line))
			}
			sb.WriteString("\t */\n")
		}

		// Simplify args for PoC
		sb.WriteString(fmt.Sprintf("\texport function %s(...args: any[]): any;\n", fn.Name))
	}

	sb.WriteString("}\n")
	sb.WriteString(fmt.Sprintf("// END: go:%s\n", info.Name))
	return sb.String()
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var rootDir string

func main() {
	target := flag.String("target", "", "Target to generate repo map for: backend or frontend")
	root := flag.String("root", ".", "Project root directory")
	flag.Parse()

	if *target == "" {
		fmt.Println("Usage: go run . --target=backend|frontend [--root=/path/to/project]")
		os.Exit(1)
	}

	rootDir = *root

	switch *target {
	case "backend":
		if err := generateBackendRepoMap(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "frontend":
		if err := generateFrontendRepoMap(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown target: %s\n", *target)
		os.Exit(1)
	}
}

// Backend (Go) repo map generation

type GoPackage struct {
	Path        string
	IsMain      bool
	Types       []string
	Funcs       []string
	Methods     map[string][]string // receiver -> methods
	Consts      []string
	Vars        []string
}

func generateBackendRepoMap() error {
	packages := make(map[string]*GoPackage)

	backendDir := filepath.Join(rootDir, "backend")
	err := filepath.WalkDir(backendDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if shouldSkipDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		return parseGoFile(path, packages)
	})

	if err != nil {
		return err
	}

	return writeBackendRepoMap(packages)
}

func shouldSkipDir(name string) bool {
	skip := []string{".git", "vendor", "node_modules", "testdata"}
	for _, s := range skip {
		if name == s {
			return true
		}
	}
	return false
}

func parseGoFile(path string, packages map[string]*GoPackage) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", path, err)
	}

	dir := filepath.Dir(path)
	pkgName := node.Name.Name

	pkg, ok := packages[dir]
	if !ok {
		pkg = &GoPackage{
			Path:    dir,
			IsMain:  pkgName == "main",
			Methods: make(map[string][]string),
		}
		packages[dir] = pkg
	}

	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if s.Name.IsExported() {
						pkg.Types = append(pkg.Types, s.Name.Name)
					}
				case *ast.ValueSpec:
					for _, name := range s.Names {
						if name.IsExported() {
							if d.Tok == token.CONST {
								pkg.Consts = append(pkg.Consts, name.Name)
							} else {
								pkg.Vars = append(pkg.Vars, name.Name)
							}
						}
					}
				}
			}
		case *ast.FuncDecl:
			if d.Name.IsExported() {
				if d.Recv != nil && len(d.Recv.List) > 0 {
					recv := getReceiverName(d.Recv.List[0].Type)
					pkg.Methods[recv] = append(pkg.Methods[recv], d.Name.Name)
				} else {
					pkg.Funcs = append(pkg.Funcs, d.Name.Name)
				}
			}
		}
	}

	return nil
}

func getReceiverName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return getReceiverName(t.X)
	}
	return ""
}

func writeBackendRepoMap(packages map[string]*GoPackage) error {
	var b strings.Builder

	b.WriteString("# Backend Repo Map\n\n")
	b.WriteString("> Auto-generated. Do not edit manually.\n\n")

	// Sort packages
	var paths []string
	for path := range packages {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	// Entry points first
	var entryPoints []string
	var regularPackages []string

	for _, path := range paths {
		if packages[path].IsMain {
			entryPoints = append(entryPoints, path)
		} else {
			regularPackages = append(regularPackages, path)
		}
	}

	if len(entryPoints) > 0 {
		b.WriteString("## Entry Points\n\n")
		for _, path := range entryPoints {
			files, _ := filepath.Glob(filepath.Join(path, "*.go"))
			for _, f := range files {
				if !strings.HasSuffix(f, "_test.go") {
					relPath, _ := filepath.Rel(rootDir, f)
					b.WriteString(fmt.Sprintf("- %s\n", relPath))
				}
			}
		}
		b.WriteString("\n")
	}

	if len(regularPackages) > 0 {
		b.WriteString("## Packages\n\n")
		for _, path := range regularPackages {
			pkg := packages[path]
			relPath, _ := filepath.Rel(rootDir, path)
			b.WriteString(fmt.Sprintf("### %s\n\n", relPath))

			// Types
			sort.Strings(pkg.Types)
			for _, t := range pkg.Types {
				b.WriteString(fmt.Sprintf("- type %s\n", t))
				// Methods for this type
				if methods, ok := pkg.Methods[t]; ok {
					sort.Strings(methods)
					for _, m := range methods {
						b.WriteString(fmt.Sprintf("  - func (%s) %s\n", t, m))
					}
				}
			}

			// Functions
			sort.Strings(pkg.Funcs)
			for _, f := range pkg.Funcs {
				b.WriteString(fmt.Sprintf("- func %s\n", f))
			}

			// Constants
			sort.Strings(pkg.Consts)
			for _, c := range pkg.Consts {
				b.WriteString(fmt.Sprintf("- const %s\n", c))
			}

			// Variables
			sort.Strings(pkg.Vars)
			for _, v := range pkg.Vars {
				b.WriteString(fmt.Sprintf("- var %s\n", v))
			}

			b.WriteString("\n")
		}
	}

	outputPath := filepath.Join(rootDir, "backend", "REPO_MAP.md")
	return os.WriteFile(outputPath, []byte(b.String()), 0644)
}

// Frontend (TypeScript/React) repo map generation

type FrontendModule struct {
	Path           string
	DefaultExport  string
	NamedExports   []string
	IsRoute        bool
}

func generateFrontendRepoMap() error {
	modules := make(map[string]*FrontendModule)

	frontendDir := filepath.Join(rootDir, "frontend")
	err := filepath.WalkDir(frontendDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if shouldSkipDir(d.Name()) || d.Name() == ".next" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".ts" && ext != ".tsx" {
			return nil
		}

		// Skip config files
		base := filepath.Base(path)
		if base == "next.config.ts" || base == "next-env.d.ts" || strings.HasSuffix(base, ".d.ts") {
			return nil
		}

		return parseTSFile(path, modules)
	})

	if err != nil {
		return err
	}

	return writeFrontendRepoMap(modules)
}

func parseTSFile(path string, modules map[string]*FrontendModule) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	relPath, _ := filepath.Rel(rootDir, path)
	module := &FrontendModule{
		Path:    relPath,
		IsRoute: strings.Contains(path, "app/") && (strings.HasSuffix(path, "page.tsx") || strings.HasSuffix(path, "layout.tsx")),
	}

	// Parse exports using regex
	text := string(content)

	// Default export: export default function Name or export default Name
	defaultFuncRe := regexp.MustCompile(`export\s+default\s+(?:async\s+)?function\s+(\w+)`)
	if matches := defaultFuncRe.FindStringSubmatch(text); len(matches) > 1 {
		module.DefaultExport = matches[1]
	}

	// Named exports: export function Name, export const Name, export type Name, export interface Name
	namedExportRe := regexp.MustCompile(`export\s+(?:async\s+)?(?:function|const|let|type|interface)\s+(\w+)`)
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		if matches := namedExportRe.FindStringSubmatch(line); len(matches) > 1 {
			module.NamedExports = append(module.NamedExports, matches[1])
		}
	}

	// Export { Name } patterns
	reExportRe := regexp.MustCompile(`export\s*\{\s*([^}]+)\s*\}`)
	if matches := reExportRe.FindAllStringSubmatch(text, -1); len(matches) > 0 {
		for _, m := range matches {
			exports := strings.Split(m[1], ",")
			for _, e := range exports {
				e = strings.TrimSpace(e)
				if e != "" && !strings.Contains(e, " as ") {
					module.NamedExports = append(module.NamedExports, e)
				}
			}
		}
	}

	if module.DefaultExport != "" || len(module.NamedExports) > 0 || module.IsRoute {
		modules[path] = module
	}

	return nil
}

func writeFrontendRepoMap(modules map[string]*FrontendModule) error {
	var b strings.Builder

	b.WriteString("# Frontend Repo Map\n\n")
	b.WriteString("> Auto-generated. Do not edit manually.\n\n")

	// Categorize modules
	var routes, components, hooks, libs, others []*FrontendModule

	for _, m := range modules {
		switch {
		case m.IsRoute:
			routes = append(routes, m)
		case strings.Contains(m.Path, "/components/"):
			components = append(components, m)
		case strings.Contains(m.Path, "/hooks/"):
			hooks = append(hooks, m)
		case strings.Contains(m.Path, "/lib/") || strings.Contains(m.Path, "/utils/"):
			libs = append(libs, m)
		default:
			others = append(others, m)
		}
	}

	// Sort each category
	sortModules := func(ms []*FrontendModule) {
		sort.Slice(ms, func(i, j int) bool {
			return ms[i].Path < ms[j].Path
		})
	}

	sortModules(routes)
	sortModules(components)
	sortModules(hooks)
	sortModules(libs)
	sortModules(others)

	// Routes
	if len(routes) > 0 {
		b.WriteString("## Routes (app/)\n\n")
		for _, m := range routes {
			exportName := m.DefaultExport
			if exportName == "" {
				exportName = "(default)"
			}
			b.WriteString(fmt.Sprintf("- %s → %s\n", m.Path, exportName))
		}
		b.WriteString("\n")
	}

	// Components
	if len(components) > 0 {
		b.WriteString("## Components\n\n")
		currentDir := ""
		for _, m := range components {
			dir := filepath.Dir(m.Path)
			if dir != currentDir {
				b.WriteString(fmt.Sprintf("### %s\n\n", dir))
				currentDir = dir
			}
			writeModuleExports(&b, m)
		}
		b.WriteString("\n")
	}

	// Hooks
	if len(hooks) > 0 {
		b.WriteString("## Hooks\n\n")
		for _, m := range hooks {
			writeModuleExports(&b, m)
		}
		b.WriteString("\n")
	}

	// Utilities
	if len(libs) > 0 {
		b.WriteString("## Utilities (lib/)\n\n")
		for _, m := range libs {
			writeModuleExports(&b, m)
		}
		b.WriteString("\n")
	}

	// Others
	if len(others) > 0 {
		b.WriteString("## Other\n\n")
		for _, m := range others {
			writeModuleExports(&b, m)
		}
		b.WriteString("\n")
	}

	outputPath := filepath.Join(rootDir, "frontend", "REPO_MAP.md")
	return os.WriteFile(outputPath, []byte(b.String()), 0644)
}

func writeModuleExports(b *strings.Builder, m *FrontendModule) {
	if m.DefaultExport != "" {
		b.WriteString(fmt.Sprintf("- %s (default: %s)\n", m.Path, m.DefaultExport))
	} else if len(m.NamedExports) > 0 {
		b.WriteString(fmt.Sprintf("- %s\n", m.Path))
		for _, e := range m.NamedExports {
			b.WriteString(fmt.Sprintf("  - %s\n", e))
		}
	}
}

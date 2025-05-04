package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type Issue struct {
	RuleID      string
	Description string
	Severity    string
	Position    token.Position
}

type AnalyzerConfig struct {
	IncludeTests bool
}

type Analyzer struct {
	config AnalyzerConfig
	fset   *token.FileSet
}

func NewAnalyzer(config AnalyzerConfig) *Analyzer {
	return &Analyzer{
		config: config,
		fset:   token.NewFileSet(),
	}
}

func (a *Analyzer) ParseFile(filepath string, src []byte) (*ast.File, error) {
	mode := parser.ParseComments
	if a.config.IncludeTests {
		mode |= parser.AllErrors
	}
	
	file, err := parser.ParseFile(a.fset, filepath, src, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}
	
	return file, nil
}

func (a *Analyzer) ParseString(filename, content string) (*ast.File, error) {
	return parser.ParseFile(token.NewFileSet(), filename, content, parser.AllErrors)
}

func (a *Analyzer) GetPositionOf(node ast.Node) token.Position {
	return a.fset.Position(node.Pos())
}

func (a *Analyzer) AnalyzeErrorHandling(file *ast.File) []Issue {
	var issues []Issue
	
	ast.Inspect(file, func(n ast.Node) bool {
		if node, ok := n.(*ast.AssignStmt); ok {
			// Only check statements that might be assigning an error
			hasErrorVar := false
			for _, lhs := range node.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok && ident.Name == "err" {
					hasErrorVar = true
					break
				}
			}
			
			if hasErrorVar {
				// Find the nearest containing block statement
				var containingBlock *ast.BlockStmt
				parent := a.findParentNode(file, node)
				for parent != nil {
					if block, ok := parent.(*ast.BlockStmt); ok {
						containingBlock = block
						break
					}
					parent = a.findParentNode(file, parent)
				}
				
				// If we found a containing block, check if it contains error handling
				if containingBlock != nil {
					if !a.blockContainsErrorCheck(containingBlock, node) {
						issues = append(issues, Issue{
							RuleID:      "error_handling",
							Description: "Missing error check after error assignment",
							Severity:    "warning",
							Position:    a.GetPositionOf(node),
						})
					}
				}
			}
			
			if a.isErrorIgnored(node) {
				issues = append(issues, Issue{
					RuleID:      "error_handling",
					Description: "Error is being ignored with underscore assignment",
					Severity:    "error",
					Position:    a.GetPositionOf(node),
				})
			}
		}
		return true
	})
	
	return issues
}

func (a *Analyzer) findParentNode(file *ast.File, target ast.Node) ast.Node {
	var parent ast.Node
	
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		
		// Check if any child nodes match our target
		switch x := n.(type) {
		case *ast.BlockStmt:
			for _, stmt := range x.List {
				if stmt == target {
					parent = n
					return false
				}
			}
		case *ast.IfStmt:
			if x.Cond == target || x.Body == target || x.Else == target {
				parent = n
				return false
			}
		case *ast.AssignStmt:
			for _, expr := range x.Lhs {
				if expr == target {
					parent = n
					return false
				}
			}
			for _, expr := range x.Rhs {
				if expr == target {
					parent = n
					return false
				}
			}
		}
		
		return true
	})
	
	return parent
}

func (a *Analyzer) blockContainsErrorCheck(block *ast.BlockStmt, errAssign ast.Node) bool {
	foundAssignment := false
	
	// Iterate through statements to find error checks that come after our assignment
	for _, stmt := range block.List {
		if !foundAssignment {
			if stmt == errAssign {
				foundAssignment = true
			}
			continue
		}
		
		// Check if this statement is an error check
		if ifStmt, ok := stmt.(*ast.IfStmt); ok {
			if binExpr, ok := ifStmt.Cond.(*ast.BinaryExpr); ok {
				if binExpr.Op == token.NEQ {
					if ident, ok := binExpr.X.(*ast.Ident); ok && ident.Name == "err" {
						return true
					}
					if ident, ok := binExpr.Y.(*ast.Ident); ok && ident.Name == "err" {
						return true
					}
				}
			}
		}
	}
	
	return false
}

func (a *Analyzer) AnalyzeAPIDesign(file *ast.File) []Issue {
	var issues []Issue
	
	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			if a.isReceiverMethod(funcDecl) && !a.hasContextParameter(funcDecl) {
				issues = append(issues, Issue{
					RuleID:      "api_design",
					Description: "API method missing context.Context as first parameter",
					Severity:    "warning",
					Position:    a.GetPositionOf(funcDecl),
				})
			}
		}
		
		return true
	})
	
	return issues
}

func (a *Analyzer) AnalyzeConcurrencySafety(file *ast.File) []Issue {
	var issues []Issue
	
	// Find maps and their associated mutexes
	mapVars := a.findMapVariables(file)
	protectedMaps := make(map[string]bool)
	
	// Find mutex-protected maps
	ast.Inspect(file, func(n ast.Node) bool {
		if structType, ok := n.(*ast.StructType); ok {
			mapFields := make(map[string]bool)
			hasMutex := false
			
			// Check if struct has both a map and a mutex
			for _, field := range structType.Fields.List {
				if _, ok := field.Type.(*ast.MapType); ok {
					// Capture map field names
					for _, name := range field.Names {
						mapFields[name.Name] = true
					}
				}
				
				if fieldType, ok := field.Type.(*ast.SelectorExpr); ok {
					if x, ok := fieldType.X.(*ast.Ident); ok && x.Name == "sync" {
						if fieldType.Sel.Name == "Mutex" || fieldType.Sel.Name == "RWMutex" {
							hasMutex = true
						}
					}
				}
			}
			
			// If struct has both map and mutex, consider the maps protected
			if hasMutex {
				for mapName := range mapFields {
					protectedMaps[mapName] = true
				}
			}
		}
		return true
	})
	
	// Check for concurrent map access outside of mutex lock/unlock
	ast.Inspect(file, func(n ast.Node) bool {
		if goStmt, ok := n.(*ast.GoStmt); ok {
			// Track map accesses and mutex locks within goroutines
			mapAccesses := make(map[string]ast.Node)
			mutexLocked := make(map[string]bool)
			
			ast.Inspect(goStmt.Call.Fun, func(innerNode ast.Node) bool {
				// Track mutex lock operations
				if callExpr, ok := innerNode.(*ast.CallExpr); ok {
					if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
						if ident, ok := selExpr.X.(*ast.Ident); ok {
							if selExpr.Sel.Name == "Lock" || selExpr.Sel.Name == "RLock" {
								mutexLocked[ident.Name] = true
							}
						}
					}
				}
				
				// Track map access
				if indexExpr, ok := innerNode.(*ast.IndexExpr); ok {
					if ident, ok := indexExpr.X.(*ast.Ident); ok {
						if mapVars[ident.Name] {
							mapAccesses[ident.Name] = indexExpr
						}
					}
				}
				
				return true
			})
			
			// Check if map accesses are protected
			for mapName, accessNode := range mapAccesses {
				if !mutexLocked["mu"] && !mutexLocked["mutex"] && !mutexLocked["rwmu"] && 
				   !mutexLocked["rwmutex"] && !mutexLocked["lock"] && !protectedMaps[mapName] {
					issues = append(issues, Issue{
						RuleID:      "concurrent_map_access",
						Description: fmt.Sprintf("Map '%s' accessed in goroutine without mutex protection", mapName),
						Severity:    "error",
						Position:    a.GetPositionOf(accessNode),
					})
				}
			}
		}
		return true
	})
	
	return issues
}

func (a *Analyzer) AnalyzeSecurityIssues(file *ast.File) []Issue {
	var issues []Issue
	
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ImportSpec:
			importPath := strings.Trim(node.Path.Value, "\"")
			if importPath == "crypto/md5" || importPath == "crypto/sha1" {
				issues = append(issues, Issue{
					RuleID:      "secure_coding",
					Description: "Using weak cryptographic algorithm: " + importPath,
					Severity:    "error",
					Position:    a.GetPositionOf(node),
				})
			}
		case *ast.CallExpr:
			if a.isSQLInjectionVulnerable(node) {
				issues = append(issues, Issue{
					RuleID:      "secure_coding",
					Description: "Potential SQL injection vulnerability detected",
					Severity:    "error",
					Position:    a.GetPositionOf(node),
				})
			}
		case *ast.ValueSpec:
			if a.hasHardcodedCredentials(node) {
				issues = append(issues, Issue{
					RuleID:      "secure_coding",
					Description: "Hardcoded credentials detected",
					Severity:    "error",
					Position:    a.GetPositionOf(node),
				})
			}
		}
		
		return true
	})
	
	return issues
}

func (a *Analyzer) AnalyzeOrganizationStandards(file *ast.File) []Issue {
	var issues []Issue
	
	// Check for global variables
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if genDecl.Tok == token.VAR {
				for _, spec := range genDecl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for _, name := range valueSpec.Names {
							if ast.IsExported(name.Name) {
								issues = append(issues, Issue{
									RuleID:      "org_coding_standards",
									Description: fmt.Sprintf("Global variable '%s' violates organizational standards", name.Name),
									Severity:    "warning",
									Position:    a.GetPositionOf(name),
								})
							}
						}
					}
				}
			}
		}
	}
	
	// Check for snake_case function names
	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			if strings.Contains(funcDecl.Name.Name, "_") {
				issues = append(issues, Issue{
					RuleID:      "org_coding_standards",
					Description: fmt.Sprintf("Function '%s' uses snake_case which violates naming conventions", funcDecl.Name.Name),
					Severity:    "warning",
					Position:    a.GetPositionOf(funcDecl.Name),
				})
			}
		}
		return true
	})
	
	return issues
}

func (a *Analyzer) hasErrorAssignment(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.AssignStmt:
		for _, rhs := range n.Rhs {
			if _, ok := rhs.(*ast.CallExpr); ok {
				// Check if one of the left-hand variables is named "err"
				for _, lhs := range n.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok && ident.Name == "err" {
						return true
					}
				}
			}
		}
		// Check for multiple return values where one is an error
		for _, lhs := range n.Lhs {
			if ident, ok := lhs.(*ast.Ident); ok && ident.Name == "err" {
				return true
			}
		}
	}
	return false
}

func (a *Analyzer) hasErrorCheck(node ast.Node) bool {
	var hasCheck bool
	ast.Inspect(node, func(n ast.Node) bool {
		if ifStmt, ok := n.(*ast.IfStmt); ok {
			if binExpr, ok := ifStmt.Cond.(*ast.BinaryExpr); ok {
				if binExpr.Op == token.NEQ {
					if ident, ok := binExpr.X.(*ast.Ident); ok && ident.Name == "err" {
						hasCheck = true
						return false
					}
					if ident, ok := binExpr.Y.(*ast.Ident); ok && ident.Name == "err" {
						hasCheck = true
						return false
					}
				}
			}
		}
		return true
	})
	return hasCheck
}

func (a *Analyzer) isErrorIgnored(node *ast.AssignStmt) bool {
	for _, lhs := range node.Lhs {
		if ident, ok := lhs.(*ast.Ident); ok {
			if ident.Name == "_" {
				// Check if the right side might be an error
				// This is also a simplification
				for _, rhs := range node.Rhs {
					if callExpr, ok := rhs.(*ast.CallExpr); ok {
						// Check if function is likely to return an error (ends with Err or Error)
						// This is a heuristic and not 100% accurate
						if funcName := a.getFunctionName(callExpr); strings.Contains(funcName, "Error") {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func (a *Analyzer) getFunctionName(call *ast.CallExpr) string {
	switch expr := call.Fun.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.SelectorExpr:
		if ident, ok := expr.X.(*ast.Ident); ok {
			return ident.Name + "." + expr.Sel.Name
		}
	}
	return ""
}

func (a *Analyzer) isReceiverMethod(node *ast.FuncDecl) bool {
	return node.Recv != nil && len(node.Recv.List) > 0
}

func (a *Analyzer) hasContextParameter(node *ast.FuncDecl) bool {
	if node.Type.Params == nil || len(node.Type.Params.List) == 0 {
		return false
	}
	
	firstParam := node.Type.Params.List[0]
	
	// Check if param type is context.Context
	if sel, ok := firstParam.Type.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "context" && sel.Sel.Name == "Context"
		}
	}
	return false
}

func (a *Analyzer) findMapVariables(file *ast.File) map[string]bool {
	mapVars := make(map[string]bool)
	
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GenDecl:
			for _, spec := range node.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					if valueSpec.Type != nil {
						if _, ok := valueSpec.Type.(*ast.MapType); ok {
							for _, name := range valueSpec.Names {
								mapVars[name.Name] = true
							}
						}
					}
				}
			}
		case *ast.AssignStmt:
			for i, lhs := range node.Lhs {
				if i < len(node.Rhs) {
					if ident, ok := lhs.(*ast.Ident); ok {
						if callExpr, ok := node.Rhs[i].(*ast.CallExpr); ok {
							if len(callExpr.Args) > 0 {
								if _, ok := callExpr.Args[0].(*ast.MapType); ok {
									mapVars[ident.Name] = true
								}
							}
						}
					}
				}
			}
		}
		return true
	})
	
	return mapVars
}

func (a *Analyzer) hasMutexProtection(file *ast.File, mapVar string) bool {
	// This is a simplification. In a real implementation, we would need to analyze
	// the synchronization patterns more carefully.
	hasMutex := false
	
	ast.Inspect(file, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok {
			if genDecl.Tok == token.TYPE {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if structType, ok := typeSpec.Type.(*ast.StructType); ok {
							for _, field := range structType.Fields.List {
								if sel, ok := field.Type.(*ast.SelectorExpr); ok {
									if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "sync" {
										if sel.Sel.Name == "Mutex" || sel.Sel.Name == "RWMutex" {
											hasMutex = true
										}
									}
								}
							}
						}
					}
				}
			}
		}
		return true
	})
	
	return hasMutex
}

func (a *Analyzer) isSQLInjectionVulnerable(node *ast.CallExpr) bool {
	if sel, ok := node.Fun.(*ast.SelectorExpr); ok {
		if fmt, ok := sel.X.(*ast.Ident); ok && fmt.Name == "fmt" && sel.Sel.Name == "Sprintf" {
			// Check if first argument contains SQL syntax
			if len(node.Args) > 0 {
				if lit, ok := node.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
					sqlStr := lit.Value
					return strings.Contains(strings.ToUpper(sqlStr), "SELECT") ||
						strings.Contains(strings.ToUpper(sqlStr), "INSERT") ||
						strings.Contains(strings.ToUpper(sqlStr), "UPDATE") ||
						strings.Contains(strings.ToUpper(sqlStr), "DELETE") ||
						strings.Contains(strings.ToUpper(sqlStr), "WHERE")
				}
			}
		}
	}
	return false
}

func (a *Analyzer) hasHardcodedCredentials(node *ast.ValueSpec) bool {
	sensitiveNames := []string{"password", "secret", "key", "token", "auth", "credential"}
	
	for _, name := range node.Names {
		nameStr := strings.ToLower(name.Name)
		for _, sensitive := range sensitiveNames {
			if strings.Contains(nameStr, sensitive) {
				// Check if value is a string literal
				for _, value := range node.Values {
					if lit, ok := value.(*ast.BasicLit); ok && lit.Kind == token.STRING {
						// String literal assigned to a sensitive variable name
						return true
					}
				}
			}
		}
	}
	
	return false
}

func (a *Analyzer) detectConcurrentMapAccess(file *ast.File) []Issue {
	var issues []Issue
	
	// Track maps declared in the file
	mapDeclarations := make(map[string]bool)
	
	// Find all map declarations
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GenDecl:
			if node.Tok == token.VAR || node.Tok == token.CONST {
				for _, spec := range node.Specs {
					if valSpec, ok := spec.(*ast.ValueSpec); ok {
						for i, name := range valSpec.Names {
							if i < len(valSpec.Values) {
								if _, ok := valSpec.Type.(*ast.MapType); ok {
									mapDeclarations[name.Name] = true
								}
							}
						}
					}
				}
			}
		case *ast.AssignStmt:
			for i, lhs := range node.Lhs {
				if i < len(node.Rhs) {
					if ident, ok := lhs.(*ast.Ident); ok {
						if callExpr, ok := node.Rhs[i].(*ast.CallExpr); ok {
							if funcIdent, ok := callExpr.Fun.(*ast.Ident); ok && funcIdent.Name == "make" {
								if len(callExpr.Args) > 0 {
									if _, ok := callExpr.Args[0].(*ast.MapType); ok {
										mapDeclarations[ident.Name] = true
									}
								}
							}
						}
					}
				}
			}
		}
		return true
	})
	
	// Check for concurrent map access
	ast.Inspect(file, func(n ast.Node) bool {
		if goStmt, ok := n.(*ast.GoStmt); ok {
			ast.Inspect(goStmt.Call.Fun, func(inner ast.Node) bool {
				if assignStmt, ok := inner.(*ast.AssignStmt); ok {
					for _, rhs := range assignStmt.Rhs {
						if indexExpr, ok := rhs.(*ast.IndexExpr); ok {
							if ident, ok := indexExpr.X.(*ast.Ident); ok {
								if mapDeclarations[ident.Name] {
									issues = append(issues, Issue{
										RuleID:      "concurrent_map_access",
										Description: "Concurrent map read without synchronization",
										Severity:    "error",
										Position:    a.GetPositionOf(indexExpr),
									})
								}
							}
						}
					}
				} else if indexExpr, ok := inner.(*ast.IndexExpr); ok {
					if ident, ok := indexExpr.X.(*ast.Ident); ok {
						if mapDeclarations[ident.Name] {
							issues = append(issues, Issue{
								RuleID:      "concurrent_map_access",
								Description: "Concurrent map access without synchronization",
								Severity:    "error",
								Position:    a.GetPositionOf(indexExpr),
							})
						}
					}
				}
				return true
			})
		}
		return true
	})
	
	return issues
}

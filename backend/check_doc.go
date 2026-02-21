package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// DocCommentCheckResult はドキュメントコメントのチェック結果を保持する
type DocCommentCheckResult struct {
	FileName      string
	HasPackageDoc bool
	ExportedDecls []ExportedDeclaration
	Issues        []string
}

// ExportedDeclaration エクスポートされた宣言情報
type ExportedDeclaration struct {
	Name    string
	Type    string // function, type, const, var
	Line    int
	Comment string
	HasDoc  bool
}

// checkDocComments は指定されたディレクトリ内のGoファイルのドキュメントコメントをチェックする
func checkDocComments(dir string) ([]DocCommentCheckResult, error) {
	fset := token.NewFileSet()
	var results []DocCommentCheckResult

	// ディレクトリ内の全Goファイルを再帰的に探索
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Goファイルのみ処理
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			// testデータを除外
			if strings.Contains(path, "test_data") {
				return nil
			}

			result, err := checkSingleFile(path, fset)
			if err != nil {
				return fmt.Errorf("failed to check %s: %v", path, err)
			}

			results = append(results, result)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

// checkSingleFile は単一のGoファイルをチェックする
func checkSingleFile(filePath string, fset *token.FileSet) (DocCommentCheckResult, error) {
	var result DocCommentCheckResult
	result.FileName = filePath

	// ファイルをパース
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return result, err
	}

	// パッケージコメントのチェック
	result.HasPackageDoc = node.Doc != nil && node.Doc.Text() != ""

	// エクスポートされた宣言を収集
	ast.Inspect(node, func(n ast.Node) bool {
		switch decl := n.(type) {
		case *ast.FuncDecl:
			if decl.Name.IsExported() {
				exported := ExportedDeclaration{
					Name: decl.Name.Name,
					Type: "function",
					Line: fset.Position(decl.Pos()).Line,
					Comment: "",
					HasDoc: decl.Doc != nil,
				}
				if decl.Doc != nil {
					exported.Comment = decl.Doc.Text()
				}
				result.ExportedDecls = append(result.ExportedDecls, exported)
			}
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if s.Name.IsExported() {
						exported := ExportedDeclaration{
							Name: s.Name.Name,
							Type: "type",
							Line: fset.Position(decl.Pos()).Line,
							Comment: "",
							HasDoc: decl.Doc != nil,
						}
						if decl.Doc != nil {
							exported.Comment = decl.Doc.Text()
						}
						result.ExportedDecls = append(result.ExportedDecls, exported)
					}
				case *ast.ValueSpec:
					for _, name := range s.Names {
						if name.IsExported() {
							exported := ExportedDeclaration{
								Name: name.Name,
								Type: getDeclarationType(decl),
								Line: fset.Position(decl.Pos()).Line,
								Comment: "",
								HasDoc: decl.Doc != nil,
							}
							if decl.Doc != nil {
								exported.Comment = decl.Doc.Text()
							}
							result.ExportedDecls = append(result.ExportedDecls, exported)
						}
					}
				}
			}
		}
		return true
	})

	// 問題点をチェック
	checkIssues(&result)

	return result, nil
}

// getDeclarationType は宣言の種類を返す
func getDeclarationType(decl *ast.GenDecl) string {
	switch decl.Tok {
	case token.CONST:
		return "const"
	case token.VAR:
		return "var"
	default:
		return "other"
	}
}

// checkIssues はドキュメントコメントの問題点をチェックする
func checkIssues(result *DocCommentCheckResult) {
	// パッケージコメントがない場合の問題点
	if !result.HasPackageDoc && result.FileName != "doc.go" {
		result.Issues = append(result.Issues, "パッケージコメントがありません")
	}

	// エクスポートされた宣言のドキュメントコメントをチェック
	for _, decl := range result.ExportedDecls {
		if !decl.HasDoc {
			result.Issues = append(result.Issues,
				fmt.Sprintf("エクスポートされた%s `%s` (行 %d) にドキュメントコメントがありません",
					decl.Type, decl.Name, decl.Line))
			continue
		}

		// コメントの内容をチェック
		comment := strings.TrimSpace(decl.Comment)
		if comment == "" {
			result.Issues = append(result.Issues,
				fmt.Sprintf("エクスポートされた%s `%s` (行 %d) のドキュメントコメントが空です",
					decl.Type, decl.Name, decl.Line))
			continue
		}

		// WhatよりもWhyの原則チェック
		if isWhatInsteadOfWhy(comment, decl.Type) {
			result.Issues = append(result.Issues,
				fmt.Sprintf("エクスポートされた%s `%s` (行 %d) のコメントが「What」中心です。「Why」を中心に記述してください",
					decl.Type, decl.Name, decl.Line))
		}

		// 完全な文になっていないチェック
		if !isCompleteSentence(comment) {
			result.Issues = append(result.Issues,
				fmt.Sprintf("エクスポートされた%s `%s` (行 %d) のコメントが完全な文ではありません",
					decl.Type, decl.Name, decl.Line))
		}
	}
}

// isWhatInsteadOfWhy コメントが「What」中心かチェック
func isWhatInsteadOfWhy(comment, declType string) bool {
	whatKeywords := []string{
		"を取得する", "を設定する", "を追加する", "を削除する", "を更新する",
		"Get", "Set", "Add", "Delete", "Update", "Create", "Do",
	}

	comment = strings.ToLower(comment)
	for _, keyword := range whatKeywords {
		if strings.Contains(comment, keyword) {
			return true
		}
	}
	return false
}

// isCompleteSentence 文章が完全な文かチェック
func isCompleteSentence(comment string) bool {
	if len(comment) == 0 {
		return false
	}

	// 文末がピリオドで終わっているか
	if comment[len(comment)-1] != '.' {
		return false
	}

	// 完全な文の基本チェック（簡略化）
	if len(comment) < 10 {
		return false
	}

	return true
}

// calculateCoverage ドキュメントコメントの網羅率を計算する
func calculateCoverage(results []DocCommentCheckResult) (float64, int, int) {
	totalDecls := 0
	documentedDecls := 0

	for _, result := range results {
		totalDecls += len(result.ExportedDecls)
		for _, decl := range result.ExportedDecls {
			if decl.HasDoc && isCompleteSentence(decl.Comment) {
				documentedDecls++
			}
		}
	}

	if totalDecls == 0 {
		return 0, 0, 0
	}

	coverage := float64(documentedDecls) / float64(totalDecls) * 100
	return coverage, documentedDecls, totalDecls
}

// printResults チェック結果を整形して出力する
func printResults(results []DocCommentCheckResult) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("ドキュメントコメントチェック結果")
	fmt.Println(strings.Repeat("=", 80))

	// 網羅率の計算と表示
	coverage, documented, total := calculateCoverage(results)
	fmt.Printf("\n【全体の網羅率】\n")
	fmt.Printf("ドキュメントコメント網羅率: %.1f%% (%d/%d)\n\n", coverage, documented, total)

	// ファイルごとの結果
	fmt.Println("\n【ファイル別の結果】")

	completeFiles := []string{}
	incompleteFiles := []string{}

	for _, result := range results {
		relPath, _ := filepath.Rel(".", result.FileName)
		fmt.Printf("\n📄 %s\n", relPath)
		fmt.Printf("   パッケージコメント: %s\n", formatStatus(result.HasPackageDoc))

		if len(result.Issues) == 0 {
			fmt.Printf("   ✅ ドキュメントコメントは完全です\n")
			completeFiles = append(completeFiles, relPath)
		} else {
			fmt.Printf("   ❌ 問題点:\n")
			for _, issue := range result.Issues {
				fmt.Printf("      - %s\n", issue)
			}
			incompleteFiles = append(incompleteFiles, relPath)
		}

		// エクスポートされた宣言の詳細
		if len(result.ExportedDecls) > 0 {
			fmt.Printf("\n   エクスポートされた宣言:\n")
			for _, decl := range result.ExportedDecls {
				status := "✅"
				if !decl.HasDoc || !isCompleteSentence(decl.Comment) {
					status = "❌"
				}
				fmt.Printf("      %s %s `%s` (行 %d)\n", status, decl.Type, decl.Name, decl.Line)
			}
		}
		fmt.Println()
	}

	// サマリー
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("サマリー")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("合計ファイル数: %d\n", len(results))
	fmt.Printf("完全なドキュメントコメントのファイル: %d\n", len(completeFiles))
	fmt.Printf("ドキュメントコメント不足のファイル: %d\n", len(incompleteFiles))

	if len(completeFiles) > 0 {
		fmt.Println("\n完全なドキュメントコメントのファイル一覧:")
		for _, file := range completeFiles {
			fmt.Printf("  - %s\n", file)
		}
	}

	if len(incompleteFiles) > 0 {
		fmt.Println("\nドキュメントコメント不足のファイル一覧:")
		for _, file := range incompleteFiles {
			fmt.Printf("  - %s\n", file)
		}
	}
}

// formatStatus 状態を文字列に変換
func formatStatus(status bool) string {
	if status {
		return "✅ あり"
	}
	return "❌ なし"
}

func main() {
	// ヘルプオプションの処理
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		printHelp()
		return
	}

	// カレントディレクトリをbackendに変更（必要に応じて）
	if _, err := os.Stat("go.mod"); err != nil {
		fmt.Println("エラー: backendディレクトリで実行してください")
		os.Exit(1)
	}

	// ドキュメントコメントチェックを実行
	fmt.Println("backendディレクトリのドキュメントコメントをチェック中...")
	results, err := checkDocComments(".")
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	printResults(results)

	// 戻り値に基づいて終了コードを返す
	hasIssues := false
	for _, result := range results {
		if len(result.Issues) > 0 {
			hasIssues = true
			break
		}
	}

	if hasIssues {
		fmt.Println("\n警告: ドキュメントコメントが不完全なファイルがあります")
		os.Exit(1)
	}

	fmt.Println("\n✅ すべてのファイルでドキュメントコメントが適切です")
}

// printHelp ヘルプメッセージを出力する
func printHelp() {
	fmt.Println("ドキュメントコメントチェッカー")
	fmt.Println()
	fmt.Println("使い方:")
	fmt.Println("  go run check_doc.go")
	fmt.Println("  make check-doc")
	fmt.Println()
	fmt.Println("機能:")
	fmt.Println("  - backendディレクトリ内のすべての.goファイルをスキャン")
	fmt.Println("  - エクスポートされる関数・型・定数がgodoc形式のコメントを持っているか確認")
	fmt.Println("  - .claude/rules/backend.mdのルールに従ったコメント形式かチェック")
	fmt.Println("  - 不足しているファイルと具体的な行番号を出力")
	fmt.Println("  - 全体のドキュメントコメントの網羅率を計算")
	fmt.Println()
	fmt.Println("ルール:")
	fmt.Println("  - エクスポートされる全ての関数・メソッドにドキュメントコメント必須")
	fmt.Println("  - エクスポートされる型・構造体・インターフェースにドキュメントコメント必須")
	fmt.Println("  - 定数（パッケージ外から参照されるもの）にドキュメントコメント必須")
	fmt.Println("  - WhatよりもWhyを中心に記述")
	fmt.Println("  - 完全な文で記述（文末にピリオド）")
	fmt.Println()
	fmt.Println("オプション:")
	fmt.Println("  -h, --help  このヘルプメッセージを表示")
	fmt.Println()
	fmt.Println("終了コード:")
	fmt.Println("  0 - すべてのファイルでドキュメントコメントが適切")
	fmt.Println("  1 - ドキュメントコメントが不完全なファイルがある")
}
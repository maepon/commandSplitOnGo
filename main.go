package main

import (
	"flag"
	"fmt"
)

func main() {
	linesPerFile := flag.Int("l", 1000, "number of lines per output file")
	numberOfFiles := flag.Int("n", 0, "number of files to split into")  // 追加
	bytesPerFile := flag.Int("b", 0, "number of bytes per output file") // -bオプションの追加
	flag.Parse()

	if *linesPerFile <= 0 && *numberOfFiles <= 0 && *bytesPerFile <= 0 {
		fmt.Println("You must specify one of -l, -n, or -b options with a positive value.")
		return
	}

	// オプションが指定された場合、デフォルト値を0に設定
	specifiedOptions := 0
	if *linesPerFile != 1000 {
		specifiedOptions++
	}
	if *numberOfFiles != 0 {
		specifiedOptions++
	}
	if *bytesPerFile != 0 {
		specifiedOptions++
	}

	// 併用のチェック
	if specifiedOptions > 1 {
		fmt.Println("Error: You cannot specify more than one of -l, -n, and -b options simultaneously.")
		return
	}

	if *numberOfFiles != 0 {
		*linesPerFile = 0 // -nが指定された場合、-lを無効にする
	}

	if *bytesPerFile != 0 {
		*linesPerFile = 0 // -bが指定された場合、-lを無効にする
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: split [-l lines] [-n parts] [-b bytes] <file>") // 更新
		return
	}
	fileName := args[0]

	if err := SplitFile(fileName, *linesPerFile, *numberOfFiles, *bytesPerFile); err != nil {
		fmt.Println("Error splitting file:", err)
	}
}

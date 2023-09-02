package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func createTestFile(fileName string, lines int) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Failed to close file:", err)
		}
	}(file)

	for i := 0; i < lines; i++ {
		_, err := fmt.Fprintln(file, "Line number:", i)
		if err != nil {
			return err
		}
	}

	return nil
}

func countLinesInFile(fileName string) (int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Failed to close file:", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	return lineCount, scanner.Err()
}

func countBytesInFile(fileName string) (int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Failed to close file:", err)
		}
	}(file)

	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return int(stat.Size()), nil
}

func TestSplitFile_NoOptions(t *testing.T) {
	// テストファイル名
	testFileName := "testfile.txt"

	// テストファイルの生成
	totalLines := 5000
	err := createTestFile(testFileName, totalLines)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
		return
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println("Failed to remove test file:", err)
		}
	}(testFileName)

	// 分割テスト
	linesPerFile := 200 // 分割後の各ファイルの行数
	err = SplitFile(testFileName, linesPerFile, 0, 0)
	if err != nil {
		t.Errorf("Failed to split file: %v", err)
		return
	}

	// 分割後のファイルを検証
	expectedNumberOfFiles := totalLines / linesPerFile
	for i := 0; i < expectedNumberOfFiles; i++ {
		fileName := fmt.Sprintf("%s%c%c", testFileName, 'a'+i/26, 'a'+i%26)
		file, err := os.Open(fileName)
		if err != nil {
			t.Errorf("Failed to open split file %s: %v", fileName, err)
			return
		}

		scanner := bufio.NewScanner(file)
		lineCount := 0
		for scanner.Scan() {
			line := scanner.Text()
			expectedLine := strconv.Itoa(i*linesPerFile + lineCount)
			if !strings.Contains(line, expectedLine) {
				t.Errorf("Unexpected content in file %s at line %d: got %s, want %s", fileName, lineCount, line, expectedLine)
			}
			lineCount++
		}
		if lineCount != linesPerFile {
			t.Errorf("Unexpected number of lines in file %s: got %d, want %d", fileName, lineCount, linesPerFile)
		}

		err = file.Close()
		if err != nil {
			t.Errorf("Failed to close file %s: %v", fileName, err)
		}
		err = os.Remove(fileName)
		if err != nil {
			fmt.Println("Failed to remove test file:", err)
		} // 必要に応じて分割ファイルを削除
	}
}

func TestSplitFile_LinesOption(t *testing.T) {
	// テストファイル名
	testFileName := "testfile_lines_option.txt"

	// テストファイルの生成
	totalLines := 5000
	err := createTestFile(testFileName, totalLines)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
		return
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println("Failed to remove test file:", err)
		}
	}(testFileName)

	// 分割テスト
	linesPerFile := 800 // 分割したい行数
	err = SplitFile(testFileName, linesPerFile, 0, 0)
	if err != nil {
		t.Errorf("Failed to split file: %v", err)
		return
	}

	// 分割後のファイルを検証
	for i := 0; i < totalLines/linesPerFile; i++ {
		outputFileName := fmt.Sprintf("%s%c%c", testFileName, 'a'+i/26, 'a'+i%26)
		lines, err := countLinesInFile(outputFileName)
		if err != nil {
			t.Errorf("Failed to read output file %s: %v", outputFileName, err)
			return
		}
		if lines != linesPerFile {
			t.Errorf("Expected %d lines in file %s, but got %d", linesPerFile, outputFileName, lines)
		}
		// 必要であれば、ファイルを削除
		err = os.Remove(outputFileName)
		if err != nil {
			return
		}
	}

	// 余りの行の検証
	if totalLines%linesPerFile != 0 {
		outputFileName := fmt.Sprintf("%s%c%c", testFileName, 'a'+(totalLines/linesPerFile)/26, 'a'+(totalLines/linesPerFile)%26)
		lines, err := countLinesInFile(outputFileName)
		if err != nil {
			t.Errorf("Failed to read output file %s: %v", outputFileName, err)
			return
		}
		if lines != totalLines%linesPerFile {
			t.Errorf("Expected %d lines in file %s, but got %d", totalLines%linesPerFile, outputFileName, lines)
		}
		// 必要であれば、ファイルを削除
		err = os.Remove(outputFileName)
		if err != nil {
			return
		}
	}
}

func TestSplitFile_NumberOfFilesOption(t *testing.T) {
	// テストファイル名
	testFileName := "testfile_number_of_files_option.txt"

	// テストファイルの生成
	totalLines := 4927
	err := createTestFile(testFileName, totalLines)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
		return
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println("Failed to remove test file:", err)
		}
	}(testFileName)

	// 分割テスト
	numberOfFiles := 6 // 分割したいファイル数
	err = SplitFile(testFileName, 0, numberOfFiles, 0)
	if err != nil {
		t.Errorf("Failed to split file: %v", err)
		return
	}

	linesPerFile := totalLines / numberOfFiles

	// 分割後のファイルを検証
	for i := 0; i < numberOfFiles; i++ {
		outputFileName := fmt.Sprintf("%s%c%c", testFileName, 'a'+i/26, 'a'+i%26)
		lines, err := countLinesInFile(outputFileName)
		if err != nil {
			t.Errorf("Failed to read output file %s: %v", outputFileName, err)
			return
		}
		expectedLines := linesPerFile
		if i == numberOfFiles-1 && totalLines%numberOfFiles != 0 {
			expectedLines += totalLines % numberOfFiles
		}
		if lines != expectedLines {
			t.Errorf("Expected %d lines in file %s, but got %d", expectedLines, outputFileName, lines)
		}
		// 必要であれば、ファイルを削除
		err = os.Remove(outputFileName)
		if err != nil {
			return
		}
	}
}

func TestSplitFile_BytesOption(t *testing.T) {
	// テストファイル名
	testFileName := "testfile_bytes_option.txt"

	// テストファイルの生成
	totalLines := 1000
	err := createTestFile(testFileName, totalLines)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
		return
	}
	defer os.Remove(testFileName)

	// 分割テスト
	bytesPerFile := 800 // 分割したいバイト数
	err = SplitFile(testFileName, 0, 0, bytesPerFile)
	if err != nil {
		t.Errorf("Failed to split file: %v", err)
		return
	}

	// 分割後のファイルを検証
	for i := 0; ; i++ {
		outputFileName := fmt.Sprintf("%s%c%c", testFileName, 'a'+i/26, 'a'+i%26)
		bytes, err := countBytesInFile(outputFileName)
		if err != nil {
			if os.IsNotExist(err) {
				break // すべてのファイルを処理した
			}
			t.Errorf("Failed to read output file %s: %v", outputFileName, err)
			return
		}
		if i < totalLines*len(strconv.Itoa(totalLines))/bytesPerFile && bytes != bytesPerFile {
			t.Errorf("Expected %d bytes in file %s, but got %d", bytesPerFile, outputFileName, bytes)
		}
		// 必要であれば、ファイルを削除
		os.Remove(outputFileName)
	}
}

func TestSplitFile_EmptyFile(t *testing.T) {
	testFileName := "testfile_empty.txt"
	file, err := os.Create(testFileName)
	if err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
		return
	}
	file.Close()
	defer os.Remove(testFileName)

	err = SplitFile(testFileName, 10, 0, 0)
	if err != nil {
		t.Errorf("Failed to split empty file: %v", err)
		return
	}

	// 空のファイルの場合、分割ファイルは生成されないはず
	outputFileName := fmt.Sprintf("%s%c%c", testFileName, 'a', 'a')
	if _, err := os.Stat(outputFileName); !os.IsNotExist(err) {
		t.Errorf("Expected no split files, but found: %s", outputFileName)
	}
}

func TestSplitFile_OneLineFile(t *testing.T) {
	testFileName := "testfile_one_line.txt"
	err := createTestFile(testFileName, 1)
	if err != nil {
		t.Fatalf("Failed to create one line test file: %v", err)
		return
	}
	defer os.Remove(testFileName)

	err = SplitFile(testFileName, 10, 0, 0)
	if err != nil {
		t.Errorf("Failed to split one line file: %v", err)
		return
	}

	// 1行だけのファイルの場合、1つの分割ファイルが生成されるはず
	outputFileName := fmt.Sprintf("%s%c%c", testFileName, 'a', 'a')
	lines, err := countLinesInFile(outputFileName)
	if err != nil {
		t.Errorf("Failed to read output file %s: %v", outputFileName, err)
		return
	}
	if lines != 1 {
		t.Errorf("Expected 1 line in file %s, but got %d", outputFileName, lines)
	}
	// 必要であれば、ファイルを削除
	os.Remove(outputFileName)
}

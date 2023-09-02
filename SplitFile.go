package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func SplitFile(inputFileName string, linesPerFile, numberOfFiles, bytesPerFile int) error {
	if linesPerFile < 0 || numberOfFiles < 0 || bytesPerFile < 0 {
		return fmt.Errorf("linesPerFile, numberOfFiles, and bytesPerFile must be non-negative")
	}
	if bytesPerFile == 0 && linesPerFile == 0 && numberOfFiles == 0 {
		return fmt.Errorf("At least one of linesPerFile, numberOfFiles, or bytesPerFile must be greater than zero")
	}

	if bytesPerFile > 0 {
		return splitByBytes(inputFileName, bytesPerFile)
	} else {
		return splitByLines(inputFileName, linesPerFile, numberOfFiles)
	}
}

func createOutputFileName(inputFileName string, fileIndex int) string {
	return fmt.Sprintf("%s%c%c", inputFileName, 'a'+fileIndex/26, 'a'+fileIndex%26)
}

func splitByBytes(inputFileName string, bytesPerFile int) error {
	file, err := os.Open(inputFileName)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Failed to close file:", err)
		}
	}(file)

	const bufferSize = 4096 // 4KB buffer size
	buffer := make([]byte, bufferSize)

	var fileIndex, totalBytes int
	var outputFile *os.File

	var closeErr error

	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		writeBytes := 0
		for writeBytes < n {
			if outputFile == nil {
				outputFileName := createOutputFileName(inputFileName, fileIndex)
				outputFile, err = os.Create(outputFileName)
				if err != nil {
					return fmt.Errorf("Failed to create output file %s: %v", outputFileName, err)
				}
				totalBytes = 0
			}

			remainingBytes := n - writeBytes
			bytesToWrite := remainingBytes
			if totalBytes+remainingBytes > bytesPerFile {
				bytesToWrite = bytesPerFile - totalBytes
			}

			_, err = outputFile.Write(buffer[writeBytes : writeBytes+bytesToWrite])
			if err != nil {
				return err
			}

			writeBytes += bytesToWrite
			totalBytes += bytesToWrite

			if totalBytes == bytesPerFile {
				if err := outputFile.Close(); err != nil && closeErr == nil {
					closeErr = err
				}
				outputFile = nil
				fileIndex++
			}
		}
	}

	if outputFile != nil {
		if err := outputFile.Close(); err != nil && closeErr == nil {
			closeErr = err
		}
	}

	if closeErr != nil {
		return fmt.Errorf("Failed to close file: %v", closeErr) // Return the first close error
	}

	return nil
}

func splitByLines(inputFileName string, linesPerFile, numberOfFiles int) error {
	file, err := os.Open(inputFileName)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Failed to close file:", err)
		}
	}(file)

	totalLines := 0
	scanner := bufio.NewScanner(file)

	if numberOfFiles > 0 {
		for scanner.Scan() {
			totalLines++
		}
		_, err := file.Seek(0, 0)
		if err != nil {
			return err
		}
		scanner = bufio.NewScanner(file)
		linesPerFile = totalLines / numberOfFiles
	}

	var lineIndex, fileIndex int
	var outputFile *os.File
	var bufferedWriter *bufio.Writer

	var closeErr error

	for scanner.Scan() {
		if lineIndex%linesPerFile == 0 {
			if bufferedWriter != nil {
				if err := bufferedWriter.Flush(); err != nil {
					return err
				}
			}
			if outputFile != nil {
				if err := outputFile.Close(); err != nil && closeErr == nil {
					closeErr = err
				}
			}
			outputFileName := createOutputFileName(inputFileName, fileIndex)
			outputFile, err = os.Create(outputFileName) // エラーをチェックするために変更
			if err != nil {
				return fmt.Errorf("Failed to create output file %s: %v", outputFileName, err)
			}
			bufferedWriter = bufio.NewWriter(outputFile)
			fileIndex++
		}
		_, err := fmt.Fprintln(bufferedWriter, scanner.Text())
		if err != nil {
			return err
		}
		lineIndex++

		// Add remaining lines to the last file if the number of files is specified
		if numberOfFiles > 0 && fileIndex == numberOfFiles {
			for scanner.Scan() {
				_, err := fmt.Fprintln(bufferedWriter, scanner.Text())
				if err != nil {
					return err
				}
			}
		}
	}

	if bufferedWriter != nil {
		if err := bufferedWriter.Flush(); err != nil {
			return err
		}
	}

	if outputFile != nil {
		if err := outputFile.Close(); err != nil && closeErr == nil {
			closeErr = err
		}
	}

	if closeErr != nil {
		return fmt.Errorf("Failed to close file: %v", closeErr) // Return the first close error
	}

	return scanner.Err()
}

# commandSplit

## Overview

`File Splitter` is a simple Go command-line utility that splits a file into smaller parts. It provides options to split by lines, bytes, or the number of files.

## Installation

1. Clone this repository:

```
git clone https://github.com/maepon/commandSplitOnGo.git
```

2. Navigate to the project directory:

```
cd commandSplitOnGo
```

3. Build the project:

```
go build
```

## Usage

To split a file named `largefile.txt`:

- By lines (splits the file into parts with 1000 lines each by default):

```
./commandSplit largefile.txt
```

or

```
./commandSplit -l 1000 largefile.txt
```

- By number of files (splits the file into 5 smaller files):

```
./commandSplit -n 5 largefile.txt
```

- By bytes (splits the file into parts with 50 bytes each):

```
./commandSplit -b 50 largefile.txt
```

**Note**: You can only specify one of the `-l`, `-n`, or `-b` options at a time.

## Flags

- `-l`: Number of lines per output file
- `-n`: Number of files to split into
- `-b`: Number of bytes per output file

## Error Handling

The program will exit and display an appropriate error message for various types of errors such as invalid arguments or file I/O issues.



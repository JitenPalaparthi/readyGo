package helper

import "os"

// IsWindows is to check whether os is windows or not
func IsWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}

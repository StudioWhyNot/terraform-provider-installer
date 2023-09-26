package system

import (
	"strings"
)

func FindExecutablePath(paths []string, programName string) (string, error) {
	binPath := "bin/" + programName
	sbinPath := "sbin/" + programName
	for _, path := range paths {
		found := strings.HasSuffix(path, binPath) || strings.HasSuffix(path, sbinPath)
		if found {
			return path, nil
		}
	}
	//Not finding an executable should not be an error, such as with apt.
	return "", nil
}

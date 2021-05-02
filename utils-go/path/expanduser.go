package path

import (
	"os"
	"runtime"
	"strings"
)

// ExpandUser takes in a path, and replaces the ~ if it is the first character
// in that path. On UNIX systems, it replaces it with the HOME varaible, on
// Windows systems it replaces it with a USERPROFILE variable.
func ExpandUser(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	return strings.Replace(path, "~", Home(), 1)
}

// Home fetches the users home directory, with HOME env var for UNIX systems and
// USERPROFILE for Windows systems.
//
// TODO: handle ~user
func Home() string {
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("USERPROFILE")
	default:
		return os.Getenv("HOME")
	}
}

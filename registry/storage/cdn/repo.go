package cdn

import (
	"strings"
)

func UseCDN(repo string) bool {
	return strings.HasPrefix(repo, "library/")
}

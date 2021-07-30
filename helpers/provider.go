package helpers

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func ProviderName() (providerName *string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not determine current working directory to evaluate provider name: %+v", err)
	}
	ps := "/"
	if runtime.GOOS == "windows" {
		ps = "\\"
	}
	wdParts := strings.Split(cwd, ps)
	fullProviderName := wdParts[len(wdParts)-1]
	if !strings.HasPrefix(fullProviderName, "terraform-provider-") {
		return nil, fmt.Errorf("current Working Directory path does not appear to be a terraform provider: %+v", cwd)
	}
	nameParts := strings.Split(fullProviderName, "-")

	providerName = &nameParts[len(nameParts)-1]

	return
}

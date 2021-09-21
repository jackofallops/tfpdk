package helpers

import (
	"fmt"
	"os/exec"
)

func CallTerraform(opts ...string) ([]byte, error) {
	cmd := exec.Command("terraform", opts...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func GoFmt(file string) error {
	cmd := exec.Command("gofmt", "-w", fmt.Sprintf("./%s", file))
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

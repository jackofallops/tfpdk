package helpers

import (
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

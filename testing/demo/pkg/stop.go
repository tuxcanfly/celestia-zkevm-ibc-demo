package pkg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func Stop(ctx context.Context) error {
	fmt.Println("Stopping demo environment...")
	if err := runDockerCompose(ctx, "rm", "-s", "-f"); err != nil {
		return err
	}

	return os.RemoveAll("./.tmp")
}

func runDockerCompose(ctx context.Context, args ...string) error {
	args = append([]string{"compose"}, args...)
	return exec.CommandContext(ctx, "docker", args...).Run()
}

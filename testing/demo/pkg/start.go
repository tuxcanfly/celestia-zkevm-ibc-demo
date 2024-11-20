package pkg

import (
	"context"
	"fmt"
	"log"
	"time"

	client "github.com/celestiaorg/celestia-openrpc"
)

func Start(ctx context.Context) error {
	fmt.Println("Starting demo environment...")

	errCh := make(chan error)
	go func() {
		if err := runDockerCompose(ctx, "up"); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Second):
		go func() {
			select {
			case err := <-errCh:
				log.Println("error running docker compose up:", err)
			case <-ctx.Done():
				return
			}
		}()
	}

	if err := waitForDA(ctx, 3, time.Minute); err != nil {
		return err
	}

	return nil
}

func waitForDA(ctx context.Context, height uint64, timeout time.Duration) error {
	client, err := client.NewClient(ctx, "http://localhost:26658", "")
	if err != nil {
		return fmt.Errorf("failed to create celestia DA client: %w", err)
	}

	var latestHeight uint64
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()
	for {
		select {
		case <-timeoutTimer.C:
			return fmt.Errorf("timeout waiting for DA height %d (latest: %d)", height, latestHeight)
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			header, err := client.Header.NetworkHead(ctx)
			if err != nil {
				return err
			}
			latestHeight = header.Height()
			if latestHeight >= height {
				return nil
			}

		}
	}
}

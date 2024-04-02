package bekup

import (
	"context"

	"github.com/pandeptwidyaop/bekup/internal/cleanup"
	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/dump"
	"github.com/pandeptwidyaop/bekup/internal/upload"
	"github.com/pandeptwidyaop/bekup/internal/zip"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, config config.Config, worker int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	dumpCh, err := dump.Run(ctx, worker, config.Sources...)
	if err != nil {
		cancel()
		return err
	}

	zipCh := zip.Run(ctx, dumpCh, worker)

	uploadCh := upload.Run(ctx, zipCh, worker, config.Destinations...)

	cleanupCh := cleanup.Run(ctx, uploadCh)

	g.Go(func() error {
		for m := range cleanupCh {
			if m.Err != nil {
				return m.Err
			}
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		cancel()
		return err
	}

	return nil
}
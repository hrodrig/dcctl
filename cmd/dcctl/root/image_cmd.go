package root

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
)

// newImageCmd adds "image ls", "image pull", "image remove" using manifests from the current environment.
// Optional service argument limits to that service/manifest.
func newImageCmd(opts *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Inspect or update images from environment manifests",
		Long:  "Lists, pulls, or removes images referenced in the manifests of the current environment. Optional [service] limits to that service/manifest.",
	}

	cmd.AddCommand(
		newImageLsCmd(opts),
		newImagePullCmd(opts),
		newImageRemoveCmd(opts),
	)
	return cmd
}

func newImageLsCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "ls [service]",
		Short: "List images used by environment manifests",
		Long:  "Shows images defined in the manifests (and their local ID if present). Add [service] to limit to one service/manifest.",
		Example: `  dcctl image ls
  dcctl image ls app`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, images, err := resolveImagesFromEnv(opts, args)
			if err != nil {
				return err
			}
			if len(images) == 0 {
				return nil
			}
			printImages(images)
			return nil
		},
	}
}

func newImagePullCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "pull [service]",
		Short: "Pull images from environment manifests",
		Long:  "Runs docker pull for each image referenced in the manifests so you have the latest. Add [service] to limit to one service/manifest.",
		Example: `  dcctl image pull
  dcctl image pull app`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := withSignalCancel(context.Background())
			defer cancel()
			if err := ensureDockerAvailable(); err != nil {
				return err
			}
			_, images, err := resolveImagesFromEnv(opts, args)
			if err != nil {
				return err
			}
			if len(images) == 0 {
				return nil
			}
			return pullImages(ctx, images)
		},
	}
}

func newImageRemoveCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "remove [service]",
		Short: "Remove images from environment manifests",
		Long:  "Removes local images referenced in the manifests (after confirmation). Add [service] to limit to one service/manifest.",
		Example: `  dcctl image remove
  dcctl image remove app`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := withSignalCancel(context.Background())
			defer cancel()
			if err := ensureDockerAvailable(); err != nil {
				return err
			}
			_, images, err := resolveImagesFromEnv(opts, args)
			if err != nil {
				return err
			}
			if len(images) == 0 {
				return nil
			}
			printImages(images)
			if !confirmImageRemoval(images) {
				log.Warn("cancelled")
				return nil
			}
			return removeImages(ctx, images)
		},
	}
}

// resolveImagesFromEnv loads config for current environment, resolves manifest paths (all or for service), returns paths and unique image list.
func resolveImagesFromEnv(opts *options, serviceArgs []string) ([]string, []string, error) {
	loaded, err := loadOrBootstrap(opts)
	if err != nil {
		return nil, nil, err
	}
	if loaded == nil {
		return nil, nil, nil
	}
	if err := ensureEnvironmentDir(loaded); err != nil {
		return nil, nil, err
	}
	var manifestPaths []string
	if len(serviceArgs) > 0 {
		service := strings.TrimSpace(serviceArgs[0])
		if service != "" {
			manifestPaths, err = resolveTargetManifestPaths(loaded, []string{service})
			if err != nil {
				return nil, nil, err
			}
		}
	}
	if len(manifestPaths) == 0 {
		manifestPaths = resolveManifestPaths(loaded)
	}
	if len(manifestPaths) == 0 {
		log.Warn("no manifests found")
		return nil, nil, nil
	}
	images, err := collectImages(manifestPaths)
	if err != nil {
		return nil, nil, err
	}
	if len(images) == 0 {
		log.Warn("no images found in manifests")
		return nil, nil, nil
	}
	return manifestPaths, images, nil
}

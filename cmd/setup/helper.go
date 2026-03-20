package setup

import "context"

// findOrCreate lists items, finds a match by name, and creates one if not found.
// Returns (id, created, error). In dry-run mode, missing entities
// return (0, true, nil) without calling create.
func findOrCreate[T any](
	ctx context.Context,
	list func(context.Context) ([]T, error),
	match func(T) bool,
	getID func(T) uint32,
	create func(context.Context) (uint32, error),
	dryRun bool,
) (uint32, bool, error) {
	items, err := list(ctx)
	if err != nil {
		return 0, false, err
	}

	for _, item := range items {
		if match(item) {
			return getID(item), false, nil
		}
	}

	if dryRun {
		return 0, true, nil
	}

	id, err := create(ctx)
	if err != nil {
		return 0, false, err
	}

	return id, true, nil
}

package go_swagger_ui

import (
	"embed"
	"io/fs"
	"slices"
	"strings"
)

func walkFS(ignorePrefix string, efs *embed.FS, ignoreFiles ...string) (map[string]struct{}, error) {
	fileMap := make(map[string]struct{})

	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			filePath := strings.TrimPrefix(path, ignorePrefix)

			if !slices.Contains(ignoreFiles, filePath) {
				fileMap[filePath] = struct{}{}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return fileMap, nil
}

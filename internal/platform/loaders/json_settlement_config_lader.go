package loaders

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type JSONSettlementConfigLoader struct {
	Dir string
}

func NewJSONSettlementConfigLoader(dir string) JSONSettlementConfigLoader {
	return JSONSettlementConfigLoader{Dir: dir}
}

func (l JSONSettlementConfigLoader) LoadSettlementNames() ([]string, error) {
	files, err := os.ReadDir(l.Dir)
	if err != nil {
		return nil, err
	}

	var settlementNames []string

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		fullPath := filepath.Join(l.Dir, file.Name())
		raw, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, err
		}

		var names []string
		if err := json.Unmarshal(raw, &names); err != nil {
			return nil, fmt.Errorf("invalid JSON in %s: expected []string: %w", fullPath, err)
		}

		settlementNames = append(settlementNames, names...)
	}

	return settlementNames, nil
}

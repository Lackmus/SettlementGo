package loaders

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/lackmus/settlementgengo/pkg/model"
)

type JSONSettlementLoader struct {
	Dir string
}

func NewJSONSettlementStorage(dir string) JSONSettlementLoader {
	return JSONSettlementLoader{Dir: dir}
}

func (l JSONSettlementLoader) LoadSettlement(name string) (model.Settlement, error) {
	filename := l.Dir + "/" + name + ".json"
	file, err := os.Open(filename)
	if err != nil {
		return model.Settlement{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return model.Settlement{}, err
	}
	var settlement model.Settlement
	if err := json.Unmarshal(data, &settlement); err != nil {
		return model.Settlement{}, err
	}
	return settlement, nil
}

func (l JSONSettlementLoader) LoadAllSettlements() ([]model.Settlement, error) {
	files, err := os.ReadDir(l.Dir)
	if err != nil {
		return nil, err
	}
	var settlements []model.Settlement
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}
		settlement, err := l.LoadSettlement(file.Name()[:len(file.Name())-5])
		if err != nil {
			return nil, err
		}
		settlements = append(settlements, settlement)
	}
	return settlements, nil
}

func (l JSONSettlementLoader) SaveSettlement(nsettlement model.Settlement) error {
	filename := l.Dir + "/" + nsettlement.Name + ".json"
	data, err := json.MarshalIndent(nsettlement, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (l JSONSettlementLoader) SaveAllSettlements(settlements []model.Settlement) error {
	for _, settlement := range settlements {
		if err := l.SaveSettlement(settlement); err != nil {
			return err
		}
	}
	return nil
}

func (l JSONSettlementLoader) DeleteSettlement(name string) error {
	filename := l.Dir + "/" + name + ".json"
	return os.Remove(filename)
}

func (l JSONSettlementLoader) DeleteAllSettlements() error {
	files, err := os.ReadDir(l.Dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}
		if err := os.Remove(l.Dir + "/" + file.Name()); err != nil {
			return err
		}
	}
	return nil
}

package shared

type SettlementConfigLoader interface {
	LoadSettlementNames() ([]string, error)
}

package shared

type SettlementViewer interface {
	SettlementObserver
	DisplaySettlements()
}

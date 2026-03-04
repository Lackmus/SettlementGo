package shared

type SettlementObservable interface {
	RegisterObserver(observer SettlementObserver)
	RemoveObserver(observer SettlementObserver)
	NotifyObservers()
}

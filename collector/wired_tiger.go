package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	wtAvailableReadTicketsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "wiredtiger",
		Name:      "read_tickets",
		Help:      "Available concurrent read operation tickets by the WiredTiger storage engine",
	})
	wtAvailableWriteTicketsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "wiredtiger",
		Name:      "write_tickets",
		Help:      "Available concurrent write operation tickets by the WiredTiger storage engine",
	})
)

// WiredTiger are some wired Tiger storage engine specific metrics
type WiredTiger struct {
	ConcurrentTransactions *WTConcurrentTransactions `bson:"concurrentTransactions"`
}

type WTConcurrentTransactions struct {
	Read  WTConcurrentTransactionsInfo `bson:"read"`
	Write WTConcurrentTransactionsInfo `bson:"write"`
}

type WTConcurrentTransactionsInfo struct {
	Out          float64 `bson:"out"`
	Available    float64 `bson:"available"`
	TotalTickets float64 `bson:"totalTickets"`
}

// Export exports the data to prometheus.
func (wiredTiger *WiredTiger) Export(ch chan<- prometheus.Metric) {
	wtAvailableReadTicketsGauge.Set(wiredTiger.ConcurrentTransactions.Read.Available)
	wtAvailableWriteTicketsGauge.Set(wiredTiger.ConcurrentTransactions.Write.Available)

	wtAvailableReadTicketsGauge.Collect(ch)
	wtAvailableWriteTicketsGauge.Collect(ch)
}

// Describe describes the metrics for prometheus
func (wiredTiger *WiredTiger) Describe(ch chan<- *prometheus.Desc) {
	wtAvailableReadTicketsGauge.Describe(ch)
	wtAvailableWriteTicketsGauge.Describe(ch)
}

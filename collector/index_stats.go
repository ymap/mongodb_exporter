package collector

import (
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	indexUsage = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "index_usage_count",
		Help:      "Contains a usage count of the given index",
	}, []string{"collection", "index"})
)

// IndexStats represents index usage information
type IndexStats struct {
	Collection string
	Items      []IndexStatsItem
}

// IndexStatsItem represents stats about an Index
type IndexStatsItem struct {
	Name     string         `bson:"name"`
	Accesses IndexUsageInfo `bson:"accesses"`
}

type IndexUsageInfo struct {
	Ops float64 `bson:"ops"`
}

// Export exports database stats to prometheus
func (indexStats *IndexStats) Export(ch chan<- prometheus.Metric) {
	indexUsage.Reset()
	for _, indexStat := range indexStats.Items {
		indexUsage.WithLabelValues(indexStats.Collection, indexStat.Name).Add(indexStat.Accesses.Ops)
	}
	indexUsage.Collect(ch)
}

// Describe describes database stats for prometheus
func (indexStats *IndexStats) Describe(ch chan<- *prometheus.Desc) {
	objectCount.Describe(ch)
	collDataSize.Describe(ch)
	collStorageSize.Describe(ch)
	collTotalIndexSize.Describe(ch)
}

// GetIndexStats returns stats for a given collection in a database
func GetIndexStats(session *mgo.Session, db, collection string) *IndexStats {
	indexStats := IndexStats{Collection: collection}
	err := session.DB(db).C(collection).Pipe([]bson.M{{"$indexStats": bson.M{}}}).All(&indexStats.Items)
	if err != nil {
		glog.Error(err)
		return nil
	}

	return &indexStats
}

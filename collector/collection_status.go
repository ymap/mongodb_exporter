package collector

import (
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	objectCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "collection",
		Name:      "objects_count",
		Help:      "Contains a count of the number of objects (i.e. documents) in this collection",
	}, []string{"collection"})
	collDataSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "collection",
		Name:      "data_size_bytes",
		Help:      "The total size in bytes of the uncompressed data held in this collection",
	}, []string{"collection"})
	collStorageSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "collection",
		Name:      "storage_size_bytes",
		Help:      "The total amount of storage allocated to this collection for document storage",
	}, []string{"collection"})
	collTotalIndexSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "collection",
		Name:      "total_index_size_bytes",
		Help:      "The total size of all indexes",
	}, []string{"collection"})
)

// CollectionStatus represents stats about a collection
type CollectionStatus struct {
	Name           string `bson:"ns,omitempty"`
	ObjectCount    int    `bson:"count,omitempty"`
	DataSize       int    `bson:"size,omitempty"`
	StorageSize    int    `bson:"storageSize,omitempty"`
	TotalIndexSize int    `bson:"totalIndexSize,omitempty"`
}

// Export exports database stats to prometheus
func (collStatus *CollectionStatus) Export(ch chan<- prometheus.Metric) {
	objectCount.Reset()
	collDataSize.Reset()
	collStorageSize.Reset()
	collTotalIndexSize.Reset()

	objectCount.WithLabelValues(collStatus.Name).Set(float64(collStatus.ObjectCount))
	collDataSize.WithLabelValues(collStatus.Name).Set(float64(collStatus.DataSize))
	collStorageSize.WithLabelValues(collStatus.Name).Set(float64(collStatus.StorageSize))
	collTotalIndexSize.WithLabelValues(collStatus.Name).Set(float64(collStatus.TotalIndexSize))

	objectCount.Collect(ch)
	collDataSize.Collect(ch)
	collStorageSize.Collect(ch)
	collTotalIndexSize.Collect(ch)
}

// Describe describes database stats for prometheus
func (collStatus *CollectionStatus) Describe(ch chan<- *prometheus.Desc) {
	objectCount.Describe(ch)
	collDataSize.Describe(ch)
	collStorageSize.Describe(ch)
	collTotalIndexSize.Describe(ch)
}

// GetCollectionStatus returns stats for a given collection in a database
func GetCollectionStatus(session *mgo.Session, db, collection string) *CollectionStatus {
	var collStatus CollectionStatus
	err := session.DB(db).Run(bson.D{{"collStats", collection}, {"scale", 1}}, &collStatus)
	if err != nil {
		glog.Error(err)
		return nil
	}

	return &collStatus
}

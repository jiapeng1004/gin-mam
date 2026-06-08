// Package asset 媒资检索索引抽象，解耦 asset 与 search 包循环依赖。
package asset

import "time"

// Indexer 媒资变更后异步更新 Elasticsearch 索引的接口。
type Indexer interface {
	// ScheduleIndex 异步调度索引写入，不阻塞调用方。
	ScheduleIndex(doc IndexDoc)
}

// IndexDoc 写入 ES 的媒资索引文档最小字段集。
type IndexDoc struct {
	// ID 媒资 ID。
	ID string
	// Title 媒资标题。
	Title string
	// Type 媒资类型。
	Type string
	// CatalogID 编目 ID，可为空。
	CatalogID *string
	// TenantID 租户标识。
	TenantID string
	// CreatedAt 媒资创建时间，用于 ES 排序。
	CreatedAt time.Time
}

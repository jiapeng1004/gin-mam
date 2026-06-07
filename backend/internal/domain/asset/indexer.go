package asset

import "time"

// Indexer schedules asynchronous search indexing after asset mutations.
type Indexer interface {
	ScheduleIndex(doc IndexDoc)
}

type IndexDoc struct {
	ID        string
	Title     string
	Type      string
	CatalogID *string
	TenantID  string
	CreatedAt time.Time
}

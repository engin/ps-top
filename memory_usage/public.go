// Package memory_usage manages collecting data from performance_schema which holds
// information about memory usage
package memory_usage

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // keep golint happy

	"github.com/sjmudd/ps-top/p_s"
)

const (
	description = "Memory Usage (memory_summary_global_by_event_name)"
)

// Object represents a table of rows
type Object struct {
	p_s.RelativeStats
	p_s.CollectionTime
	current Rows // last loaded values
	results Rows // results (maybe with subtraction)
	totals  Row  // totals of results
}

// Collect data from the db, no merging needed
func (t *Object) Collect(dbh *sql.DB) {
	t.current = selectRows(dbh)

	t.makeResults()
}

// SetInitialFromCurrent resets the statistics to current values
func (t *Object) SetInitialFromCurrent() {
	t.SetCollected()

	t.makeResults()
}

// Headings returns the headings for a table
func (t Object) Headings() string {
	var r Row

	return r.headings()
}

// RowContent returns the rows we need for displaying
func (t Object) RowContent(maxRows int) []string {
	rows := make([]string, 0, maxRows)

	for i := range t.results {
		if i < maxRows {
			rows = append(rows, t.results[i].rowContent(t.totals))
		}
	}

	return rows
}

// Rows() returns the rows we have which are interesting
func (t Object) Rows() []Row {
	rows := make([]Row, 0, len(t.results))

	for i := range t.results {
		rows = append(rows, t.results[i])
	}

	return rows
}

// Totals return the row of totals
func (t Object) Totals() Row {
	return t.totals
}

// TotalRowContent returns all the totals
func (t Object) TotalRowContent() string {
	return t.totals.rowContent(t.totals)
}

// EmptyRowContent returns an empty string of data (for filling in)
func (t Object) EmptyRowContent() string {
	var empty Row
	return empty.rowContent(empty)
}

// Description provides a description of the table
func (t Object) Description() string {
	return description
}

// Len returns the length of the result set
func (t Object) Len() int {
	return len(t.results)
}
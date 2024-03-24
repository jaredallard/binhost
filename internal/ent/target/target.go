// Code generated by ent, DO NOT EDIT.

package target

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the target type in the database.
	Label = "target"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// EdgePkgs holds the string denoting the pkgs edge name in mutations.
	EdgePkgs = "pkgs"
	// Table holds the table name of the target in the database.
	Table = "targets"
	// PkgsTable is the table that holds the pkgs relation/edge.
	PkgsTable = "targets"
	// PkgsInverseTable is the table name for the Pkg entity.
	// It exists in this package in order to avoid circular dependency with the "pkg" package.
	PkgsInverseTable = "pkgs"
	// PkgsColumn is the table column denoting the pkgs relation/edge.
	PkgsColumn = "target_pkgs"
)

// Columns holds all SQL columns for target fields.
var Columns = []string{
	FieldID,
	FieldName,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "targets"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"target_pkgs",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Target queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByPkgsField orders the results by pkgs field.
func ByPkgsField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPkgsStep(), sql.OrderByField(field, opts...))
	}
}
func newPkgsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PkgsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, PkgsTable, PkgsColumn),
	)
}

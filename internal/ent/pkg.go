// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/jaredallard/binhost/internal/ent/pkg"
)

// Pkg is the model entity for the Pkg schema.
type Pkg struct {
	config
	// ID of the ent.
	ID           uuid.UUID `json:"id,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Pkg) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case pkg.FieldID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Pkg fields.
func (pk *Pkg) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case pkg.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				pk.ID = *value
			}
		default:
			pk.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Pkg.
// This includes values selected through modifiers, order, etc.
func (pk *Pkg) Value(name string) (ent.Value, error) {
	return pk.selectValues.Get(name)
}

// Update returns a builder for updating this Pkg.
// Note that you need to call Pkg.Unwrap() before calling this method if this Pkg
// was returned from a transaction, and the transaction was committed or rolled back.
func (pk *Pkg) Update() *PkgUpdateOne {
	return NewPkgClient(pk.config).UpdateOne(pk)
}

// Unwrap unwraps the Pkg entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pk *Pkg) Unwrap() *Pkg {
	_tx, ok := pk.config.driver.(*txDriver)
	if !ok {
		panic("ent: Pkg is not a transactional entity")
	}
	pk.config.driver = _tx.drv
	return pk
}

// String implements the fmt.Stringer.
func (pk *Pkg) String() string {
	var builder strings.Builder
	builder.WriteString("Pkg(")
	builder.WriteString(fmt.Sprintf("id=%v", pk.ID))
	builder.WriteByte(')')
	return builder.String()
}

// Pkgs is a parsable slice of Pkg.
type Pkgs []*Pkg

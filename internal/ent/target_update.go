// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/jaredallard/binhost/internal/ent/pkg"
	"github.com/jaredallard/binhost/internal/ent/predicate"
	"github.com/jaredallard/binhost/internal/ent/target"
)

// TargetUpdate is the builder for updating Target entities.
type TargetUpdate struct {
	config
	hooks    []Hook
	mutation *TargetMutation
}

// Where appends a list predicates to the TargetUpdate builder.
func (tu *TargetUpdate) Where(ps ...predicate.Target) *TargetUpdate {
	tu.mutation.Where(ps...)
	return tu
}

// SetName sets the "name" field.
func (tu *TargetUpdate) SetName(s string) *TargetUpdate {
	tu.mutation.SetName(s)
	return tu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (tu *TargetUpdate) SetNillableName(s *string) *TargetUpdate {
	if s != nil {
		tu.SetName(*s)
	}
	return tu
}

// SetPkgsID sets the "pkgs" edge to the Pkg entity by ID.
func (tu *TargetUpdate) SetPkgsID(id uuid.UUID) *TargetUpdate {
	tu.mutation.SetPkgsID(id)
	return tu
}

// SetNillablePkgsID sets the "pkgs" edge to the Pkg entity by ID if the given value is not nil.
func (tu *TargetUpdate) SetNillablePkgsID(id *uuid.UUID) *TargetUpdate {
	if id != nil {
		tu = tu.SetPkgsID(*id)
	}
	return tu
}

// SetPkgs sets the "pkgs" edge to the Pkg entity.
func (tu *TargetUpdate) SetPkgs(p *Pkg) *TargetUpdate {
	return tu.SetPkgsID(p.ID)
}

// Mutation returns the TargetMutation object of the builder.
func (tu *TargetUpdate) Mutation() *TargetMutation {
	return tu.mutation
}

// ClearPkgs clears the "pkgs" edge to the Pkg entity.
func (tu *TargetUpdate) ClearPkgs() *TargetUpdate {
	tu.mutation.ClearPkgs()
	return tu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tu *TargetUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, tu.sqlSave, tu.mutation, tu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TargetUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TargetUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TargetUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tu *TargetUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(target.Table, target.Columns, sqlgraph.NewFieldSpec(target.FieldID, field.TypeUUID))
	if ps := tu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tu.mutation.Name(); ok {
		_spec.SetField(target.FieldName, field.TypeString, value)
	}
	if tu.mutation.PkgsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   target.PkgsTable,
			Columns: []string{target.PkgsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(pkg.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.PkgsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   target.PkgsTable,
			Columns: []string{target.PkgsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(pkg.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{target.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	tu.mutation.done = true
	return n, nil
}

// TargetUpdateOne is the builder for updating a single Target entity.
type TargetUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *TargetMutation
}

// SetName sets the "name" field.
func (tuo *TargetUpdateOne) SetName(s string) *TargetUpdateOne {
	tuo.mutation.SetName(s)
	return tuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (tuo *TargetUpdateOne) SetNillableName(s *string) *TargetUpdateOne {
	if s != nil {
		tuo.SetName(*s)
	}
	return tuo
}

// SetPkgsID sets the "pkgs" edge to the Pkg entity by ID.
func (tuo *TargetUpdateOne) SetPkgsID(id uuid.UUID) *TargetUpdateOne {
	tuo.mutation.SetPkgsID(id)
	return tuo
}

// SetNillablePkgsID sets the "pkgs" edge to the Pkg entity by ID if the given value is not nil.
func (tuo *TargetUpdateOne) SetNillablePkgsID(id *uuid.UUID) *TargetUpdateOne {
	if id != nil {
		tuo = tuo.SetPkgsID(*id)
	}
	return tuo
}

// SetPkgs sets the "pkgs" edge to the Pkg entity.
func (tuo *TargetUpdateOne) SetPkgs(p *Pkg) *TargetUpdateOne {
	return tuo.SetPkgsID(p.ID)
}

// Mutation returns the TargetMutation object of the builder.
func (tuo *TargetUpdateOne) Mutation() *TargetMutation {
	return tuo.mutation
}

// ClearPkgs clears the "pkgs" edge to the Pkg entity.
func (tuo *TargetUpdateOne) ClearPkgs() *TargetUpdateOne {
	tuo.mutation.ClearPkgs()
	return tuo
}

// Where appends a list predicates to the TargetUpdate builder.
func (tuo *TargetUpdateOne) Where(ps ...predicate.Target) *TargetUpdateOne {
	tuo.mutation.Where(ps...)
	return tuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tuo *TargetUpdateOne) Select(field string, fields ...string) *TargetUpdateOne {
	tuo.fields = append([]string{field}, fields...)
	return tuo
}

// Save executes the query and returns the updated Target entity.
func (tuo *TargetUpdateOne) Save(ctx context.Context) (*Target, error) {
	return withHooks(ctx, tuo.sqlSave, tuo.mutation, tuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TargetUpdateOne) SaveX(ctx context.Context) *Target {
	node, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tuo *TargetUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TargetUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tuo *TargetUpdateOne) sqlSave(ctx context.Context) (_node *Target, err error) {
	_spec := sqlgraph.NewUpdateSpec(target.Table, target.Columns, sqlgraph.NewFieldSpec(target.FieldID, field.TypeUUID))
	id, ok := tuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Target.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := tuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, target.FieldID)
		for _, f := range fields {
			if !target.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != target.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tuo.mutation.Name(); ok {
		_spec.SetField(target.FieldName, field.TypeString, value)
	}
	if tuo.mutation.PkgsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   target.PkgsTable,
			Columns: []string{target.PkgsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(pkg.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.PkgsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   target.PkgsTable,
			Columns: []string{target.PkgsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(pkg.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Target{config: tuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{target.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	tuo.mutation.done = true
	return _node, nil
}
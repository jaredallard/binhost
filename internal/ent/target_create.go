// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/jaredallard/binhost/internal/ent/pkg"
	"github.com/jaredallard/binhost/internal/ent/target"
)

// TargetCreate is the builder for creating a Target entity.
type TargetCreate struct {
	config
	mutation *TargetMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (tc *TargetCreate) SetName(s string) *TargetCreate {
	tc.mutation.SetName(s)
	return tc
}

// SetID sets the "id" field.
func (tc *TargetCreate) SetID(u uuid.UUID) *TargetCreate {
	tc.mutation.SetID(u)
	return tc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (tc *TargetCreate) SetNillableID(u *uuid.UUID) *TargetCreate {
	if u != nil {
		tc.SetID(*u)
	}
	return tc
}

// SetPkgsID sets the "pkgs" edge to the Pkg entity by ID.
func (tc *TargetCreate) SetPkgsID(id uuid.UUID) *TargetCreate {
	tc.mutation.SetPkgsID(id)
	return tc
}

// SetNillablePkgsID sets the "pkgs" edge to the Pkg entity by ID if the given value is not nil.
func (tc *TargetCreate) SetNillablePkgsID(id *uuid.UUID) *TargetCreate {
	if id != nil {
		tc = tc.SetPkgsID(*id)
	}
	return tc
}

// SetPkgs sets the "pkgs" edge to the Pkg entity.
func (tc *TargetCreate) SetPkgs(p *Pkg) *TargetCreate {
	return tc.SetPkgsID(p.ID)
}

// Mutation returns the TargetMutation object of the builder.
func (tc *TargetCreate) Mutation() *TargetMutation {
	return tc.mutation
}

// Save creates the Target in the database.
func (tc *TargetCreate) Save(ctx context.Context) (*Target, error) {
	tc.defaults()
	return withHooks(ctx, tc.sqlSave, tc.mutation, tc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TargetCreate) SaveX(ctx context.Context) *Target {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tc *TargetCreate) Exec(ctx context.Context) error {
	_, err := tc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tc *TargetCreate) ExecX(ctx context.Context) {
	if err := tc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tc *TargetCreate) defaults() {
	if _, ok := tc.mutation.ID(); !ok {
		v := target.DefaultID()
		tc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tc *TargetCreate) check() error {
	if _, ok := tc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Target.name"`)}
	}
	return nil
}

func (tc *TargetCreate) sqlSave(ctx context.Context) (*Target, error) {
	if err := tc.check(); err != nil {
		return nil, err
	}
	_node, _spec := tc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	tc.mutation.id = &_node.ID
	tc.mutation.done = true
	return _node, nil
}

func (tc *TargetCreate) createSpec() (*Target, *sqlgraph.CreateSpec) {
	var (
		_node = &Target{config: tc.config}
		_spec = sqlgraph.NewCreateSpec(target.Table, sqlgraph.NewFieldSpec(target.FieldID, field.TypeUUID))
	)
	if id, ok := tc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := tc.mutation.Name(); ok {
		_spec.SetField(target.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if nodes := tc.mutation.PkgsIDs(); len(nodes) > 0 {
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
		_node.target_pkgs = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// TargetCreateBulk is the builder for creating many Target entities in bulk.
type TargetCreateBulk struct {
	config
	err      error
	builders []*TargetCreate
}

// Save creates the Target entities in the database.
func (tcb *TargetCreateBulk) Save(ctx context.Context) ([]*Target, error) {
	if tcb.err != nil {
		return nil, tcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(tcb.builders))
	nodes := make([]*Target, len(tcb.builders))
	mutators := make([]Mutator, len(tcb.builders))
	for i := range tcb.builders {
		func(i int, root context.Context) {
			builder := tcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TargetMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, tcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, tcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tcb *TargetCreateBulk) SaveX(ctx context.Context) []*Target {
	v, err := tcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcb *TargetCreateBulk) Exec(ctx context.Context) error {
	_, err := tcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcb *TargetCreateBulk) ExecX(ctx context.Context) {
	if err := tcb.Exec(ctx); err != nil {
		panic(err)
	}
}
// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/jaredallard/binhost/internal/ent/pkg"
	"github.com/jaredallard/binhost/internal/ent/predicate"
	"github.com/jaredallard/binhost/internal/ent/target"
)

// TargetQuery is the builder for querying Target entities.
type TargetQuery struct {
	config
	ctx        *QueryContext
	order      []target.OrderOption
	inters     []Interceptor
	predicates []predicate.Target
	withPkgs   *PkgQuery
	withFKs    bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the TargetQuery builder.
func (tq *TargetQuery) Where(ps ...predicate.Target) *TargetQuery {
	tq.predicates = append(tq.predicates, ps...)
	return tq
}

// Limit the number of records to be returned by this query.
func (tq *TargetQuery) Limit(limit int) *TargetQuery {
	tq.ctx.Limit = &limit
	return tq
}

// Offset to start from.
func (tq *TargetQuery) Offset(offset int) *TargetQuery {
	tq.ctx.Offset = &offset
	return tq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (tq *TargetQuery) Unique(unique bool) *TargetQuery {
	tq.ctx.Unique = &unique
	return tq
}

// Order specifies how the records should be ordered.
func (tq *TargetQuery) Order(o ...target.OrderOption) *TargetQuery {
	tq.order = append(tq.order, o...)
	return tq
}

// QueryPkgs chains the current query on the "pkgs" edge.
func (tq *TargetQuery) QueryPkgs() *PkgQuery {
	query := (&PkgClient{config: tq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := tq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := tq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(target.Table, target.FieldID, selector),
			sqlgraph.To(pkg.Table, pkg.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, target.PkgsTable, target.PkgsColumn),
		)
		fromU = sqlgraph.SetNeighbors(tq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Target entity from the query.
// Returns a *NotFoundError when no Target was found.
func (tq *TargetQuery) First(ctx context.Context) (*Target, error) {
	nodes, err := tq.Limit(1).All(setContextOp(ctx, tq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{target.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (tq *TargetQuery) FirstX(ctx context.Context) *Target {
	node, err := tq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Target ID from the query.
// Returns a *NotFoundError when no Target ID was found.
func (tq *TargetQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = tq.Limit(1).IDs(setContextOp(ctx, tq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{target.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (tq *TargetQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := tq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Target entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Target entity is found.
// Returns a *NotFoundError when no Target entities are found.
func (tq *TargetQuery) Only(ctx context.Context) (*Target, error) {
	nodes, err := tq.Limit(2).All(setContextOp(ctx, tq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{target.Label}
	default:
		return nil, &NotSingularError{target.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (tq *TargetQuery) OnlyX(ctx context.Context) *Target {
	node, err := tq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Target ID in the query.
// Returns a *NotSingularError when more than one Target ID is found.
// Returns a *NotFoundError when no entities are found.
func (tq *TargetQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = tq.Limit(2).IDs(setContextOp(ctx, tq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{target.Label}
	default:
		err = &NotSingularError{target.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (tq *TargetQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := tq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Targets.
func (tq *TargetQuery) All(ctx context.Context) ([]*Target, error) {
	ctx = setContextOp(ctx, tq.ctx, "All")
	if err := tq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Target, *TargetQuery]()
	return withInterceptors[[]*Target](ctx, tq, qr, tq.inters)
}

// AllX is like All, but panics if an error occurs.
func (tq *TargetQuery) AllX(ctx context.Context) []*Target {
	nodes, err := tq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Target IDs.
func (tq *TargetQuery) IDs(ctx context.Context) (ids []uuid.UUID, err error) {
	if tq.ctx.Unique == nil && tq.path != nil {
		tq.Unique(true)
	}
	ctx = setContextOp(ctx, tq.ctx, "IDs")
	if err = tq.Select(target.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (tq *TargetQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := tq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (tq *TargetQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, tq.ctx, "Count")
	if err := tq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, tq, querierCount[*TargetQuery](), tq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (tq *TargetQuery) CountX(ctx context.Context) int {
	count, err := tq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (tq *TargetQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, tq.ctx, "Exist")
	switch _, err := tq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (tq *TargetQuery) ExistX(ctx context.Context) bool {
	exist, err := tq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the TargetQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (tq *TargetQuery) Clone() *TargetQuery {
	if tq == nil {
		return nil
	}
	return &TargetQuery{
		config:     tq.config,
		ctx:        tq.ctx.Clone(),
		order:      append([]target.OrderOption{}, tq.order...),
		inters:     append([]Interceptor{}, tq.inters...),
		predicates: append([]predicate.Target{}, tq.predicates...),
		withPkgs:   tq.withPkgs.Clone(),
		// clone intermediate query.
		sql:  tq.sql.Clone(),
		path: tq.path,
	}
}

// WithPkgs tells the query-builder to eager-load the nodes that are connected to
// the "pkgs" edge. The optional arguments are used to configure the query builder of the edge.
func (tq *TargetQuery) WithPkgs(opts ...func(*PkgQuery)) *TargetQuery {
	query := (&PkgClient{config: tq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	tq.withPkgs = query
	return tq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Target.Query().
//		GroupBy(target.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (tq *TargetQuery) GroupBy(field string, fields ...string) *TargetGroupBy {
	tq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &TargetGroupBy{build: tq}
	grbuild.flds = &tq.ctx.Fields
	grbuild.label = target.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.Target.Query().
//		Select(target.FieldName).
//		Scan(ctx, &v)
func (tq *TargetQuery) Select(fields ...string) *TargetSelect {
	tq.ctx.Fields = append(tq.ctx.Fields, fields...)
	sbuild := &TargetSelect{TargetQuery: tq}
	sbuild.label = target.Label
	sbuild.flds, sbuild.scan = &tq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a TargetSelect configured with the given aggregations.
func (tq *TargetQuery) Aggregate(fns ...AggregateFunc) *TargetSelect {
	return tq.Select().Aggregate(fns...)
}

func (tq *TargetQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range tq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, tq); err != nil {
				return err
			}
		}
	}
	for _, f := range tq.ctx.Fields {
		if !target.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if tq.path != nil {
		prev, err := tq.path(ctx)
		if err != nil {
			return err
		}
		tq.sql = prev
	}
	return nil
}

func (tq *TargetQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Target, error) {
	var (
		nodes       = []*Target{}
		withFKs     = tq.withFKs
		_spec       = tq.querySpec()
		loadedTypes = [1]bool{
			tq.withPkgs != nil,
		}
	)
	if tq.withPkgs != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, target.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Target).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Target{config: tq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, tq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := tq.withPkgs; query != nil {
		if err := tq.loadPkgs(ctx, query, nodes, nil,
			func(n *Target, e *Pkg) { n.Edges.Pkgs = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (tq *TargetQuery) loadPkgs(ctx context.Context, query *PkgQuery, nodes []*Target, init func(*Target), assign func(*Target, *Pkg)) error {
	ids := make([]uuid.UUID, 0, len(nodes))
	nodeids := make(map[uuid.UUID][]*Target)
	for i := range nodes {
		if nodes[i].target_pkgs == nil {
			continue
		}
		fk := *nodes[i].target_pkgs
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(pkg.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "target_pkgs" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (tq *TargetQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := tq.querySpec()
	_spec.Node.Columns = tq.ctx.Fields
	if len(tq.ctx.Fields) > 0 {
		_spec.Unique = tq.ctx.Unique != nil && *tq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, tq.driver, _spec)
}

func (tq *TargetQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(target.Table, target.Columns, sqlgraph.NewFieldSpec(target.FieldID, field.TypeUUID))
	_spec.From = tq.sql
	if unique := tq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if tq.path != nil {
		_spec.Unique = true
	}
	if fields := tq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, target.FieldID)
		for i := range fields {
			if fields[i] != target.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := tq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := tq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := tq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := tq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (tq *TargetQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(tq.driver.Dialect())
	t1 := builder.Table(target.Table)
	columns := tq.ctx.Fields
	if len(columns) == 0 {
		columns = target.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if tq.sql != nil {
		selector = tq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if tq.ctx.Unique != nil && *tq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range tq.predicates {
		p(selector)
	}
	for _, p := range tq.order {
		p(selector)
	}
	if offset := tq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := tq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// TargetGroupBy is the group-by builder for Target entities.
type TargetGroupBy struct {
	selector
	build *TargetQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (tgb *TargetGroupBy) Aggregate(fns ...AggregateFunc) *TargetGroupBy {
	tgb.fns = append(tgb.fns, fns...)
	return tgb
}

// Scan applies the selector query and scans the result into the given value.
func (tgb *TargetGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, tgb.build.ctx, "GroupBy")
	if err := tgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*TargetQuery, *TargetGroupBy](ctx, tgb.build, tgb, tgb.build.inters, v)
}

func (tgb *TargetGroupBy) sqlScan(ctx context.Context, root *TargetQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(tgb.fns))
	for _, fn := range tgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*tgb.flds)+len(tgb.fns))
		for _, f := range *tgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*tgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := tgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// TargetSelect is the builder for selecting fields of Target entities.
type TargetSelect struct {
	*TargetQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ts *TargetSelect) Aggregate(fns ...AggregateFunc) *TargetSelect {
	ts.fns = append(ts.fns, fns...)
	return ts
}

// Scan applies the selector query and scans the result into the given value.
func (ts *TargetSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ts.ctx, "Select")
	if err := ts.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*TargetQuery, *TargetSelect](ctx, ts.TargetQuery, ts, ts.inters, v)
}

func (ts *TargetSelect) sqlScan(ctx context.Context, root *TargetQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ts.fns))
	for _, fn := range ts.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ts.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ts.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

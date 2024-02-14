// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"gosfV2/src/ent/file"
	"gosfV2/src/ent/note"
	"gosfV2/src/ent/user"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// FileCreate is the builder for creating a File entity.
type FileCreate struct {
	config
	mutation *FileMutation
	hooks    []Hook
}

// SetCreatedAt sets the "created_at" field.
func (fc *FileCreate) SetCreatedAt(t time.Time) *FileCreate {
	fc.mutation.SetCreatedAt(t)
	return fc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (fc *FileCreate) SetNillableCreatedAt(t *time.Time) *FileCreate {
	if t != nil {
		fc.SetCreatedAt(*t)
	}
	return fc
}

// SetUpdatedAt sets the "updated_at" field.
func (fc *FileCreate) SetUpdatedAt(t time.Time) *FileCreate {
	fc.mutation.SetUpdatedAt(t)
	return fc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (fc *FileCreate) SetNillableUpdatedAt(t *time.Time) *FileCreate {
	if t != nil {
		fc.SetUpdatedAt(*t)
	}
	return fc
}

// SetFilename sets the "filename" field.
func (fc *FileCreate) SetFilename(s string) *FileCreate {
	fc.mutation.SetFilename(s)
	return fc
}

// SetIsDir sets the "is_dir" field.
func (fc *FileCreate) SetIsDir(b bool) *FileCreate {
	fc.mutation.SetIsDir(b)
	return fc
}

// SetNillableIsDir sets the "is_dir" field if the given value is not nil.
func (fc *FileCreate) SetNillableIsDir(b *bool) *FileCreate {
	if b != nil {
		fc.SetIsDir(*b)
	}
	return fc
}

// SetIsShared sets the "is_shared" field.
func (fc *FileCreate) SetIsShared(b bool) *FileCreate {
	fc.mutation.SetIsShared(b)
	return fc
}

// SetNillableIsShared sets the "is_shared" field if the given value is not nil.
func (fc *FileCreate) SetNillableIsShared(b *bool) *FileCreate {
	if b != nil {
		fc.SetIsShared(*b)
	}
	return fc
}

// SetID sets the "id" field.
func (fc *FileCreate) SetID(u uint) *FileCreate {
	fc.mutation.SetID(u)
	return fc
}

// SetOwnerID sets the "owner" edge to the User entity by ID.
func (fc *FileCreate) SetOwnerID(id uint) *FileCreate {
	fc.mutation.SetOwnerID(id)
	return fc
}

// SetOwner sets the "owner" edge to the User entity.
func (fc *FileCreate) SetOwner(u *User) *FileCreate {
	return fc.SetOwnerID(u.ID)
}

// AddSharedWithIDs adds the "shared_with" edge to the User entity by IDs.
func (fc *FileCreate) AddSharedWithIDs(ids ...uint) *FileCreate {
	fc.mutation.AddSharedWithIDs(ids...)
	return fc
}

// AddSharedWith adds the "shared_with" edges to the User entity.
func (fc *FileCreate) AddSharedWith(u ...*User) *FileCreate {
	ids := make([]uint, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return fc.AddSharedWithIDs(ids...)
}

// SetParentID sets the "parent" edge to the File entity by ID.
func (fc *FileCreate) SetParentID(id uint) *FileCreate {
	fc.mutation.SetParentID(id)
	return fc
}

// SetNillableParentID sets the "parent" edge to the File entity by ID if the given value is not nil.
func (fc *FileCreate) SetNillableParentID(id *uint) *FileCreate {
	if id != nil {
		fc = fc.SetParentID(*id)
	}
	return fc
}

// SetParent sets the "parent" edge to the File entity.
func (fc *FileCreate) SetParent(f *File) *FileCreate {
	return fc.SetParentID(f.ID)
}

// AddChildIDs adds the "children" edge to the File entity by IDs.
func (fc *FileCreate) AddChildIDs(ids ...uint) *FileCreate {
	fc.mutation.AddChildIDs(ids...)
	return fc
}

// AddChildren adds the "children" edges to the File entity.
func (fc *FileCreate) AddChildren(f ...*File) *FileCreate {
	ids := make([]uint, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return fc.AddChildIDs(ids...)
}

// SetNotesID sets the "notes" edge to the Note entity by ID.
func (fc *FileCreate) SetNotesID(id uint) *FileCreate {
	fc.mutation.SetNotesID(id)
	return fc
}

// SetNillableNotesID sets the "notes" edge to the Note entity by ID if the given value is not nil.
func (fc *FileCreate) SetNillableNotesID(id *uint) *FileCreate {
	if id != nil {
		fc = fc.SetNotesID(*id)
	}
	return fc
}

// SetNotes sets the "notes" edge to the Note entity.
func (fc *FileCreate) SetNotes(n *Note) *FileCreate {
	return fc.SetNotesID(n.ID)
}

// Mutation returns the FileMutation object of the builder.
func (fc *FileCreate) Mutation() *FileMutation {
	return fc.mutation
}

// Save creates the File in the database.
func (fc *FileCreate) Save(ctx context.Context) (*File, error) {
	fc.defaults()
	return withHooks(ctx, fc.sqlSave, fc.mutation, fc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (fc *FileCreate) SaveX(ctx context.Context) *File {
	v, err := fc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fc *FileCreate) Exec(ctx context.Context) error {
	_, err := fc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fc *FileCreate) ExecX(ctx context.Context) {
	if err := fc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (fc *FileCreate) defaults() {
	if _, ok := fc.mutation.CreatedAt(); !ok {
		v := file.DefaultCreatedAt()
		fc.mutation.SetCreatedAt(v)
	}
	if _, ok := fc.mutation.IsDir(); !ok {
		v := file.DefaultIsDir
		fc.mutation.SetIsDir(v)
	}
	if _, ok := fc.mutation.IsShared(); !ok {
		v := file.DefaultIsShared
		fc.mutation.SetIsShared(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fc *FileCreate) check() error {
	if _, ok := fc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "File.created_at"`)}
	}
	if _, ok := fc.mutation.Filename(); !ok {
		return &ValidationError{Name: "filename", err: errors.New(`ent: missing required field "File.filename"`)}
	}
	if v, ok := fc.mutation.Filename(); ok {
		if err := file.FilenameValidator(v); err != nil {
			return &ValidationError{Name: "filename", err: fmt.Errorf(`ent: validator failed for field "File.filename": %w`, err)}
		}
	}
	if _, ok := fc.mutation.IsDir(); !ok {
		return &ValidationError{Name: "is_dir", err: errors.New(`ent: missing required field "File.is_dir"`)}
	}
	if _, ok := fc.mutation.IsShared(); !ok {
		return &ValidationError{Name: "is_shared", err: errors.New(`ent: missing required field "File.is_shared"`)}
	}
	if _, ok := fc.mutation.OwnerID(); !ok {
		return &ValidationError{Name: "owner", err: errors.New(`ent: missing required edge "File.owner"`)}
	}
	return nil
}

func (fc *FileCreate) sqlSave(ctx context.Context) (*File, error) {
	if err := fc.check(); err != nil {
		return nil, err
	}
	_node, _spec := fc.createSpec()
	if err := sqlgraph.CreateNode(ctx, fc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = uint(id)
	}
	fc.mutation.id = &_node.ID
	fc.mutation.done = true
	return _node, nil
}

func (fc *FileCreate) createSpec() (*File, *sqlgraph.CreateSpec) {
	var (
		_node = &File{config: fc.config}
		_spec = sqlgraph.NewCreateSpec(file.Table, sqlgraph.NewFieldSpec(file.FieldID, field.TypeUint))
	)
	if id, ok := fc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := fc.mutation.CreatedAt(); ok {
		_spec.SetField(file.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := fc.mutation.UpdatedAt(); ok {
		_spec.SetField(file.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = &value
	}
	if value, ok := fc.mutation.Filename(); ok {
		_spec.SetField(file.FieldFilename, field.TypeString, value)
		_node.Filename = value
	}
	if value, ok := fc.mutation.IsDir(); ok {
		_spec.SetField(file.FieldIsDir, field.TypeBool, value)
		_node.IsDir = value
	}
	if value, ok := fc.mutation.IsShared(); ok {
		_spec.SetField(file.FieldIsShared, field.TypeBool, value)
		_node.IsShared = value
	}
	if nodes := fc.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   file.OwnerTable,
			Columns: []string{file.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.user_files = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fc.mutation.SharedWithIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   file.SharedWithTable,
			Columns: file.SharedWithPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fc.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   file.ParentTable,
			Columns: []string{file.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(file.FieldID, field.TypeUint),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.file_children = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fc.mutation.ChildrenIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.ChildrenTable,
			Columns: []string{file.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(file.FieldID, field.TypeUint),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fc.mutation.NotesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   file.NotesTable,
			Columns: []string{file.NotesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(note.FieldID, field.TypeUint),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.note_files = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// FileCreateBulk is the builder for creating many File entities in bulk.
type FileCreateBulk struct {
	config
	err      error
	builders []*FileCreate
}

// Save creates the File entities in the database.
func (fcb *FileCreateBulk) Save(ctx context.Context) ([]*File, error) {
	if fcb.err != nil {
		return nil, fcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(fcb.builders))
	nodes := make([]*File, len(fcb.builders))
	mutators := make([]Mutator, len(fcb.builders))
	for i := range fcb.builders {
		func(i int, root context.Context) {
			builder := fcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*FileMutation)
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
					_, err = mutators[i+1].Mutate(root, fcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, fcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = uint(id)
				}
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
		if _, err := mutators[0].Mutate(ctx, fcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (fcb *FileCreateBulk) SaveX(ctx context.Context) []*File {
	v, err := fcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fcb *FileCreateBulk) Exec(ctx context.Context) error {
	_, err := fcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fcb *FileCreateBulk) ExecX(ctx context.Context) {
	if err := fcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// Code generated by ent, DO NOT EDIT.

package file

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the file type in the database.
	Label = "file"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldFilename holds the string denoting the filename field in the database.
	FieldFilename = "filename"
	// FieldIsDir holds the string denoting the is_dir field in the database.
	FieldIsDir = "is_dir"
	// FieldIsShared holds the string denoting the is_shared field in the database.
	FieldIsShared = "is_shared"
	// EdgeOwner holds the string denoting the owner edge name in mutations.
	EdgeOwner = "owner"
	// EdgeSharedWith holds the string denoting the shared_with edge name in mutations.
	EdgeSharedWith = "shared_with"
	// EdgeParent holds the string denoting the parent edge name in mutations.
	EdgeParent = "parent"
	// EdgeChildren holds the string denoting the children edge name in mutations.
	EdgeChildren = "children"
	// EdgeNotes holds the string denoting the notes edge name in mutations.
	EdgeNotes = "notes"
	// Table holds the table name of the file in the database.
	Table = "files"
	// OwnerTable is the table that holds the owner relation/edge.
	OwnerTable = "files"
	// OwnerInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	OwnerInverseTable = "users"
	// OwnerColumn is the table column denoting the owner relation/edge.
	OwnerColumn = "user_files"
	// SharedWithTable is the table that holds the shared_with relation/edge. The primary key declared below.
	SharedWithTable = "file_shared_with"
	// SharedWithInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	SharedWithInverseTable = "users"
	// ParentTable is the table that holds the parent relation/edge.
	ParentTable = "files"
	// ParentColumn is the table column denoting the parent relation/edge.
	ParentColumn = "file_children"
	// ChildrenTable is the table that holds the children relation/edge.
	ChildrenTable = "files"
	// ChildrenColumn is the table column denoting the children relation/edge.
	ChildrenColumn = "file_children"
	// NotesTable is the table that holds the notes relation/edge.
	NotesTable = "files"
	// NotesInverseTable is the table name for the Note entity.
	// It exists in this package in order to avoid circular dependency with the "note" package.
	NotesInverseTable = "notes"
	// NotesColumn is the table column denoting the notes relation/edge.
	NotesColumn = "note_files"
)

// Columns holds all SQL columns for file fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldFilename,
	FieldIsDir,
	FieldIsShared,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "files"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"file_children",
	"note_files",
	"user_files",
}

var (
	// SharedWithPrimaryKey and SharedWithColumn2 are the table columns denoting the
	// primary key for the shared_with relation (M2M).
	SharedWithPrimaryKey = []string{"file_id", "user_id"}
)

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
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// FilenameValidator is a validator for the "filename" field. It is called by the builders before save.
	FilenameValidator func(string) error
	// DefaultIsDir holds the default value on creation for the "is_dir" field.
	DefaultIsDir bool
	// DefaultIsShared holds the default value on creation for the "is_shared" field.
	DefaultIsShared bool
)

// OrderOption defines the ordering options for the File queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByFilename orders the results by the filename field.
func ByFilename(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFilename, opts...).ToFunc()
}

// ByIsDir orders the results by the is_dir field.
func ByIsDir(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsDir, opts...).ToFunc()
}

// ByIsShared orders the results by the is_shared field.
func ByIsShared(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsShared, opts...).ToFunc()
}

// ByOwnerField orders the results by owner field.
func ByOwnerField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOwnerStep(), sql.OrderByField(field, opts...))
	}
}

// BySharedWithCount orders the results by shared_with count.
func BySharedWithCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newSharedWithStep(), opts...)
	}
}

// BySharedWith orders the results by shared_with terms.
func BySharedWith(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSharedWithStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByParentField orders the results by parent field.
func ByParentField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newParentStep(), sql.OrderByField(field, opts...))
	}
}

// ByChildrenCount orders the results by children count.
func ByChildrenCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newChildrenStep(), opts...)
	}
}

// ByChildren orders the results by children terms.
func ByChildren(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newChildrenStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByNotesField orders the results by notes field.
func ByNotesField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newNotesStep(), sql.OrderByField(field, opts...))
	}
}
func newOwnerStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OwnerInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, OwnerTable, OwnerColumn),
	)
}
func newSharedWithStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SharedWithInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, SharedWithTable, SharedWithPrimaryKey...),
	)
}
func newParentStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(Table, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ParentTable, ParentColumn),
	)
}
func newChildrenStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(Table, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ChildrenTable, ChildrenColumn),
	)
}
func newNotesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(NotesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, NotesTable, NotesColumn),
	)
}

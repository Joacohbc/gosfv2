// Code generated by ent, DO NOT EDIT.

package ent

import (
	"gosfV2/src/ent/file"
	"gosfV2/src/ent/note"
	"gosfV2/src/ent/schema"
	"gosfV2/src/ent/user"
	"time"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	fileMixin := schema.File{}.Mixin()
	fileMixinFields0 := fileMixin[0].Fields()
	_ = fileMixinFields0
	fileFields := schema.File{}.Fields()
	_ = fileFields
	// fileDescCreatedAt is the schema descriptor for created_at field.
	fileDescCreatedAt := fileMixinFields0[0].Descriptor()
	// file.DefaultCreatedAt holds the default value on creation for the created_at field.
	file.DefaultCreatedAt = fileDescCreatedAt.Default.(func() time.Time)
	// fileDescUpdatedAt is the schema descriptor for updated_at field.
	fileDescUpdatedAt := fileMixinFields0[1].Descriptor()
	// file.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	file.UpdateDefaultUpdatedAt = fileDescUpdatedAt.UpdateDefault.(func() time.Time)
	// fileDescFilename is the schema descriptor for filename field.
	fileDescFilename := fileFields[1].Descriptor()
	// file.FilenameValidator is a validator for the "filename" field. It is called by the builders before save.
	file.FilenameValidator = func() func(string) error {
		validators := fileDescFilename.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(filename string) error {
			for _, fn := range fns {
				if err := fn(filename); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// fileDescIsDir is the schema descriptor for is_dir field.
	fileDescIsDir := fileFields[2].Descriptor()
	// file.DefaultIsDir holds the default value on creation for the is_dir field.
	file.DefaultIsDir = fileDescIsDir.Default.(bool)
	// fileDescIsShared is the schema descriptor for is_shared field.
	fileDescIsShared := fileFields[3].Descriptor()
	// file.DefaultIsShared holds the default value on creation for the is_shared field.
	file.DefaultIsShared = fileDescIsShared.Default.(bool)
	noteMixin := schema.Note{}.Mixin()
	noteMixinFields0 := noteMixin[0].Fields()
	_ = noteMixinFields0
	noteFields := schema.Note{}.Fields()
	_ = noteFields
	// noteDescCreatedAt is the schema descriptor for created_at field.
	noteDescCreatedAt := noteMixinFields0[0].Descriptor()
	// note.DefaultCreatedAt holds the default value on creation for the created_at field.
	note.DefaultCreatedAt = noteDescCreatedAt.Default.(func() time.Time)
	// noteDescUpdatedAt is the schema descriptor for updated_at field.
	noteDescUpdatedAt := noteMixinFields0[1].Descriptor()
	// note.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	note.UpdateDefaultUpdatedAt = noteDescUpdatedAt.UpdateDefault.(func() time.Time)
	// noteDescTitle is the schema descriptor for title field.
	noteDescTitle := noteFields[1].Descriptor()
	// note.TitleValidator is a validator for the "title" field. It is called by the builders before save.
	note.TitleValidator = func() func(string) error {
		validators := noteDescTitle.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(title string) error {
			for _, fn := range fns {
				if err := fn(title); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// noteDescContent is the schema descriptor for content field.
	noteDescContent := noteFields[2].Descriptor()
	// note.ContentValidator is a validator for the "content" field. It is called by the builders before save.
	note.ContentValidator = noteDescContent.Validators[0].(func(string) error)
	userMixin := schema.User{}.Mixin()
	userMixinFields0 := userMixin[0].Fields()
	_ = userMixinFields0
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreatedAt is the schema descriptor for created_at field.
	userDescCreatedAt := userMixinFields0[0].Descriptor()
	// user.DefaultCreatedAt holds the default value on creation for the created_at field.
	user.DefaultCreatedAt = userDescCreatedAt.Default.(func() time.Time)
	// userDescUpdatedAt is the schema descriptor for updated_at field.
	userDescUpdatedAt := userMixinFields0[1].Descriptor()
	// user.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	user.UpdateDefaultUpdatedAt = userDescUpdatedAt.UpdateDefault.(func() time.Time)
	// userDescUsername is the schema descriptor for username field.
	userDescUsername := userFields[1].Descriptor()
	// user.UsernameValidator is a validator for the "username" field. It is called by the builders before save.
	user.UsernameValidator = userDescUsername.Validators[0].(func(string) error)
	// userDescPassword is the schema descriptor for password field.
	userDescPassword := userFields[2].Descriptor()
	// user.PasswordValidator is a validator for the "password" field. It is called by the builders before save.
	user.PasswordValidator = userDescPassword.Validators[0].(func(string) error)
}
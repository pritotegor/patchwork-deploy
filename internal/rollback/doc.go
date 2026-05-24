// Package rollback provides a lightweight rollback stack for patchwork-deploy.
//
// As patches are successfully applied during a deployment run, callers
// Record each patch along with an optional undo script path. If a later
// patch fails, Execute walks the stack in reverse order and invokes the
// caller-supplied exec function for each entry that has an undo script.
//
// Usage:
//
//	rb := rollback.New(os.Stdout)
//	rb.Record("001-create-table.sql", "001-drop-table.sql")
//	rb.Record("002-add-index.sql", "")
//
//	if err := rb.Execute(func(undoPath string) error {
//		return sshClient.RunScript(undoPath)
//	}); err != nil {
//		log.Fatal(err)
//	}
package rollback

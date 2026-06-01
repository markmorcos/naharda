// Package migrations embeds the goose SQL migration files.
package migrations

import "embed"

// FS holds the embedded migration SQL files.
//
//go:embed *.sql
var FS embed.FS

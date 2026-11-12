// Package web holds pages that are built into the app,
// so they always ship with it.
package web

import _ "embed"

//go:embed key-test.html
var KeyTest []byte

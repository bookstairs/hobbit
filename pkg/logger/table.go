package logger

import (
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
)

type TableLogger interface {
	Title(title string) TableLogger // Title adds the table title.
	Head(heads ...any) TableLogger  // Head adds the table head.
	Row(fields ...any) TableLogger  // Row add a row to the table.
	AllowZeroValue() TableLogger    // AllowZeroValue The row will be printed if it contains zero value.
	Log()                           // Print would print a table-like message from the given config.
}

type tableLogger struct {
	title     string
	heads     []any
	rows      [][]any
	allowZero bool
}

func (t *tableLogger) Title(title string) TableLogger {
	t.title = title
	return t
}

func (t *tableLogger) Head(heads ...any) TableLogger {
	t.heads = heads
	return t
}

func (t *tableLogger) Row(fields ...any) TableLogger {
	if len(fields) > 0 {
		t.rows = append(t.rows, fields)
	}
	return t
}

func (t *tableLogger) AllowZeroValue() TableLogger {
	t.allowZero = true
	return t
}

func (t *tableLogger) Log() {
	w := table.NewWriter()
	w.SetOutputMirror(os.Stdout)
	w.SetTitle(t.title)
	if len(t.heads) > 0 {
		w.AppendHeader(t.heads)
	}

	for _, row := range t.rows {
		if len(row) == 1 {
			w.AppendRow([]any{row[0]})
		} else {
			zero := true
			for _, r := range row[1:] {
				v := reflect.ValueOf(r)
				if !v.IsZero() {
					zero = false
					break
				}
			}

			if !zero || t.allowZero {
				w.AppendRow(row)
			}
		}
	}

	w.Render()
}

// NewTableLogger will return a printer for table-like logs.
func NewTableLogger() TableLogger {
	return &tableLogger{}
}

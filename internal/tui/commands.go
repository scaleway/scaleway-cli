package tui

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	zoneType   = reflect.TypeFor[scw.Zone]()
	regionType = reflect.TypeFor[scw.Region]()
)

// ResourcesFromCommands builds one Resource per `list` command in the
// CLI registry, so every listable Scaleway resource is browsable in the
// TUI. On a name collision the first command wins, matching how the
// cobra builder resolves duplicate command paths. Fetches read the
// session zone/region from loc at fetch time, so the locality can
// change during the session.
func ResourcesFromCommands(
	cmds []*core.Command,
	meta *core.Meta,
	loc *Locality,
) map[string]*Resource {
	resources := map[string]*Resource{}
	for _, cmd := range cmds {
		if cmd.Run == nil || !cmd.IsList() {
			continue
		}
		name := cmd.Namespace
		if cmd.Resource != "" {
			name += " " + cmd.Resource
		}
		if _, ok := resources[name]; ok {
			continue
		}
		resources[name] = &Resource{
			Name:  name,
			Title: name,
			Fetch: commandFetch(cmd, meta, loc),
		}
	}

	return resources
}

func commandFetch(
	cmd *core.Command,
	meta *core.Meta,
	loc *Locality,
) func(context.Context) (listResult, error) {
	return func(ctx context.Context) (listResult, error) {
		zone, region := loc.Get()
		var args any
		if cmd.ArgsType != nil {
			args = reflect.New(cmd.ArgsType).Interface()
			fillLocalities(args, zone, region)
		}
		if meta != nil {
			ctx = core.InjectMeta(ctx, meta)
		}
		// Run through the command interceptor chain, exactly like the
		// CLI does: interceptors may rewrite the args before Run sees
		// them (e.g. instance server list passes its embedded SDK
		// request). Localities are filled again on the final args.
		runner := func(ctx context.Context, argsI any) (any, error) {
			fillLocalities(argsI, zone, region)

			return cmd.Run(ctx, argsI)
		}
		var resp any
		var err error
		if cmd.Interceptor != nil {
			resp, err = cmd.Interceptor(ctx, args, runner)
		} else {
			resp, err = runner(ctx, args)
		}
		if err != nil {
			return listResult{}, err
		}

		return toRows(resp)
	}
}

// fillLocalities sets the default zone/region on the request args. The
// cobra layer does this from the config when running the CLI; here the
// zone/region are the ones chosen for the TUI session. An explicit
// region takes precedence; otherwise it is derived from the zone.
func fillLocalities(args any, zone scw.Zone, region scw.Region) {
	v := reflect.ValueOf(args)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}
	t := v.Type()
	for i := range t.NumField() {
		f := v.Field(i)
		if !f.CanSet() {
			continue
		}
		switch f.Type() {
		case zoneType:
			if f.String() == "" && zone != "" {
				f.Set(reflect.ValueOf(zone))
			}
		case regionType:
			if f.String() == "" {
				if region != "" {
					f.Set(reflect.ValueOf(region))
				} else if r, err := zone.Region(); err == nil {
					f.Set(reflect.ValueOf(r))
				}
			}
		}
	}
}

// toRows turns a list command result into table data. Most generated
// commands return a slice (`return resp.Servers`); some wrap it in a
// struct, in which case the first exported slice field is used.
func toRows(resp any) (listResult, error) {
	v := reflect.ValueOf(resp)
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return listResult{}, errors.New("empty response")
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return sliceToRows(v)
	case reflect.Struct:
		for i := range v.Type().NumField() {
			ft := v.Type().Field(i)
			f := v.Field(i)
			if !ft.IsExported() || (f.Kind() != reflect.Slice && f.Kind() != reflect.Array) {
				continue
			}

			return sliceToRows(f)
		}
	}

	return listResult{}, fmt.Errorf("no list found in response of type %T", resp)
}

// preferredColumns is the column order of the generic table: the fields
// every Scaleway resource exposes, in the k9s-ish ID/NAME/STATE order.
var preferredColumns = []struct {
	header string
	fields []string
}{
	{"ID", []string{"ID"}},
	{"NAME", []string{"Name"}},
	{"STATE", []string{"State", "Status"}},
	{"ZONE", []string{"Zone"}},
	{"REGION", []string{"Region"}},
	{"CREATED", []string{"CreationDate", "CreatedAt", "Created"}},
}

// pickColumns selects the table columns for an element type: the
// preferred fields that exist, falling back to the first scalar fields.
func pickColumns(t reflect.Type) (columns []string, fields []int, stateCol int) {
	stateCol = -1
	if t == nil || t.Kind() != reflect.Struct {
		return nil, nil, stateCol
	}
	fieldIdx := func(name string) int {
		for i := range t.NumField() {
			f := t.Field(i)
			if f.IsExported() && f.Name == name {
				return i
			}
		}

		return -1
	}
	for _, pc := range preferredColumns {
		for _, fn := range pc.fields {
			if i := fieldIdx(fn); i >= 0 {
				columns = append(columns, pc.header)
				fields = append(fields, i)
				if pc.header == "STATE" {
					stateCol = len(columns) - 1
				}

				break
			}
		}
	}
	if len(columns) > 0 {
		return columns, fields, stateCol
	}
	// Fallback: first scalar-ish fields (types without ID/Name/State).
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() || !isScalarish(f.Type) {
			continue
		}
		columns = append(columns, strings.ToUpper(f.Name))
		fields = append(fields, i)
		if len(columns) >= 5 {
			break
		}
	}

	return columns, fields, stateCol
}

func isScalarish(t reflect.Type) bool {
	for t.Kind() == reflect.Pointer || t.Kind() == reflect.Interface {
		if t.Kind() == reflect.Interface {
			return false
		}
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}

	return t == timeType
}

func sliceToRows(v reflect.Value) (listResult, error) {
	elemT := v.Type().Elem()
	for elemT.Kind() == reflect.Pointer || elemT.Kind() == reflect.Interface {
		elemT = elemT.Elem()
	}
	columns, fields, stateCol := pickColumns(elemT)
	if len(columns) == 0 {
		columns, fields = []string{"VALUE"}, []int{-1}
	}
	idCol := -1
	for c, h := range columns {
		if h == "ID" {
			idCol = c
		}
	}
	n := v.Len()
	rows := make([]Row, 0, n)
	for i := range n {
		el := v.Index(i)
		row := Row{Obj: el.Interface(), Cells: make([]string, len(columns))}
		ev := el
		for ev.Kind() == reflect.Pointer {
			if ev.IsNil() {
				break
			}
			ev = ev.Elem()
		}
		for c, fi := range fields {
			fv := ev
			if fi >= 0 {
				fv = ev.Field(fi)
			}
			row.Cells[c] = formatCell(columns[c], fv)
		}
		if idCol >= 0 {
			row.ID = rawString(ev.Field(fields[idCol]))
		}
		if stateCol >= 0 && row.Cells[stateCol] != "-" {
			row.CellStyles = map[int]lipgloss.Style{stateCol: stateStyle(row.Cells[stateCol])}
		}
		rows = append(rows, row)
	}

	return listResult{columns: columns, rows: rows}, nil
}

// formatCell renders one table cell.
func formatCell(header string, fv reflect.Value) string {
	for fv.Kind() == reflect.Pointer || fv.Kind() == reflect.Interface {
		if fv.IsNil() {
			return "-"
		}
		fv = fv.Elem()
	}
	if fv.Type() == timeType {
		return fv.Interface().(time.Time).Format("2006-01-02 15:04")
	}
	cell := formatValue(fv)
	if header == "ID" && cell != "-" {
		cell = shortID(cell)
	}

	return cell
}

// rawString returns the string content of a (possibly pointer) field.
func rawString(fv reflect.Value) string {
	for fv.Kind() == reflect.Pointer || fv.Kind() == reflect.Interface {
		if fv.IsNil() {
			return ""
		}
		fv = fv.Elem()
	}
	if fv.Kind() == reflect.String {
		return fv.String()
	}

	return ""
}

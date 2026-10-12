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
	// Index every command by its "ns res verb" path (first wins, like
	// the cobra builder) so a resource can find its get/update sibling.
	byName := map[string]*core.Command{}
	for _, cmd := range cmds {
		if cmd.Run == nil {
			continue
		}
		if name := commandName(cmd); name != "" {
			if _, ok := byName[name]; !ok {
				byName[name] = cmd
			}
		}
	}

	resources := map[string]*Resource{}
	listCmds := map[string]*core.Command{}
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
		res := &Resource{
			Name:  name,
			Title: name,
			Fetch: commandFetch(cmd, meta, loc),
		}
		if get, ok := byName[name+" get"]; ok {
			res.GetCmd = get
		}
		if update, ok := byName[name+" update"]; ok {
			res.UpdateCmd = update
		}
		if del, ok := byName[name+" delete"]; ok {
			res.DeleteCmd = del
		}
		resources[name] = res
		listCmds[name] = cmd
	}
	linkSubResources(resources, listCmds, meta, loc)

	return resources
}

// subResources maps a child list to its parent: the child's list
// command requires the parent row's ID in the request field named here.
// Children are hidden from the goto menu and the `:` prompt; the parent
// browser gains the drill key ('s') to list the child for the selected
// row.
//
// ponytail: by hand, one sub per parent (a second entry for the same
// parent would overwrite the first); auto-detect from the list
// commands' required *-id arg specs if the table stops scaling.
var subResources = map[string]struct {
	parent string
	field  string
}{
	"k8s pool": {parent: "k8s cluster", field: "ClusterID"},
}

// linkSubResources wires the subResources table onto the built
// resources: the child learns its parent (menu/prompt filtering), the
// parent learns how to fetch the child list scoped to one row.
func linkSubResources(
	resources map[string]*Resource,
	listCmds map[string]*core.Command,
	meta *core.Meta,
	loc *Locality,
) {
	for name, spec := range subResources {
		child, ok := resources[name]
		parent, ok2 := resources[spec.parent]
		cmd, ok3 := listCmds[name]
		if !ok || !ok2 || !ok3 {
			continue
		}
		child.Parent = spec.parent
		parent.Sub = &SubResource{
			Res: child,
			Scoped: func(parentID string) func(context.Context) (listResult, error) {
				return fetchWith(cmd, meta, loc, func(args any) {
					setStringField(args, spec.field, parentID)
				})
			},
		}
	}
}

// setStringField sets a string field on the (pointer to) request args.
func setStringField(args any, field, value string) {
	v := reflect.ValueOf(args)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	if f := v.FieldByName(field); f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
		f.SetString(value)
	}
}

// commandName is the full command path "ns res verb" (empty parts
// skipped), the way the cobra builder assembles command paths.
func commandName(cmd *core.Command) string {
	var parts []string
	for _, p := range []string{cmd.Namespace, cmd.Resource, cmd.Verb} {
		if p != "" {
			parts = append(parts, p)
		}
	}

	return strings.Join(parts, " ")
}

func commandFetch(
	cmd *core.Command,
	meta *core.Meta,
	loc *Locality,
) func(context.Context) (listResult, error) {
	return fetchWith(cmd, meta, loc, nil)
}

// fetchWith builds a list fetch for cmd; fill, if non-nil, adjusts the
// request args after the localities are filled (the drill-down uses it
// to pin the list to one parent row).
func fetchWith(
	cmd *core.Command,
	meta *core.Meta,
	loc *Locality,
	fill func(any),
) func(context.Context) (listResult, error) {
	return func(ctx context.Context) (listResult, error) {
		zone, region := loc.Get()
		var args any
		if cmd.ArgsType != nil {
			args = reflect.New(cmd.ArgsType).Interface()
			fillLocalities(args, zone, region)
			if fill != nil {
				fill(args)
			}
		}
		if meta != nil {
			ctx = core.InjectMeta(ctx, meta)
		}
		resp, err := runCommand(cmd, meta, loc, ctx, args)
		if err != nil {
			return listResult{}, err
		}

		return toRows(resp)
	}
}

// runCommand runs a CLI command the TUI way: through the interceptor
// chain (exactly like the CLI does: interceptors may rewrite the args
// before Run sees them) with the session locality re-filled on the
// final args. The caller fills the localities on args beforehand:
// interceptors may read them before the runner re-fills them.
func runCommand(
	cmd *core.Command,
	meta *core.Meta,
	loc *Locality,
	ctx context.Context,
	args any,
) (any, error) {
	zone, region := loc.Get()
	if meta != nil {
		ctx = core.InjectMeta(ctx, meta)
	}
	runner := func(ctx context.Context, argsI any) (any, error) {
		fillLocalities(argsI, zone, region)

		return cmd.Run(ctx, argsI)
	}
	if cmd.Interceptor != nil {
		return cmd.Interceptor(ctx, args, runner)
	}

	return runner(ctx, args)
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

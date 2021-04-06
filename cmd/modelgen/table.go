package main

import (
	"fmt"
	"strings"

	"github.com/ebay/libovsdb"
)

/*TableGenerator is a code generator capable of generating a file such as this:
package main

import (
	goovn "github.com/ebay/go-ovn"
)

//Chassis struct defines an object in Chassis table
type LogicalRouter struct {
	goovn.BaseModel
	Name         string            `ovs:"name"`
	StaticRoutes []string          `ovs:"static_routes"`
	Nat          []string          `ovs:"nat"`
	ExternalIds  map[string]string `ovs:"external_ids"`
	Ports        []string          `ovs:"ports"`
	LoadBalancer []string          `ovs:"load_balancer"`
	Options      map[string]string `ovs:"options"`
	Policies     []string          `ovs:"policies"`
	Enabled      []bool            `ovs:"enabled"`
}

func (lr *LogicalRouter) Table() goovn.TableName {
	return "Logical_Router"
}

func NewRouter(uuid string) goovn.Model {
	return &LogicalRouter{
		BaseModel: goovn.BaseModel{UUID: uuid},
	}
}
*/
type TableGenerator struct {
	Generator
	name  string
	table libovsdb.TableSchema
}

func (g *TableGenerator) Generate() error {
	g.PrintHeader()
	g.Printf("\n")
	// Add imports
	imports := [][]string{
		{"goovn", "github.com/ebay/go-ovn"},
	}
	g.PrintImports(imports)
	g.Printf("\n")

	vars := [][]string{
		{g.TableVarName(), "goovn.TableName", "=", fmt.Sprintf("\"%s\"", g.name)},
	}
	g.PrintVars(vars)
	g.Printf("\n")
	g.PrintTableStruct()
	//g.Printf("\n")
	//g.PrintApi()
	return nil
}

/*PrintTableStruct prints a struct called `TableName` that implements
goovn.Model and that inherits from BaseModel

type LogicalRouter struct {
	goovn.BaseModel
	Name         string            `ovs:"name"`
	StaticRoutes []string          `ovs:"static_routes"`
	Nat          []string          `ovs:"nat"`
	ExternalIds  map[string]string `ovs:"external_ids"`
	Ports        []string          `ovs:"ports"`
	LoadBalancer []string          `ovs:"load_balancer"`
	Options      map[string]string `ovs:"options"`
	Policies     []string          `ovs:"policies"`
	Enabled      []bool            `ovs:"enabled"`
}

func (*LogicalRouter) Table() goovn.TableName {
	return "Logical_Router"
}

func NewRouter(uuid string) goovn.Model {
	return &LogicalRouter{
		BaseModel: goovn.BaseModel{UUID: uuid},
	}
}
*/
func (g *TableGenerator) PrintTableStruct() {
	structName := g.TableTypeName()
	g.Printf("//%s struct defines an object in %s table\n",
		g.TableTypeName(), g.name)

	g.Printf("type %s struct {\n", structName)
	// TODO: Inherit from BaseModel
	g.Printf("\tUUID\tstring\t%s\n", g.UidTagString())
	for colName, colSchema := range g.table.Columns {
		g.Printf("\t%s\t%s\t%s\n", g.ColumnFieldName(colName),
			g.ColumnTypeString(colSchema),
			g.TagString(colName))
	}
	g.Printf("}\n")
	g.Printf("\n")
	// Print the Table() func
	g.Printf("// Table returns the table name. It's part of the Model interface\n")
	g.Printf("func (*%s) Table() goovn.TableName {\n", structName)
	g.Printf("\treturn %s", g.TableVarName())
	g.Printf("}\n")
}

func (g *TableGenerator) TableApiStruct() string {
	return g.TableTypeName() + "Api"
}
func (g *TableGenerator) TableVarName() string {
	return fmt.Sprintf("%sTable", g.TableTypeName())
}
func (g *TableGenerator) FileName() string {
	return g.TableFileName(g.name)
}
func (g *TableGenerator) TagString(column string) string {
	return fmt.Sprintf("`ovs:\"%s\"`", column)
}
func (g *TableGenerator) UidTagString() string {
	return fmt.Sprintf("`ovs:\"_uuid\"`")
}

/* e.g
DatapathBindingApi {
    client: c,
},
*/
func (g *TableGenerator) ApiInitString(tabLevel int) string {
	return strings.Join([]string{
		fmt.Sprintf("%s {", g.TableApiStruct()),
		fmt.Sprintf("%sclient: c,", strings.Repeat("\t", tabLevel+1)),
		fmt.Sprintf("%s}", strings.Repeat("\t", tabLevel)),
	}, "\n")
}

func (g *TableGenerator) ExtendedTypeString(etype libovsdb.ExtendedType) string {
	switch etype {
	case libovsdb.TypeInteger:
		return "int"
	case libovsdb.TypeReal:
		return "float64"
	case libovsdb.TypeBoolean:
		return "bool"
	case libovsdb.TypeString:
		return "string"
	case libovsdb.TypeUUID:
		return "string"
	}
	return ""
}

func (g *TableGenerator) ColumnTypeString(column *libovsdb.ColumnSchema) string {
	// Type can be either a string, such as "string" or "float"
	// or a map describing a more complex type
	//log.Printf(reflect.TypeOf(column.Type).String())

	switch column.Type {
	case libovsdb.TypeEnum:
		return g.ExtendedTypeString(column.TypeObj.Key.Type)
	case libovsdb.TypeMap:
		return fmt.Sprintf("map[%s]%s", g.ExtendedTypeString(column.TypeObj.Key.Type),
			g.ExtendedTypeString(column.TypeObj.Value.Type))
	case libovsdb.TypeSet:
		return fmt.Sprintf("[]%s", g.ExtendedTypeString(column.TypeObj.Key.Type))
	default:
		return g.ExtendedTypeString(column.Type)
	}
}

/* tunnel_key -> TunnelKey */
func (g *TableGenerator) ColumnFieldName(field string) string {
	return camelCase(field)
}

/*Logical_Flow -> LogicalFlow*/
func (g *TableGenerator) TableTypeName() string {
	return strings.ReplaceAll(g.name, "_", "")
}

/* Logical_Flow -> logical_flow.go */
func (g *TableGenerator) TableFileName(table string) string {
	return fmt.Sprintf("%s.go", strings.ToLower(table))
}

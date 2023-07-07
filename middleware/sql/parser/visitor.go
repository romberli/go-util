package parser

import (
	"reflect"
	"strings"

	"github.com/pingcap/parser/ast"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"

	driver "github.com/pingcap/parser/test_driver"
)

const (
	CreateTableStmtString = "*ast.CreateTableStmt"
	AlterTableStmtString  = "*ast.AlterTableStmt"
	DropTableStmtString   = "*ast.DropTableStmt"
	SelectStmtString      = "*ast.SelectStmt"
	UnionStmtString       = "*ast.UnionStmt"
	InsertStmtString      = "*ast.InsertStmt"
	ReplaceStmtString     = "*ast.ReplaceStmt"
	UpdateStmtString      = "*ast.UpdateStmt"
	DeleteStmtString      = "*ast.DeleteStmt"

	FuncCallExprString      = "*ast.FuncCallExpr"
	AggregateFuncExprString = "*ast.AggregateFuncExpr"
	WindowFuncExprString    = "*ast.WindowFuncExpr"
)

var (
	DefaultSQLList = []string{
		CreateTableStmtString,
		AlterTableStmtString,
		DropTableStmtString,
		SelectStmtString,
		UnionStmtString,
		InsertStmtString,
		ReplaceStmtString,
		UpdateStmtString,
		DeleteStmtString,
	}
	DefaultFuncList = []string{
		FuncCallExprString,
		AggregateFuncExprString,
		WindowFuncExprString,
	}
)

type Visitor struct {
	toParse  bool
	sqlList  []string
	funcList []string
	result   *Result
}

// NewVisitor returns a new *Visitor
func NewVisitor(sqlList, funcList []string) *Visitor {
	return &Visitor{
		sqlList:  sqlList,
		funcList: funcList,
		result:   NewEmptyResult(),
	}
}

// NewVisitorWithDefault returns a new *Visitor with default sql list and function list
func NewVisitorWithDefault() *Visitor {
	return &Visitor{
		sqlList:  DefaultSQLList,
		funcList: DefaultFuncList,
		result:   NewEmptyResult(),
	}
}

// GetSQLList returns the sql list
func (v *Visitor) GetSQLList() []string {
	return v.sqlList
}

// GetFuncList returns the function list
func (v *Visitor) GetFuncList() []string {
	return v.funcList
}

// GetResult returns the result
func (v *Visitor) GetResult() *Result {
	return v.result
}

// Enter enters into the given node, it will traverse each child node to find useful information such as table name, column name...
// note that it only traverses some kinds of node types, see the constants at the top of this file
func (v *Visitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	astType := reflect.TypeOf(in).String()

	if common.StringInSlice(v.sqlList, astType) {
		v.toParse = true
		// set sql type
		v.result.SetSQLType(strings.Split(astType, constant.DotString)[1])
	}

	if v.toParse {
		switch node := in.(type) {
		case *ast.TableName:
			v.visitTableName(node)
		case *ast.CreateTableStmt:
			v.visitCreateTableStmt(node)
		case *ast.AlterTableStmt:
			v.visitAlterTableStmt(node)
		case *ast.SelectField:
			v.visitSelectField(node)
		case *ast.ColumnDef:
			v.visitColumnDef(node)
		case *ast.ColumnName:
			v.visitColumnName(node)
		}
	}

	return in, false
}

// Leave leaves the given node, traversal is over
func (v *Visitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

// visitTableName visits the given node which type is *ast.TableName
func (v *Visitor) visitTableName(node *ast.TableName) {
	tableName := node.Name.L
	dbName := node.Schema.L

	v.result.AddTableDBListMap(tableName, dbName)
	v.result.AddDBName(dbName)
	v.result.AddTableName(tableName)
}

// visitCreateTableStmt visits the given node which type is *ast.CreateTableStmt
func (v *Visitor) visitCreateTableStmt(node *ast.CreateTableStmt) {
	for _, tableOption := range node.Options {
		if tableOption.Tp == ast.TableOptionComment {
			v.result.SetTableComment(node.Table.Name.L, tableOption.StrValue)
			break
		}
	}
}

// visitAlterTableStmt visits the given node which type is *ast.AlterTableStmt
func (v *Visitor) visitAlterTableStmt(node *ast.AlterTableStmt) {
	for _, tableSpec := range node.Specs {
		for _, tableOption := range tableSpec.Options {
			if tableOption.Tp == ast.TableOptionComment {
				v.result.SetTableComment(node.Table.Name.L, tableOption.StrValue)
				break
			}
		}
	}
}

// visitSelectField visits the given node which type is *ast.SelectField
func (v *Visitor) visitSelectField(node *ast.SelectField) {
	var funcArgs []ast.ExprNode

	expr := node.Expr
	if expr == nil && node.WildCard != nil {
		v.result.AddColumn(constant.AsteriskString)
	} else if expr != nil {
		switch e := expr.(type) {
		case *ast.AggregateFuncExpr:
			funcArgs = e.Args
		case *ast.FuncCallExpr:
			funcArgs = e.Args
		case *ast.WindowFuncExpr:
			funcArgs = e.Args
		case *ast.ColumnNameExpr:
			v.result.AddColumn(e.Name.Name.L)
		}

		for _, arg := range funcArgs {
			switch e := arg.(type) {
			case *ast.ColumnNameExpr:
				v.result.AddColumn(e.Name.Name.L)
			}
		}
	}
}

// visitColumnDef visits the given node which type is *ast.ColumnDef
func (v *Visitor) visitColumnDef(node *ast.ColumnDef) {
	var columnComment string

	columnName := node.Name.Name.L

	v.result.AddColumn(columnName)
	v.result.SetColumnType(columnName, node.Tp.InfoSchemaStr())

	for _, columnOption := range node.Options {
		if columnOption.Tp == ast.ColumnOptionComment {
			columnComment = columnOption.Expr.(*driver.ValueExpr).GetDatumString()
		}
	}

	v.result.SetColumnComment(columnName, columnComment)
}

// visitColumnName visits the given node which type is *ast.ColumnName
func (v *Visitor) visitColumnName(node *ast.ColumnName) {
	v.result.AddColumn(node.Name.L)
}

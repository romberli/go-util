package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pingcap/errors"
	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pingcap/tidb/pkg/parser/types"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"

	driver "github.com/pingcap/tidb/pkg/parser/test_driver"
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
	GrantStmtString       = "*ast.GrantStmt"

	FuncCallExprString      = "*ast.FuncCallExpr"
	AggregateFuncExprString = "*ast.AggregateFuncExpr"
	WindowFuncExprString    = "*ast.WindowFuncExpr"

	CurrentTimeStampFuncName = "current_timestamp"
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
		GrantStmtString,
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

	parseTableDefinition bool
	tableDefinition      *TableDefinition
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

// SetParseTableDefinition sets the flag to parse table definition
func (v *Visitor) SetParseTableDefinition(parseTableDefinition bool) {
	v.parseTableDefinition = parseTableDefinition
	v.tableDefinition = NewEmptyTableDefinition()
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

	if common.ElementInSlice(v.sqlList, astType) {
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
		case *ast.GrantStmt:
			v.visitGrantStmt(node)
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
	if v.parseTableDefinition {
		tableSchema := node.Table.Schema.L
		tableName := node.Table.Name.L
		v.tableDefinition.TableSchema = tableSchema
		v.tableDefinition.TableName = tableName
		// table options
		for _, tableOption := range node.Options {
			switch tableOption.Tp {
			case ast.TableOptionEngine:
				v.tableDefinition.TableEngine = tableOption.StrValue
			case ast.TableOptionCharset:
				v.tableDefinition.Charset = tableOption.StrValue
			case ast.TableOptionCollate:
				v.tableDefinition.Collation = tableOption.StrValue
			case ast.TableOptionComment:
				v.tableDefinition.TableComment = tableOption.StrValue
			case ast.TableOptionRowFormat:
				v.tableDefinition.RowFormat = tableOption.StrValue
			}
		}
		// column definitions
		for i, column := range node.Cols {
			columnName := column.Name.Name.L
			cd := NewColumnDefinition(tableSchema, tableName, columnName)
			cd.OrdinalPosition = i + constant.OneInt
			cd.CharacterSetName = column.Tp.GetCharset()
			cd.CollationName = column.Tp.GetCollate()
			fieldType := column.Tp.GetType()
			cd.DataType = types.TypeToStr(fieldType, cd.CharacterSetName)
			cd.ColumnType = column.Tp.InfoSchemaStr()

			for _, option := range column.Options {
				switch option.Tp {
				case ast.ColumnOptionDefaultValue:
					switch expr := option.Expr.(type) {
					case *driver.ValueExpr:
						cd.DefaultValue = common.ConvertInterfaceToString(expr.GetValue())
					case *ast.FuncCallExpr:
						if expr.FnName.L == CurrentTimeStampFuncName {
							var args []string
							for _, arg := range expr.Args {
								strVal := common.ConvertInterfaceToString(arg.(*driver.ValueExpr).GetValue())
								args = append(args, strVal)
							}
							cd.DefaultValue = GetFullFuncName(expr.FnName.L, args...)
						}
					default:
						err := errors.Errorf("unknown default value expression type. columnName: %s", columnName)
						cd.AddError(err)
					}
				case ast.ColumnOptionCollate:
					cd.CollationName = option.StrValue
				case ast.ColumnOptionNotNull:
					cd.NotNull = true
				case ast.ColumnOptionAutoIncrement:
					cd.IsAutoIncrement = true
					cd.NotNull = true
				case ast.ColumnOptionComment:
					cd.ColumnComment = option.Expr.(*driver.ValueExpr).GetDatumString()
				case ast.ColumnOptionOnUpdate:
					switch expr := option.Expr.(type) {
					case *driver.ValueExpr:
						cd.OnUpdateValue = common.ConvertInterfaceToString(expr.GetValue())
					case *ast.FuncCallExpr:
						if expr.FnName.L == CurrentTimeStampFuncName {
							var args []string
							for _, arg := range expr.Args {
								strVal := common.ConvertInterfaceToString(arg.(*driver.ValueExpr).GetValue())
								args = append(args, strVal)
							}
							cd.OnUpdateValue = GetFullFuncName(expr.FnName.L, args...)
						}
					default:
						err := errors.Errorf("unknown default value expression type. columnName: %s", columnName)
						cd.AddError(err)
					}
				default:
					err := errors.Errorf("unknown column option. columnName: %s, optionType: %d", columnName, option.Tp)
					cd.AddError(err)
				}
			}

			v.tableDefinition.AddColumn(cd)
		}
		// index definition
		for _, constraint := range node.Constraints {
			indexName := constraint.Name
			id := NewIndexDefinition(tableSchema, tableName, indexName)
			id.HandleOption(constraint.Option)

			switch constraint.Tp {
			case ast.ConstraintPrimaryKey:
				id.IndexName = PrimaryKeyName
				id.IsPrimary = true
				id.IsUnique = true

				for _, column := range constraint.Keys {
					columnName := column.Column.Name.L
					cd := v.tableDefinition.GetColumnDefinition(columnName)
					cd.NotNull = true
					is := NewIndexSpec(cd, column.Desc, column.Length)
					id.AddIndexSpec(is)
				}
			case ast.ConstraintIndex, ast.ConstraintUniq:
				id.IndexName = constraint.Name
				if constraint.Tp == ast.ConstraintUniq {
					id.IsUnique = true
				}

				for _, column := range constraint.Keys {
					columnName := column.Column.Name.L
					cd := v.tableDefinition.GetColumnDefinition(columnName)
					is := NewIndexSpec(cd, column.Desc, column.Length)
					id.AddIndexSpec(is)
				}
			default:
				err := errors.Errorf("unknown index type. indexName: %s, indexType: %d", indexName, constraint.Tp)
				id.AddError(err)
			}

			v.tableDefinition.AddIndex(id)
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

// visitGrantStmt visits the given node which type is *ast.GrantStmt
func (v *Visitor) visitGrantStmt(node *ast.GrantStmt) {
	if len(node.Users) > constant.ZeroInt {
		user := node.Users[constant.ZeroInt]
		fullUserName := fmt.Sprintf("'%s'@'%s'", user.User.Username, user.User.Hostname)
		v.result.SetUser(fullUserName)
	}

	for _, priv := range node.Privs {
		v.result.AddPrivilege(priv.Priv, node.WithGrant)
	}

	v.result.AddDBName(node.Level.DBName)

	if node.Level.TableName != constant.EmptyString {
		v.result.AddTableName(node.Level.TableName)
		v.result.AddTableDBListMap(node.Level.TableName, node.Level.DBName)
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

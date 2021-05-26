package parser

import (
	"reflect"
	"strings"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/tidb/types/parser_driver"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

var (
	DefaultSQLList = []string{
		"*ast.CreateTableStmt",
		"*ast.AlterTableStmt",
		"*ast.DropTableStmt",
		"*ast.SelectStmt",
		"*ast.UnionStmt",
		"*ast.InsertStmt",
		"*ast.ReplaceStmt",
		"*ast.InsertStmt",
		"*ast.UpdateStmt",
		"*ast.DeleteStmt",
	}
	DefaultFuncList = []string{
		"*ast.FuncCallExpr",
		"*ast.AggregateFuncExpr",
		"*ast.WindowFuncExpr",
	}
)

type Visitor struct {
	ToParse  bool
	SQLList  []string
	FuncList []string
	Result   *Result
}

func NewVisitor() *Visitor {
	return &Visitor{
		ToParse:  false,
		SQLList:  DefaultSQLList,
		FuncList: DefaultFuncList,
		Result:   NewEmptyResult(),
	}
}

func (v *Visitor) AddDB(dbName string) {
	for _, db := range v.Result.DBNames {
		if dbName == db {
			return
		}
	}

	v.Result.DBNames = append(v.Result.DBNames, dbName)
}

func (v *Visitor) AddTable(tableName string) {
	if !common.StringInSlice(v.Result.TableNames, tableName) {
		v.Result.TableNames = append(v.Result.TableNames, tableName)
	}
}

func (v *Visitor) AddTableComment(tableName string, tableComment string) {
	v.Result.TableComments[tableName] = tableComment
}

func (v *Visitor) AddColumn(columnName string) {
	if !common.StringInSlice(v.Result.ColumnNames, columnName) {
		v.Result.ColumnNames = append(v.Result.ColumnNames, columnName)
	}
}

func (v *Visitor) AddColumnComment(columnName string, columnComment string) {
	v.Result.ColumnComments[columnName] = columnComment
}

func (v *Visitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	var (
		funcArgs      []ast.ExprNode
		dbName        string
		tableName     string
		columnName    string
		tableComment  string
		columnComment string
	)

	astType := reflect.TypeOf(in).String()

	if common.StringInSlice(v.SQLList, astType) {
		v.ToParse = true
		// get sql type
		v.Result.SQLType = strings.Split(astType, ".")[1]
	}

	if v.ToParse {
		switch in.(type) {
		case *ast.TableName:
			// get table name
			tableName = in.(*ast.TableName).Name.L
			v.AddTable(tableName)
			// get db name
			dbName = in.(*ast.TableName).Schema.L
			if dbName != "" {
				v.AddDB(dbName)
			}
		case *ast.CreateTableStmt:
			// get table comment
			tableName = in.(*ast.CreateTableStmt).Table.Name.L

			for _, tableOption := range in.(*ast.CreateTableStmt).Options {
				if tableOption.Tp == ast.TableOptionComment {
					tableComment = tableOption.StrValue
					v.AddTableComment(tableName, tableComment)
					break
				}
			}
		case *ast.AlterTableStmt:
			// get table comment
			tableName = in.(*ast.AlterTableStmt).Table.Name.L

			for _, tableSpec := range in.(*ast.AlterTableStmt).Specs {
				for _, tableOption := range tableSpec.Options {
					if tableOption.Tp == ast.TableOptionComment {
						tableComment = tableOption.StrValue
						v.AddTableComment(tableName, tableComment)
						break
					}
				}
			}
		case *ast.SelectField:
			// get column name
			expr := in.(*ast.SelectField).Expr
			if expr == nil && in.(*ast.SelectField).WildCard != nil {
				columnName = "*"
				v.AddColumn(columnName)
			} else if expr != nil {
				switch expr.(type) {
				case *ast.AggregateFuncExpr:
					funcArgs = expr.(*ast.AggregateFuncExpr).Args
				case *ast.FuncCallExpr:
					funcArgs = expr.(*ast.FuncCallExpr).Args
				case *ast.WindowFuncExpr:
					funcArgs = expr.(*ast.WindowFuncExpr).Args
				case *ast.ColumnNameExpr:
					columnName = expr.(*ast.ColumnNameExpr).Name.Name.L
					v.AddColumn(columnName)
				}

				for _, arg := range funcArgs {
					switch arg.(type) {
					case *ast.ColumnNameExpr:
						columnName = arg.(*ast.ColumnNameExpr).Name.Name.L
						v.AddColumn(columnName)
					}
				}
			}
		case *ast.ColumnDef:
			// get column comment
			columnName = in.(*ast.ColumnDef).Name.Name.L
			v.AddColumn(columnName)

			for _, columnOption := range in.(*ast.ColumnDef).Options {
				if columnOption.Tp == ast.ColumnOptionComment {
					columnComment = columnOption.Expr.(*driver.ValueExpr).Datum.GetString()
				}
			}

			v.AddColumnComment(columnName, columnComment)
		case *ast.ColumnName:
			columnName = in.(*ast.ColumnName).Name.L
			v.AddColumn(columnName)
		}
	}

	return in, false
}

func (v *Visitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

type Parser struct {
	Parser  *parser.Parser
	Visitor *Visitor
}

// NewParser returns a new *Parser
func NewParser() *Parser {
	return &Parser{
		Parser:  parser.New(),
		Visitor: NewVisitor(),
	}
}

// Parse parses sql and returns the result
func (p *Parser) Parse(sql string) (*Result, []error, error) {
	stmtNodes, warns, err := p.Parser.Parse(sql, constant.EmptyString, constant.EmptyString)
	if warns != nil || err != nil {
		return nil, warns, err
	}

	for _, stmtNode := range stmtNodes {
		stmtNode.Accept(p.Visitor)
	}

	return p.Visitor.Result, nil, nil
}

// Split splits multiple sqls into a slice
func (p *Parser) Split(sqls string) ([]string, []error, error) {
	var sqlList []string

	stmtNodes, warns, err := p.Parser.Parse(sqls, constant.EmptyString, constant.EmptyString)
	if warns != nil || err != nil {
		return nil, warns, err
	}

	for _, stmtNode := range stmtNodes {
		sqlList = append(sqlList, stmtNode.Text())
	}

	return sqlList, nil, nil
}

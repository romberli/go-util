package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/percona/go-mysql/query"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/romberli/go-util/constant"
)

const (
	alterTablePrefix = "alter table"
	addIndexKeyword  = "add index"

	backTicks            = "(`[^\\s]*`)"
	alterTableExpString  = `(?i)alter\s*table\s*(` + backTicks + `|([^\s]*))\s*`
	createIndexExpString = `(?i)create((unique)|(fulltext)|(spatial)|(primary)|(\s*)\s*)((index)|(key))\s*`
	indexNameExpString   = `(?i)(` + backTicks + `|([^\s]*))\s*`
)

type Parser struct {
	parser  *parser.Parser
	visitor *Visitor
}

// NewParser returns a new *Parser
func NewParser(visitor *Visitor) *Parser {
	return &Parser{
		parser:  parser.New(),
		visitor: visitor,
	}
}

// NewParserWithDefault returns a new *Parser with default visitor
func NewParserWithDefault() *Parser {
	return &Parser{
		parser:  parser.New(),
		visitor: NewVisitorWithDefault(),
	}
}

// GetTiDBParser returns the TiDB parser
func (p *Parser) GetTiDBParser() *parser.Parser {
	return p.parser
}

// GetVisitor returns the visitor
func (p *Parser) GetVisitor() *Visitor {
	return p.visitor
}

// GetFingerprint returns fingerprint of the given sql
func (p *Parser) GetFingerprint(sql string) string {
	return query.Fingerprint(sql)
}

// GetSQLID returns the sql id of the given sql
func (p *Parser) GetSQLID(sql string) string {
	return query.Id(p.GetFingerprint(sql))
}

// Parse parses sql and returns the result,
// not that only some kinds of statements will be parsed,
// see the constants defined at the top of visitor.go file
func (p *Parser) Parse(sql string) (*Result, error) {
	stmtNodes, err := p.GetStatementNodes(sql)
	if err != nil {
		return nil, err
	}

	for _, stmtNode := range stmtNodes {
		stmtNode.Accept(p.visitor)
	}

	return p.visitor.result, nil
}

// GetStatementNodes gets the statement nodes of the given sql
func (p *Parser) GetStatementNodes(sql string) ([]ast.StmtNode, error) {
	stmtNodes, warns, err := p.GetTiDBParser().Parse(sql, constant.EmptyString, constant.EmptyString)
	if err != nil || warns != nil {
		merr := &multierror.Error{}

		if err != nil {
			merr = multierror.Append(merr, err)
		}
		if warns != nil {
			merr = multierror.Append(merr, warns...)
		}

		return nil, merr.ErrorOrNil()
	}

	return stmtNodes, nil
}

// Split splits multiple sql statements into a slice
func (p *Parser) Split(multiSQL string) ([]string, error) {
	var sqlList []string

	stmtNodes, err := p.GetStatementNodes(multiSQL)
	if err != nil {
		return nil, err

	}

	for _, stmtNode := range stmtNodes {
		sqlList = append(sqlList, stmtNode.Text())
	}

	return sqlList, nil
}

// MergeDDLStatements merges ddl statements by table names.
// note that only alter table statement and create index statement will be merged,
// inputting other sql statements will return error,
// each argument in the input sqls could contain multiple sql statements
func (p *Parser) MergeDDLStatements(sqls ...string) ([]string, error) {
	var (
		mergedSQLList []string
		stmtNodes     []ast.StmtNode
	)

	alterTableClauses := make(map[string][]string)

	// init regular expression variables
	alterTableExp := regexp.MustCompile(alterTableExpString)
	createIndexExp := regexp.MustCompile(createIndexExpString)
	indexNameExp := regexp.MustCompile(indexNameExpString)

	for _, sqlOrig := range sqls {
		// try to split sql
		sqlList, err := p.Split(sqlOrig)
		if err != nil {
			return nil, err
		}

		for _, sql := range sqlList {
			sql = strings.Trim(sql, constant.SemicolonString)
			// get statement nodes
			stmtNodes, err = p.GetStatementNodes(sql)
			if err != nil {
				return nil, err
			}

			var (
				dbName           string
				tableName        string
				alterTableClause string
			)

			for _, stmtNode := range stmtNodes {
				switch node := stmtNode.(type) {
				case *ast.AlterTableStmt:
					dbName = node.Table.Schema.L
					tableName = node.Table.Name.L

					if alterTableExp.MatchString(sql) {
						// get alter table clause
						alterTableClause = fmt.Sprint(alterTableExp.ReplaceAllString(sql, constant.EmptyString))
					}
				case *ast.CreateIndexStmt:
					dbName = node.Table.Schema.L
					tableName = node.Table.Name.L

					sqlExp := createIndexExp.ReplaceAllString(sql, constant.EmptyString)
					indexName := strings.TrimSpace(indexNameExp.FindString(sqlExp))
					sqlExp = string([]byte(sqlExp)[strings.Index(sqlExp, constant.LeftParenthesis):])
					// get alter table clause
					alterTableClause = fmt.Sprintf("%s %s %s", addIndexKeyword, indexName, sqlExp)
				default:
					return nil, errors.New(fmt.Sprintf(
						"sql statement must be either alter table statement or create index statement, this is not valid. sql:%s\n", sql))
				}

				if alterTableClause != constant.EmptyString && tableName != constant.EmptyString {
					fullTableName := tableName
					if dbName != constant.EmptyString {
						fullTableName = fmt.Sprintf("%s.%s", dbName, tableName)
					}

					// add the alter table clause to the map
					alterTableClauses[fullTableName] = append(alterTableClauses[fullTableName], alterTableClause)
				}
			}
		}
	}

	for fullTableName, alterClauseList := range alterTableClauses {
		// get merged alter table statement
		alterTableStatement := fmt.Sprintf("%s %s %s%s",
			alterTablePrefix, fullTableName, strings.Join(alterClauseList, fmt.Sprintf("%s ", constant.CommaString)), constant.SemicolonString)
		// add to the result
		mergedSQLList = append(mergedSQLList, alterTableStatement)
	}

	return mergedSQLList, nil
}

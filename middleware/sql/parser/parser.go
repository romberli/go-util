package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

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

// Parse parses sql and returns the result,
// not that only some kinds of statements will be parsed,
// see the constants defined at the top of visitor.go file
func (p *Parser) Parse(sql string) (*Result, []error, error) {
	stmtNodes, warns, err := p.parser.Parse(sql, constant.EmptyString, constant.EmptyString)
	if warns != nil || err != nil {
		return nil, warns, err
	}

	for _, stmtNode := range stmtNodes {
		stmtNode.Accept(p.visitor)
	}

	return p.visitor.result, nil, nil
}

// Split splits multiple sqls into a slice
func (p *Parser) Split(sqls string) ([]string, []error, error) {
	var sqlList []string

	stmtNodes, warns, err := p.parser.Parse(sqls, constant.EmptyString, constant.EmptyString)
	if warns != nil || err != nil {
		return nil, warns, err
	}

	for _, stmtNode := range stmtNodes {
		sqlList = append(sqlList, stmtNode.Text())
	}

	return sqlList, nil, nil
}

// GetFingerprint returns fingerprint of the given sql
func (p *Parser) GetFingerprint(sql string) string {
	return query.Fingerprint(sql)
}

// GetSQLID returns the sql id of the given sql
func (p *Parser) GetSQLID(sql string) string {
	return query.Id(p.GetFingerprint(sql))
}

// MergeDDLStatements merges ddl statements by table names,
// note that only alter table statement and create index statement will be merged,
// inputting other sql statements will return error
func (p *Parser) MergeDDLStatements(sqls ...string) ([]string, []error, error) {
	var result []string

	alterTableClauseMap := make(map[string][]string)

	// init regular expression
	alterTableExp := regexp.MustCompile(alterTableExpString)
	createIndexExp := regexp.MustCompile(createIndexExpString)
	indexNameExp := regexp.MustCompile(indexNameExpString)

	for _, sql := range sqls {
		sql = strings.Trim(sql, constant.SemicolonString)
		// parse sql
		stmtNodes, warns, err := p.GetTiDBParser().Parse(sql, constant.EmptyString, constant.EmptyString)
		if err != nil || warns != nil {
			return nil, warns, err
		}

		var (
			alterTableClause string
			dbName           string
			tableName        string
		)

		for _, stmtNode := range stmtNodes {
			switch node := stmtNode.(type) {
			case *ast.AlterTableStmt:
				tableName = node.Table.Name.L
				dbName = node.Table.Schema.L

				if alterTableExp.MatchString(sql) {
					// get alter table clause
					alterTableClause = fmt.Sprint(alterTableExp.ReplaceAllString(sql, constant.EmptyString))
				}
			case *ast.CreateIndexStmt:
				tableName = node.Table.Name.L
				dbName = node.Table.Schema.L

				sqlExp := createIndexExp.ReplaceAllString(sql, constant.EmptyString)
				indexName := strings.TrimSpace(indexNameExp.FindString(sqlExp))
				sqlExp = string([]byte(sqlExp)[strings.Index(sqlExp, constant.LeftParenthesis):])
				// get alter table clause
				alterTableClause = fmt.Sprintf("%s %s %s", addIndexKeyword, indexName, sqlExp)
			default:
				return nil, nil, errors.New(fmt.Sprintf(
					"sql statement must be either alter table statement or create index statement, this is not valid. sql:%s\n", sql))
			}

			if alterTableClause != constant.EmptyString && tableName != constant.EmptyString {
				fullTableName := tableName
				if dbName != constant.EmptyString {
					fullTableName = fmt.Sprintf("%s.%s", dbName, tableName)
				}

				// add the alter table clause to the map
				alterTableClauseMap[fullTableName] = append(alterTableClauseMap[fullTableName], alterTableClause)
			}
		}
	}

	for fullTableName, alterClauseList := range alterTableClauseMap {
		// get merged alter table statement
		alterTableStatement := fmt.Sprintf("%s %s %s%s",
			alterTablePrefix, fullTableName, strings.Join(alterClauseList, fmt.Sprintf("%s ", constant.CommaString)), constant.SemicolonString)

		result = append(result, alterTableStatement)
	}

	return result, nil, nil
}

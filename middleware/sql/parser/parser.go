package parser

import (
	"github.com/percona/go-mysql/query"
	"github.com/pingcap/parser"
	"github.com/romberli/go-util/constant"
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

// GetFingerprint returns fingerprint fo given sql
func (p *Parser) GetFingerprint(sql string) string {
	return query.Fingerprint(sql)
}

func (p *Parser) GetSQLID(sql string) string {
	return query.Id(p.GetFingerprint(sql))
}

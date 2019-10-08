package parser

type StmtType uint8

const (
	StmtTypeFunctionCall = iota
	StmtTypeImport
	StmtTypePrint
	StmtTypeComment
	StmtTypeExpr
	StmtTypeTypeDecl
	StmtTypeVarDecl
	StmtTypeFuncDecl
	StmtUnknown
)

type ParseResult struct {
	StmtType StmtType
}

func Parse(code string) (StmtType, error) {
	if isComment(code) {
		return StmtTypeComment, nil
	} else if isImport(code) {
		return StmtTypeImport, nil
	} else if isFunc(code) {
		return StmtTypeFuncDecl, nil
	} else if isTypeDecl(code) {
		return StmtTypeTypeDecl, nil
	} else if isPrint(code) {
		return StmtTypePrint, nil
	} else if isComment(code) {
		return StmtTypeComment, nil
	} else if isExpr(code) {
		return StmtTypeExpr, nil
	} else {
		return StmtUnknown, nil
	}
}

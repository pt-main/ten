package ten

const (
	_CONFIG_BLOCK_SYNTAX = `\(\(\?CONFIG\s([\s\S]*?)\s*\)\)`
	_CODE_BLOCK_SYNTAX   = `\{\{\?CODE\s([\s\S]*?)\s*\}\}`
	_PLACEHOLDER_SYNTAX  = `\[\[\?\s\n*(.*?)\n*\s\]\]`
	_RAW_SYNTAX          = `(?s)(.+?)(?=\[\[\?|\{\{\?|\(\(\?|$)`
)

var (
	_CONFIG_START_END      = []string{"((?CONFIG", "))"}
	_CODE_START_END        = []string{"{{?CODE", "}}"}
	_PLACEHOLDER_START_END = []string{"[[?", "]]"}
)

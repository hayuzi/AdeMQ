package handler

import "strings"

// ParsedCmd 解析之后的命令结构
type ParsedCmd struct {
	Cmd    string
	Params []string
	Origin string
}

// Parser 命令解析器
type Parser struct{}

// NewParser 创建命令解析器
func NewParser() *Parser {
	return &Parser{}
}

// Parse 解析命令并返回结果
func (p *Parser) Parse(cmdStr string) *ParsedCmd {
	parsedCmd := &ParsedCmd{
		Cmd:    "",
		Params: make([]string, 0),
		Origin: cmdStr,
	}
	tmpSlice := strings.Split(cmdStr, " ")
	cmdStored := false
	for _, str := range tmpSlice {
		if str == "" || str == " " {
			continue
		}
		// 如果命令字段未解析，先解析命令
		if !cmdStored {
			parsedCmd.Cmd = str
			cmdStored = true
			continue
		}
		// 如果是参数, 注入到参数中去
		parsedCmd.Params = append(parsedCmd.Params, str)
	}
	return parsedCmd
}

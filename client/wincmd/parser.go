package wincmd

import "strings"

// Parser 命令解析器
type Parser struct {
	Cmd    string
	Params []string
	Origin string
}

// New 创建命令解析器
func New(cmdStr string) *Parser {
	p := &Parser{
		Cmd:    "",
		Params: make([]string, 4),
		Origin: cmdStr,
	}
	tmpSlice := strings.Split(cmdStr, " ")
	cmdStored := false
	for _, str := range tmpSlice {
		if str == "" {
			continue
		}
		// 如果命令字段未解析，先解析命令
		if !cmdStored {
			p.Cmd = str
			cmdStored = false
		}
		// 如果是参数, 注入到参数中去
		p.Params = append(p.Params, str)
	}
	return p
}

func (p *Parser) Exec() {

}

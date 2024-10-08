package raw

import (
	"fmt"
	"strings"
)

type Raw struct {
}

func New() Raw {
	return Raw{}
}

func (r Raw) AddLine(line ...string) {
	fmt.Println(strings.Join(line, " "))
}

func (r Raw) Render() {

}

package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/talitore/smartass/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}

package generator

import . "gopkg.in/check.v1"

type TplSuite struct{}

var _ = Suite(&TplSuite{})

var typeText = `type {{.Name}}Iter []{{.Type}}

func New{{.Name}}Iter(items ...{{.Type}}) {{.Name}}Iter {
  return {{.Name}}Iter(items)
}
`

func (s *TplSuite) TestLoadTemplateText(c *C) {
	text := loadTemplateText("type")
	c.Assert(text, Equals, typeText)
}

func (s *TplSuite) TestLoadTemplate(c *C) {
	// No panic means it worked
	t := loadTemplate("type")
	c.Assert(t, Not(IsNil))
}

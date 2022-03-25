package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"text/template"
)

type Conditioner interface {
	Argument() string
	ArgumentName() string
	Condition() string
	SQL() string
}

type FieldEqual struct {
	Field *ModelField
	Val   string
}

func (fe *FieldEqual) SQL() string {
	return fmt.Sprintf("%s = ?", fe.Field.SqlName())
}

func (fe *FieldEqual) Condition() string {
	var sb = &strings.Builder{}

	text := `
	SQL.WriteString(" {{.COLUMN}}=?")	
	ARGS = append(ARGS, {{.NAME}})
`
	if fe.Val == "" {
		text = `
	SQL.WriteString(" {{.COLUMN}}={{.VAL}}")
`
	}
	tpl := template.New("FieldEqual")
	tpl.Parse(text)

	tpl.Execute(sb, struct {
		COLUMN string
		NAME   string
		VAL    string
	}{
		COLUMN: fe.Field.SqlName(),
		NAME:   fe.ArgumentName(),
		VAL:    fe.Val,
	})

	return sb.String()
}

func (fe *FieldEqual) ArgumentName() string {
	return fe.Field.Name
}

func (fe *FieldEqual) Argument() string {
	if fe.Val == "" {
		return ""
	}
	return fe.ArgumentName() + " " + fe.Field.Type
}

type FieldIn struct {
	Field *ModelField
}

func (in *FieldIn) ArgumentName() string {
	return in.Field.Name + "List"
}

func (in *FieldIn) ArgumentType() string {
	return "[]" + in.Field.Type
}

func (in *FieldIn) Argument() string {
	return in.ArgumentName() + " " + in.ArgumentType()
}

func (in *FieldIn) SQL() string {
	return fmt.Sprintf("%s IN (%%s)", in.Field.SqlName())
}

func (in *FieldIn) Condition() string {

	var sb = &strings.Builder{}

	text := `
	if len({{.LIST}}) == 0 {
		return false
	}
	SQL.WriteString(" {{.COLUMN}} IN (")
	SQL.WriteString(strings.Repeat("?, ", len({{.LIST}})-1))
	SQL.WriteString("?)")
	for _, v := range {{.LIST}} {
		ARGS = append(ARGS, v)
	}
`
	tpl := template.New("sqlIn")
	tpl.Parse(text)

	tpl.Execute(sb, struct {
		COLUMN string
		LIST   string
	}{
		COLUMN: in.Field.SqlName(),
		LIST:   in.ArgumentName(),
	})

	return sb.String()
}

func PrivateFieldCase(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func (m *Model) ParseCondition(s string) Conditioner {

	if X, Y, ok := MatchEqualCond(s); ok {
		field := m.GetField(X)
		if field == nil {
			log.Fatalf("invalid cond '%s' get field '%s' nil", s, X)
		}
		return &FieldEqual{
			Field: field,
			Val:   Y,
		}
	}
	if X, ok := MatchFieldInCond(s); ok {
		field := m.GetField(X)
		if field == nil {
			log.Fatalf("invalid conditon '%s' get field '%s' nil", s, X)
		}
		return &FieldIn{
			Field: field,
		}
	}
	log.Fatalf("ParseCondition %v", s)
	return nil
}

// MatchEqualCond 'key=?'
func MatchEqualCond(s string) (X string, Y string, Ok bool) {
	re := regexp.MustCompile(`^(\w+)(?:\s*=\s*(\w+|\?))?$`)
	matches := re.FindAllStringSubmatch(s, -1)
	if len(matches) == 1 {
		X, Y = matches[0][1], matches[0][2]
		if Y == "" {
			Y = "?"
		}
		return X, Y, true
	}

	return "", "", false
}

//MatchFieldInCond 'key IN (?)'
func MatchFieldInCond(s string) (X string, Ok bool) {
	re := regexp.MustCompile(`^(\w+)\s* IN \s*\(\?\)$`)
	matches := re.FindAllStringSubmatch(s, -1)
	if len(matches) == 1 {
		return matches[0][1], true
	}

	return "", false
}

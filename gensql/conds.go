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
}

type FieldEqual struct {
	Field *ModelField
	Val   string
}

func (fe *FieldEqual) Condition() string {
	return fmt.Sprintf("%s = ?", snakeCase(fe.Field.Name))
}

func (fe *FieldEqual) ArgumentName() string {
	return PrivateFieldCase(fe.Field.Name)
}

func (fe *FieldEqual) Argument() string {
	return fe.ArgumentName() + " " + fe.Field.Type
}

type FieldIn struct {
	Field *ModelField
}

func (in *FieldIn) ArgumentName() string {
	return strings.ToLower(in.Field.Name[:1]) + in.Field.Name[1:] + "List"
}

func (in *FieldIn) ArgumentType() string {
	return "[]" + in.Field.Type
}

func (in *FieldIn) Argument() string {
	return in.ArgumentName() + " " + in.ArgumentType()
}

func (in *FieldIn) Condition() string {

	var sb = &strings.Builder{}

	text := `
	sqlSb.WriteString("{{.COLUMN}} IN(")
	for i, v := range {{.LIST}} {
		if i != 0 {
			sqlSb.WriteString(", ")
		}
		sqlSb.WriteString("?")
		args = append(args, v)
	}
	sqlSb.WriteString(")")
`
	tpl := template.New("sqlIn")
	tpl.Parse(text)

	tpl.Execute(sb, struct {
		COLUMN string
		LIST   string
	}{
		COLUMN: snakeCase(in.Field.Name),
		LIST:   in.ArgumentName(),
	})

	return sb.String()
}

func PrivateFieldCase(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func (m *Model) ParseCondition(s string) Conditioner {

	if X, Y, ok := MatchEqualCond(s); ok {
		return &FieldEqual{
			Field: m.GetField(X),
			Val:   Y,
		}
	}
	if X, ok := MatchFieldInCond(s); ok {
		return &FieldIn{
			Field: m.GetField(X),
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

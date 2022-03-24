package main

import (
	"testing"
)

func TestParseCondition(t *testing.T) {
	t.Log(MatchEqualCond("player_sn=2"))
}

func TestMatchFiledInCond(t *testing.T) {
	t.Log(MatchFieldInCond("player_sn IN (?)"))
}

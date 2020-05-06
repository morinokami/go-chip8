#!/usr/bin/env python

import os

template = '''
package games

import (
	"fmt"
	"strings"
)

type Game struct {{
	Name   string
	Binary []byte
}}

var Games = []Game{{
	{}
}}

func AvailableGames() string {{
	var desc []string
	for i, g := range Games {{
		desc = append(desc, fmt.Sprintf("%d. %s", i, g.Name))
	}}
	return strings.Join(desc, ", ")
}}
'''

game_template = '''
	{{
                Name: "{}",
                Binary: []byte{{{}}},
	}},
'''

games = ''
filter_func = lambda x: not x.endswith('.py') and not x.endswith('.go')
for game in filter(filter_func, sorted(os.listdir('.'))):
    with open(game, 'rb') as f:
        binary = ', '.join('0x{:02X}'.format(b) for b in f.read())
        games += game_template.format(game, binary)

with open('games.go', 'w') as f:
    f.write(template.format(games))

%{

package sgf;

import (
        __yyfmt__ "fmt"
       )

%}

%union {
    g *GameTree
    gs []*GameTree
    s *Sequence
    n *Node
    ps []Property
    vs []PropValue
    v PropValue
    name string
}

%type   <gs>            game_trees
%type   <g>             game_tree
%type   <s>             sequence
%type   <n>             node
%type   <ps>            properties
%type   <vs>            prop_values

%token '(' ')' ';'

%token  <v>             TokPropValue
%token  <name>          TokPropName

%%

collection:
                game_trees
                {
                    yylex.(*lexer).c = &Collection{
                      Trees: $1,
                    }
                }

game_trees:
                game_tree
                {
                    $$ = []*GameTree{$1}
                }
        |       game_trees game_tree {
                    $$ = append($1, $2)
                }

game_tree:
                '(' sequence ')'
                {
                    $$ = &GameTree{
                      Principal: *$2,
                    }
                }
        |       '(' sequence game_trees ')'
                {
                    $$ = &GameTree{
                      Principal: *$2,
                      Children: $3,
                    }
                }

sequence:
                node
                {
                    $$ = &Sequence{[]Node{*$1}}
                }
        |       sequence node
                {
                    $1.Nodes = append($1.Nodes, *$2)
                    $$ = $1
                }

node:
                ';' properties
                {
                    $$ = &Node{$2}
                }

properties:
                properties TokPropName prop_values
                {
                    $$ = append($1, Property{$2, $3})
                }
        |
                {
                    $$ = nil
                }

prop_values:
                TokPropValue
                {
                    $$ = []PropValue{$1}
                }
        |       prop_values TokPropValue
                {
                    $$ = append($1, $2)
                }
%%

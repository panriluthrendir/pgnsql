package main

import (
    "database/sql/driver"
    "errors"
    "strings"
)

func applicable(game *Game, cond *condition, value driver.Value) bool {
    if cond == nil {
        return true
    }
    switch cond.ctype {
        case EQ:
            return game.headers[cond.col] == value 
    }
    return false
}

func parseCondition(str string) (*condition, error) {
    for s, ctype := range CTYPES {
        if strings.Index(str, s) != -1 {
            split := strings.Split(str, s)
            return &condition{
                col: split[0], 
                ctype: ctype}, nil
        }
    }
    return nil, errors.New("Cannot parse " + str)
}

package main

import (
    "database/sql/driver"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "strings"
)

type PgnDriver struct{}

type PgnConn struct {
    content []byte
}

type PgnRows struct {
    content []byte
    cols    []string
    reader  *PgnReader
    cond    *condition
    value   driver.Value
}

type condition struct {
    col   string
    ctype ctype
}

type ctype int

const (
    EQ ctype = iota
)

var CTYPES = map[string]ctype{
    "=": EQ,
}

func (driver *PgnDriver) Open(name string) (driver.Conn, error) {
    content, err := ioutil.ReadFile(name)

    if err != nil {
        return nil, err
    }

    return &PgnConn{content: content}, nil
}

func (c *PgnConn) Query(query string, args []driver.Value) (driver.Rows, error) {
    parts := strings.Split(query, " FROM")
    colPart := strings.Split(parts[0], "SELECT ")[1]
    cols := strings.Split(colPart, ",")

    if len(args) > 1 {
        return nil, errors.New("Multiple conditions are not supported")
    }

    if len(args) == 1 {
        condPart := strings.Split(parts[1], "WHERE ")[1]
        cond, err := parseCondition(condPart)

        if err != nil {
            return nil, err
        }

        return &PgnRows{
            cols:    cols,
            reader:  NewPgnReader(),
            content: c.content,
            cond:    cond,
            value:   args[0]}, nil
    }

    return &PgnRows{
        cols:    cols,
        reader:  NewPgnReader(),
        content: c.content,
        cond:    nil,
        value:   nil}, nil

}

func (c *PgnConn) Close() error {
    return nil
}

func (r *PgnRows) Next(dest []driver.Value) error {
    for {
        if len(r.content) == 0 {
            return io.EOF
        }
        read, err := r.reader.Read(r.content)

        if err != nil {
            return err
        }

        game := r.reader.games[len(r.reader.games)-1]
        r.content = r.content[read:]

        if applicable(&game, r.cond, r.value) {
            for i := 0; i < len(r.cols); i++ {
                for k, v := range game.headers {
                    if r.cols[i] == k {
                        dest[i] = driver.Value(v)
                    }
                }
            }
            return nil
        }
    }
}

func (r *PgnRows) Columns() []string {
    return r.cols
}

func (r *PgnRows) Close() error {
    return nil
}

// Not implemented
func (c *PgnConn) Prepare(query string) (driver.Stmt, error) {
    return nil, fmt.Errorf("Prepare method not implemented")
}

func (c *PgnConn) Begin() (driver.Tx, error) {
    return c, fmt.Errorf("Begin method not implemented")
}

func (c *PgnConn) Commit() error {
    return fmt.Errorf("Commit method not implemented")
}

func (c *PgnConn) Rollback() error {
    return fmt.Errorf("Rollback method not implemented")
}

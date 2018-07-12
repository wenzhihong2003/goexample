package main

import (
	"database/sql"
	"fmt"
	"io"
)

type Command struct {
	ID     int
	Result string
}

func (c Command) Error() string {
	return fmt.Sprintf("%s %d", c.Result, c.ID)
}

func main() {
	for i := 0; i < 5; i++ {
		err := process(i)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
				continue
			}
			switch e := err.(type) {
			case Command:
				fmt.Println("Command Error")
			default:
				fmt.Printf("### e -> %T\n", e)
				fmt.Printf("### e -> %+v\n", e)
			}
		}
	}
}

func process(i int) error {
	switch i {
	case 1:
		return fmt.Errorf("error %d", i)
	case 2:
		return io.EOF
	case 3:
		return sql.ErrNoRows
	case 4:
		return Command{}
	}
	return nil
}

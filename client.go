package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

func (cl *client) readInput() {
	for {
		msg, err := bufio.NewReader(cl.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			cl.commands <- command{
				id:     CMD_NICK,
				client: cl,
				args:   args,
			}
		case "/join":
			cl.commands <- command{
				id:     CMD_JOIN,
				client: cl,
				args:   args,
			}
		case "/rooms":
			cl.commands <- command{
				id:     CMD_ROOMS,
				client: cl,
				args:   args,
			}
		case "/msg":
			cl.commands <- command{
				id:     CMD_MSG,
				client: cl,
				args:   args,
			}
		case "/quit":
			cl.commands <- command{
				id:     CMD_QUIT,
				client: cl,
				args:   args,
			}
		default:
			cl.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}

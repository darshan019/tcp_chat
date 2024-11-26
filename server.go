package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listOfRooms(cmd.client)
		case CMD_QUIT:
			s.quit(cmd.client)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		}
	}
}

func (server *server) newClient(conn net.Conn) {
	log.Printf("new client: %s", conn.RemoteAddr().String())

	client := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: server.commands,
	}

	client.readInput()

}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) nick(c *client, args []string) {
	c.nick = args[1]
	c.msg(fmt.Sprintf("You are now: %s", c.nick))
}

func (s *server) join(c *client, args []string) {
	roomName := args[1]
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	r.members[c.conn.RemoteAddr()] = c
	s.quitCurrentRoom(c)

	c.room = r
	r.broadCast(c, fmt.Sprintf("%s has joined the room", c.nick))
	c.msg(fmt.Sprintf("Welcome to %s", r.name))
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadCast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}

func (s *server) listOfRooms(c *client) {
	var rooms []string
	for roomName := range s.rooms {
		rooms = append(rooms, roomName)
	}

	c.msg(fmt.Sprintf("Rooms available are: %s", strings.Join(rooms, ", ")))
}

func (s *server) quit(c *client) {
	log.Printf("%s has left the room: %s", c.nick, c.conn.RemoteAddr())
	s.quitCurrentRoom(c)
	c.msg(fmt.Sprintf("You have left room %s", c.room.name))
	c.conn.Close()
}

func (s *server) msg(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("join a room first"))
		return
	}

	c.room.broadCast(c, c.nick+": "+strings.Join(args[1:], " "))
}

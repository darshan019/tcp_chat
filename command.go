package main

type commandID int

const (
	CMD_NICK commandID = iota // 0
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
	CMD_MSGP
	CMD_MEMS
)

type command struct {
	id     commandID
	client *client
	args   []string
}

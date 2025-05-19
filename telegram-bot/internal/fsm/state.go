package fsm

import "github.com/vitaliy-ukiru/fsm-telebot/v2"

var TicketGroup = fsm.NewStateGroup("ticket")

var (
	Type    = TicketGroup.New("type")
	Content = TicketGroup.New("content")
	Confirm = TicketGroup.New("confirm")
)

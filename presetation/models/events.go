package models

type AIResponseMsg struct {
	Content string
}

type WaitingMsg struct{}

type ErrorMsg error

package agent

import (
	"errors"
	"io"

	"github.com/cloudwego/eino/schema"
)

type AgentStreamResponse struct {
	msgReader *schema.StreamReader[*schema.Message]
	content   string
	finished  bool
}

func (r *AgentStreamResponse) Content() (string, bool, error) {
	msg, err := r.msgReader.Recv()
	if err != nil && errors.Is(err, io.EOF) {
		r.finished = true
		return r.content, r.finished, nil
	}

	if err != nil {
		r.finished = true
		return "Error on generating response", r.finished, err
	}

	r.content = msg.Content

	return r.content, r.finished, nil
}

func NewAgentStreamResponse(msgReader *schema.StreamReader[*schema.Message]) *AgentStreamResponse {
	return &AgentStreamResponse{msgReader: msgReader}
}

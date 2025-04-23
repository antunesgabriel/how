package agent

import (
	"errors"
	"io"

	"github.com/cloudwego/eino/schema"
)

type AgentStreamResponse struct {
	msgReader *schema.StreamReader[*schema.Message]
	content   string
}

func (r *AgentStreamResponse) Content() (string, bool, error) {
	msg, err := r.msgReader.Recv()
	if err != nil && errors.Is(err, io.EOF) {
		return r.content, false, nil
	}

	if err != nil {
		return "Error on generating response", false, err
	}

	r.content = msg.Content

	return r.content, true, nil
}

func NewAgentStreamResponse(msgReader *schema.StreamReader[*schema.Message]) *AgentStreamResponse {
	return &AgentStreamResponse{msgReader: msgReader}
}

package agent

import (
	"github.com/cloudwego/eino/schema"
)

type AgentStreamResponse struct {
	msgReader *schema.StreamReader[*schema.Message]
	content   string
}

func (r *AgentStreamResponse) Content() (string, bool, error) {
	msg, err := r.msgReader.Recv()

	if err != nil {
		// caller should verify if error is io.EOF
		return r.content, false, err
	}

	r.content = msg.Content

	return r.content, true, nil
}

func NewAgentStreamResponse(msgReader *schema.StreamReader[*schema.Message]) *AgentStreamResponse {
	return &AgentStreamResponse{msgReader: msgReader}
}

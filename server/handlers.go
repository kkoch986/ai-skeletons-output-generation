package server

import (
	"context"

	"github.com/juju/zaputil/zapctx"
	"github.com/kkoch986/ai-skeletons-output-generation/viseme"
	"github.com/kkoch986/ai-skeletons-playback/output"
	"go.uber.org/zap"
)

// GenerateRequest is the struct that can parse the POST request sent to /generate
type GenerateRequest struct {
	Prompt         string          `json:"prompt"`
	VisemeSequence viseme.Sequence `json:"visemes"`
}

// GenerateResponse contains the respost to POST /generate calls
type GenerateResponse struct {
	OutputSequence output.Sequence `json:"sequence"`
}

// HandleGenerate handles the POST /generate calls
func HandleGenerate(ctx context.Context, r *GenerateRequest) (*GenerateResponse, error) {
	logger := zapctx.Logger(ctx)
	logger.Info("generate request", zap.Any("request", r))
	return &GenerateResponse{
		OutputSequence: r.VisemeSequence.ToJawOutputSequence("/skeleton/jaw"),
	}, nil
}

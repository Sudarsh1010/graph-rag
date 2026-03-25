package di

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sudarsh1010/graph-rag/internal/api"
	"github.com/sudarsh1010/graph-rag/internal/graph"
	"github.com/sudarsh1010/graph-rag/internal/nlquery"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewGraphService(db *Database, logger *zap.Logger) *graph.Service {
	return graph.NewService(db.DB, logger)
}

func NewNLQueryService(db *Database, gs *graph.Service, logger *zap.Logger) *nlquery.NLQueryService {
	llm := nlquery.DefaultMockProvider()
	return nlquery.NewService(llm, gs, db.DB, logger)
}

func NewRouter(db *Database, gs *graph.Service, nls *nlquery.NLQueryService) http.Handler {
	nlHandler := api.NLQueryHandlerFunc(func(ctx context.Context, req api.NLQueryRequest) (*api.NLQueryResponse, error) {
		nlReq := nlquery.NLQueryRequest{
			Question:       req.Question,
			ConversationID: req.ConversationID,
		}
		nlResp, err := nls.Query(ctx, nlReq)
		if err != nil {
			return nil, err
		}

		var queryPlan interface{}
		if nlResp.QueryPlan != nil {
			raw, _ := json.Marshal(nlResp.QueryPlan)
			json.Unmarshal(raw, &queryPlan)
		}

		return &api.NLQueryResponse{
			Answer:        nlResp.Answer,
			Data:          nlResp.Data,
			Visualization: (*api.VisualizationHint)(nlResp.Visualization),
			Sources:       nlResp.Sources,
			QueryPlan:     queryPlan,
		}, nil
	})
	return api.NewRouter(db.DB, gs, api.NewNLQueryAdapter(nlHandler))
}

var Module = fx.Options(
	fx.Provide(NewConfig),
	fx.Provide(NewLogger),
	fx.Provide(NewDatabase),
	fx.Provide(NewPipelineService),
	fx.Provide(NewGraphService),
	fx.Provide(NewNLQueryService),
	fx.Provide(NewRouter),
	fx.Provide(NewServer),
	fx.Invoke(RunPipeline),
)

package nlquery

import (
	"context"
	"fmt"

	"github.com/sudarsh1010/graph-rag/internal/api"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

// LLMProvider is the interface for generating LLM completions.
// Implementations wrap specific LLM APIs (OpenAI, Anthropic, etc.).
type LLMProvider interface {
	GenerateCompletion(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}

// NLQueryService is the main orchestrator for natural language queries.
// It translates questions into query plans, executes them, and formats results.
type NLQueryService struct {
	llm    LLMProvider
	graph  api.GraphQuerier
	db     *bun.DB
	logger *zap.Logger
}

// NewService creates a new NLQueryService.
func NewService(llm LLMProvider, graph api.GraphQuerier, db *bun.DB, logger *zap.Logger) *NLQueryService {
	return &NLQueryService{
		llm:    llm,
		graph:  graph,
		db:     db,
		logger: logger,
	}
}

// Query processes a natural language question end-to-end:
//  1. Translates the question into a structured query plan
//  2. Executes each step in the plan
//  3. Formats the results into a natural language answer
func (s *NLQueryService) Query(ctx context.Context, req NLQueryRequest) (*NLQueryResponse, error) {
	if req.Question == "" {
		return nil, fmt.Errorf("question is required")
	}

	s.logger.Info("processing NL query",
		zap.String("question", req.Question),
		zap.String("conversation_id", req.ConversationID),
	)

	// Step 1: Translate the natural language question into a query plan.
	plan, err := translate(ctx, s.llm, req.Question)
	if err != nil {
		return nil, fmt.Errorf("translate question: %w", err)
	}

	s.logger.Info("generated query plan",
		zap.Int("steps", len(plan.Steps)),
		zap.String("answer_format", plan.AnswerFormat),
	)

	// Step 2: Execute the query plan steps.
	results, sources := executePlan(ctx, s.graph, s.db, s.logger, plan)

	s.logger.Info("executed query plan",
		zap.Int("results", len(results)),
		zap.Strings("sources", sources),
	)

	// Step 3: Format the execution results into a natural language answer.
	answer, err := formatResults(ctx, s.llm, req.Question, results, plan.AnswerFormat)
	if err != nil {
		s.logger.Warn("failed to format results with LLM, returning raw data",
			zap.Error(err),
		)
		// Graceful fallback: return the data without a formatted answer.
		answer = fmt.Sprintf("Query executed successfully. %d results found.", len(results))
	}

	// Step 4: Determine a visualization hint from the plan.
	viz := planToVisualization(plan)

	return &NLQueryResponse{
		Answer:        answer,
		Data:          results,
		Visualization: viz,
		Sources:       sources,
		QueryPlan:     plan,
	}, nil
}

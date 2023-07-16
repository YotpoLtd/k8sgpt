package server

import (
	"context"
	json "encoding/json"
	"github.com/k8sgpt-ai/k8sgpt/pkg/analysis"
	schemav1 "github.com/k8sgpt-ai/k8sgpt/schema/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (h *handler) Analyze(ctx context.Context, i *schemav1.AnalyzeRequest) (
	*schemav1.AnalyzeResponse,
	error,
) {
	if i.Output == "" {
		i.Output = "json"
	}

	if i.Backend == "" {
		i.Backend = "openai"
	}

	if int(i.MaxConcurrency) == 0 {
		i.MaxConcurrency = 10
	}

	var listopt v1.ListOptions
	err := json.Unmarshal([]byte(i.Listoptions), &listopt)
	if err != nil {
		return &schemav1.AnalyzeResponse{}, err
	}

	config, err := analysis.NewAnalysis(
		i.Backend,
		i.Language,
		i.Filters,
		i.Namespace,
		listopt,
		i.Nocache,
		i.Explain,
		int(i.MaxConcurrency),
		false, // Kubernetes Doc disabled in server mode
	)
	if err != nil {
		return &schemav1.AnalyzeResponse{}, err
	}
	config.RunAnalysis()

	if i.Explain {
		err := config.GetAIResults(i.Output, i.Anonymize)
		if err != nil {
			return &schemav1.AnalyzeResponse{}, err
		}
	}

	out, err := config.PrintOutput(i.Output)
	if err != nil {
		return &schemav1.AnalyzeResponse{}, err
	}
	var obj schemav1.AnalyzeResponse

	err = json.Unmarshal(out, &obj)
	if err != nil {
		return &schemav1.AnalyzeResponse{}, err
	}

	return &obj, nil
}

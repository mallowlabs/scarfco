package output

import "encoding/json"

type sarifOutput struct {
	Schema  string    `json:"$schema"`
	Version string    `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string      `json:"name"`
	InformationURI string      `json:"informationUri"`
	Rules          []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	ShortDescription sarifMessage `json:"shortDescription"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   sarifMessage    `json:"message"`
	Locations []sarifLocation `json:"locations"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           sarifRegion           `json:"region"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine int `json:"startLine"`
}

func toSARIF(r *Result) string {
	seen := map[string]bool{}
	rules := []sarifRule{}
	for _, f := range r.Files {
		for _, e := range f.Errors {
			if !seen[e.Source] {
				seen[e.Source] = true
				rules = append(rules, sarifRule{
					ID:               e.Source,
					Name:             e.Source,
					ShortDescription: sarifMessage{Text: e.Source},
				})
			}
		}
	}

	results := []sarifResult{}
	for _, f := range r.Files {
		for _, e := range f.Errors {
			results = append(results, sarifResult{
				RuleID:  e.Source,
				Level:   sarifLevel(e.Severity),
				Message: sarifMessage{Text: e.Message},
				Locations: []sarifLocation{
					{
						PhysicalLocation: sarifPhysicalLocation{
							ArtifactLocation: sarifArtifactLocation{URI: f.Name},
							Region:           sarifRegion{StartLine: e.Line},
						},
					},
				},
			})
		}
	}

	out := sarifOutput{
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Version: "2.1.0",
		Runs: []sarifRun{
			{
				Tool: sarifTool{
					Driver: sarifDriver{
						Name:           "scarfco",
						InformationURI: "https://github.com/mallowlabs/scarfco",
						Rules:          rules,
					},
				},
				Results: results,
			},
		},
	}

	b, _ := json.MarshalIndent(out, "", "  ")
	return string(b)
}

func sarifLevel(severity string) string {
	switch severity {
	case "error":
		return "error"
	case "warning":
		return "warning"
	default:
		return "note"
	}
}

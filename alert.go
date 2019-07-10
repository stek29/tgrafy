package main

import "fmt"

type nullFloat float32

type EvalMatch struct {
	Value  nullFloat         `json:"value"`
	Metric string            `json:"metric"`
	Tags   map[string]string `json:"tags"`
}

type AlertBody struct {
	Title       string      `json:"title"`
	RuleID      int64       `json:"ruleId"`
	RuleName    string      `json:"ruleName"`
	State       string      `json:"state"`
	EvalMatches []EvalMatch `json:"evalMatches"`
	RuleURL     string      `json:"ruleUrl"`
	ImageURL    string      `json:"imageUrl"`
	Message     string      `json:"message"`
}

// Based on https://github.com/grafana/grafana/blob/master/pkg/services/alerting/notifiers/telegram.go
func generateMetricsMessage(alert *AlertBody) string {
	metrics := ""
	fieldLimitCount := 5
	for index, evt := range alert.EvalMatches {
		metrics += fmt.Sprintf("\n%s: %f", evt.Metric, evt.Value)
		if index == fieldLimitCount {
			break
		}
	}
	if len(alert.EvalMatches) > fieldLimitCount {
		metrics += "\n<i>some metrics were hidden</i>"
	}
	return metrics
}

// Based on https://github.com/grafana/grafana/blob/master/pkg/services/alerting/notifiers/telegram.go
func buildMessage(alert *AlertBody) string {
	message := fmt.Sprintf("<b>%s</b>\nState: %s\nMessage: %s\n", alert.Title, alert.RuleName, alert.Message)

	if ruleURL := alert.RuleURL; ruleURL != "" {
		message = message + fmt.Sprintf("URL: %s\n", ruleURL)
	}

	if imageURL := alert.ImageURL; imageURL != "" {
		message = message + fmt.Sprintf(`<a href="%s">Image</a>\n`, imageURL)
	}

	if metrics := generateMetricsMessage(alert); metrics != "" {
		message = message + fmt.Sprintf("\n<i>Metrics:</i>%s", metrics)
	}

	return message
}

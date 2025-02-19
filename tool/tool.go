package tool

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/url"
	"regexp"
	"strings"
)

func GetDeepSeekResp(input string) (string, string) {
	log.Info("原始回答>> ", input)
	reThink := regexp.MustCompile(`(?s)<think>(.*?)</think>`)
	reAnswer := regexp.MustCompile(`</think>\s*(.*)`)

	input = strings.ReplaceAll(input, "\n", "")
	thinkMatch := reThink.FindStringSubmatch(input)
	answerMatch := reAnswer.FindStringSubmatch(input)

	var auThink, aiAnswer string
	if len(thinkMatch) > 1 {
		auThink = thinkMatch[1]
	}
	if len(answerMatch) > 1 {
		aiAnswer = answerMatch[1]
	} else {
		aiAnswer = input
	}

	return auThink, aiAnswer
}

func extractBaseDomain(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	return parsedURL.Host, nil
}

func GetSessionID() string {
	return uuid.New().String()
}

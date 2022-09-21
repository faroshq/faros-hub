package yaml

import (
	"io"
	"regexp"
	"strings"
)

// YAMLDocuments returns a collection of documents from the reader
func YAMLDocuments(reader io.Reader) ([]string, error) {

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	splitter := regexp.MustCompile("(?m)^---\n")

	var list []string

	for _, document := range splitter.Split(string(content), -1) {
		if strings.TrimSpace(document) == "" {
			continue
		}
		list = append(list, document)
	}

	return list, nil
}

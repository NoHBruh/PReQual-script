package helper

import (
	"PReQual/model"
	"encoding/json"
	"os"
	"path/filepath"
)

func WriteMetaDataFile(path string, pr model.PullRequest) {
	writePullRequestMetaFile(path, pr)
	writeCommentsMetaFile(path, pr)
}

func writeCommentsMetaFile(path string, pr model.PullRequest) {
	type CommentJSON struct {
		CreatedAt string `json:"created_at"`
		Body      string `json:"body"`
	}

	comments := make([]CommentJSON, 0, len(pr.Comments))
	for _, c := range pr.Comments {
		comments = append(comments, CommentJSON{
			CreatedAt: c.CreatedAt,
			Body:      c.Body,
		})
	}

	commentsData, err := json.MarshalIndent(comments, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filepath.Join(path, "comments.json"), commentsData, 0644)
	if err != nil {
		panic(err)
	}
}

func writePullRequestMetaFile(path string, pr model.PullRequest) {
	metadata := map[string]interface{}{
		"title":      pr.Title,
		"body":       pr.Body,
		"created_at": pr.CreateAt,
		"closed_at":  pr.ClosedAt,
		"state":      pr.State,
	}

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		panic(err)
	}

	filePath := filepath.Join(path, "metadata.json")
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		panic(err)
	}
}

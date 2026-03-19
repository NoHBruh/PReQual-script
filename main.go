package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"PReQual/client"
	"PReQual/helper"
	"PReQual/metric"
)

const defaultWorkspace = "tmp"
const defaultMetric = "complexity,cognitive_complexity"

var repos = []string{
	"ReViSE-EuroSpaceCenter/ReViSE-backend",
	"sipeed/picoclaw",
	"iluwatar/java-design-patterns",
	"TheAlgorithms/Java",
	"google/guava",
	"dbeaver/dbeaver",
	"apache/dubbo",
	"netty/netty",
	"keycloak/keycloak",
}

func main() {
	// CLI flags
	reposArg := flag.String("repos", "", "GitHub repositories in the form <owner>/<repo>(,<owner>/<repo>)* (required)")
	workspace := flag.String("workspace", defaultWorkspace, "Workspace directory (default: tmp)")
	metricsArg := flag.String("metrics", defaultMetric, "Comma-separated list of metrics to analyze (default: complexity,cognitive_complexity)")

	flag.Parse()

	if *reposArg == "" {
		fmt.Println("Error: -repos argument is required")
		flag.Usage()
		os.Exit(1)
	}

	repos := strings.Split(*reposArg, ",")

	metrics := strings.Split(*metricsArg, ",")

	var prClient client.PullRequestClient
	prClient = &client.GhClient{}

	var analyzer metric.ProjectAnalyser
	analyzer = &metric.SonarQubeAnalyzer{}

	for _, repo := range repos {
		fmt.Printf("\n===== Traitement du repo: %s =====\n", repo)

		prs, err := prClient.GetPullRequests(repo)
		if err != nil {
			fmt.Printf("Error fetching pull requests: %v\n", err)
			return
		}

		for _, pr := range prs {
			var startTime = time.Now()
			fmt.Printf("PR #%d: %s (Base: %s, Head: %s)\n", pr.Number, pr.Title, pr.BaseRefOid, pr.HeadRefOid)

			var path = fmt.Sprintf("%s/%s/pr_%d", *workspace, repo, pr.Number)

			if err := prClient.RetrieveBranchZip(repo, pr.HeadRefOid, path, "head.zip"); err != nil {
				return
			}
			if err = prClient.RetrieveBranchZip(repo, pr.BaseRefOid, path, "base.zip"); err != nil {
				return
			}

			helper.WriteMetaDataFile(path, pr)

			formattedRepo := strings.Replace(repo, "/", "-", -1)

			err := analyzer.AnalyzeProject(formattedRepo, path, metrics)
			var totalDuration = helper.FormatDuration(time.Since(startTime))
			var totalSize = helper.FormatSizeRounded([]string{path + "/head.zip", path + "/base.zip"})
			var baseSize = helper.FormatSizeRounded([]string{path + "/base.zip"})
			var headSize = helper.FormatSizeRounded([]string{path + "/head.zip"})
			fmt.Printf("End of the PR #%d's analysis.\n", pr.Number)
			fmt.Printf("Total ZIP files size : %s (base.zip: %s ; head.zip %s).\n", totalSize, baseSize, headSize)
			fmt.Printf("Total duration of the PR's analysis : %s.\n", totalDuration)
			if err != nil {
				fmt.Printf("Error analyzing pull requests: %v\n", err)
				return
			}
		}
	}
}

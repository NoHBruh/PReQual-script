package model

type SonarMeasures struct {
	Component struct {
		Measures []struct {
			Metric string `json:"metric"`
			Value  string `json:"value"`
		} `json:"measures"`
	} `json:"component"`
}

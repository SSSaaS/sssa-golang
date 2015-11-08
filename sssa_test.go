package sssa

import (
	"testing"
)

func TestCreateCombine(t *testing.T) {
	// Short, medium, and long tests
	strings := []string{
		"N17FigASkL6p1EOgJhRaIquQLGvYV0",
		"0y10VAfmyH7GLQY6QccCSLKJi8iFgpcSBTLyYOGbiYPqOpStAf1OYuzEBzZR",
		"KjRHO1nHmIDidf6fKvsiXWcTqNYo2U9U8juO94EHXVqgearRISTQe0zAjkeUYYBvtcB8VWzZHYm6ktMlhOXXCfRFhbJzBUsXaHb5UDQAvs2GKy6yq0mnp8gCj98ksDlUultqygybYyHvjqR7D7EAWIKPKUVz4of8OzSjZlYg7YtCUMYhwQDryESiYabFID1PKBfKn5WSGgJBIsDw5g2HB2AqC1r3K8GboDN616Swo6qjvSFbseeETCYDB3ikS7uiK67ErIULNqVjf7IKoOaooEhQACmZ5HdWpr34tstg18rO",
	}

	minimum := []int{4, 6, 20}
	shares := []int{5, 100, 100}

	for i := range strings {
		if Combine(Create(minimum[i], shares[i], strings[i])) != strings[i] {
			t.Fatal("Fatal: creating and combining returned invalid data")
		}
	}
}

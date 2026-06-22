//go:build unit

package newrelic

import (
	"testing"
)

func TestFormatEntitySearchQueryTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple single tag",
			input:    "tags.accountId = '12345678'",
			expected: "`tags.accountId` = '12345678'",
		},
		{
			name:     "Multiple tags with AND",
			input:    "tags.accountId = '12345678' AND tags.environment='production' AND tags.language='java'",
			expected: "`tags.accountId` = '12345678' AND `tags.environment`='production' AND `tags.language`='java'",
		},
		{
			name:     "Tags with OR",
			input:    "tags.env='prod' OR tags.env='staging'",
			expected: "`tags.env`='prod' OR `tags.env`='staging'",
		},
		{
			name:     "Already backticked tags should not be modified",
			input:    "`tags.accountId` = '12345678' AND tags.environment='production'",
			expected: "`tags.accountId` = '12345678' AND `tags.environment`='production'",
		},
		{
			name:     "All tags already backticked",
			input:    "`tags.accountId` = '12345678' AND `tags.environment`='production'",
			expected: "`tags.accountId` = '12345678' AND `tags.environment`='production'",
		},
		{
			name:     "Tags with underscores",
			input:    "tags.account_id = '123' AND tags.env_name='test'",
			expected: "`tags.account_id` = '123' AND `tags.env_name`='test'",
		},
		{
			name:     "Tags with numbers",
			input:    "tags.version2 = '1.0' AND tags.cluster3='prod'",
			expected: "`tags.version2` = '1.0' AND `tags.cluster3`='prod'",
		},
		{
			name:     "Complex query with parentheses",
			input:    "(tags.env='prod' OR tags.env='staging') AND tags.app='myapp'",
			expected: "(`tags.env`='prod' OR `tags.env`='staging') AND `tags.app`='myapp'",
		},
		{
			name:     "Tag at the beginning of query",
			input:    "tags.environment='prod'",
			expected: "`tags.environment`='prod'",
		},
		{
			name:     "Tag at the end of query",
			input:    "type='APPLICATION' AND tags.env='prod'",
			expected: "type='APPLICATION' AND `tags.env`='prod'",
		},
		{
			name:     "Mixed conditions with non-tag fields",
			input:    "type='APPLICATION' AND tags.accountId='123' AND name LIKE 'prod%'",
			expected: "type='APPLICATION' AND `tags.accountId`='123' AND name LIKE 'prod%'",
		},
		{
			name:     "Empty query",
			input:    "",
			expected: "",
		},
		{
			name:     "Query without tags",
			input:    "type='APPLICATION' AND name='myapp'",
			expected: "type='APPLICATION' AND name='myapp'",
		},
		{
			name:     "Nested tags - tags.tags.REPOSITORY",
			input:    "tags.tags.REPOSITORY = 'myrepo'",
			expected: "`tags.tags.REPOSITORY` = 'myrepo'",
		},
		{
			name:     "Multiple nested tags",
			input:    "tags.tags.REPOSITORY = 'repo1' AND tags.tags.BRANCH='main'",
			expected: "`tags.tags.REPOSITORY` = 'repo1' AND `tags.tags.BRANCH`='main'",
		},
		{
			name:     "Mixed single and nested tags",
			input:    "tags.environment='prod' AND tags.tags.REPOSITORY='myrepo'",
			expected: "`tags.environment`='prod' AND `tags.tags.REPOSITORY`='myrepo'",
		},
		{
			name:     "Three-level nested tags",
			input:    "tags.foo.bar.baz = 'value'",
			expected: "`tags.foo.bar.baz` = 'value'",
		},
		{
			name:     "Nested tags with underscores",
			input:    "tags.tags.REPOSITORY_NAME = 'test_repo'",
			expected: "`tags.tags.REPOSITORY_NAME` = 'test_repo'",
		},
		{
			name:     "Already backticked nested tags",
			input:    "`tags.tags.REPOSITORY` = 'myrepo' AND tags.tags.BRANCH='main'",
			expected: "`tags.tags.REPOSITORY` = 'myrepo' AND `tags.tags.BRANCH`='main'",
		},
		{
			name:     "Complex query with nested tags and parentheses",
			input:    "(tags.tags.REPOSITORY='repo1' OR tags.tags.REPOSITORY='repo2') AND tags.env='prod'",
			expected: "(`tags.tags.REPOSITORY`='repo1' OR `tags.tags.REPOSITORY`='repo2') AND `tags.env`='prod'",
		},
		{
			name:     "Four-level nested tags",
			input:    "tags.tags.tags.tags.REPO = 'test'",
			expected: "`tags.tags.tags.tags.REPO` = 'test'",
		},
		{
			name:     "Five-level nested tags",
			input:    "tags.a.b.c.d.e = 'value'",
			expected: "`tags.a.b.c.d.e` = 'value'",
		},
		{
			name:     "Deep nesting with multiple conditions",
			input:    "tags.tags.tags.REPO='r1' AND tags.tags.BRANCH='main' AND tags.ENV='prod'",
			expected: "`tags.tags.tags.REPO`='r1' AND `tags.tags.BRANCH`='main' AND `tags.ENV`='prod'",
		},
		{
			name:     "Mixed depth nested tags in one query",
			input:    "tags.a='1' AND tags.b.c='2' AND tags.d.e.f='3' AND tags.g.h.i.j='4'",
			expected: "`tags.a`='1' AND `tags.b.c`='2' AND `tags.d.e.f`='3' AND `tags.g.h.i.j`='4'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatEntitySearchQueryTags(tt.input)
			if result != tt.expected {
				t.Errorf("formatEntitySearchQueryTags() = %v, want %v", result, tt.expected)
			}
		})
	}
}
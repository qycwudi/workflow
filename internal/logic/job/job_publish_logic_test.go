package job

import "testing"

func TestValidateCronExpression(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		shouldError bool
	}{
		{
			name:        "Valid expression with all stars",
			expression:  "* * * * *",
			shouldError: false,
		},
		{
			name:        "Valid expression with numbers",
			expression:  "0 0 1 1 0",
			shouldError: false,
		},
		{
			name:        "Valid expression with ranges",
			expression:  "1-5 0-23 1-31 1-12 0-6",
			shouldError: false,
		},
		{
			name:        "Valid expression with intervals",
			expression:  "*/15 */6 */10 */3 */2",
			shouldError: false,
		},
		{
			name:        "Valid expression with lists",
			expression:  "1,15,30 0,12,23 1,15,31 1,6,12 0,3,6",
			shouldError: false,
		},
		{
			name:        "Invalid number of fields",
			expression:  "* * * *",
			shouldError: true,
		},
		{
			name:        "Invalid minute value",
			expression:  "60 * * * *",
			shouldError: true,
		},
		{
			name:        "Invalid hour value",
			expression:  "* 24 * * *",
			shouldError: true,
		},
		{
			name:        "Invalid day value",
			expression:  "* * 32 * *",
			shouldError: true,
		},
		{
			name:        "Invalid month value",
			expression:  "* * * 13 *",
			shouldError: true,
		},
		{
			name:        "Invalid week value",
			expression:  "* * * * 7",
			shouldError: true,
		},
		{
			name:        "Invalid interval format",
			expression:  "1/15/30 * * * *",
			shouldError: true,
		},
		{
			name:        "Invalid range format",
			expression:  "1-5-10 * * * *",
			shouldError: true,
		},
		{
			name:        "Invalid range values",
			expression:  "5-1 * * * *",
			shouldError: true,
		},
		{
			name:        "Valid expression with question mark",
			expression:  "* * 1 * ?",
			shouldError: false,
		},
		{
			name:        "Valid expression with question mark",
			expression:  "1",
			shouldError: false,
		},
		{
			name:        "Valid expression with question mark",
			expression:  "* * ?",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCronExpression(tt.expression)
			if (err != nil) != tt.shouldError {
				t.Errorf("ValidateCronExpression() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

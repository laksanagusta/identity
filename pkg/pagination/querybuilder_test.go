package pagination

import (
	"testing"
)

func TestQueryBuilder_AddSearch(t *testing.T) {
	tests := []struct {
		name          string
		initialQuery  string
		searchTerm    string
		searchFields  []string
		expectedQuery string
		expectedArgs  []interface{}
		expectError   bool
	}{
		{
			name:          "Single field search",
			initialQuery:  "SELECT * FROM users",
			searchTerm:    "john",
			searchFields:  []string{"name"},
			expectedQuery: "SELECT * FROM users WHERE (LOWER(name) ILIKE LOWER($1))",
			expectedArgs:  []interface{}{"%john%"},
			expectError:   false,
		},
		{
			name:          "Multiple field search",
			initialQuery:  "SELECT * FROM users",
			searchTerm:    "john",
			searchFields:  []string{"name", "uuid"}, // uuid is in validFields
			expectedQuery: "SELECT * FROM users WHERE (LOWER(name) ILIKE LOWER($1) OR LOWER(uuid) ILIKE LOWER($1))",
			expectedArgs:  []interface{}{"%john%"},
			expectError:   false,
		},
		{
			name:          "Empty search term",
			initialQuery:  "SELECT * FROM users",
			searchTerm:    "",
			searchFields:  []string{"name"},
			expectedQuery: "SELECT * FROM users",
			expectedArgs:  nil,
			expectError:   false,
		},
		{
			name:          "Invalid field",
			initialQuery:  "SELECT * FROM users",
			searchTerm:    "john",
			searchFields:  []string{"invalid_field"},
			expectedQuery: "",
			expectedArgs:  nil,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder(tt.initialQuery)
			err := qb.AddSearch(tt.searchTerm, tt.searchFields)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			query, args := qb.Build()

			if query != tt.expectedQuery {
				t.Errorf("expected query %q, got %q", tt.expectedQuery, query)
			}

			if len(args) != len(tt.expectedArgs) {
				t.Errorf("expected %d args, got %d", len(tt.expectedArgs), len(args))
			} else {
				for i, arg := range args {
					if arg != tt.expectedArgs[i] {
						t.Errorf("expected arg %v, got %v", tt.expectedArgs[i], arg)
					}
				}
			}
		})
	}
}

package formatter

import "testing"

func TestFormatSessionName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic cases
		{"simple name", "myproject", "myproject"},
		{"dots to underscores", "my.project.name", "my_project_name"},

		// Whitespace handling
		{"leading whitespace", "  myproject", "myproject"},
		{"trailing whitespace", "myproject  ", "myproject"},
		{"both whitespace", "  my.project  ", "my_project"},
		{"tabs and spaces", "\t my.project \n", "my_project"},

		// Leading/trailing special characters
		{"leading hyphens", "--my-project", "my-project"},
		{"trailing hyphens", "my-project--", "my-project"},
		{"leading underscores", "__my_project", "my_project"},
		{"trailing underscores", "my_project__", "my_project"},
		{"leading dots", "...my.project", "my_project"},
		{"trailing dots", "my.project...", "my_project"},
		{"mixed leading chars", "-_.my.project", "my_project"},
		{"mixed trailing chars", "my.project.-_", "my_project"},

		// Repository name patterns
		{"github style", "user/repo-name", "user/repo-name"},
		{"namespaced project", "org.example.project", "org_example_project"},
		{"version in name", "project-v1.2.3", "project-v1_2_3"},
		{"scoped package", "@angular/core", "@angular/core"},
		{"js framework", "react.js", "react_js"},
		{"vue router", "vue-router", "vue-router"},
		{"dotnet project", "My.Company.ProjectName", "My_Company_ProjectName"},

		// Complex combinations
		{"mixed separators", "my-project_name.api", "my-project_name_api"},
		{"multiple dots", "com.example.my.project", "com_example_my_project"},
		{"camelCase with dots", "myProject.apiService", "myProject_apiService"},

		// Edge cases
		{"empty string", "", ""},
		{"whitespace only", "   ", ""},
		{"special chars only", "-_.", ""},
		{"mixed whitespace and special", "  -_.  ", ""},
		{"single character", "a", "a"},
		{"single dot", ".", ""},

		// Real-world examples
		{"kubernetes", "k8s.io/client-go", "k8s_io/client-go"},
		{"spring boot", "spring-boot.starter", "spring-boot_starter"},
		{"nodejs package", "lodash.debounce", "lodash_debounce"},
		{"maven artifact", "com.fasterxml.jackson", "com_fasterxml_jackson"},
		{"python package", "django-rest-framework", "django-rest-framework"},

		// Stress cases
		{"many leading chars", "---___...my.project", "my_project"},
		{"many trailing chars", "my.project...___---", "my_project"},
		{"scattered dots", "a.b.c.d.e.f", "a_b_c_d_e_f"},
		{"alternating chars", "-_.-_.-_.project.-_.-_", "project"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatSessionName(tt.input)
			if result != tt.expected {
				t.Errorf("FormatSessionName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type TestCase struct {
	Name       string
	FilePath   string
	Command    string
	WantOutput string
	WantError  []string
	WantExit   int
}

func TestMain(m *testing.M) {
	// Build interpreter before tests
	build := exec.Command("go", "build", "-o", "phaeton", "./main.go")
	if err := build.Run(); err != nil {
		panic("Build failed: " + err.Error())
	}
	code := m.Run()
	os.Remove("phaeton")
	os.Exit(code)
}

func TestIntegration(t *testing.T) {
	tests := []TestCase{
		{
			Name:       "IfStatement",
			FilePath:   "tests/control_flow/if/test1.phn",
			Command:    "run",
			WantOutput: "tests/control_flow/if/expected_test1.txt",
			WantExit:   0,
		},
		{
			Name:       "nested_Complex",
			FilePath:   "tests/control_flow/nested/test1.phn",
			Command:    "run",
			WantOutput: "tests/control_flow/nested/expected_test1.txt",
			WantExit:   0,
		},
		{
			Name:       "nested_block",
			FilePath:   "tests/control_flow/nested/test2.phn",
			Command:    "run",
			WantOutput: "tests/control_flow/nested/expected_test2.txt",
			WantExit:   0,
		},
		{
			Name:       "nested_statement",
			FilePath:   "tests/control_flow/nested/test3.phn",
			Command:    "run",
			WantOutput: "tests/control_flow/nested/expected_test3.txt",
			WantExit:   0,
		},
		{
			Name:       "tokenize_basic",
			FilePath:   "tests/tokenize/test1.phn",
			Command:    "tokenize",
			WantOutput: "tests/tokenize/expected_test1.txt",
			WantExit:   0,
		},
		// {
		// 	Name:     "Parsing_InvalidSyntax",
		// 	FilePath: "parsing/invalid_syntax.phn",
		// 	Command:  "parse",
		// 	WantError: []string{
		// 		"Unterminated string",
		// 	},
		// 	WantExit: 65,
		// },
		// // Add more test cases
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			fullPath := filepath.Join("../", tc.FilePath)
			wantFullPath := filepath.Join("../", tc.WantOutput)

			cmd := exec.Command("./phaeton", tc.Command, fullPath)
			output, _ := cmd.CombinedOutput()

			gotExit := cmd.ProcessState.ExitCode()
			if gotExit != tc.WantExit {
				t.Errorf("Exit code mismatch\nWant: %d\nGot:  %d", tc.WantExit, gotExit)
			}

			want, err := os.ReadFile(wantFullPath)
			if err != nil {
				t.Fatalf("Failed to read expected output file: %v", err)
			}

			wantStr := strings.TrimSpace(string(want))
			outStr := strings.TrimSpace(string(output))

			if wantStr != outStr {
				t.Errorf("Output mismatch\n--- Want ---\n%s\n\n--- Got ---\n%s", wantStr, outStr)
			}

			for _, wantErr := range tc.WantError {
				if !strings.Contains(outStr, wantErr) {
					t.Errorf("Missing expected error: %q", wantErr)
				}
			}
		})
	}
}

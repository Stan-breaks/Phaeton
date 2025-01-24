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
	WantOutput []string
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
			Name:     "IfStatement_Basic",
			FilePath: "tests/control_flow/if/test1.phn",
			Command:  "run",
			WantOutput: []string{
				"if branch",
			},
			WantExit: 0,
		},
		{
			Name:     "nested_Complex",
			FilePath: "tests/control_flow/nested/test1.phn",
			Command:  "run",
			WantOutput: []string{
				"young adult\neligible for voting: true\nfirst-time voter: likely\neligible for driving: full license\neligible for drinking (US): true\nremember: drink responsibly",
			},
			WantExit: 0,
		},
		{
			Name:     "nested_Basic",
			FilePath: "tests/control_flow/nested/test2.phn",
			Command:  "run",
			WantOutput: []string{
				"world",
			},
			WantExit: 0,
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

			cmd := exec.Command("./phaeton", tc.Command, fullPath)
			output, _ := cmd.CombinedOutput()
			gotExit := cmd.ProcessState.ExitCode()

			// Validate exit code
			if gotExit != tc.WantExit {
				t.Errorf("Exit code mismatch\nWant: %d\nGot:  %d", tc.WantExit, gotExit)
			}

			// Validate output
			outStr := string(output)
			for _, want := range tc.WantOutput {
				if !strings.Contains(outStr, want) {
					t.Errorf("expected output: %q\n received output: %q", want, outStr)
				}
			}

			// Validate errors
			for _, wantErr := range tc.WantError {
				if !strings.Contains(outStr, wantErr) {
					t.Errorf("Missing expected error: %q", wantErr)
				}
			}
		})
	}
}

package ritual

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRitualCommand(t *testing.T) {
	cmd := RitualCommand()
	
	if cmd == nil {
		t.Fatal("RitualCommand() returned nil")
	}
	
	if cmd.Use != "ritual" {
		t.Errorf("Expected Use='ritual', got %q", cmd.Use)
	}
	
	if cmd.Short == "" {
		t.Error("Short description should not be empty")
	}
}

func TestRitualInitCommand(t *testing.T) {
	cmd := RitualCommand()
	
	// Find init subcommand
	initCmd := findSubcommand(cmd, "init")
	if initCmd == nil {
		t.Fatal("init subcommand not found")
	}
	
	if initCmd.Use != "init <ritual-name>" {
		t.Errorf("Expected Use='init <ritual-name>', got %q", initCmd.Use)
	}
}

func TestRitualListCommand(t *testing.T) {
	cmd := RitualCommand()
	
	// Find list subcommand
	listCmd := findSubcommand(cmd, "list")
	if listCmd == nil {
		t.Fatal("list subcommand not found")
	}
	
	if listCmd.Use != "list" {
		t.Errorf("Expected Use='list', got %q", listCmd.Use)
	}
}

func TestRitualInfoCommand(t *testing.T) {
	cmd := RitualCommand()
	
	// Find info subcommand
	infoCmd := findSubcommand(cmd, "info")
	if infoCmd == nil {
		t.Fatal("info subcommand not found")
	}
	
	if infoCmd.Use != "info <ritual-name>" {
		t.Errorf("Expected Use='info <ritual-name>', got %q", infoCmd.Use)
	}
}

func TestRitualValidateCommand(t *testing.T) {
	cmd := RitualCommand()
	
	// Find validate subcommand
	validateCmd := findSubcommand(cmd, "validate")
	if validateCmd == nil {
		t.Fatal("validate subcommand not found")
	}
	
	if validateCmd.Use != "validate" {
		t.Errorf("Expected Use='validate', got %q", validateCmd.Use)
	}
}

func TestRitualCreateCommand(t *testing.T) {
	cmd := RitualCommand()
	
	// Find create subcommand
	createCmd := findSubcommand(cmd, "create")
	if createCmd == nil {
		t.Fatal("create subcommand not found")
	}
	
	if createCmd.Use != "create <name>" {
		t.Errorf("Expected Use='create <name>', got %q", createCmd.Use)
	}
}

func TestInitRitual_CreatesProjectStructure(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "ritual-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	
	projectDir := filepath.Join(tmpDir, "test-project")
	
	// This test will fail initially until we implement the function
	// For now, we're just defining the expected behavior
	
	// TODO: Implement this test once initRitual function is created
	// err = initRitual("basic-site", projectDir)
	// if err != nil {
	// 	t.Fatalf("initRitual failed: %v", err)
	// }
	
	// // Verify project structure was created
	// expectedDirs := []string{
	// 	"handlers",
	// 	"templates",
	// 	"static",
	// 	"config",
	// }
	// 
	// for _, dir := range expectedDirs {
	// 	path := filepath.Join(projectDir, dir)
	// 	if _, err := os.Stat(path); os.IsNotExist(err) {
	// 		t.Errorf("Expected directory %s to exist", dir)
	// 	}
	// }
	
	t.Skip("Test not yet implemented - waiting for initRitual function")
}

func TestListRituals_ShowsAvailableRituals(t *testing.T) {
	// Create temp directory for ritual registry
	tmpDir, err := os.MkdirTemp("", "ritual-registry-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	
	// TODO: Implement this test once listRituals function is created
	// Create a mock ritual in the registry
	// ritualDir := filepath.Join(tmpDir, "basic-site")
	// if err := os.MkdirAll(ritualDir, 0755); err != nil {
	// 	t.Fatal(err)
	// }
	
	// // Create ritual.yaml
	// ritualYAML := `name: basic-site
	// version: 1.0.0
	// description: A basic website template
	// `
	// if err := os.WriteFile(filepath.Join(ritualDir, "ritual.yaml"), []byte(ritualYAML), 0644); err != nil {
	// 	t.Fatal(err)
	// }
	
	// // List rituals
	// var buf bytes.Buffer
	// err = listRituals(tmpDir, &buf)
	// if err != nil {
	// 	t.Fatalf("listRituals failed: %v", err)
	// }
	
	// output := buf.String()
	// if !strings.Contains(output, "basic-site") {
	// 	t.Errorf("Expected output to contain 'basic-site', got: %s", output)
	// }
	
	t.Skip("Test not yet implemented - waiting for listRituals function")
}

func TestValidateRitual_ValidRitualYAML(t *testing.T) {
	// Create temp directory with a valid ritual.yaml
	tmpDir, err := os.MkdirTemp("", "ritual-validate-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	
	// TODO: Implement this test once validateRitual function is created
	t.Skip("Test not yet implemented - waiting for validateRitual function")
}

func TestValidateRitual_InvalidRitualYAML(t *testing.T) {
	// Create temp directory with an invalid ritual.yaml
	tmpDir, err := os.MkdirTemp("", "ritual-validate-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Write invalid YAML
	invalidYAML := `name: test
version: invalid-version
`
	if err := os.WriteFile(filepath.Join(tmpDir, "ritual.yaml"), []byte(invalidYAML), 0644); err != nil {
		t.Fatal(err)
	}
	
	// TODO: Implement this test once validateRitual function is created
	// err = validateRitual(tmpDir)
	// if err == nil {
	// 	t.Error("Expected validation to fail for invalid YAML")
	// }
	
	t.Skip("Test not yet implemented - waiting for validateRitual function")
}

// Helper function to find subcommand by name
func findSubcommand(parent interface{ Commands() []interface{ Use() string } }, name string) interface{ Use() string } {
	for _, cmd := range parent.Commands() {
		if len(cmd.Use()) > 0 && cmd.Use()[:len(name)] == name {
			return cmd
		}
	}
	return nil
}

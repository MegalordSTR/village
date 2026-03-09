package deployment

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func projectRoot() string {
	// Walk up from this test file location to find go.mod
	dir, _ := filepath.Abs(filepath.Dir("."))
	for dir != "/" {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	return "."
}

func TestDockerfileExists(t *testing.T) {
	root := projectRoot()
	dockerfile := filepath.Join(root, "Dockerfile")
	if _, err := os.Stat(dockerfile); os.IsNotExist(err) {
		t.Fatalf("Dockerfile does not exist at %s", dockerfile)
	}
}

func TestDockerBuildBackendTarget(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker not available")
	}
	root := projectRoot()
	cmd := exec.Command("docker", "build", "--target", "backend", "-t", "village-backend-test", root) //nolint:gosec
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("docker build backend target failed: %v", err)
	}
}

func TestDockerBuildFrontendTarget(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker not available")
	}
	root := projectRoot()
	cmd := exec.Command("docker", "build", "--target", "frontend", "-t", "village-frontend-test", root) //nolint:gosec
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("docker build frontend target failed: %v", err)
	}
}

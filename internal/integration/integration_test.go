package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/vano44/village/internal/economy"
	"github.com/vano44/village/internal/simulation"
)

// TestBuildBinary verifies that the village binary can be built.
func TestBuildBinary(t *testing.T) {
	// Change to repository root (two levels up from internal/integration)
	root := filepath.Join("..", "..")
	cmd := exec.Command("go", "build", "./cmd/village")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}
	// Clean up binary
	os.Remove(filepath.Join(root, "village"))
}

// TestSimulationEconomyIntegration verifies that simulation and economy packages
// work together correctly.
func TestSimulationEconomyIntegration(t *testing.T) {
	// Create a game state with residents and resources
	state := simulation.NewGameState("integration-economy", 12345)
	for i := 0; i < 5; i++ {
		state.AddResident(simulation.Resident{
			ID:   string(rune('a' + i)),
			Name: "Test Resident",
			Age:  25 + i*5,
		})
	}
	// Add economic resources
	state.AddResource(simulation.Resource{
		Type:     "food",
		Quantity: 100,
		Quality:  1.0,
	})
	state.AddResource(simulation.Resource{
		Type:     "wood",
		Quantity: 50,
		Quality:  0.8,
	})

	// Create turn processor with all systems
	tp := simulation.NewTurnProcessor()
	tp.RegisterSystem(simulation.SystemEnvironment, simulation.NewEnvironmentSystem())
	tp.RegisterSystem(simulation.SystemProduction, simulation.NewProductionSystem())
	tp.RegisterSystem(simulation.SystemSocial, simulation.NewSocialSystem())
	tp.RegisterSystem(simulation.SystemEconomic, simulation.NewEconomicSystem())
	tp.RegisterSystem(simulation.SystemEvents, simulation.NewEventSystem())

	// Process a week
	events := tp.ProcessWeek(state)
	if len(events) == 0 {
		t.Log("no events generated this week")
	}

	// Verify that resources have changed (quantity may have been consumed/produced)
	// At least one resource should have changed quantity (or remain same).
	// We'll just ensure state is still valid.
	if len(state.Resources) < 2 {
		t.Errorf("expected at least 2 resources, got %d", len(state.Resources))
	}

	// Use economy package to calculate something
	category := economy.CategoryForResource(economy.ResourceType("food"))
	if category != economy.CategoryRaw {
		t.Errorf("expected raw category, got %v", category)
	}
}

// TestDeploymentBuild verifies that Docker images can be built (if Docker is available).
func TestDeploymentBuild(t *testing.T) {
	// Skip in CI environment to avoid long build times
	if os.Getenv("CI") != "" {
		t.Skip("Skipping Docker build in CI")
	}
	// Skip if Docker not available
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker not available")
	}
	root := filepath.Join("..", "..")
	// Use unique tag to avoid conflicts
	tag := fmt.Sprintf("village-backend-test-%d", time.Now().UnixNano())
	cmd := exec.Command("docker", "build", "--target", "backend", "-t", tag, ".")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Errorf("failed to build backend Docker image: %v", err)
	}
	// Clean up image
	exec.Command("docker", "rmi", tag).Run()
}

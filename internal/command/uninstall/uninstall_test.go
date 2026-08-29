package uninstall

import (
	"reflect"
	"testing"

	_ "github.com/jamesawo/mdev/internal/tools/gradle"
	_ "github.com/jamesawo/mdev/internal/tools/java"
	_ "github.com/jamesawo/mdev/internal/tools/maven"
	_ "github.com/jamesawo/mdev/internal/tools/podman"
	_ "github.com/jamesawo/mdev/internal/tools/podmancompose"
	_ "github.com/jamesawo/mdev/internal/tools/podmandesktop"
	_ "github.com/jamesawo/mdev/internal/tools/sdkman"
)

func TestInstalledPlanExcludesMissingDependents(t *testing.T) {
	status := map[string]bool{"sdkman": true}
	plan, err := installedPlan([]string{"maven", "gradle", "java", "sdkman"}, func(name string) (bool, error) {
		return status[name], nil
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"sdkman"}
	if !reflect.DeepEqual(plan, want) {
		t.Fatalf("plan = %v, want %v", plan, want)
	}

}

func TestBuildPlanIncludesTransitiveDependents(t *testing.T) {
	plan, err := BuildPlan("sdkman")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{"maven": true, "gradle": true, "java": true, "sdkman": true}
	for _, name := range plan {
		delete(want, name)
	}
	if len(want) != 0 {
		t.Fatalf("plan = %v, missing %v", plan, want)
	}
}

func TestBuildPlanRemovesPodmanAddOnsBeforeBaseTool(t *testing.T) {
	plan, err := BuildPlan("podman")
	if err != nil {
		t.Fatal(err)
	}
	positions := map[string]int{}
	for index, name := range plan {
		positions[name] = index
	}
	base, hasBase := positions["podman"]
	desktop, hasDesktop := positions["podman-desktop"]
	compose, hasCompose := positions["podman-compose"]
	if !hasBase || !hasDesktop || !hasCompose {
		t.Fatalf("plan = %v, want all Podman tools", plan)
	}
	if desktop >= base || compose >= base {
		t.Fatalf("plan = %v, add-ons must precede podman", plan)
	}
}

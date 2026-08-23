package uninstall

import (
	"reflect"
	"testing"

	_ "github.com/jamesawo/mdev/internal/tools/gradle"
	_ "github.com/jamesawo/mdev/internal/tools/java"
	_ "github.com/jamesawo/mdev/internal/tools/maven"
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

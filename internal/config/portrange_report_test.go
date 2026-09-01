package config

import "testing"

// TestAPortOutsideItsRangeDoesNotBlockLoading is the check behind a defect that
// locked a user out. Widening the ranges made an existing configuration invalid,
// every command that reads it failed, and ports reassign, which repairs it,
// reads the same file.
func TestAPortOutsideItsRangeDoesNotBlockLoading(t *testing.T) {
	t.Parallel()

	value := Config{
		Version: 1,
		Project: Project{Name: "example"},
		Runtime: DefaultRuntime(),
		Services: []Service{{
			Name: "developer", Command: "npm run dev:developer", Role: RoleDeveloper,
			// Inside the range this role had before it was widened.
			Port: 3100, Host: "127.0.0.1", HealthPath: "/",
		}},
	}

	if err := value.Validate(); err != nil {
		t.Fatalf("a port outside its range must not stop the file loading: %v", err)
	}

	outside := value.PortsOutsideTheirRange()
	if len(outside) != 1 {
		t.Fatalf("reported %+v, want the one service named", outside)
	}
	if outside[0].Service != "developer" || outside[0].Port != 3100 {
		t.Fatalf("reported %+v, want the service and its port", outside[0])
	}
	allowed, _ := RoleDeveloper.PortRange()
	if outside[0].Allowed != allowed {
		t.Fatalf("reported range %+v, want the role's own %+v", outside[0].Allowed, allowed)
	}
}

func TestAPortInsideItsRangeIsNotReported(t *testing.T) {
	t.Parallel()

	allowed, _ := RoleDeveloper.PortRange()
	value := Config{
		Version: 1,
		Project: Project{Name: "example"},
		Runtime: DefaultRuntime(),
		Services: []Service{
			{Name: "developer", Command: "x", Role: RoleDeveloper, Port: allowed.First, Host: "127.0.0.1", HealthPath: "/"},
			// A worker has no port and no range, so it is never reported.
			{Name: "queue", Command: "y", Role: RoleWorker, Port: 0},
		},
	}
	if outside := value.PortsOutsideTheirRange(); len(outside) != 0 {
		t.Fatalf("reported %+v, want nothing", outside)
	}
}

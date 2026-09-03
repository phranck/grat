package detect

import (
	"path/filepath"
	"strings"
)

// webProjectSDK marks a .NET project as an ASP.NET Core application. A project
// file without it is a library or a console program, neither of which serves.
const webProjectSDK = `Sdk="Microsoft.NET.Sdk.Web"`

// dotnetRunCommand starts the application with the port on the command line.
//
// The bare -- separator is what makes this reliable. `dotnet run` has no --urls
// option of its own and forwards unrecognised tokens to the application anyway,
// but its own documentation warns that a recognised option interleaved among
// them is removed and shifts the rest. After the separator every token reaches
// the application unread and unreordered.
//
// The address has to be given, rather than left to ASPNETCORE_URLS or to the
// default: a project scaffolded from a current template carries a generated
// endpoint in Properties/launchSettings.json on a random port, so what a fresh
// project binds without this is neither 5000 nor the port grat assigned.
const dotnetRunCommand = `dotnet run -- --urls "http://127.0.0.1:$PORT"`

// detectDotnet recognises an ASP.NET Core application by its project SDK.
func detectDotnet(root string) ([]Service, []Unresolved) {
	projects := make([]string, 0, 1)
	for _, entry := range entries(root) {
		if entry.IsDir() {
			continue
		}
		extension := filepath.Ext(entry.Name())
		if extension != ".csproj" && extension != ".fsproj" {
			continue
		}
		data, ok := readBounded(join(root, entry.Name()))
		if !ok {
			continue
		}
		if strings.Contains(string(data), webProjectSDK) {
			projects = append(projects, entry.Name())
		}
	}

	switch len(projects) {
	case 0:
		return nil, nil
	case 1:
		return []Service{service("backend", dotnetRunCommand, projects[0])}, nil
	default:
		return nil, []Unresolved{{
			Marker: projects[0],
			Reason: "several web projects sit beside each other (" + strings.Join(projects, ", ") +
				"), so which one dotnet run would start is not decidable from the directory",
		}}
	}
}

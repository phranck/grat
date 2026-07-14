# Graph Report - grat  (2026-07-14)

## Corpus Check
- 53 files · ~36,974 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 761 nodes · 1849 edges · 34 communities (29 shown, 5 thin omitted)
- Extraction: 83% EXTRACTED · 17% INFERRED · 0% AMBIGUOUS · INFERRED: 313 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `32ae4317`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- [[_COMMUNITY_CLI Command Orchestration|CLI Command Orchestration]]
- [[_COMMUNITY_Lifecycle TUI Engine|Lifecycle TUI Engine]]
- [[_COMMUNITY_Runtime Process Control|Runtime Process Control]]
- [[_COMMUNITY_Configuration Schema|Configuration Schema]]
- [[_COMMUNITY_Port CLI Integration Tests|Port CLI Integration Tests]]
- [[_COMMUNITY_Terminal Presentation|Terminal Presentation]]
- [[_COMMUNITY_Presentation Tests|Presentation Tests]]
- [[_COMMUNITY_Port Registry|Port Registry]]
- [[_COMMUNITY_Service Lifecycle Tests|Service Lifecycle Tests]]
- [[_COMMUNITY_Service Lifecycle Manager|Service Lifecycle Manager]]
- [[_COMMUNITY_Logging and Registry Locking|Logging and Registry Locking]]
- [[_COMMUNITY_README Project Guide|README Project Guide]]
- [[_COMMUNITY_Interactive Project Init|Interactive Project Init]]
- [[_COMMUNITY_Managed State Storage|Managed State Storage]]
- [[_COMMUNITY_Security and Support|Security and Support]]
- [[_COMMUNITY_Linux Listener Inspection|Linux Listener Inspection]]
- [[_COMMUNITY_Release Distribution|Release Distribution]]
- [[_COMMUNITY_Contribution Guidelines|Contribution Guidelines]]
- [[_COMMUNITY_Service Safety Model|Service Safety Model]]
- [[_COMMUNITY_Project Root Discovery|Project Root Discovery]]
- [[_COMMUNITY_Version Reporting|Version Reporting]]
- [[_COMMUNITY_CI Quality Gates|CI Quality Gates]]
- [[_COMMUNITY_Graphify Governance|Graphify Governance]]
- [[_COMMUNITY_CLI Helper Tests|CLI Helper Tests]]
- [[_COMMUNITY_README Verification|README Verification]]
- [[_COMMUNITY_Darwin Listener Tests|Darwin Listener Tests]]
- [[_COMMUNITY_Process Inspection Tests|Process Inspection Tests]]
- [[_COMMUNITY_Listener Lookup Interface|Listener Lookup Interface]]
- [[_COMMUNITY_Code of Conduct|Code of Conduct]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
- [[_COMMUNITY_Community 30|Community 30]]
- [[_COMMUNITY_Community 31|Community 31]]
- [[_COMMUNITY_Community 32|Community 32]]
- [[_COMMUNITY_Community 33|Community 33]]

## God Nodes (most connected - your core abstractions)
1. `Contains()` - 51 edges
2. `grat` - 46 edges
3. `Run()` - 35 edges
4. `New()` - 35 edges
5. `runWithEnvironment()` - 29 edges
6. `T` - 25 edges
7. `T` - 25 edges
8. `Renderer` - 23 edges
9. `Config` - 22 edges
10. `Manager` - 22 edges

## Surprising Connections (you probably didn't know these)
- `Configuration compatibility` --semantically_similar_to--> `Approved Declarative Tasks`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → README.md
- `Local Quality Gate` --semantically_similar_to--> `Verify Job`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → .github/workflows/ci.yml
- `Cross-Platform Compatibility` --semantically_similar_to--> `Platform Verification Matrix`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → .github/workflows/ci.yml
- `main()` --calls--> `Run()`  [INFERRED]
  cmd/grat/main.go → internal/cli/cli.go
- `grat` --references--> `Code of conduct`  [EXTRACTED]
  README.md → CONTRIBUTING.md

## Import Cycles
- None detected.

## Communities (34 total, 5 thin omitted)

### Community 0 - "CLI Command Orchestration"
Cohesion: 0.06
Nodes (81): assignReassignedPorts(), configuredRoots(), copyReservations(), defaultEnvironment(), detectServices(), ensureValidRegistry(), executeLifecycle(), fileExists() (+73 more)

### Community 1 - "Lifecycle TUI Engine"
Cohesion: 0.10
Nodes (34): CancelFunc, Cmd, Context, Reader, Style, Writer, lifecycleRow, Model (+26 more)

### Community 2 - "Runtime Process Control"
Cohesion: 0.11
Nodes (18): asset, installation, Client, Context, Service, Client, Context, Service (+10 more)

### Community 3 - "Configuration Schema"
Cohesion: 0.12
Nodes (32): Config, containsControl(), DefaultRuntime(), InferRole(), Load(), prepareWrite(), replaceFile(), rollbackWrites() (+24 more)

### Community 4 - "Port CLI Integration Tests"
Cohesion: 0.15
Nodes (36): exitCode(), Run(), assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort(), loadPortFixtureConfig(), runWithConfiguredRoots() (+28 more)

### Community 5 - "Terminal Presentation"
Cohesion: 0.10
Nodes (18): Renderer, Style, Renderer, Style, Writer, ColorMode, CommandGroup, helpUsageWidth() (+10 more)

### Community 6 - "Presentation Tests"
Cohesion: 0.19
Nodes (33): T, NewLifecycleModel(), DividerLine(), New(), columnOf(), displayColumn(), stripANSI(), TestHelpRendersBorderlessGroupsWithGloballyAlignedDescriptions() (+25 more)

### Community 7 - "Port Registry"
Cohesion: 0.13
Nodes (28): Config, Listener, T, fakeLookup, Listener, ListenerLookup, Problem, ProjectConfig (+20 more)

### Community 8 - "Service Lifecycle Tests"
Cohesion: 0.08
Nodes (42): Listener, Config, Listener, Manager, Service, T, processState, Service (+34 more)

### Community 9 - "Service Lifecycle Manager"
Cohesion: 0.11
Nodes (22): main(), mustGetwd(), Client, Config, Context, processState, ProgressObserver, Manager (+14 more)

### Community 10 - "Logging and Registry Locking"
Cohesion: 0.11
Nodes (18): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, File, T, Context, T, T, Mutex (+10 more)

### Community 11 - "README Project Guide"
Cohesion: 0.08
Nodes (28): CI Workflow, Code of Conduct, Command contract, Commands, Complete Development Stacks, Configuration reference, Contents, Contributing and support (+20 more)

### Community 12 - "Interactive Project Init"
Cohesion: 0.28
Nodes (15): collectInitInterview(), parseServiceDefinition(), promptDefault(), promptRequired(), promptServiceName(), TestCollectInitInterviewAcceptsEditsAndAdditionalServices(), TestCollectInitInterviewAcceptsSuppliedProjectName(), TestCollectInitInterviewRequiresAtLeastOneService() (+7 more)

### Community 13 - "Managed State Storage"
Cohesion: 0.28
Nodes (4): Manager, loadedState, processState, Time

### Community 14 - "Security and Support"
Cohesion: 0.21
Nodes (10): Fixed System Tool Paths, Private Vulnerability Reporting, Reporting a vulnerability, Security policy, Documented Shell Semantics, Supported versions, Trusted Configuration Boundary, Diagnostic Support Request (+2 more)

### Community 15 - "Linux Listener Inspection"
Cohesion: 0.25
Nodes (8): Listener, T, systemListener(), linuxListeningSocketInodes(), linuxSocketOwnerPIDs(), socketInode(), TestLinuxListenerParsingFindsOnlyTCPListenInodesForTheRequestedPort(), TestLinuxSocketOwnersFindsPIDsFromProcFileDescriptors()

### Community 16 - "Release Distribution"
Cohesion: 0.29
Nodes (10): Cross-Platform Compatibility, GitHub Releases, Release Binary Installation, Platform Verification Matrix, Release Build Job, Checksum Generation, Cross-Platform Binary Matrix, GitHub Release Publication (+2 more)

### Community 17 - "Contribution Guidelines"
Cohesion: 0.22
Nodes (8): Code of conduct, Configuration compatibility, Contributing to grat, Development setup, Focused Pull Requests, Pull requests, Approved Declarative Tasks, Declarative Non-Executable Configuration

### Community 18 - "Service Safety Model"
Cohesion: 0.25
Nodes (8): Bounded Local Logs, Cancellation Recovery, Listener Ownership and HTTP Readiness, Managed Process Identity, Process Group Isolation, Process-Owned Service Readiness, Safe Local Service Manager, Service Recovery Workflow

### Community 19 - "Project Root Discovery"
Cohesion: 0.47
Nodes (4): T, FindRoot(), TestFindRootReturnsNotFoundOutsideProject(), TestFindRootUsesNearestConfig()

### Community 20 - "Version Reporting"
Cohesion: 0.13
Nodes (32): repeatedValue, Service, Store, T, Context, Service, T, commandResponse (+24 more)

### Community 21 - "CI Quality Gates"
Cohesion: 0.40
Nodes (5): Local Quality Gate, CI Workflow, Race-Enabled Tests, Verify Job, Vulnerability Scan

### Community 22 - "Graphify Governance"
Cohesion: 0.67
Nodes (3): Graphify Before Push, Project Agent Rules, Versioned Graphify Artifacts

### Community 23 - "CLI Helper Tests"
Cohesion: 0.67
Nodes (3): TestListenerOwnerLabelHandlesUnknownPID(), TestLogFollowUsesTrustedExecutable(), T

### Community 24 - "README Verification"
Cohesion: 0.83
Nodes (3): require(), require_in(), check-readme.sh script

### Community 30 - "Community 30"
Cohesion: 0.15
Nodes (22): Store, T, Settings, canonicalExistingDirectory(), canonicalExistingPath(), ConfigDirectory(), canonicalPath(), equalStrings() (+14 more)

### Community 31 - "Community 31"
Cohesion: 0.21
Nodes (23): canonicalCLITestPath(), environmentForTest(), newCLITestStore(), sameStringSlices(), TestDirectoriesAddDoesNotPromptForInitialSetup(), TestDirectoriesCommandsPersistAndListConfiguredRoots(), TestFirstUseAcceptsExistingSitesDefault(), TestFirstUseFallsBackToWorkingDirectory() (+15 more)

### Community 32 - "Community 32"
Cohesion: 0.19
Nodes (14): Context, Service, Reader, Result, Store, Writer, installation, installationKind (+6 more)

### Community 33 - "Community 33"
Cohesion: 0.29
Nodes (7): Go HTTP API, Laravel, Project examples, Python with FastAPI, React, Laravel, and a queue worker, React with Vite, Swift with Vapor

## Knowledge Gaps
- **101 isolated node(s):** `Store`, `updateService`, `uninstallService`, `updateService`, `uninstallService` (+96 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `New()` connect `Presentation Tests` to `CLI Command Orchestration`, `Community 32`, `Runtime Process Control`, `Port CLI Integration Tests`, `Terminal Presentation`, `Service Lifecycle Manager`, `Interactive Project Init`, `Version Reporting`, `Community 30`?**
  _High betweenness centrality (0.185) - this node is a cross-community bridge._
- **Why does `Contains()` connect `Presentation Tests` to `Community 32`, `Configuration Schema`, `Port CLI Integration Tests`, `Service Lifecycle Tests`, `Interactive Project Init`, `Version Reporting`, `Community 30`, `Community 31`?**
  _High betweenness centrality (0.175) - this node is a cross-community bridge._
- **Why does `Run()` connect `Port CLI Integration Tests` to `CLI Command Orchestration`, `Service Lifecycle Manager`, `Presentation Tests`, `Lifecycle TUI Engine`?**
  _High betweenness centrality (0.098) - this node is a cross-community bridge._
- **Are the 48 inferred relationships involving `Contains()` (e.g. with `TestInitRejectsDeprecatedAppFlag()` and `TestInitRejectsInvalidGlobalRegistry()`) actually correct?**
  _`Contains()` has 48 INFERRED edges - model-reasoned connections that need verification._
- **Are the 20 inferred relationships involving `Run()` (e.g. with `New()` and `Current()`) actually correct?**
  _`Run()` has 20 INFERRED edges - model-reasoned connections that need verification._
- **Are the 29 inferred relationships involving `New()` (e.g. with `configuredRoots()` and `Run()`) actually correct?**
  _`New()` has 29 INFERRED edges - model-reasoned connections that need verification._
- **Are the 13 inferred relationships involving `runWithEnvironment()` (e.g. with `New()` and `Current()`) actually correct?**
  _`runWithEnvironment()` has 13 INFERRED edges - model-reasoned connections that need verification._
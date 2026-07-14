# Graph Report - .  (2026-07-14)

## Corpus Check
- 4 files · ~26,978 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 557 nodes · 1233 edges · 30 communities (25 shown, 5 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 202 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

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

## God Nodes (most connected - your core abstractions)
1. `grat` - 37 edges
2. `Run()` - 33 edges
3. `T` - 25 edges
4. `T` - 24 edges
5. `Renderer` - 22 edges
6. `runPortAssignLocked()` - 21 edges
7. `Config` - 19 edges
8. `LifecycleModel` - 19 edges
9. `runInitWithInput()` - 18 edges
10. `NewLifecycleModel()` - 18 edges

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

## Hyperedges (group relationships)
- **Safe Service Lifecycle** — readme_safe_local_service_manager, readme_approved_declarative_tasks, readme_process_group_isolation, readme_listener_and_http_readiness, readme_managed_process_identity, readme_cancellation_recovery [INFERRED 0.95]
- **Coordinated Port Management** — readme_cross_project_port_coordination, readme_role_based_port_registry, readme_serialized_port_allocation, readme_safe_port_reassignment [INFERRED 0.95]
- **Project Participation Guidance** — readme_contributing_guide, readme_security_policy, readme_code_of_conduct, readme_support_guide [EXTRACTED 1.00]

## Communities (30 total, 5 thin omitted)

### Community 0 - "CLI Command Orchestration"
Cohesion: 0.07
Nodes (66): assignReassignedPorts(), copyReservations(), detectServices(), ensureValidRegistry(), executeLifecycle(), fileExists(), globalScanRoots(), hasConfiguredCollision() (+58 more)

### Community 1 - "Lifecycle TUI Engine"
Cohesion: 0.10
Nodes (33): CancelFunc, Cmd, Context, Reader, Style, Writer, lifecycleRow, Model (+25 more)

### Community 2 - "Runtime Process Control"
Cohesion: 0.08
Nodes (28): Listener, Renderer, Style, Context, Duration, processState, Manager, Service (+20 more)

### Community 3 - "Configuration Schema"
Cohesion: 0.11
Nodes (32): repeatedValue, Config, containsControl(), DefaultRuntime(), InferRole(), Load(), replaceFile(), rollbackWrites() (+24 more)

### Community 4 - "Port CLI Integration Tests"
Cohesion: 0.14
Nodes (34): exitCode(), isHelp(), Run(), assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort(), loadPortFixtureConfig() (+26 more)

### Community 5 - "Terminal Presentation"
Cohesion: 0.14
Nodes (12): Renderer, ColorMode, formatProjectRows(), fprint(), fprintf(), fprintln(), pad(), ParseColorMode() (+4 more)

### Community 6 - "Presentation Tests"
Cohesion: 0.17
Nodes (32): T, NewLifecycleModel(), DividerLine(), New(), columnOf(), displayColumn(), stripANSI(), TestHelpRendersBorderlessGroupsWithGloballyAlignedDescriptions() (+24 more)

### Community 7 - "Port Registry"
Cohesion: 0.14
Nodes (27): Config, T, fakeLookup, Listener, ListenerLookup, Problem, ProjectConfig, FirstFree() (+19 more)

### Community 8 - "Service Lifecycle Tests"
Cohesion: 0.15
Nodes (26): Config, Service, T, Manager, Service, Listener, Manager, fixtureConfig() (+18 more)

### Community 9 - "Service Lifecycle Manager"
Cohesion: 0.17
Nodes (14): Client, main(), mustGetwd(), Config, Context, processState, ProgressObserver, Manager (+6 more)

### Community 10 - "Logging and Registry Locking"
Cohesion: 0.13
Nodes (14): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, File, T, T, T, Mutex, TestRegistryLockHonorsContextWhileContended() (+6 more)

### Community 11 - "README Project Guide"
Cohesion: 0.11
Nodes (20): CI Workflow, Code of Conduct, Commands, Complete Development Stacks, Configuration reference, Contributing and support, Contributing Guide, Cross-Project Port Coordination (+12 more)

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
Cohesion: 0.47
Nodes (4): T, Current(), TestCurrentPrefixesLinkerOverrideWithV(), TestCurrentReturnsSourceVersion()

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

## Knowledge Gaps
- **67 isolated node(s):** `Reader`, `ProgressObserver`, `StepKind`, `LifecycleEvent`, `ProgressStage` (+62 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Run()` connect `Port CLI Integration Tests` to `CLI Command Orchestration`, `Lifecycle TUI Engine`, `Presentation Tests`, `Service Lifecycle Manager`, `Version Reporting`?**
  _High betweenness centrality (0.175) - this node is a cross-community bridge._
- **Why does `New()` connect `Presentation Tests` to `CLI Command Orchestration`, `Port CLI Integration Tests`, `Terminal Presentation`, `Service Lifecycle Manager`, `Interactive Project Init`?**
  _High betweenness centrality (0.146) - this node is a cross-community bridge._
- **Why does `processAlive()` connect `Service Lifecycle Manager` to `Service Lifecycle Tests`, `Runtime Process Control`, `Configuration Schema`?**
  _High betweenness centrality (0.093) - this node is a cross-community bridge._
- **Are the 20 inferred relationships involving `Run()` (e.g. with `New()` and `Current()`) actually correct?**
  _`Run()` has 20 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Reader`, `ProgressObserver`, `StepKind` to the rest of the system?**
  _69 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `CLI Command Orchestration` be split into smaller, more focused modules?**
  _Cohesion score 0.06862745098039216 - nodes in this community are weakly interconnected._
- **Should `Lifecycle TUI Engine` be split into smaller, more focused modules?**
  _Cohesion score 0.09725490196078432 - nodes in this community are weakly interconnected._
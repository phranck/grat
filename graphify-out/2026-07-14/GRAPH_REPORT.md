# Graph Report - grat  (2026-07-14)

## Corpus Check
- 45 files · ~25,274 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 553 nodes · 1226 edges · 25 communities (20 shown, 5 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 198 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `28a814f0`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- [[_COMMUNITY_CLI Command Orchestration|CLI Command Orchestration]]
- [[_COMMUNITY_Lifecycle TUI Engine|Lifecycle TUI Engine]]
- [[_COMMUNITY_Runtime Process Control|Runtime Process Control]]
- [[_COMMUNITY_Configuration Schema|Configuration Schema]]
- [[_COMMUNITY_Terminal Rendering|Terminal Rendering]]
- [[_COMMUNITY_CLI Integration Tests|CLI Integration Tests]]
- [[_COMMUNITY_Documentation and Security|Documentation and Security]]
- [[_COMMUNITY_Presentation Tests|Presentation Tests]]
- [[_COMMUNITY_Port Registry|Port Registry]]
- [[_COMMUNITY_Service Lifecycle Manager|Service Lifecycle Manager]]
- [[_COMMUNITY_Runtime Lifecycle Tests|Runtime Lifecycle Tests]]
- [[_COMMUNITY_Interactive Project Init|Interactive Project Init]]
- [[_COMMUNITY_Managed State Storage|Managed State Storage]]
- [[_COMMUNITY_Logging and Locking|Logging and Locking]]
- [[_COMMUNITY_Linux Listener Inspection|Linux Listener Inspection]]
- [[_COMMUNITY_Project Root Discovery|Project Root Discovery]]
- [[_COMMUNITY_Version Reporting Tests|Version Reporting Tests]]
- [[_COMMUNITY_CLI Helper Tests|CLI Helper Tests]]
- [[_COMMUNITY_Listener Lookup Interface|Listener Lookup Interface]]
- [[_COMMUNITY_README Verification|README Verification]]
- [[_COMMUNITY_Graphify Governance|Graphify Governance]]
- [[_COMMUNITY_Darwin Listener Tests|Darwin Listener Tests]]
- [[_COMMUNITY_Process Inspection Tests|Process Inspection Tests]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
- [[_COMMUNITY_Code of Conduct|Code of Conduct]]

## God Nodes (most connected - your core abstractions)
1. `Run()` - 33 edges
2. `T` - 25 edges
3. `T` - 24 edges
4. `Renderer` - 23 edges
5. `Config` - 22 edges
6. `Manager` - 22 edges
7. `grat` - 20 edges
8. `LifecycleModel` - 19 edges
9. `runInitWithInput()` - 18 edges
10. `NewLifecycleModel()` - 18 edges

## Surprising Connections (you probably didn't know these)
- `Local Quality Gate` --semantically_similar_to--> `Verify Job`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → .github/workflows/ci.yml
- `Configuration compatibility` --semantically_similar_to--> `Approved Declarative Tasks`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → README.md
- `Cross-Platform Compatibility` --semantically_similar_to--> `Platform Verification Matrix`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → .github/workflows/ci.yml
- `main()` --calls--> `Run()`  [INFERRED]
  cmd/grat/main.go → internal/cli/cli.go
- `grat` --references--> `Code of conduct`  [EXTRACTED]
  README.md → CONTRIBUTING.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Release Artifact Delivery** — workflows_release_build_job, workflows_release_publish_job, workflows_release_checksum_generation, workflows_release_github_release_publication, readme_release_binary_installation [INFERRED 0.95]
- **Project Governance Guidance** — readme_grat, contributing_contributing_to_grat, security_security_policy, support_support, contributing_code_of_conduct [EXTRACTED 1.00]
- **Safe Local Task Execution** — readme_approved_declarative_tasks, readme_process_owned_readiness, readme_managed_process_identity, security_trusted_configuration_boundary, security_fixed_system_tool_paths [INFERRED 0.85]

## Communities (25 total, 5 thin omitted)

### Community 0 - "CLI Command Orchestration"
Cohesion: 0.06
Nodes (69): assignReassignedPorts(), copyReservations(), detectServices(), ensureValidRegistry(), executeLifecycle(), fileExists(), globalScanRoots(), hasConfiguredCollision() (+61 more)

### Community 1 - "Lifecycle TUI Engine"
Cohesion: 0.09
Nodes (34): CancelFunc, Cmd, Context, Reader, Style, Writer, lifecycleRow, Model (+26 more)

### Community 2 - "Runtime Process Control"
Cohesion: 0.07
Nodes (35): File, Listener, Renderer, Style, T, Context, Duration, processState (+27 more)

### Community 3 - "Configuration Schema"
Cohesion: 0.11
Nodes (32): Config, containsControl(), DefaultRuntime(), InferRole(), Load(), prepareWrite(), replaceFile(), rollbackWrites() (+24 more)

### Community 4 - "Terminal Rendering"
Cohesion: 0.12
Nodes (14): Renderer, Style, Writer, ColorMode, formatProjectRows(), fprint(), fprintf(), fprintln() (+6 more)

### Community 5 - "CLI Integration Tests"
Cohesion: 0.14
Nodes (34): exitCode(), isHelp(), Run(), assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort(), loadPortFixtureConfig() (+26 more)

### Community 6 - "Documentation and Security"
Cohesion: 0.06
Nodes (44): Code of conduct, Configuration compatibility, Contributing to grat, Cross-Platform Compatibility, Development setup, Focused Pull Requests, Local Quality Gate, Pull requests (+36 more)

### Community 7 - "Presentation Tests"
Cohesion: 0.17
Nodes (32): T, NewLifecycleModel(), DividerLine(), New(), columnOf(), displayColumn(), stripANSI(), TestHelpRendersBorderlessGroupsWithGloballyAlignedDescriptions() (+24 more)

### Community 8 - "Port Registry"
Cohesion: 0.13
Nodes (28): Config, Listener, T, fakeLookup, Listener, ListenerLookup, Problem, ProjectConfig (+20 more)

### Community 9 - "Service Lifecycle Manager"
Cohesion: 0.16
Nodes (14): Client, main(), mustGetwd(), Config, Context, processState, ProgressObserver, Manager (+6 more)

### Community 10 - "Runtime Lifecycle Tests"
Cohesion: 0.16
Nodes (25): Config, Listener, Manager, Service, T, Manager, Service, fixtureConfig() (+17 more)

### Community 11 - "Interactive Project Init"
Cohesion: 0.28
Nodes (15): collectInitInterview(), parseServiceDefinition(), promptDefault(), promptRequired(), promptServiceName(), TestCollectInitInterviewAcceptsEditsAndAdditionalServices(), TestCollectInitInterviewAcceptsSuppliedProjectName(), TestCollectInitInterviewRequiresAtLeastOneService() (+7 more)

### Community 12 - "Managed State Storage"
Cohesion: 0.28
Nodes (4): Manager, loadedState, processState, Time

### Community 13 - "Logging and Locking"
Cohesion: 0.25
Nodes (7): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, T, T, TestRegistryLockHonorsContextWhileContended(), TestRegistryLockReleasesAfterCallbackPanic(), TestRegistryLockSerializesCallbacks()

### Community 14 - "Linux Listener Inspection"
Cohesion: 0.25
Nodes (8): Listener, T, systemListener(), linuxListeningSocketInodes(), linuxSocketOwnerPIDs(), socketInode(), TestLinuxListenerParsingFindsOnlyTCPListenInodesForTheRequestedPort(), TestLinuxSocketOwnersFindsPIDsFromProcFileDescriptors()

### Community 15 - "Project Root Discovery"
Cohesion: 0.47
Nodes (4): T, FindRoot(), TestFindRootReturnsNotFoundOutsideProject(), TestFindRootUsesNearestConfig()

### Community 16 - "Version Reporting Tests"
Cohesion: 0.47
Nodes (4): T, Current(), TestCurrentPrefixesLinkerOverrideWithV(), TestCurrentReturnsSourceVersion()

### Community 17 - "CLI Helper Tests"
Cohesion: 0.67
Nodes (3): TestListenerOwnerLabelHandlesUnknownPID(), TestLogFollowUsesTrustedExecutable(), T

### Community 19 - "README Verification"
Cohesion: 0.83
Nodes (3): require(), require_in(), check-readme.sh script

### Community 20 - "Graphify Governance"
Cohesion: 0.67
Nodes (3): Graphify Before Push, Project Agent Rules, Versioned Graphify Artifacts

## Knowledge Gaps
- **63 isolated node(s):** `Reader`, `ProgressObserver`, `StepKind`, `LifecycleEvent`, `ProgressStage` (+58 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Run()` connect `CLI Integration Tests` to `CLI Command Orchestration`, `Lifecycle TUI Engine`, `Presentation Tests`, `Service Lifecycle Manager`, `Version Reporting Tests`?**
  _High betweenness centrality (0.194) - this node is a cross-community bridge._
- **Why does `New()` connect `Presentation Tests` to `CLI Command Orchestration`, `Terminal Rendering`, `CLI Integration Tests`, `Service Lifecycle Manager`, `Interactive Project Init`?**
  _High betweenness centrality (0.160) - this node is a cross-community bridge._
- **Why does `processAlive()` connect `Service Lifecycle Manager` to `Runtime Process Control`, `Runtime Lifecycle Tests`?**
  _High betweenness centrality (0.098) - this node is a cross-community bridge._
- **Are the 20 inferred relationships involving `Run()` (e.g. with `New()` and `Current()`) actually correct?**
  _`Run()` has 20 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Reader`, `ProgressObserver`, `StepKind` to the rest of the system?**
  _65 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `CLI Command Orchestration` be split into smaller, more focused modules?**
  _Cohesion score 0.0625694187338023 - nodes in this community are weakly interconnected._
- **Should `Lifecycle TUI Engine` be split into smaller, more focused modules?**
  _Cohesion score 0.09071117561683599 - nodes in this community are weakly interconnected._
# Graph Report - .  (2026-07-14)

## Corpus Check
- 5 files · ~37,800 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 780 nodes · 1847 edges · 39 communities (32 shown, 7 thin omitted)
- Extraction: 83% EXTRACTED · 17% INFERRED · 0% AMBIGUOUS · INFERRED: 313 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_CLI Service Lifecycle|CLI Service Lifecycle]]
- [[_COMMUNITY_Runtime Process Management|Runtime Process Management]]
- [[_COMMUNITY_Terminal Lifecycle Display|Terminal Lifecycle Display]]
- [[_COMMUNITY_CLI Presentation Framework|CLI Presentation Framework]]
- [[_COMMUNITY_Release Update Client|Release Update Client]]
- [[_COMMUNITY_CLI Command Tests|CLI Command Tests]]
- [[_COMMUNITY_Presentation Rendering Tests|Presentation Rendering Tests]]
- [[_COMMUNITY_Project Configuration Validation|Project Configuration Validation]]
- [[_COMMUNITY_Maintenance Test Helpers|Maintenance Test Helpers]]
- [[_COMMUNITY_README Product Guide|README Product Guide]]
- [[_COMMUNITY_Port Registry and Scanning|Port Registry and Scanning]]
- [[_COMMUNITY_Directory Command Tests|Directory Command Tests]]
- [[_COMMUNITY_Runtime Log Tests|Runtime Log Tests]]
- [[_COMMUNITY_CLI Bootstrap|CLI Bootstrap]]
- [[_COMMUNITY_Update Installation|Update Installation]]
- [[_COMMUNITY_Interactive Project Setup|Interactive Project Setup]]
- [[_COMMUNITY_Settings Safety Tests|Settings Safety Tests]]
- [[_COMMUNITY_Settings Persistence|Settings Persistence]]
- [[_COMMUNITY_Contributing and CI|Contributing and CI]]
- [[_COMMUNITY_Runtime State Storage|Runtime State Storage]]
- [[_COMMUNITY_Release Distribution|Release Distribution]]
- [[_COMMUNITY_Security Documentation|Security Documentation]]
- [[_COMMUNITY_Linux Listener Inspection|Linux Listener Inspection]]
- [[_COMMUNITY_Safe Service Design|Safe Service Design]]
- [[_COMMUNITY_Framework Configuration Examples|Framework Configuration Examples]]
- [[_COMMUNITY_Project Root Discovery|Project Root Discovery]]
- [[_COMMUNITY_Bottle Packaging Tests|Bottle Packaging Tests]]
- [[_COMMUNITY_Bottle Packaging Script|Bottle Packaging Script]]
- [[_COMMUNITY_Published Bottle Verification|Published Bottle Verification]]
- [[_COMMUNITY_Graphify Project Rules|Graphify Project Rules]]
- [[_COMMUNITY_README Contract Script|README Contract Script]]
- [[_COMMUNITY_Darwin Listener Tests|Darwin Listener Tests]]
- [[_COMMUNITY_Readiness Inspection Tests|Readiness Inspection Tests]]
- [[_COMMUNITY_Safe Self-Uninstall|Safe Self-Uninstall]]
- [[_COMMUNITY_System Listener Lookup|System Listener Lookup]]
- [[_COMMUNITY_Code of Conduct|Code of Conduct]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
- [[_COMMUNITY_Bottle Release Verification Tests|Bottle Release Verification Tests]]
- [[_COMMUNITY_Help Contract Tests|Help Contract Tests]]

## God Nodes (most connected - your core abstractions)
1. `Contains()` - 52 edges
2. `grat` - 51 edges
3. `Run()` - 35 edges
4. `New()` - 34 edges
5. `runWithEnvironment()` - 32 edges
6. `T` - 25 edges
7. `T` - 25 edges
8. `Renderer` - 22 edges
9. `runPortAssignLocked()` - 21 edges
10. `Service` - 20 edges

## Surprising Connections (you probably didn't know these)
- `Configuration compatibility` --semantically_similar_to--> `Approved Declarative Tasks`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → README.md
- `Local Quality Gate` --semantically_similar_to--> `Verify Job`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → .github/workflows/ci.yml
- `grat` --references--> `WithRegistryLock()`  [EXTRACTED]
  README.md → internal/ports/lock.go
- `Contributing to grat` --references--> `README Contract Check`  [EXTRACTED]
  CONTRIBUTING.md → scripts/check-readme.sh
- `grat` --references--> `Service.Uninstall`  [EXTRACTED]
  README.md → internal/maintenance/uninstall.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Cross-Platform Binary Delivery** — workflows_release_cross_platform_build_matrix, workflows_release_versioned_binary_build, workflows_release_binary_artifacts, workflows_release_artifact_download [EXTRACTED 1.00]
- **Verified Homebrew Release** — workflows_release_homebrew_bottle_packaging, workflows_release_release_checksums, workflows_release_github_release_publication, workflows_release_published_bottle_verification [EXTRACTED 1.00]

## Communities (39 total, 7 thin omitted)

### Community 0 - "CLI Service Lifecycle"
Cohesion: 0.06
Nodes (73): assignReassignedPorts(), configuredRoots(), copyReservations(), defaultEnvironment(), detectServices(), ensureValidRegistry(), executeLifecycle(), fileExists() (+65 more)

### Community 1 - "Runtime Process Management"
Cohesion: 0.07
Nodes (52): repeatedValue, Listener, Config, Listener, Manager, Service, T, Context (+44 more)

### Community 2 - "Terminal Lifecycle Display"
Cohesion: 0.10
Nodes (32): CancelFunc, Cmd, Context, Reader, Style, Writer, lifecycleRow, Model (+24 more)

### Community 3 - "CLI Presentation Framework"
Cohesion: 0.10
Nodes (17): Renderer, Style, Renderer, ColorMode, CommandGroup, helpUsageWidth(), formatProjectRows(), fprint() (+9 more)

### Community 4 - "Release Update Client"
Cohesion: 0.09
Nodes (24): asset, installation, Client, Context, Service, Client, Context, Service (+16 more)

### Community 5 - "CLI Command Tests"
Cohesion: 0.15
Nodes (36): exitCode(), Run(), assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort(), loadPortFixtureConfig(), runWithConfiguredRoots() (+28 more)

### Community 6 - "Presentation Rendering Tests"
Cohesion: 0.17
Nodes (35): T, NewLifecycleModel(), DividerLine(), New(), columnOf(), displayColumn(), stripANSI(), TestHelpRendersBorderlessGroupsWithGloballyAlignedDescriptions() (+27 more)

### Community 7 - "Project Configuration Validation"
Cohesion: 0.12
Nodes (31): Config, containsControl(), DefaultRuntime(), InferRole(), Load(), replaceFile(), rollbackWrites(), safeServiceName() (+23 more)

### Community 8 - "Maintenance Test Helpers"
Cohesion: 0.14
Nodes (31): Service, Store, T, Context, Service, T, commandResponse, fakeCommands (+23 more)

### Community 9 - "README Product Guide"
Cohesion: 0.07
Nodes (31): Bounded Local Logs, CI Workflow, Code of Conduct, Command contract, Commands, Complete Development Stacks, Configuration reference, Contents (+23 more)

### Community 10 - "Port Registry and Scanning"
Cohesion: 0.14
Nodes (27): Config, T, fakeLookup, Listener, ListenerLookup, Problem, ProjectConfig, FirstFree() (+19 more)

### Community 11 - "Directory Command Tests"
Cohesion: 0.20
Nodes (26): isHelp(), runWithEnvironment(), CLI Command Routing Tests, canonicalCLITestPath(), environmentForTest(), newCLITestStore(), sameStringSlices(), TestDirectoriesAddDoesNotPromptForInitialSetup() (+18 more)

### Community 12 - "Runtime Log Tests"
Cohesion: 0.11
Nodes (19): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, File, T, Context, T, T, Mutex (+11 more)

### Community 13 - "CLI Bootstrap"
Cohesion: 0.17
Nodes (13): main(), mustGetwd(), Client, Config, Context, processState, ProgressObserver, Manager (+5 more)

### Community 14 - "Update Installation"
Cohesion: 0.19
Nodes (14): Context, Service, Reader, Result, Store, Writer, installation, installationKind (+6 more)

### Community 15 - "Interactive Project Setup"
Cohesion: 0.26
Nodes (17): initServiceSuggestions(), runInitWithInput(), collectInitInterview(), parseServiceDefinition(), promptDefault(), promptRequired(), promptServiceName(), TestCollectInitInterviewAcceptsEditsAndAdditionalServices() (+9 more)

### Community 16 - "Settings Safety Tests"
Cohesion: 0.31
Nodes (17): Store, T, canonicalPath(), equalStrings(), newTestStore(), TestContainsAcceptsRegularFileBelowRoot(), TestContainsRejectsPathsOutsideRootAndThroughSymlinks(), TestStoreAddCanonicalizesAndDeduplicatesDirectories() (+9 more)

### Community 17 - "Settings Persistence"
Cohesion: 0.30
Nodes (4): Settings, canonicalExistingDirectory(), canonicalExistingPath(), Store

### Community 18 - "Contributing and CI"
Cohesion: 0.14
Nodes (14): Code of conduct, Contributing to grat, Cross-Platform Compatibility, Development setup, Focused Pull Requests, Local Quality Gate, Pull requests, README Contract Check (+6 more)

### Community 19 - "Runtime State Storage"
Cohesion: 0.28
Nodes (4): Manager, loadedState, processState, Time

### Community 20 - "Release Distribution"
Cohesion: 0.20
Nodes (15): GitHub Releases, Release Binary Installation, Homebrew Bottle Builder, Homebrew Bottle Tests, Artifact Download, Binary Artifacts, Release Build Job, Checksum Generation (+7 more)

### Community 21 - "Security Documentation"
Cohesion: 0.21
Nodes (10): Fixed System Tool Paths, Private Vulnerability Reporting, Reporting a vulnerability, Security policy, Documented Shell Semantics, Supported versions, Trusted Configuration Boundary, Diagnostic Support Request (+2 more)

### Community 22 - "Linux Listener Inspection"
Cohesion: 0.25
Nodes (8): Listener, T, systemListener(), linuxListeningSocketInodes(), linuxSocketOwnerPIDs(), socketInode(), TestLinuxListenerParsingFindsOnlyTCPListenInodesForTheRequestedPort(), TestLinuxSocketOwnersFindsPIDsFromProcFileDescriptors()

### Community 23 - "Safe Service Design"
Cohesion: 0.25
Nodes (8): Configuration compatibility, Approved Declarative Tasks, Cancellation Recovery, Declarative Non-Executable Configuration, Managed Process Identity, Process Group Isolation, Process-Owned Service Readiness, Safe Local Service Manager

### Community 24 - "Framework Configuration Examples"
Cohesion: 0.29
Nodes (7): Go HTTP API, Laravel, Project examples, Python with FastAPI, React, Laravel, and a queue worker, React with Vite, Swift with Vapor

### Community 25 - "Project Root Discovery"
Cohesion: 0.47
Nodes (4): T, FindRoot(), TestFindRootReturnsNotFoundOutsideProject(), TestFindRootUsesNearestConfig()

### Community 26 - "Bottle Packaging Tests"
Cohesion: 0.60
Nodes (5): assert_archive_contains(), assert_binary(), assert_file(), assert_mode(), test-homebrew-bottles.sh script

### Community 27 - "Bottle Packaging Script"
Cohesion: 0.70
Nodes (4): package(), usage(), write_formula(), build-homebrew-bottles.sh script

### Community 28 - "Published Bottle Verification"
Cohesion: 0.60
Nodes (3): usage(), verify_bottle(), verify-homebrew-bottles.sh script

### Community 29 - "Graphify Project Rules"
Cohesion: 0.67
Nodes (3): Graphify Before Push, Project Agent Rules, Versioned Graphify Artifacts

### Community 30 - "README Contract Script"
Cohesion: 0.83
Nodes (3): require(), require_in(), check-readme.sh script

### Community 33 - "Safe Self-Uninstall"
Cohesion: 0.67
Nodes (3): Safe Uninstall Tests, Service.Uninstall, Settings Store

## Knowledge Gaps
- **107 isolated node(s):** `Store`, `StepKind`, `LifecycleEvent`, `ProgressStage`, `LifecycleStage` (+102 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `New()` connect `Presentation Rendering Tests` to `CLI Service Lifecycle`, `CLI Presentation Framework`, `Release Update Client`, `CLI Command Tests`, `Maintenance Test Helpers`, `Directory Command Tests`, `CLI Bootstrap`, `Update Installation`, `Interactive Project Setup`, `Settings Safety Tests`, `Settings Persistence`?**
  _High betweenness centrality (0.210) - this node is a cross-community bridge._
- **Why does `Contains()` connect `Presentation Rendering Tests` to `Runtime Process Management`, `CLI Command Tests`, `Project Configuration Validation`, `Maintenance Test Helpers`, `Directory Command Tests`, `Update Installation`, `Interactive Project Setup`, `Settings Safety Tests`, `Settings Persistence`?**
  _High betweenness centrality (0.181) - this node is a cross-community bridge._
- **Why does `grat` connect `README Product Guide` to `Safe Self-Uninstall`, `Release Update Client`, `Runtime Log Tests`, `Contributing and CI`, `Release Distribution`, `Security Documentation`, `Safe Service Design`, `Framework Configuration Examples`?**
  _High betweenness centrality (0.176) - this node is a cross-community bridge._
- **Are the 48 inferred relationships involving `Contains()` (e.g. with `TestInitRejectsDeprecatedAppFlag()` and `TestInitRejectsInvalidGlobalRegistry()`) actually correct?**
  _`Contains()` has 48 INFERRED edges - model-reasoned connections that need verification._
- **Are the 20 inferred relationships involving `Run()` (e.g. with `New()` and `Current()`) actually correct?**
  _`Run()` has 20 INFERRED edges - model-reasoned connections that need verification._
- **Are the 29 inferred relationships involving `New()` (e.g. with `configuredRoots()` and `Run()`) actually correct?**
  _`New()` has 29 INFERRED edges - model-reasoned connections that need verification._
- **Are the 12 inferred relationships involving `runWithEnvironment()` (e.g. with `New()` and `runWithConfiguredRoots()`) actually correct?**
  _`runWithEnvironment()` has 12 INFERRED edges - model-reasoned connections that need verification._
# Graph Report - .  (2026-07-14)

## Corpus Check
- 23 files · ~37,524 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 767 nodes · 1827 edges · 37 communities (31 shown, 6 thin omitted)
- Extraction: 83% EXTRACTED · 17% INFERRED · 0% AMBIGUOUS · INFERRED: 313 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_CLI Service Lifecycle|CLI Service Lifecycle]]
- [[_COMMUNITY_Runtime Process Management|Runtime Process Management]]
- [[_COMMUNITY_Terminal Lifecycle Display|Terminal Lifecycle Display]]
- [[_COMMUNITY_Maintenance Test Helpers|Maintenance Test Helpers]]
- [[_COMMUNITY_CLI Presentation Framework|CLI Presentation Framework]]
- [[_COMMUNITY_CLI Command Tests|CLI Command Tests]]
- [[_COMMUNITY_Update Release Verification|Update Release Verification]]
- [[_COMMUNITY_Presentation Rendering Tests|Presentation Rendering Tests]]
- [[_COMMUNITY_Project Configuration Validation|Project Configuration Validation]]
- [[_COMMUNITY_Port Registry and Scanning|Port Registry and Scanning]]
- [[_COMMUNITY_README Product Guide|README Product Guide]]
- [[_COMMUNITY_Directory Command Tests|Directory Command Tests]]
- [[_COMMUNITY_Safe Self-Uninstall|Safe Self-Uninstall]]
- [[_COMMUNITY_Settings Safety Tests|Settings Safety Tests]]
- [[_COMMUNITY_Settings Persistence|Settings Persistence]]
- [[_COMMUNITY_Interactive Project Setup|Interactive Project Setup]]
- [[_COMMUNITY_Registry Locking|Registry Locking]]
- [[_COMMUNITY_Security Documentation|Security Documentation]]
- [[_COMMUNITY_Runtime State Storage|Runtime State Storage]]
- [[_COMMUNITY_Contributing and CI|Contributing and CI]]
- [[_COMMUNITY_Release Distribution|Release Distribution]]
- [[_COMMUNITY_Linux Listener Inspection|Linux Listener Inspection]]
- [[_COMMUNITY_Runtime Progress Reporting|Runtime Progress Reporting]]
- [[_COMMUNITY_Readiness and Recovery|Readiness and Recovery]]
- [[_COMMUNITY_Framework Configuration Examples|Framework Configuration Examples]]
- [[_COMMUNITY_Project Root Discovery|Project Root Discovery]]
- [[_COMMUNITY_Bottle Packaging Tests|Bottle Packaging Tests]]
- [[_COMMUNITY_Bottle Packaging Script|Bottle Packaging Script]]
- [[_COMMUNITY_Graphify Project Rules|Graphify Project Rules]]
- [[_COMMUNITY_CLI Trusted Utilities|CLI Trusted Utilities]]
- [[_COMMUNITY_README Contract Script|README Contract Script]]
- [[_COMMUNITY_Darwin Listener Tests|Darwin Listener Tests]]
- [[_COMMUNITY_Readiness Inspection Tests|Readiness Inspection Tests]]
- [[_COMMUNITY_System Listener Lookup|System Listener Lookup]]
- [[_COMMUNITY_Code of Conduct|Code of Conduct]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
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
- **Safe Self-Maintenance** — cli_cli_runwithenvironment, maintenance_update_update, maintenance_update_replaceverifiedbinary, maintenance_uninstall_uninstall, settings_settings_store [INFERRED 0.95]
- **Serialized Port Configuration** — cli_cli_runinitwithinput, cli_cli_runportassign, cli_cli_runportreassign, ports_lock_withregistrylock [EXTRACTED 1.00]
- **Release Distribution Pipeline** — workflows_release_release_workflow, scripts_build_release_release_builder, scripts_build_homebrew_bottles_homebrew_bottle_builder, scripts_test_homebrew_bottles_homebrew_bottle_tests [INFERRED 0.95]

## Communities (37 total, 6 thin omitted)

### Community 0 - "CLI Service Lifecycle"
Cohesion: 0.06
Nodes (84): assignReassignedPorts(), configuredRoots(), copyReservations(), defaultEnvironment(), detectServices(), ensureValidRegistry(), executeLifecycle(), fileExists() (+76 more)

### Community 1 - "Runtime Process Management"
Cohesion: 0.06
Nodes (59): main(), mustGetwd(), Listener, Client, Config, Context, processState, ProgressObserver (+51 more)

### Community 2 - "Terminal Lifecycle Display"
Cohesion: 0.10
Nodes (33): CancelFunc, Cmd, Context, Reader, Style, Writer, lifecycleRow, Model (+25 more)

### Community 3 - "Maintenance Test Helpers"
Cohesion: 0.09
Nodes (39): repeatedValue, File, Service, Store, T, Context, Service, T (+31 more)

### Community 4 - "CLI Presentation Framework"
Cohesion: 0.10
Nodes (16): Renderer, Style, Renderer, ColorMode, CommandGroup, helpUsageWidth(), formatProjectRows(), fprint() (+8 more)

### Community 5 - "CLI Command Tests"
Cohesion: 0.15
Nodes (36): exitCode(), Run(), assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort(), loadPortFixtureConfig(), runWithConfiguredRoots() (+28 more)

### Community 6 - "Update Release Verification"
Cohesion: 0.12
Nodes (18): asset, installation, Client, Context, Service, Client, Context, Service (+10 more)

### Community 7 - "Presentation Rendering Tests"
Cohesion: 0.17
Nodes (35): T, NewLifecycleModel(), DividerLine(), New(), columnOf(), displayColumn(), stripANSI(), TestHelpRendersBorderlessGroupsWithGloballyAlignedDescriptions() (+27 more)

### Community 8 - "Project Configuration Validation"
Cohesion: 0.12
Nodes (31): Config, containsControl(), DefaultRuntime(), InferRole(), Load(), replaceFile(), rollbackWrites(), safeServiceName() (+23 more)

### Community 9 - "Port Registry and Scanning"
Cohesion: 0.14
Nodes (27): Config, T, fakeLookup, Listener, ListenerLookup, Problem, ProjectConfig, FirstFree() (+19 more)

### Community 10 - "README Product Guide"
Cohesion: 0.08
Nodes (28): CI Workflow, Code of Conduct, Command contract, Commands, Complete Development Stacks, Configuration reference, Contents, Contributing and support (+20 more)

### Community 11 - "Directory Command Tests"
Cohesion: 0.21
Nodes (23): canonicalCLITestPath(), environmentForTest(), newCLITestStore(), sameStringSlices(), TestDirectoriesAddDoesNotPromptForInitialSetup(), TestDirectoriesCommandsPersistAndListConfiguredRoots(), TestFirstUseAcceptsExistingSitesDefault(), TestFirstUseFallsBackToWorkingDirectory() (+15 more)

### Community 12 - "Safe Self-Uninstall"
Cohesion: 0.19
Nodes (14): Context, Service, Reader, Result, Store, Writer, installation, installationKind (+6 more)

### Community 13 - "Settings Safety Tests"
Cohesion: 0.31
Nodes (17): Store, T, canonicalPath(), equalStrings(), newTestStore(), TestContainsAcceptsRegularFileBelowRoot(), TestContainsRejectsPathsOutsideRootAndThroughSymlinks(), TestStoreAddCanonicalizesAndDeduplicatesDirectories() (+9 more)

### Community 14 - "Settings Persistence"
Cohesion: 0.30
Nodes (4): Settings, canonicalExistingDirectory(), canonicalExistingPath(), Store

### Community 15 - "Interactive Project Setup"
Cohesion: 0.28
Nodes (15): collectInitInterview(), parseServiceDefinition(), promptDefault(), promptRequired(), promptServiceName(), TestCollectInitInterviewAcceptsEditsAndAdditionalServices(), TestCollectInitInterviewAcceptsSuppliedProjectName(), TestCollectInitInterviewRequiresAtLeastOneService() (+7 more)

### Community 16 - "Registry Locking"
Cohesion: 0.19
Nodes (12): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, T, Context, T, Registry Lock Tests, TestRegistryLockHonorsContextWhileContended(), TestRegistryLockReleasesAfterCallbackPanic() (+4 more)

### Community 17 - "Security Documentation"
Cohesion: 0.16
Nodes (13): Configuration compatibility, Approved Declarative Tasks, Declarative Non-Executable Configuration, Fixed System Tool Paths, Private Vulnerability Reporting, Reporting a vulnerability, Security policy, Documented Shell Semantics (+5 more)

### Community 18 - "Runtime State Storage"
Cohesion: 0.28
Nodes (4): Manager, loadedState, processState, Time

### Community 19 - "Contributing and CI"
Cohesion: 0.17
Nodes (11): Code of conduct, Contributing to grat, Development setup, Focused Pull Requests, Local Quality Gate, Pull requests, README Contract Check, CI Workflow (+3 more)

### Community 20 - "Release Distribution"
Cohesion: 0.23
Nodes (12): Cross-Platform Compatibility, GitHub Releases, Release Binary Installation, Homebrew Bottle Builder, Homebrew Bottle Tests, Platform Verification Matrix, Release Build Job, Checksum Generation (+4 more)

### Community 21 - "Linux Listener Inspection"
Cohesion: 0.25
Nodes (8): Listener, T, systemListener(), linuxListeningSocketInodes(), linuxSocketOwnerPIDs(), socketInode(), TestLinuxListenerParsingFindsOnlyTCPListenInodesForTheRequestedPort(), TestLinuxSocketOwnersFindsPIDsFromProcFileDescriptors()

### Community 22 - "Runtime Progress Reporting"
Cohesion: 0.36
Nodes (5): Manager, Service, ProgressEvent, ProgressObserver, ProgressStage

### Community 23 - "Readiness and Recovery"
Cohesion: 0.25
Nodes (8): Bounded Local Logs, Cancellation Recovery, Listener Ownership and HTTP Readiness, Managed Process Identity, Process Group Isolation, Process-Owned Service Readiness, Safe Local Service Manager, Service Recovery Workflow

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

### Community 28 - "Graphify Project Rules"
Cohesion: 0.67
Nodes (3): Graphify Before Push, Project Agent Rules, Versioned Graphify Artifacts

### Community 29 - "CLI Trusted Utilities"
Cohesion: 0.67
Nodes (3): TestListenerOwnerLabelHandlesUnknownPID(), TestLogFollowUsesTrustedExecutable(), T

### Community 30 - "README Contract Script"
Cohesion: 0.83
Nodes (3): require(), require_in(), check-readme.sh script

## Knowledge Gaps
- **105 isolated node(s):** `Store`, `StepKind`, `LifecycleEvent`, `ProgressStage`, `LifecycleStage` (+100 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `New()` connect `Presentation Rendering Tests` to `CLI Service Lifecycle`, `Runtime Process Management`, `Maintenance Test Helpers`, `CLI Presentation Framework`, `CLI Command Tests`, `Update Release Verification`, `Safe Self-Uninstall`, `Settings Safety Tests`, `Settings Persistence`, `Interactive Project Setup`?**
  _High betweenness centrality (0.214) - this node is a cross-community bridge._
- **Why does `grat` connect `README Product Guide` to `CLI Service Lifecycle`, `Registry Locking`, `Security Documentation`, `Contributing and CI`, `Release Distribution`, `Readiness and Recovery`, `Framework Configuration Examples`?**
  _High betweenness centrality (0.195) - this node is a cross-community bridge._
- **Why does `Contains()` connect `Presentation Rendering Tests` to `Runtime Process Management`, `Maintenance Test Helpers`, `CLI Command Tests`, `Project Configuration Validation`, `Directory Command Tests`, `Safe Self-Uninstall`, `Settings Safety Tests`, `Settings Persistence`, `Interactive Project Setup`?**
  _High betweenness centrality (0.189) - this node is a cross-community bridge._
- **Are the 48 inferred relationships involving `Contains()` (e.g. with `TestInitRejectsDeprecatedAppFlag()` and `TestInitRejectsInvalidGlobalRegistry()`) actually correct?**
  _`Contains()` has 48 INFERRED edges - model-reasoned connections that need verification._
- **Are the 20 inferred relationships involving `Run()` (e.g. with `New()` and `Current()`) actually correct?**
  _`Run()` has 20 INFERRED edges - model-reasoned connections that need verification._
- **Are the 29 inferred relationships involving `New()` (e.g. with `configuredRoots()` and `Run()`) actually correct?**
  _`New()` has 29 INFERRED edges - model-reasoned connections that need verification._
- **Are the 12 inferred relationships involving `runWithEnvironment()` (e.g. with `New()` and `runWithConfiguredRoots()`) actually correct?**
  _`runWithEnvironment()` has 12 INFERRED edges - model-reasoned connections that need verification._
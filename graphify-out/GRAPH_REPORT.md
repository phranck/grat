# Graph Report - .  (2026-07-14)

## Corpus Check
- 2 files · ~37,937 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 793 nodes · 1872 edges · 41 communities (33 shown, 8 thin omitted)
- Extraction: 83% EXTRACTED · 17% INFERRED · 0% AMBIGUOUS · INFERRED: 321 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_CLI Orchestration|CLI Orchestration]]
- [[_COMMUNITY_Runtime Terminal Models|Runtime Terminal Models]]
- [[_COMMUNITY_CLI Test Utilities|CLI Test Utilities]]
- [[_COMMUNITY_Update Maintenance|Update Maintenance]]
- [[_COMMUNITY_Terminal Presentation|Terminal Presentation]]
- [[_COMMUNITY_Uninstall Maintenance|Uninstall Maintenance]]
- [[_COMMUNITY_Project Documentation|Project Documentation]]
- [[_COMMUNITY_CLI Command Tests|CLI Command Tests]]
- [[_COMMUNITY_Presentation Tests|Presentation Tests]]
- [[_COMMUNITY_Configuration Validation|Configuration Validation]]
- [[_COMMUNITY_Port Registry|Port Registry]]
- [[_COMMUNITY_Directory CLI Tests|Directory CLI Tests]]
- [[_COMMUNITY_Runtime Process Manager|Runtime Process Manager]]
- [[_COMMUNITY_Process Shutdown|Process Shutdown]]
- [[_COMMUNITY_Runtime Manager Tests|Runtime Manager Tests]]
- [[_COMMUNITY_Settings Tests|Settings Tests]]
- [[_COMMUNITY_Interactive Setup|Interactive Setup]]
- [[_COMMUNITY_Logging And Locks|Logging And Locks]]
- [[_COMMUNITY_Contributing And CI|Contributing And CI]]
- [[_COMMUNITY_Runtime State Storage|Runtime State Storage]]
- [[_COMMUNITY_Release Distribution|Release Distribution]]
- [[_COMMUNITY_Product Documentation|Product Documentation]]
- [[_COMMUNITY_Security Documentation|Security Documentation]]
- [[_COMMUNITY_Linux Listener Inspection|Linux Listener Inspection]]
- [[_COMMUNITY_Service Safety Architecture|Service Safety Architecture]]
- [[_COMMUNITY_Runtime Progress Reporting|Runtime Progress Reporting]]
- [[_COMMUNITY_Runtime Readiness|Runtime Readiness]]
- [[_COMMUNITY_Project Root Discovery|Project Root Discovery]]
- [[_COMMUNITY_Bottle Packaging Tests|Bottle Packaging Tests]]
- [[_COMMUNITY_Bottle Packaging Script|Bottle Packaging Script]]
- [[_COMMUNITY_Bottle Verification|Bottle Verification]]
- [[_COMMUNITY_Project Agent Rules|Project Agent Rules]]
- [[_COMMUNITY_README Contract Script|README Contract Script]]
- [[_COMMUNITY_CLI Entrypoint|CLI Entrypoint]]
- [[_COMMUNITY_Darwin Listener Tests|Darwin Listener Tests]]
- [[_COMMUNITY_Readiness Tests|Readiness Tests]]
- [[_COMMUNITY_Listener Lookup|Listener Lookup]]
- [[_COMMUNITY_Code Of Conduct|Code Of Conduct]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
- [[_COMMUNITY_Bottle Verification Tests|Bottle Verification Tests]]
- [[_COMMUNITY_Help Tests|Help Tests]]

## God Nodes (most connected - your core abstractions)
1. `grat` - 61 edges
2. `Contains()` - 52 edges
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
- `grat` --references--> `Service.Update`  [EXTRACTED]
  README.md → internal/maintenance/update.go

## Import Cycles
- None detected.

## Communities (41 total, 8 thin omitted)

### Community 0 - "CLI Orchestration"
Cohesion: 0.06
Nodes (75): assignReassignedPorts(), configuredRoots(), copyReservations(), defaultEnvironment(), detectServices(), ensureValidRegistry(), executeLifecycle(), fileExists() (+67 more)

### Community 1 - "Runtime Terminal Models"
Cohesion: 0.10
Nodes (33): CancelFunc, Cmd, Context, Reader, Style, Writer, lifecycleRow, Model (+25 more)

### Community 2 - "CLI Test Utilities"
Cohesion: 0.09
Nodes (39): repeatedValue, File, Service, Store, T, Context, Service, T (+31 more)

### Community 3 - "Update Maintenance"
Cohesion: 0.09
Nodes (24): asset, installation, Client, Context, Service, Client, Context, Service (+16 more)

### Community 4 - "Terminal Presentation"
Cohesion: 0.10
Nodes (18): Renderer, Style, ColorMode, CommandGroup, helpUsageWidth(), formatProjectRows(), fprint(), fprintf() (+10 more)

### Community 5 - "Uninstall Maintenance"
Cohesion: 0.12
Nodes (18): Context, Service, Reader, Result, Store, Writer, installation, installationKind (+10 more)

### Community 6 - "Project Documentation"
Cohesion: 0.06
Nodes (39): Safe Uninstall Tests, Service.Uninstall, Bounded Local Logs, Bounded Service Logs, CI Workflow, Code of Conduct, Commands, Complete Development Stacks (+31 more)

### Community 7 - "CLI Command Tests"
Cohesion: 0.15
Nodes (36): exitCode(), Run(), assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort(), loadPortFixtureConfig(), runWithConfiguredRoots() (+28 more)

### Community 8 - "Presentation Tests"
Cohesion: 0.17
Nodes (35): NewLifecycleModel(), DividerLine(), New(), columnOf(), displayColumn(), stripANSI(), TestHelpRendersBorderlessGroupsWithGloballyAlignedDescriptions(), TestLifecycleModelAlignsHeadersWithDataColumns() (+27 more)

### Community 9 - "Configuration Validation"
Cohesion: 0.12
Nodes (31): Config, containsControl(), DefaultRuntime(), InferRole(), Load(), replaceFile(), rollbackWrites(), safeServiceName() (+23 more)

### Community 10 - "Port Registry"
Cohesion: 0.14
Nodes (27): Config, T, fakeLookup, Listener, ListenerLookup, Problem, ProjectConfig, FirstFree() (+19 more)

### Community 11 - "Directory CLI Tests"
Cohesion: 0.20
Nodes (26): isHelp(), runWithEnvironment(), CLI Command Routing Tests, canonicalCLITestPath(), environmentForTest(), newCLITestStore(), sameStringSlices(), TestDirectoriesAddDoesNotPromptForInitialSetup() (+18 more)

### Community 12 - "Runtime Process Manager"
Cohesion: 0.20
Nodes (12): Client, Config, Context, processState, ProgressObserver, Manager, State, Service (+4 more)

### Community 13 - "Process Shutdown"
Cohesion: 0.15
Nodes (17): Listener, Context, Duration, processState, Manager, Service, loadedState, systemListener() (+9 more)

### Community 14 - "Runtime Manager Tests"
Cohesion: 0.25
Nodes (20): Config, Listener, Manager, Service, T, fixtureConfig(), fixtureService(), freeTCPPort() (+12 more)

### Community 15 - "Settings Tests"
Cohesion: 0.31
Nodes (17): Store, T, canonicalPath(), equalStrings(), newTestStore(), TestContainsAcceptsRegularFileBelowRoot(), TestContainsRejectsPathsOutsideRootAndThroughSymlinks(), TestStoreAddCanonicalizesAndDeduplicatesDirectories() (+9 more)

### Community 16 - "Interactive Setup"
Cohesion: 0.28
Nodes (15): collectInitInterview(), parseServiceDefinition(), promptDefault(), promptRequired(), promptServiceName(), TestCollectInitInterviewAcceptsEditsAndAdditionalServices(), TestCollectInitInterviewAcceptsSuppliedProjectName(), TestCollectInitInterviewRequiresAtLeastOneService() (+7 more)

### Community 17 - "Logging And Locks"
Cohesion: 0.19
Nodes (12): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, T, Context, T, Registry Lock Tests, TestRegistryLockHonorsContextWhileContended(), TestRegistryLockReleasesAfterCallbackPanic() (+4 more)

### Community 18 - "Contributing And CI"
Cohesion: 0.14
Nodes (14): Code of conduct, Contributing to grat, Cross-Platform Compatibility, Development setup, Focused Pull Requests, Local Quality Gate, Pull requests, README Contract Check (+6 more)

### Community 19 - "Runtime State Storage"
Cohesion: 0.28
Nodes (4): Manager, loadedState, processState, Time

### Community 20 - "Release Distribution"
Cohesion: 0.20
Nodes (15): GitHub Releases, Release Binary Installation, Homebrew Bottle Builder, Homebrew Bottle Tests, Artifact Download, Binary Artifacts, Release Build Job, Checksum Generation (+7 more)

### Community 21 - "Product Documentation"
Cohesion: 0.15
Nodes (13): Configuration compatibility, Approved Declarative Tasks, Approved Task Management, Command contract, Declarative Non-Executable Configuration, Go HTTP API, Interactive Project Initialization, Laravel (+5 more)

### Community 22 - "Security Documentation"
Cohesion: 0.21
Nodes (10): Fixed System Tool Paths, Private Vulnerability Reporting, Reporting a vulnerability, Security policy, Documented Shell Semantics, Supported versions, Trusted Configuration Boundary, Diagnostic Support Request (+2 more)

### Community 23 - "Linux Listener Inspection"
Cohesion: 0.25
Nodes (8): Listener, T, systemListener(), linuxListeningSocketInodes(), linuxSocketOwnerPIDs(), socketInode(), TestLinuxListenerParsingFindsOnlyTCPListenInodesForTheRequestedPort(), TestLinuxSocketOwnersFindsPIDsFromProcFileDescriptors()

### Community 24 - "Service Safety Architecture"
Cohesion: 0.24
Nodes (10): Cancellation Recovery, HTTP Service Readiness, Managed Process Identity, Process Group Isolation, Process Group Shutdown, Process-Owned Service Readiness, Role-Based Port Registry, Safe Local Service Manager (+2 more)

### Community 25 - "Runtime Progress Reporting"
Cohesion: 0.33
Nodes (6): Manager, Service, ProgressEvent, ProgressObserver, ProgressObserverFunc, ProgressStage

### Community 26 - "Runtime Readiness"
Cohesion: 0.31
Nodes (7): Context, processState, Manager, Service, readiness, isInProcessTree(), parentProcessID()

### Community 27 - "Project Root Discovery"
Cohesion: 0.47
Nodes (4): T, FindRoot(), TestFindRootReturnsNotFoundOutsideProject(), TestFindRootUsesNearestConfig()

### Community 28 - "Bottle Packaging Tests"
Cohesion: 0.60
Nodes (5): assert_archive_contains(), assert_binary(), assert_file(), assert_mode(), test-homebrew-bottles.sh script

### Community 29 - "Bottle Packaging Script"
Cohesion: 0.70
Nodes (4): package(), usage(), write_formula(), build-homebrew-bottles.sh script

### Community 30 - "Bottle Verification"
Cohesion: 0.60
Nodes (3): usage(), verify_bottle(), verify-homebrew-bottles.sh script

### Community 31 - "Project Agent Rules"
Cohesion: 0.67
Nodes (3): Graphify Before Push, Project Agent Rules, Versioned Graphify Artifacts

### Community 32 - "README Contract Script"
Cohesion: 0.83
Nodes (3): require(), require_in(), check-readme.sh script

## Knowledge Gaps
- **107 isolated node(s):** `Store`, `StepKind`, `LifecycleEvent`, `ProgressStage`, `LifecycleStage` (+102 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **8 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `New()` connect `Presentation Tests` to `CLI Orchestration`, `CLI Test Utilities`, `Update Maintenance`, `Terminal Presentation`, `Uninstall Maintenance`, `CLI Command Tests`, `Directory CLI Tests`, `Runtime Process Manager`, `Settings Tests`, `Interactive Setup`?**
  _High betweenness centrality (0.209) - this node is a cross-community bridge._
- **Why does `grat` connect `Project Documentation` to `Update Maintenance`, `Logging And Locks`, `Contributing And CI`, `Release Distribution`, `Product Documentation`, `Security Documentation`, `Service Safety Architecture`?**
  _High betweenness centrality (0.201) - this node is a cross-community bridge._
- **Why does `Contains()` connect `Presentation Tests` to `CLI Test Utilities`, `Uninstall Maintenance`, `CLI Command Tests`, `Configuration Validation`, `Directory CLI Tests`, `Runtime Manager Tests`, `Settings Tests`, `Interactive Setup`?**
  _High betweenness centrality (0.178) - this node is a cross-community bridge._
- **Are the 48 inferred relationships involving `Contains()` (e.g. with `TestInitRejectsDeprecatedAppFlag()` and `TestInitRejectsInvalidGlobalRegistry()`) actually correct?**
  _`Contains()` has 48 INFERRED edges - model-reasoned connections that need verification._
- **Are the 20 inferred relationships involving `Run()` (e.g. with `New()` and `Current()`) actually correct?**
  _`Run()` has 20 INFERRED edges - model-reasoned connections that need verification._
- **Are the 29 inferred relationships involving `New()` (e.g. with `configuredRoots()` and `Run()`) actually correct?**
  _`New()` has 29 INFERRED edges - model-reasoned connections that need verification._
- **Are the 12 inferred relationships involving `runWithEnvironment()` (e.g. with `New()` and `runWithConfiguredRoots()`) actually correct?**
  _`runWithEnvironment()` has 12 INFERRED edges - model-reasoned connections that need verification._
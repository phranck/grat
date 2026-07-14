# Graph Report - .  (2026-07-14)

## Corpus Check
- 1 files · ~37,919 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 793 nodes · 1871 edges · 41 communities (33 shown, 8 thin omitted)
- Extraction: 83% EXTRACTED · 17% INFERRED · 0% AMBIGUOUS · INFERRED: 321 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_CLI Service Lifecycle|CLI Service Lifecycle]]
- [[_COMMUNITY_Runtime Process Management|Runtime Process Management]]
- [[_COMMUNITY_Terminal Lifecycle Display|Terminal Lifecycle Display]]
- [[_COMMUNITY_CLI Command Tests|CLI Command Tests]]
- [[_COMMUNITY_Release Update Client|Release Update Client]]
- [[_COMMUNITY_Maintenance Test Helpers|Maintenance Test Helpers]]
- [[_COMMUNITY_Presentation Rendering Tests|Presentation Rendering Tests]]
- [[_COMMUNITY_Project Configuration Validation|Project Configuration Validation]]
- [[_COMMUNITY_CLI Presentation Framework|CLI Presentation Framework]]
- [[_COMMUNITY_README Product Guide|README Product Guide]]
- [[_COMMUNITY_Port Registry and Scanning|Port Registry and Scanning]]
- [[_COMMUNITY_Directory Command Tests|Directory Command Tests]]
- [[_COMMUNITY_CLI Bootstrap|CLI Bootstrap]]
- [[_COMMUNITY_Update Installation|Update Installation]]
- [[_COMMUNITY_Settings Safety Tests|Settings Safety Tests]]
- [[_COMMUNITY_Settings Persistence|Settings Persistence]]
- [[_COMMUNITY_Interactive Project Setup|Interactive Project Setup]]
- [[_COMMUNITY_Runtime Log Tests|Runtime Log Tests]]
- [[_COMMUNITY_Contributing and CI|Contributing and CI]]
- [[_COMMUNITY_Runtime State Storage|Runtime State Storage]]
- [[_COMMUNITY_Release Distribution|Release Distribution]]
- [[_COMMUNITY_Security Documentation|Security Documentation]]
- [[_COMMUNITY_Linux Listener Inspection|Linux Listener Inspection]]
- [[_COMMUNITY_Framework Configuration Examples|Framework Configuration Examples]]
- [[_COMMUNITY_Help Rendering|Help Rendering]]
- [[_COMMUNITY_Safe Service Design|Safe Service Design]]
- [[_COMMUNITY_Project Root Discovery|Project Root Discovery]]
- [[_COMMUNITY_Bottle Packaging Tests|Bottle Packaging Tests]]
- [[_COMMUNITY_Bottle Packaging Script|Bottle Packaging Script]]
- [[_COMMUNITY_Published Bottle Verification|Published Bottle Verification]]
- [[_COMMUNITY_Graphify Project Rules|Graphify Project Rules]]
- [[_COMMUNITY_CLI Test Helpers|CLI Test Helpers]]
- [[_COMMUNITY_README Contract Script|README Contract Script]]
- [[_COMMUNITY_Darwin Listener Tests|Darwin Listener Tests]]
- [[_COMMUNITY_Readiness Inspection Tests|Readiness Inspection Tests]]
- [[_COMMUNITY_System Listener Lookup|System Listener Lookup]]
- [[_COMMUNITY_Code of Conduct|Code of Conduct]]
- [[_COMMUNITY_Update Release Verification|Update Release Verification]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
- [[_COMMUNITY_Bottle Release Verification Tests|Bottle Release Verification Tests]]
- [[_COMMUNITY_Help Contract Tests|Help Contract Tests]]

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
- `grat` --references--> `Service.Uninstall`  [EXTRACTED]
  README.md → internal/maintenance/uninstall.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Managed Service Lifecycle** — readme_managed_process_identity, readme_http_service_readiness, readme_worker_readiness, readme_process_group_shutdown, readme_cancellation_recovery [INFERRED 0.95]
- **Safe Port Coordination** — readme_registered_directory_discovery, readme_role_based_port_registry, readme_serialized_port_operations [INFERRED 0.95]
- **Safe Maintenance** — readme_update_workflow, readme_verified_binary_replacement, readme_safe_uninstall_workflow, readme_machine_local_settings [INFERRED 0.85]

## Communities (41 total, 8 thin omitted)

### Community 0 - "CLI Service Lifecycle"
Cohesion: 0.06
Nodes (84): assignReassignedPorts(), configuredRoots(), copyReservations(), defaultEnvironment(), detectServices(), ensureValidRegistry(), executeLifecycle(), fileExists() (+76 more)

### Community 1 - "Runtime Process Management"
Cohesion: 0.05
Nodes (59): repeatedValue, File, Listener, T, Config, Listener, Manager, Service (+51 more)

### Community 2 - "Terminal Lifecycle Display"
Cohesion: 0.10
Nodes (32): CancelFunc, Cmd, Context, Reader, Style, Writer, lifecycleRow, Model (+24 more)

### Community 3 - "CLI Command Tests"
Cohesion: 0.15
Nodes (36): exitCode(), Run(), assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort(), loadPortFixtureConfig(), runWithConfiguredRoots() (+28 more)

### Community 4 - "Release Update Client"
Cohesion: 0.12
Nodes (18): asset, installation, Client, Context, Service, Client, Context, Service (+10 more)

### Community 5 - "Maintenance Test Helpers"
Cohesion: 0.14
Nodes (32): Service, Store, T, Context, Service, T, commandResponse, fakeCommands (+24 more)

### Community 6 - "Presentation Rendering Tests"
Cohesion: 0.17
Nodes (35): T, NewLifecycleModel(), DividerLine(), New(), columnOf(), displayColumn(), stripANSI(), TestHelpRendersBorderlessGroupsWithGloballyAlignedDescriptions() (+27 more)

### Community 7 - "Project Configuration Validation"
Cohesion: 0.12
Nodes (31): Config, containsControl(), DefaultRuntime(), InferRole(), Load(), replaceFile(), rollbackWrites(), safeServiceName() (+23 more)

### Community 8 - "CLI Presentation Framework"
Cohesion: 0.14
Nodes (12): Renderer, ColorMode, formatProjectRows(), fprint(), fprintln(), pad(), ParseColorMode(), stepStyle() (+4 more)

### Community 9 - "README Product Guide"
Cohesion: 0.07
Nodes (34): Bounded Local Logs, Bounded Service Logs, CI Workflow, Code of Conduct, Commands, Complete Development Stacks, Configuration reference, Contents (+26 more)

### Community 10 - "Port Registry and Scanning"
Cohesion: 0.14
Nodes (27): Config, T, fakeLookup, Listener, ListenerLookup, Problem, ProjectConfig, FirstFree() (+19 more)

### Community 11 - "Directory Command Tests"
Cohesion: 0.21
Nodes (23): canonicalCLITestPath(), environmentForTest(), newCLITestStore(), sameStringSlices(), TestDirectoriesAddDoesNotPromptForInitialSetup(), TestDirectoriesCommandsPersistAndListConfiguredRoots(), TestFirstUseAcceptsExistingSitesDefault(), TestFirstUseFallsBackToWorkingDirectory() (+15 more)

### Community 12 - "CLI Bootstrap"
Cohesion: 0.17
Nodes (13): main(), mustGetwd(), Client, Config, Context, processState, ProgressObserver, Manager (+5 more)

### Community 13 - "Update Installation"
Cohesion: 0.19
Nodes (14): Context, Service, Reader, Result, Store, Writer, installation, installationKind (+6 more)

### Community 14 - "Settings Safety Tests"
Cohesion: 0.31
Nodes (17): Store, T, canonicalPath(), equalStrings(), newTestStore(), TestContainsAcceptsRegularFileBelowRoot(), TestContainsRejectsPathsOutsideRootAndThroughSymlinks(), TestStoreAddCanonicalizesAndDeduplicatesDirectories() (+9 more)

### Community 15 - "Settings Persistence"
Cohesion: 0.30
Nodes (4): Settings, canonicalExistingDirectory(), canonicalExistingPath(), Store

### Community 16 - "Interactive Project Setup"
Cohesion: 0.28
Nodes (15): collectInitInterview(), parseServiceDefinition(), promptDefault(), promptRequired(), promptServiceName(), TestCollectInitInterviewAcceptsEditsAndAdditionalServices(), TestCollectInitInterviewAcceptsSuppliedProjectName(), TestCollectInitInterviewRequiresAtLeastOneService() (+7 more)

### Community 17 - "Runtime Log Tests"
Cohesion: 0.19
Nodes (12): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, T, Context, T, Registry Lock Tests, TestRegistryLockHonorsContextWhileContended(), TestRegistryLockReleasesAfterCallbackPanic() (+4 more)

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
Cohesion: 0.18
Nodes (12): Configuration compatibility, Approved Declarative Tasks, Fixed System Tool Paths, Private Vulnerability Reporting, Reporting a vulnerability, Security policy, Documented Shell Semantics, Supported versions (+4 more)

### Community 22 - "Linux Listener Inspection"
Cohesion: 0.25
Nodes (8): Listener, T, systemListener(), linuxListeningSocketInodes(), linuxSocketOwnerPIDs(), socketInode(), TestLinuxListenerParsingFindsOnlyTCPListenInodesForTheRequestedPort(), TestLinuxSocketOwnersFindsPIDsFromProcFileDescriptors()

### Community 23 - "Framework Configuration Examples"
Cohesion: 0.18
Nodes (11): Approved Task Management, Command contract, Declarative Non-Executable Configuration, Go HTTP API, Interactive Project Initialization, Laravel, Project examples, Python with FastAPI (+3 more)

### Community 24 - "Help Rendering"
Cohesion: 0.38
Nodes (4): Renderer, Style, CommandGroup, helpUsageWidth()

### Community 25 - "Safe Service Design"
Cohesion: 0.24
Nodes (10): Cancellation Recovery, HTTP Service Readiness, Managed Process Identity, Process Group Isolation, Process Group Shutdown, Process-Owned Service Readiness, Role-Based Port Registry, Safe Local Service Manager (+2 more)

### Community 26 - "Project Root Discovery"
Cohesion: 0.47
Nodes (4): T, FindRoot(), TestFindRootReturnsNotFoundOutsideProject(), TestFindRootUsesNearestConfig()

### Community 27 - "Bottle Packaging Tests"
Cohesion: 0.60
Nodes (5): assert_archive_contains(), assert_binary(), assert_file(), assert_mode(), test-homebrew-bottles.sh script

### Community 28 - "Bottle Packaging Script"
Cohesion: 0.70
Nodes (4): package(), usage(), write_formula(), build-homebrew-bottles.sh script

### Community 29 - "Published Bottle Verification"
Cohesion: 0.60
Nodes (3): usage(), verify_bottle(), verify-homebrew-bottles.sh script

### Community 30 - "Graphify Project Rules"
Cohesion: 0.67
Nodes (3): Graphify Before Push, Project Agent Rules, Versioned Graphify Artifacts

### Community 31 - "CLI Test Helpers"
Cohesion: 0.67
Nodes (3): TestListenerOwnerLabelHandlesUnknownPID(), TestLogFollowUsesTrustedExecutable(), T

### Community 32 - "README Contract Script"
Cohesion: 0.83
Nodes (3): require(), require_in(), check-readme.sh script

## Knowledge Gaps
- **107 isolated node(s):** `Store`, `StepKind`, `LifecycleEvent`, `ProgressStage`, `LifecycleStage` (+102 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **8 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `New()` connect `Presentation Rendering Tests` to `CLI Service Lifecycle`, `CLI Command Tests`, `Release Update Client`, `Maintenance Test Helpers`, `CLI Presentation Framework`, `CLI Bootstrap`, `Update Installation`, `Settings Safety Tests`, `Settings Persistence`, `Interactive Project Setup`?**
  _High betweenness centrality (0.209) - this node is a cross-community bridge._
- **Why does `grat` connect `README Product Guide` to `CLI Service Lifecycle`, `Update Release Verification`, `Runtime Log Tests`, `Contributing and CI`, `Release Distribution`, `Security Documentation`, `Framework Configuration Examples`, `Safe Service Design`?**
  _High betweenness centrality (0.201) - this node is a cross-community bridge._
- **Why does `Contains()` connect `Presentation Rendering Tests` to `Runtime Process Management`, `CLI Command Tests`, `Maintenance Test Helpers`, `Project Configuration Validation`, `Directory Command Tests`, `Update Installation`, `Settings Safety Tests`, `Settings Persistence`, `Interactive Project Setup`?**
  _High betweenness centrality (0.178) - this node is a cross-community bridge._
- **Are the 48 inferred relationships involving `Contains()` (e.g. with `TestInitRejectsDeprecatedAppFlag()` and `TestInitRejectsInvalidGlobalRegistry()`) actually correct?**
  _`Contains()` has 48 INFERRED edges - model-reasoned connections that need verification._
- **Are the 20 inferred relationships involving `Run()` (e.g. with `New()` and `Current()`) actually correct?**
  _`Run()` has 20 INFERRED edges - model-reasoned connections that need verification._
- **Are the 29 inferred relationships involving `New()` (e.g. with `configuredRoots()` and `Run()`) actually correct?**
  _`New()` has 29 INFERRED edges - model-reasoned connections that need verification._
- **Are the 12 inferred relationships involving `runWithEnvironment()` (e.g. with `New()` and `runWithConfiguredRoots()`) actually correct?**
  _`runWithEnvironment()` has 12 INFERRED edges - model-reasoned connections that need verification._
# Graph Report - .  (2026-07-26)

## Corpus Check
- 11 files · ~47,769 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 873 nodes · 2129 edges · 50 communities (35 shown, 15 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 343 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_CLI Service Lifecycle|CLI Service Lifecycle]]
- [[_COMMUNITY_Process Execution|Process Execution]]
- [[_COMMUNITY_Log Streaming Tests|Log Streaming Tests]]
- [[_COMMUNITY_Terminal Lifecycle UI|Terminal Lifecycle UI]]
- [[_COMMUNITY_Update Release Assets|Update Release Assets]]
- [[_COMMUNITY_Command Help Rendering|Command Help Rendering]]
- [[_COMMUNITY_Configuration Model|Configuration Model]]
- [[_COMMUNITY_Runtime Manager Tests|Runtime Manager Tests]]
- [[_COMMUNITY_CLI Directory Tests|CLI Directory Tests]]
- [[_COMMUNITY_Presentation Tests|Presentation Tests]]
- [[_COMMUNITY_Maintenance Uninstall|Maintenance Uninstall]]
- [[_COMMUNITY_Port Registry|Port Registry]]
- [[_COMMUNITY_CLI Test Helpers|CLI Test Helpers]]
- [[_COMMUNITY_Application Runtime|Application Runtime]]
- [[_COMMUNITY_Project Documentation|Project Documentation]]
- [[_COMMUNITY_Settings Store Tests|Settings Store Tests]]
- [[_COMMUNITY_Settings Persistence|Settings Persistence]]
- [[_COMMUNITY_Interactive Project Setup|Interactive Project Setup]]
- [[_COMMUNITY_Uninstall Tests|Uninstall Tests]]
- [[_COMMUNITY_Runtime State Storage|Runtime State Storage]]
- [[_COMMUNITY_Recovery CLI Tests|Recovery CLI Tests]]
- [[_COMMUNITY_Linux Listener Inspection|Linux Listener Inspection]]
- [[_COMMUNITY_Port Registry Lock|Port Registry Lock]]
- [[_COMMUNITY_Terminal Text Safety|Terminal Text Safety]]
- [[_COMMUNITY_Release Workflow|Release Workflow]]
- [[_COMMUNITY_Project Root Discovery|Project Root Discovery]]
- [[_COMMUNITY_Bottle Packaging Tests|Bottle Packaging Tests]]
- [[_COMMUNITY_Homebrew Bottle Builder|Homebrew Bottle Builder]]
- [[_COMMUNITY_Bottle Publication Verification|Bottle Publication Verification]]
- [[_COMMUNITY_CI Verification Workflow|CI Verification Workflow]]
- [[_COMMUNITY_Environment Security Contract|Environment Security Contract]]
- [[_COMMUNITY_README Contract Checks|README Contract Checks]]
- [[_COMMUNITY_Security Documentation|Security Documentation]]
- [[_COMMUNITY_Graphify Project Rule|Graphify Project Rule]]
- [[_COMMUNITY_Darwin Listener Tests|Darwin Listener Tests]]
- [[_COMMUNITY_Readiness Inspection Tests|Readiness Inspection Tests]]
- [[_COMMUNITY_Listener Lookup API|Listener Lookup API]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
- [[_COMMUNITY_Bottle Verification Tests|Bottle Verification Tests]]
- [[_COMMUNITY_Code Of Conduct|Code Of Conduct]]
- [[_COMMUNITY_Contribution Guide|Contribution Guide]]
- [[_COMMUNITY_Port Allocation|Port Allocation]]
- [[_COMMUNITY_HTTP Readiness|HTTP Readiness]]
- [[_COMMUNITY_Service Management|Service Management]]
- [[_COMMUNITY_Artifact Attestation|Artifact Attestation]]
- [[_COMMUNITY_Trusted Configuration|Trusted Configuration]]
- [[_COMMUNITY_Platform Helper Safety|Platform Helper Safety]]
- [[_COMMUNITY_Release Provenance|Release Provenance]]
- [[_COMMUNITY_Settings Store|Settings Store]]

## God Nodes (most connected - your core abstractions)
1. `Contains()` - 70 edges
2. `New()` - 39 edges
3. `runWithEnvironment()` - 35 edges
4. `T` - 28 edges
5. `Service` - 27 edges
6. `T` - 25 edges
7. `Manager` - 24 edges
8. `grat` - 23 edges
9. `Renderer` - 22 edges
10. `runWithConfiguredRoots()` - 21 edges

## Surprising Connections (you probably didn't know these)
- `Registry Lock Tests` --references--> `WithRegistryLock()`  [EXTRACTED]
  /Users/phranck/Developer/tools/cli/grat/internal/ports/lock_test.go → internal/ports/lock.go
- `Settings Safety Tests` --references--> `Contains()`  [EXTRACTED]
  /Users/phranck/Developer/tools/cli/grat/internal/settings/settings_test.go → internal/settings/settings.go
- `defaultEnvironment()` --calls--> `DefaultService()`  [INFERRED]
  internal/cli/cli.go → internal/maintenance/system.go
- `runWithEnvironment()` --calls--> `New()`  [INFERRED]
  internal/cli/cli.go → internal/presentation/presentation.go
- `runWithConfiguredRoots()` --calls--> `runWithEnvironment()`  [INFERRED]
  internal/cli/cli_test.go → internal/cli/cli.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Service Orchestration Flow** — readme_service_lifecycle_management, readme_role_specific_readiness, readme_managed_state, readme_process_group_shutdown [EXTRACTED 1.00]
- **Verified Release Supply Chain** — readme_update_verification, readme_artifact_attestation, readme_homebrew_tap, readme_github_releases [EXTRACTED 1.00]
- **Fail-Closed Process Safety** — readme_managed_process_identity, readme_process_group_shutdown, readme_fail_closed_identity_validation, readme_legacy_process_recovery [EXTRACTED 1.00]

## Communities (50 total, 15 thin omitted)

### Community 0 - "CLI Service Lifecycle"
Cohesion: 0.06
Nodes (81): assignReassignedPorts(), configuredRoots(), confirmRecovery(), copyReservations(), defaultEnvironment(), detectServices(), ensureValidRegistry(), executeLifecycle() (+73 more)

### Community 1 - "Process Execution"
Cohesion: 0.09
Nodes (42): Duration, Context, loadedState, processState, Manager, Service, T, Context (+34 more)

### Community 2 - "Log Streaming Tests"
Cohesion: 0.08
Nodes (39): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, repeatedValue, File, T, Context, Service, T (+31 more)

### Community 3 - "Terminal Lifecycle UI"
Cohesion: 0.10
Nodes (32): CancelFunc, Cmd, Context, Reader, Style, Writer, lifecycleRow, Model (+24 more)

### Community 4 - "Update Release Assets"
Cohesion: 0.10
Nodes (22): asset, installation, Client, Context, Service, Client, Context, Service (+14 more)

### Community 5 - "Command Help Rendering"
Cohesion: 0.11
Nodes (19): Renderer, Style, ColorMode, Command, CommandGroup, helpUsageWidth(), formatProjectRows(), fprint() (+11 more)

### Community 6 - "Configuration Model"
Cohesion: 0.10
Nodes (39): Config, configDecodeError(), DefaultRuntime(), InferRole(), Load(), readConfigFile(), replaceFile(), rollbackWrites() (+31 more)

### Community 7 - "Runtime Manager Tests"
Cohesion: 0.11
Nodes (37): Config, Service, T, T, Manager, Service, Listener, Manager (+29 more)

### Community 8 - "CLI Directory Tests"
Cohesion: 0.12
Nodes (35): exitCode(), hasHelpFlag(), isHelp(), runWithEnvironment(), canonicalCLITestPath(), environmentForTest(), newCLITestStore(), sameStringSlices() (+27 more)

### Community 9 - "Presentation Tests"
Cohesion: 0.16
Nodes (38): T, NewLifecycleModel(), DividerLine(), New(), columnOf(), displayColumn(), stripANSI(), TestHelpRendersBorderlessGroupsWithGloballyAlignedDescriptions() (+30 more)

### Community 10 - "Maintenance Uninstall"
Cohesion: 0.13
Nodes (23): Context, Service, Reader, Result, Store, Writer, Context, T (+15 more)

### Community 11 - "Port Registry"
Cohesion: 0.13
Nodes (31): Config, T, fakeLookup, Listener, ListenerLookup, Problem, ProjectConfig, FirstFree() (+23 more)

### Community 12 - "CLI Test Helpers"
Cohesion: 0.15
Nodes (32): assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort(), loadPortFixtureConfig(), runWithConfiguredRoots(), TestCLIRuntimeHelper(), TestExitCodeMapsInterruptedOperationsTo130() (+24 more)

### Community 13 - "Application Runtime"
Cohesion: 0.14
Nodes (16): Client, main(), mustGetwd(), Config, Context, loadedState, processState, Manager (+8 more)

### Community 14 - "Project Documentation"
Cohesion: 0.09
Nodes (31): GitHub Artifact Attestation, BACKEND_URL Discovery, Code of Conduct, Contributing Guide, Directory Discovery, Fail-Closed Identity Validation, FastAPI Deployment Guide, grat GitHub Releases (+23 more)

### Community 15 - "Settings Store Tests"
Cohesion: 0.31
Nodes (17): Store, T, canonicalPath(), equalStrings(), newTestStore(), TestContainsAcceptsRegularFileBelowRoot(), TestContainsRejectsPathsOutsideRootAndThroughSymlinks(), TestStoreAddCanonicalizesAndDeduplicatesDirectories() (+9 more)

### Community 16 - "Settings Persistence"
Cohesion: 0.30
Nodes (4): Settings, canonicalExistingDirectory(), canonicalExistingPath(), Store

### Community 17 - "Interactive Project Setup"
Cohesion: 0.28
Nodes (15): collectInitInterview(), parseServiceDefinition(), promptDefault(), promptRequired(), promptServiceName(), TestCollectInitInterviewAcceptsEditsAndAdditionalServices(), TestCollectInitInterviewAcceptsSuppliedProjectName(), TestCollectInitInterviewRequiresAtLeastOneService() (+7 more)

### Community 18 - "Uninstall Tests"
Cohesion: 0.34
Nodes (16): Service, Store, T, fakeUninstallService(), newUninstallStore(), TestDiscoverUninstallArtifactsRejectsScanLimitOverrun(), TestUninstallAbortsBeforePromptsForActiveManagedService(), TestUninstallDefaultYesRemovesOnlyRegisteredProjectArtifacts() (+8 more)

### Community 19 - "Runtime State Storage"
Cohesion: 0.23
Nodes (6): Manager, Service, loadedState, processState, RecoveryCandidate, Time

### Community 20 - "Recovery CLI Tests"
Cohesion: 0.46
Nodes (14): assertCLIRecoveryState(), assertRecoveryPreview(), cliProcessAlive(), legacyCLIStartIdentity(), recoveryEnvironment(), stopCLIRecoveryGroup(), TestRecoverDeclinedConfirmationLeavesLegacyProcessAndState(), TestRecoverInteractiveConfirmationStopsLegacyProcessAndRemovesState() (+6 more)

### Community 21 - "Linux Listener Inspection"
Cohesion: 0.25
Nodes (8): Listener, T, systemListener(), linuxListeningSocketInodes(), linuxSocketOwnerPIDs(), socketInode(), TestLinuxListenerParsingFindsOnlyTCPListenInodesForTheRequestedPort(), TestLinuxSocketOwnersFindsPIDsFromProcFileDescriptors()

### Community 22 - "Port Registry Lock"
Cohesion: 0.33
Nodes (9): Context, T, Registry Lock Tests, TestRegistryLockHonorsContextWhileContended(), TestRegistryLockReleasesAfterCallbackPanic(), TestRegistryLockSerializesCallbacks(), TestRegistryLockUsesProvidedGratConfigurationDirectory(), WithRegistryLock() (+1 more)

### Community 23 - "Terminal Text Safety"
Cohesion: 0.39
Nodes (6): T, ContainsUnsafe(), Sanitize(), TestSanitizeReplacesEveryUnsafeRune(), TestUnsafeRuneRejectsControlsAndUnicodeFormatCharacters(), UnsafeRune()

### Community 24 - "Release Workflow"
Cohesion: 0.38
Nodes (7): GitHub Artifact Attestation, Cross-Platform Build Job, GitHub Release Publication, Release Publish Job, Release Workflow, SHA-256 Release Checksums, Tag-Triggered Release

### Community 25 - "Project Root Discovery"
Cohesion: 0.47
Nodes (4): T, FindRoot(), TestFindRootReturnsNotFoundOutsideProject(), TestFindRootUsesNearestConfig()

### Community 26 - "Bottle Packaging Tests"
Cohesion: 0.60
Nodes (5): assert_archive_contains(), assert_binary(), assert_file(), assert_mode(), test-homebrew-bottles.sh script

### Community 27 - "Homebrew Bottle Builder"
Cohesion: 0.70
Nodes (4): package(), usage(), write_formula(), build-homebrew-bottles.sh script

### Community 28 - "Bottle Publication Verification"
Cohesion: 0.60
Nodes (3): usage(), verify_bottle(), verify-homebrew-bottles.sh script

### Community 29 - "CI Verification Workflow"
Cohesion: 0.40
Nodes (5): CI Workflow, Platform Verification Matrix, Race-Enabled Tests, Verify Job, Vulnerability Scan

### Community 30 - "Environment Security Contract"
Cohesion: 0.50
Nodes (4): BACKEND_URL, inherit_env, Non-Secret Environment Baseline, Trusted Local Configurations

### Community 31 - "README Contract Checks"
Cohesion: 0.83
Nodes (3): require(), require_in(), check-readme.sh script

### Community 32 - "Security Documentation"
Cohesion: 0.50
Nodes (4): Security Policy, Diagnostic Support Request, Sensitive Data Redaction, Support

### Community 33 - "Graphify Project Rule"
Cohesion: 1.00
Nodes (3): Graphify Before Push, Project Agent Rules, Versioned Graphify Artifacts

## Knowledge Gaps
- **95 isolated node(s):** `Store`, `ProgressObserver`, `StepKind`, `LifecycleEvent`, `ProgressStage` (+90 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **15 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Contains()` connect `Presentation Tests` to `Process Execution`, `Log Streaming Tests`, `Configuration Model`, `Runtime Manager Tests`, `CLI Directory Tests`, `Maintenance Uninstall`, `CLI Test Helpers`, `Settings Store Tests`, `Settings Persistence`, `Interactive Project Setup`, `Uninstall Tests`, `Recovery CLI Tests`?**
  _High betweenness centrality (0.215) - this node is a cross-community bridge._
- **Why does `New()` connect `Presentation Tests` to `CLI Service Lifecycle`, `Log Streaming Tests`, `Update Release Assets`, `Command Help Rendering`, `CLI Directory Tests`, `Maintenance Uninstall`, `CLI Test Helpers`, `Application Runtime`, `Settings Store Tests`, `Settings Persistence`, `Interactive Project Setup`, `Uninstall Tests`?**
  _High betweenness centrality (0.205) - this node is a cross-community bridge._
- **Why does `processAlive()` connect `Process Execution` to `Log Streaming Tests`, `Application Runtime`, `Runtime Manager Tests`?**
  _High betweenness centrality (0.083) - this node is a cross-community bridge._
- **Are the 66 inferred relationships involving `Contains()` (e.g. with `TestInitRejectsDeprecatedAppFlag()` and `TestInitRejectsInvalidGlobalRegistry()`) actually correct?**
  _`Contains()` has 66 INFERRED edges - model-reasoned connections that need verification._
- **Are the 34 inferred relationships involving `New()` (e.g. with `configuredRoots()` and `runDirectories()`) actually correct?**
  _`New()` has 34 INFERRED edges - model-reasoned connections that need verification._
- **Are the 17 inferred relationships involving `runWithEnvironment()` (e.g. with `New()` and `Current()`) actually correct?**
  _`runWithEnvironment()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Store`, `ProgressObserver`, `StepKind` to the rest of the system?**
  _100 weakly-connected nodes found - possible documentation gaps or missing edges._
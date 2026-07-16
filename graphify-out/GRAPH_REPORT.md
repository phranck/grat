# Graph Report - .  (2026-07-16)

## Corpus Check
- Corpus is ~47,459 words - fits in a single context window. You may not need a graph.

## Summary
- 867 nodes · 2103 edges · 56 communities (34 shown, 22 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 340 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_CLI Service Lifecycle|CLI Service Lifecycle]]
- [[_COMMUNITY_Log Streaming Tests|Log Streaming Tests]]
- [[_COMMUNITY_Runtime File Utilities|Runtime File Utilities]]
- [[_COMMUNITY_Process Execution|Process Execution]]
- [[_COMMUNITY_Update Release Assets|Update Release Assets]]
- [[_COMMUNITY_Configuration Model|Configuration Model]]
- [[_COMMUNITY_Runtime Service Management|Runtime Service Management]]
- [[_COMMUNITY_Port Presentation Tests|Port Presentation Tests]]
- [[_COMMUNITY_CLI Directory Tests|CLI Directory Tests]]
- [[_COMMUNITY_Listener Test Fixtures|Listener Test Fixtures]]
- [[_COMMUNITY_Terminal Presentation|Terminal Presentation]]
- [[_COMMUNITY_Command Entrypoint|Command Entrypoint]]
- [[_COMMUNITY_Maintenance Service API|Maintenance Service API]]
- [[_COMMUNITY_CLI Test Helpers|CLI Test Helpers]]
- [[_COMMUNITY_Settings Path Validation|Settings Path Validation]]
- [[_COMMUNITY_Settings Persistence|Settings Persistence]]
- [[_COMMUNITY_Interactive Project Setup|Interactive Project Setup]]
- [[_COMMUNITY_Managed Process State|Managed Process State]]
- [[_COMMUNITY_Recovery CLI Tests|Recovery CLI Tests]]
- [[_COMMUNITY_Linux Listener Inspection|Linux Listener Inspection]]
- [[_COMMUNITY_Operation Locking|Operation Locking]]
- [[_COMMUNITY_Command Help Rendering|Command Help Rendering]]
- [[_COMMUNITY_Terminal Text Safety|Terminal Text Safety]]
- [[_COMMUNITY_Release Workflow|Release Workflow]]
- [[_COMMUNITY_Project Root Discovery|Project Root Discovery]]
- [[_COMMUNITY_Bottle Packaging Tests|Bottle Packaging Tests]]
- [[_COMMUNITY_Homebrew Bottle Builder|Homebrew Bottle Builder]]
- [[_COMMUNITY_Bottle Publication Verification|Bottle Publication Verification]]
- [[_COMMUNITY_CI Verification Workflow|CI Verification Workflow]]
- [[_COMMUNITY_Listener Lookup API|Listener Lookup API]]
- [[_COMMUNITY_Environment Security Contract|Environment Security Contract]]
- [[_COMMUNITY_README Contract Checks|README Contract Checks]]
- [[_COMMUNITY_Security Documentation|Security Documentation]]
- [[_COMMUNITY_Graphify Project Rule|Graphify Project Rule]]
- [[_COMMUNITY_Darwin Listener Tests|Darwin Listener Tests]]
- [[_COMMUNITY_Readiness Inspection Tests|Readiness Inspection Tests]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
- [[_COMMUNITY_Bottle Verification Tests|Bottle Verification Tests]]
- [[_COMMUNITY_Code Of Conduct|Code Of Conduct]]
- [[_COMMUNITY_Contribution Guide|Contribution Guide]]
- [[_COMMUNITY_Artifact Attestation|Artifact Attestation]]
- [[_COMMUNITY_grat Overview|grat Overview]]
- [[_COMMUNITY_Configuration File|Configuration File]]
- [[_COMMUNITY_Homebrew Distribution|Homebrew Distribution]]
- [[_COMMUNITY_Port Allocation|Port Allocation]]
- [[_COMMUNITY_Process Group Shutdown|Process Group Shutdown]]
- [[_COMMUNITY_HTTP Readiness|HTTP Readiness]]
- [[_COMMUNITY_Legacy Recovery|Legacy Recovery]]
- [[_COMMUNITY_Service Management|Service Management]]
- [[_COMMUNITY_Service Roles|Service Roles]]
- [[_COMMUNITY_Sigstore|Sigstore]]
- [[_COMMUNITY_Trusted Configuration|Trusted Configuration]]
- [[_COMMUNITY_Platform Helper Safety|Platform Helper Safety]]
- [[_COMMUNITY_Release Provenance|Release Provenance]]
- [[_COMMUNITY_Settings Store|Settings Store]]

## God Nodes (most connected - your core abstractions)
1. `Contains()` - 70 edges
2. `New()` - 38 edges
3. `runWithEnvironment()` - 35 edges
4. `Service` - 28 edges
5. `T` - 27 edges
6. `Manager` - 27 edges
7. `T` - 25 edges
8. `Config` - 24 edges
9. `Renderer` - 23 edges
10. `runWithConfiguredRoots()` - 21 edges

## Surprising Connections (you probably didn't know these)
- `Init Interview Tests` --references--> `runInitWithInput()`  [INFERRED]
  /Users/phranck/Developer/tools/cli/grat/internal/cli/init_interview_test.go → internal/cli/cli.go
- `Registry Lock Tests` --references--> `WithRegistryLock()`  [EXTRACTED]
  /Users/phranck/Developer/tools/cli/grat/internal/ports/lock_test.go → internal/ports/lock.go
- `Settings Safety Tests` --references--> `Contains()`  [EXTRACTED]
  /Users/phranck/Developer/tools/cli/grat/internal/settings/settings_test.go → internal/settings/settings.go
- `defaultEnvironment()` --calls--> `DefaultService()`  [INFERRED]
  internal/cli/cli.go → internal/maintenance/system.go
- `runWithEnvironment()` --calls--> `New()`  [INFERRED]
  internal/cli/cli.go → internal/presentation/presentation.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Release Supply-Chain Verification** — workflows_release_artifact_attestation, workflows_release_sha256_checksums, workflows_release_github_release_publication, readme_release_provenance_verification, security_fail_closed_update_verification [INFERRED 0.95]

## Communities (56 total, 22 thin omitted)

### Community 0 - "CLI Service Lifecycle"
Cohesion: 0.05
Nodes (89): assignReassignedPorts(), configuredRoots(), confirmRecovery(), copyReservations(), defaultEnvironment(), detectServices(), ensureValidRegistry(), executeLifecycle() (+81 more)

### Community 1 - "Log Streaming Tests"
Cohesion: 0.07
Nodes (52): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, repeatedValue, T, Service, Store, T, Context (+44 more)

### Community 2 - "Runtime File Utilities"
Cohesion: 0.08
Nodes (46): File, T, Context, Duration, loadedState, processState, Manager, Service (+38 more)

### Community 3 - "Process Execution"
Cohesion: 0.09
Nodes (34): CancelFunc, Cmd, Context, Reader, Style, Writer, lifecycleRow, Model (+26 more)

### Community 4 - "Update Release Assets"
Cohesion: 0.09
Nodes (23): asset, installation, Client, Context, Service, Client, Context, Service (+15 more)

### Community 5 - "Configuration Model"
Cohesion: 0.09
Nodes (40): Config, configDecodeError(), DefaultRuntime(), InferRole(), Load(), prepareWrite(), readConfigFile(), replaceFile() (+32 more)

### Community 6 - "Runtime Service Management"
Cohesion: 0.11
Nodes (37): Config, Listener, Manager, Service, T, T, Manager, Service (+29 more)

### Community 7 - "Port Presentation Tests"
Cohesion: 0.16
Nodes (38): TestRenderPortReassignSummaryGroupsAssignmentsByProject(), T, NewLifecycleModel(), DividerLine(), New(), columnOf(), displayColumn(), stripANSI() (+30 more)

### Community 8 - "CLI Directory Tests"
Cohesion: 0.14
Nodes (32): runWithEnvironment(), canonicalCLITestPath(), environmentForTest(), newCLITestStore(), sameStringSlices(), TestDirectoriesAddDoesNotPromptForInitialSetup(), TestDirectoriesCommandsPersistAndListConfiguredRoots(), TestFirstUseAcceptsExistingSitesDefault() (+24 more)

### Community 9 - "Listener Test Fixtures"
Cohesion: 0.12
Nodes (32): Config, Listener, T, fakeLookup, Listener, ListenerLookup, Problem, ProjectConfig (+24 more)

### Community 10 - "Terminal Presentation"
Cohesion: 0.13
Nodes (13): Renderer, Style, Writer, ColorMode, formatProjectRows(), fprint(), fprintln(), pad() (+5 more)

### Community 11 - "Command Entrypoint"
Cohesion: 0.14
Nodes (16): main(), mustGetwd(), Client, Config, Context, loadedState, processState, ProgressObserver (+8 more)

### Community 12 - "Maintenance Service API"
Cohesion: 0.13
Nodes (23): Context, Service, Reader, Result, Store, Writer, Context, T (+15 more)

### Community 13 - "CLI Test Helpers"
Cohesion: 0.16
Nodes (31): exitCode(), assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort(), loadPortFixtureConfig(), runWithConfiguredRoots(), TestCLIRuntimeHelper() (+23 more)

### Community 14 - "Settings Path Validation"
Cohesion: 0.28
Nodes (5): Settings, canonicalExistingDirectory(), canonicalExistingPath(), ConfigDirectory(), Store

### Community 15 - "Settings Persistence"
Cohesion: 0.31
Nodes (17): Store, T, canonicalPath(), equalStrings(), newTestStore(), TestContainsAcceptsRegularFileBelowRoot(), TestContainsRejectsPathsOutsideRootAndThroughSymlinks(), TestStoreAddCanonicalizesAndDeduplicatesDirectories() (+9 more)

### Community 16 - "Interactive Project Setup"
Cohesion: 0.28
Nodes (15): collectInitInterview(), parseServiceDefinition(), promptDefault(), promptRequired(), promptServiceName(), TestCollectInitInterviewAcceptsEditsAndAdditionalServices(), TestCollectInitInterviewAcceptsSuppliedProjectName(), TestCollectInitInterviewRequiresAtLeastOneService() (+7 more)

### Community 17 - "Managed Process State"
Cohesion: 0.23
Nodes (6): Manager, Service, loadedState, processState, RecoveryCandidate, Time

### Community 18 - "Recovery CLI Tests"
Cohesion: 0.46
Nodes (14): assertCLIRecoveryState(), assertRecoveryPreview(), cliProcessAlive(), legacyCLIStartIdentity(), recoveryEnvironment(), stopCLIRecoveryGroup(), TestRecoverDeclinedConfirmationLeavesLegacyProcessAndState(), TestRecoverInteractiveConfirmationStopsLegacyProcessAndRemovesState() (+6 more)

### Community 19 - "Linux Listener Inspection"
Cohesion: 0.25
Nodes (8): Listener, T, systemListener(), linuxListeningSocketInodes(), linuxSocketOwnerPIDs(), socketInode(), TestLinuxListenerParsingFindsOnlyTCPListenInodesForTheRequestedPort(), TestLinuxSocketOwnersFindsPIDsFromProcFileDescriptors()

### Community 20 - "Operation Locking"
Cohesion: 0.33
Nodes (9): Context, T, Registry Lock Tests, TestRegistryLockHonorsContextWhileContended(), TestRegistryLockReleasesAfterCallbackPanic(), TestRegistryLockSerializesCallbacks(), TestRegistryLockUsesProvidedGratConfigurationDirectory(), WithRegistryLock() (+1 more)

### Community 21 - "Command Help Rendering"
Cohesion: 0.36
Nodes (5): Renderer, Style, Command, CommandGroup, helpUsageWidth()

### Community 22 - "Terminal Text Safety"
Cohesion: 0.39
Nodes (6): T, ContainsUnsafe(), Sanitize(), TestSanitizeReplacesEveryUnsafeRune(), TestUnsafeRuneRejectsControlsAndUnicodeFormatCharacters(), UnsafeRune()

### Community 23 - "Release Workflow"
Cohesion: 0.38
Nodes (7): GitHub Artifact Attestation, Cross-Platform Build Job, GitHub Release Publication, Release Publish Job, Release Workflow, SHA-256 Release Checksums, Tag-Triggered Release

### Community 24 - "Project Root Discovery"
Cohesion: 0.47
Nodes (4): T, FindRoot(), TestFindRootReturnsNotFoundOutsideProject(), TestFindRootUsesNearestConfig()

### Community 25 - "Bottle Packaging Tests"
Cohesion: 0.60
Nodes (5): assert_archive_contains(), assert_binary(), assert_file(), assert_mode(), test-homebrew-bottles.sh script

### Community 26 - "Homebrew Bottle Builder"
Cohesion: 0.70
Nodes (4): package(), usage(), write_formula(), build-homebrew-bottles.sh script

### Community 27 - "Bottle Publication Verification"
Cohesion: 0.60
Nodes (3): usage(), verify_bottle(), verify-homebrew-bottles.sh script

### Community 28 - "CI Verification Workflow"
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
- **96 isolated node(s):** `Store`, `updateService`, `uninstallService`, `updateService`, `uninstallService` (+91 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **22 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `New()` connect `Port Presentation Tests` to `CLI Service Lifecycle`, `Log Streaming Tests`, `Update Release Assets`, `CLI Directory Tests`, `Terminal Presentation`, `Command Entrypoint`, `Maintenance Service API`, `Settings Path Validation`, `Settings Persistence`, `Interactive Project Setup`?**
  _High betweenness centrality (0.225) - this node is a cross-community bridge._
- **Why does `Contains()` connect `Port Presentation Tests` to `Log Streaming Tests`, `Runtime File Utilities`, `Configuration Model`, `Runtime Service Management`, `CLI Directory Tests`, `Maintenance Service API`, `CLI Test Helpers`, `Settings Path Validation`, `Settings Persistence`, `Interactive Project Setup`, `Recovery CLI Tests`?**
  _High betweenness centrality (0.224) - this node is a cross-community bridge._
- **Why does `processAlive()` connect `Runtime File Utilities` to `Log Streaming Tests`, `Command Entrypoint`, `Runtime Service Management`?**
  _High betweenness centrality (0.085) - this node is a cross-community bridge._
- **Are the 66 inferred relationships involving `Contains()` (e.g. with `TestInitRejectsDeprecatedAppFlag()` and `TestInitRejectsInvalidGlobalRegistry()`) actually correct?**
  _`Contains()` has 66 INFERRED edges - model-reasoned connections that need verification._
- **Are the 33 inferred relationships involving `New()` (e.g. with `configuredRoots()` and `runDirectories()`) actually correct?**
  _`New()` has 33 INFERRED edges - model-reasoned connections that need verification._
- **Are the 17 inferred relationships involving `runWithEnvironment()` (e.g. with `New()` and `Current()`) actually correct?**
  _`runWithEnvironment()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Store`, `updateService`, `uninstallService` to the rest of the system?**
  _100 weakly-connected nodes found - possible documentation gaps or missing edges._
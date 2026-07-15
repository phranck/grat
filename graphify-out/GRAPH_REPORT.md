# Graph Report - grat  (2026-07-15)

## Corpus Check
- 71 files · ~46,787 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 861 nodes · 2108 edges · 41 communities (33 shown, 8 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 346 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Application Runtime|Application Runtime]]
- [[_COMMUNITY_CLI Command Dispatch|CLI Command Dispatch]]
- [[_COMMUNITY_CLI Directory Setup|CLI Directory Setup]]
- [[_COMMUNITY_CLI Log Handling|CLI Log Handling]]
- [[_COMMUNITY_CLI Command Tests|CLI Command Tests]]
- [[_COMMUNITY_Init Interview|Init Interview]]
- [[_COMMUNITY_Legacy Recovery Tests|Legacy Recovery Tests]]
- [[_COMMUNITY_Configuration Persistence|Configuration Persistence]]
- [[_COMMUNITY_Update Installation Service|Update Installation Service]]
- [[_COMMUNITY_Maintenance Uninstall|Maintenance Uninstall]]
- [[_COMMUNITY_Uninstall Tests|Uninstall Tests]]
- [[_COMMUNITY_Settings Path Handling|Settings Path Handling]]
- [[_COMMUNITY_Listener Interface|Listener Interface]]
- [[_COMMUNITY_Darwin Listener Tests|Darwin Listener Tests]]
- [[_COMMUNITY_Linux Listener Discovery|Linux Listener Discovery]]
- [[_COMMUNITY_Port Registry Lock|Port Registry Lock]]
- [[_COMMUNITY_Port Registry|Port Registry]]
- [[_COMMUNITY_Terminal Presentation|Terminal Presentation]]
- [[_COMMUNITY_Terminal Lifecycle UI|Terminal Lifecycle UI]]
- [[_COMMUNITY_Presentation Tests|Presentation Tests]]
- [[_COMMUNITY_Project Root Discovery|Project Root Discovery]]
- [[_COMMUNITY_Runtime Manager Tests|Runtime Manager Tests]]
- [[_COMMUNITY_Runtime Service Manager|Runtime Service Manager]]
- [[_COMMUNITY_Readiness Tests|Readiness Tests]]
- [[_COMMUNITY_Runtime State Storage|Runtime State Storage]]
- [[_COMMUNITY_Settings Store Tests|Settings Store Tests]]
- [[_COMMUNITY_Text Sanitization|Text Sanitization]]
- [[_COMMUNITY_Bottle Packaging|Bottle Packaging]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
- [[_COMMUNITY_README Contract Checks|README Contract Checks]]
- [[_COMMUNITY_Bottle Package Tests|Bottle Package Tests]]
- [[_COMMUNITY_Bottle Verification Tests|Bottle Verification Tests]]
- [[_COMMUNITY_Bottle Verification|Bottle Verification]]
- [[_COMMUNITY_Settings Store|Settings Store]]
- [[_COMMUNITY_CI Verification|CI Verification]]
- [[_COMMUNITY_Security Release Documentation|Security Release Documentation]]
- [[_COMMUNITY_Graphify Release Index|Graphify Release Index]]
- [[_COMMUNITY_Code of Conduct|Code of Conduct]]
- [[_COMMUNITY_Contribution Guide|Contribution Guide]]
- [[_COMMUNITY_Project Documentation|Project Documentation]]

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
- **Managed Service Lifecycle** — readme_declarative_configuration, readme_http_readiness, readme_process_group_shutdown [INFERRED 0.85]

## Communities (41 total, 8 thin omitted)

### Community 11 - "Application Runtime"
Cohesion: 0.14
Nodes (16): main(), mustGetwd(), Manager, State, Status, Service, Config, ListenerLookup (+8 more)

### Community 0 - "CLI Command Dispatch"
Cohesion: 0.05
Nodes (89): Run(), Context, Writer, environment, Reader, Store, updateService, uninstallService (+81 more)

### Community 7 - "CLI Directory Setup"
Cohesion: 0.14
Nodes (32): runWithEnvironment(), TestDirectoriesCommandsPersistAndListConfiguredRoots(), T, TestDirectoriesAddDoesNotPromptForInitialSetup(), TestFirstUseAcceptsExistingSitesDefault(), TestFirstUseFallsBackToWorkingDirectory(), TestFunctionalCommandWithoutRootsFailsNonInteractively(), TestHelpAndVersionDoNotCreateSettings() (+24 more)

### Community 3 - "CLI Log Handling"
Cohesion: 0.08
Nodes (39): repeatedValue, notifyingWriter, TestOutputLogStreamsBeforeInputReachesEOF(), T, TestUpdateDelegatesToHomebrewOnlyForOwnedExecutable(), T, TestUpdateDoesNotDelegateToHomebrewForAnotherExecutable(), TestUpdateReplacesVerifiedReleaseForEverySupportedTarget() (+31 more)

### Community 13 - "CLI Command Tests"
Cohesion: 0.16
Nodes (32): TestVersionCommandsRenderTheToolVersion(), T, TestExitCodeMapsInterruptedOperationsTo130(), TestInitAllocatesPortsForExplicitServices(), TestInitRejectsInvalidGlobalRegistry(), TestInitRejectsDeprecatedAppFlag(), TestRunRejectsUnknownCommand(), TestRunRejectsRemovedWorkerCommand() (+24 more)

### Community 16 - "Init Interview"
Cohesion: 0.28
Nodes (15): collectInitInterview(), Reader, Writer, serviceDefinition, promptRequired(), promptServiceName(), promptDefault(), parseServiceDefinition() (+7 more)

### Community 19 - "Legacy Recovery Tests"
Cohesion: 0.46
Nodes (14): TestRecoverRequiresConfirmationOrYes(), T, TestRecoverRequiresYesForStaleLegacyState(), TestRecoverWithYesStopsLegacyProcessAndRemovesState(), TestRecoverInteractiveConfirmationStopsLegacyProcessAndRemovesState(), TestRecoverDeclinedConfirmationLeavesLegacyProcessAndState(), writeLegacyCLIRecoveryFixture(), recoveryEnvironment() (+6 more)

### Community 6 - "Configuration Persistence"
Cohesion: 0.09
Nodes (40): Config, Role, PortRange, Project, Runtime, Durations, Duration, Service (+32 more)

### Community 2 - "Update Installation Service"
Cohesion: 0.09
Nodes (23): Service, Context, Client, installation, Result, DefaultService(), runCommand(), runningBuildInfo() (+15 more)

### Community 12 - "Maintenance Uninstall"
Cohesion: 0.13
Nodes (23): installationKind, installation, uninstallArtifacts, artifactScanLimits, Service, Context, Store, Reader (+15 more)

### Community 17 - "Uninstall Tests"
Cohesion: 0.34
Nodes (16): TestUninstallDefaultYesRemovesOnlyRegisteredProjectArtifacts(), T, TestUninstallKeepsDeclinedArtifactClass(), TestUninstallRejectsNonInteractiveCleanup(), TestUninstallUsesOperationLockBeforePreflight(), TestUninstallAbortsBeforePromptsForActiveManagedService(), TestUninstallSkipsSymlinkedDirectoriesOutsideRegisteredRoots(), TestDiscoverUninstallArtifactsRejectsScanLimitOverrun() (+8 more)

### Community 14 - "Settings Path Handling"
Cohesion: 0.28
Nodes (5): Settings, Store, ConfigDirectory(), canonicalExistingPath(), canonicalExistingDirectory()

### Community 22 - "Linux Listener Discovery"
Cohesion: 0.25
Nodes (8): systemListener(), Listener, linuxListeningSocketInodes(), linuxSocketOwnerPIDs(), socketInode(), TestLinuxListenerParsingFindsOnlyTCPListenInodesForTheRequestedPort(), T, TestLinuxSocketOwnersFindsPIDsFromProcFileDescriptors()

### Community 23 - "Port Registry Lock"
Cohesion: 0.33
Nodes (9): WithRegistryLock(), Context, withRegistryLockIn(), TestRegistryLockUsesProvidedGratConfigurationDirectory(), T, TestRegistryLockHonorsContextWhileContended(), TestRegistryLockSerializesCallbacks(), TestRegistryLockReleasesAfterCallbackPanic() (+1 more)

### Community 10 - "Port Registry"
Cohesion: 0.12
Nodes (32): scanLimits, scanCounters, Source, Reservation, ProjectConfig, Config, Problem, Report (+24 more)

### Community 5 - "Terminal Presentation"
Cohesion: 0.09
Nodes (19): Command, CommandGroup, Renderer, helpUsageWidth(), Style, ColorMode, StepKind, Renderer (+11 more)

### Community 1 - "Terminal Lifecycle UI"
Cohesion: 0.09
Nodes (34): truncate(), LifecycleStage, LifecycleService, LifecycleGroup, LifecycleOperation, LifecycleEvent, LifecycleRunner, LifecycleModel (+26 more)

### Community 8 - "Presentation Tests"
Cohesion: 0.17
Nodes (37): NewLifecycleModel(), New(), DividerLine(), TestRendererUsesPlainTextForNonTerminalOutput(), T, TestRendererUsesSemanticColorWhenForced(), TestRendererSanitizesDynamicTerminalControlCharacters(), TestRendererSanitizesUnsafeUnicodeFormatCharacters() (+29 more)

### Community 25 - "Project Root Discovery"
Cohesion: 0.47
Nodes (4): FindRoot(), TestFindRootUsesNearestConfig(), T, TestFindRootReturnsNotFoundOutsideProject()

### Community 9 - "Runtime Manager Tests"
Cohesion: 0.12
Nodes (33): TestFixtureManagerUsesRuntimeBackendPortRange(), T, TestStartAndStopRequiresOwnedHealthyListener(), TestRestartEmitsOrderedLifecycleEvents(), TestStartRejectsUnhealthyHTTPResponse(), TestStartGracefullyStopsPreviouslyStartedServicesWhenCancelled(), TestStatusIgnoresLegacyPIDFiles(), TestLogTailReadsOnlyTheFinalWindow() (+25 more)

### Community 4 - "Runtime Service Manager"
Cohesion: 0.10
Nodes (43): Manager, Service, processState, commandEnvironment(), Context, loadedState, validateManagedState(), validateLegacyManagedState() (+35 more)

### Community 18 - "Runtime State Storage"
Cohesion: 0.23
Nodes (6): processState, Time, loadedState, RecoveryCandidate, Service, Manager

### Community 15 - "Settings Store Tests"
Cohesion: 0.31
Nodes (17): TestStoreLoadReportsMissingSettings(), T, TestStoreAddCanonicalizesAndDeduplicatesDirectories(), TestStoreAddResolvesRelativeDirectoriesAgainstWorkingDirectory(), TestStoreAddRejectsMissingAndNonDirectoryPaths(), TestStoreRemovePersistsRemainingDirectories(), TestStoreRejectsInvalidSettingsDocuments(), TestStoreSaveUsesRestrictivePermissions() (+9 more)

### Community 24 - "Text Sanitization"
Cohesion: 0.39
Nodes (6): UnsafeRune(), ContainsUnsafe(), Sanitize(), TestUnsafeRuneRejectsControlsAndUnicodeFormatCharacters(), T, TestSanitizeReplacesEveryUnsafeRune()

### Community 27 - "Bottle Packaging"
Cohesion: 0.70
Nodes (4): build-homebrew-bottles.sh script, usage(), write_formula(), package()

### Community 31 - "README Contract Checks"
Cohesion: 0.83
Nodes (3): check-readme.sh script, require(), require_in()

### Community 26 - "Bottle Package Tests"
Cohesion: 0.60
Nodes (5): test-homebrew-bottles.sh script, assert_file(), assert_archive_contains(), assert_binary(), assert_mode()

### Community 28 - "Bottle Verification"
Cohesion: 0.60
Nodes (3): verify-homebrew-bottles.sh script, usage(), verify_bottle()

### Community 29 - "CI Verification"
Cohesion: 0.40
Nodes (5): CI Workflow, Verify Job, Platform Verification Matrix, Race-Enabled Tests, Vulnerability Scan

### Community 20 - "Security Release Documentation"
Cohesion: 0.17
Nodes (15): Release Workflow, Tag-Triggered Release, Cross-Platform Build Job, GitHub Artifact Attestation, Release Publish Job, SHA-256 Release Checksums, GitHub Release Publication, Security Policy (+7 more)

### Community 32 - "Graphify Release Index"
Cohesion: 1.00
Nodes (3): Project Agent Rules, Graphify Before Push, Versioned Graphify Artifacts

### Community 21 - "Project Documentation"
Cohesion: 0.19
Nodes (13): grat README, grat, Declarative Configuration, Directory Discovery, Role-Based Port Allocation, HTTP Readiness, Process-Group Shutdown, Legacy Process Recovery (+5 more)

## Knowledge Gaps
- **88 isolated node(s):** `Store`, `updateService`, `uninstallService`, `updateService`, `uninstallService` (+83 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **8 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `New()` connect `Presentation Tests` to `CLI Command Dispatch`, `Update Installation Service`, `CLI Log Handling`, `Terminal Presentation`, `CLI Directory Setup`, `Application Runtime`, `Maintenance Uninstall`, `CLI Command Tests`, `Settings Path Handling`, `Settings Store Tests`, `Init Interview`, `Uninstall Tests`?**
  _High betweenness centrality (0.227) - this node is a cross-community bridge._
- **Why does `Contains()` connect `Presentation Tests` to `CLI Log Handling`, `Runtime Service Manager`, `Configuration Persistence`, `CLI Directory Setup`, `Runtime Manager Tests`, `Maintenance Uninstall`, `CLI Command Tests`, `Settings Path Handling`, `Settings Store Tests`, `Init Interview`, `Uninstall Tests`, `Legacy Recovery Tests`?**
  _High betweenness centrality (0.225) - this node is a cross-community bridge._
- **Why does `processAlive()` connect `Runtime Service Manager` to `CLI Log Handling`, `Runtime Manager Tests`, `Application Runtime`?**
  _High betweenness centrality (0.084) - this node is a cross-community bridge._
- **Are the 66 inferred relationships involving `Contains()` (e.g. with `TestInitRejectsDeprecatedAppFlag()` and `TestInitRejectsInvalidGlobalRegistry()`) actually correct?**
  _`Contains()` has 66 INFERRED edges - model-reasoned connections that need verification._
- **Are the 33 inferred relationships involving `New()` (e.g. with `configuredRoots()` and `runDirectories()`) actually correct?**
  _`New()` has 33 INFERRED edges - model-reasoned connections that need verification._
- **Are the 17 inferred relationships involving `runWithEnvironment()` (e.g. with `New()` and `Current()`) actually correct?**
  _`runWithEnvironment()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Store`, `updateService`, `uninstallService` to the rest of the system?**
  _93 weakly-connected nodes found - possible documentation gaps or missing edges._
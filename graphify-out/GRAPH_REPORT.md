# Graph Report - .  (2026-07-14)

## Corpus Check
- 34 files · ~43,468 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 803 nodes · 1817 edges · 43 communities (34 shown, 9 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 211 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_CLI Orchestration|CLI Orchestration]]
- [[_COMMUNITY_Runtime Readiness|Runtime Readiness]]
- [[_COMMUNITY_Runtime Terminal Models|Runtime Terminal Models]]
- [[_COMMUNITY_Update Maintenance|Update Maintenance]]
- [[_COMMUNITY_Configuration Validation|Configuration Validation]]
- [[_COMMUNITY_Uninstall Maintenance|Uninstall Maintenance]]
- [[_COMMUNITY_CLI Command Tests|CLI Command Tests]]
- [[_COMMUNITY_Maintenance CLI Tests|Maintenance CLI Tests]]
- [[_COMMUNITY_Terminal Presentation|Terminal Presentation]]
- [[_COMMUNITY_Port Registry|Port Registry]]
- [[_COMMUNITY_Presentation Tests|Presentation Tests]]
- [[_COMMUNITY_Configuration Tests|Configuration Tests]]
- [[_COMMUNITY_Runtime State Storage|Runtime State Storage]]
- [[_COMMUNITY_Runtime Process Manager|Runtime Process Manager]]
- [[_COMMUNITY_Settings Store|Settings Store]]
- [[_COMMUNITY_Security And Release Docs|Security And Release Docs]]
- [[_COMMUNITY_Settings Tests|Settings Tests]]
- [[_COMMUNITY_Interactive Setup|Interactive Setup]]
- [[_COMMUNITY_Uninstall Tests|Uninstall Tests]]
- [[_COMMUNITY_Project Documentation|Project Documentation]]
- [[_COMMUNITY_Help Rendering|Help Rendering]]
- [[_COMMUNITY_Linux Listener Inspection|Linux Listener Inspection]]
- [[_COMMUNITY_Operation Locking|Operation Locking]]
- [[_COMMUNITY_Runtime Progress Reporting|Runtime Progress Reporting]]
- [[_COMMUNITY_Unicode Text Safety|Unicode Text Safety]]
- [[_COMMUNITY_Project Root Discovery|Project Root Discovery]]
- [[_COMMUNITY_Bottle Packaging Tests|Bottle Packaging Tests]]
- [[_COMMUNITY_Bottle Packaging Script|Bottle Packaging Script]]
- [[_COMMUNITY_Bottle Verification|Bottle Verification]]
- [[_COMMUNITY_Project Agent Rules|Project Agent Rules]]
- [[_COMMUNITY_CLI Test Utilities|CLI Test Utilities]]
- [[_COMMUNITY_README Contract Script|README Contract Script]]
- [[_COMMUNITY_Support Documentation|Support Documentation]]
- [[_COMMUNITY_Darwin Listener Tests|Darwin Listener Tests]]
- [[_COMMUNITY_Readiness Tests|Readiness Tests]]
- [[_COMMUNITY_Listener Lookup|Listener Lookup]]
- [[_COMMUNITY_Code Of Conduct|Code Of Conduct]]
- [[_COMMUNITY_Homebrew Bottle Tooling|Homebrew Bottle Tooling]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
- [[_COMMUNITY_Bottle Verification Tests|Bottle Verification Tests]]
- [[_COMMUNITY_Help Tests|Help Tests]]
- [[_COMMUNITY_Settings Documentation|Settings Documentation]]

## God Nodes (most connected - your core abstractions)
1. `New()` - 32 edges
2. `runWithEnvironment()` - 30 edges
3. `Service` - 27 edges
4. `T` - 27 edges
5. `T` - 25 edges
6. `Renderer` - 22 edges
7. `runPortAssignLocked()` - 21 edges
8. `LifecycleModel` - 19 edges
9. `Manager` - 18 edges
10. `runWithConfiguredRoots()` - 18 edges

## Surprising Connections (you probably didn't know these)
- `Local Quality Gate` --semantically_similar_to--> `Verify Job`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → .github/workflows/ci.yml
- `Cross-Platform Compatibility` --semantically_similar_to--> `Platform Verification Matrix`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → .github/workflows/ci.yml
- `Approved Shell Command Contract` --semantically_similar_to--> `Shell Execution Trust Boundary`  [INFERRED] [semantically similar]
  README.md → SECURITY.md
- `Minimized Command Environment` --semantically_similar_to--> `Non-Secret Environment Baseline`  [INFERRED] [semantically similar]
  README.md → SECURITY.md
- `Release Provenance Verification` --semantically_similar_to--> `Fail-Closed Update Verification`  [INFERRED] [semantically similar]
  README.md → SECURITY.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Release Supply-Chain Verification** — workflows_release_artifact_attestation, workflows_release_sha256_checksums, workflows_release_github_release_publication, readme_release_provenance_verification, security_fail_closed_update_verification [INFERRED 0.95]

## Communities (43 total, 9 thin omitted)

### Community 0 - "CLI Orchestration"
Cohesion: 0.07
Nodes (73): assignReassignedPorts(), configuredRoots(), copyReservations(), defaultEnvironment(), detectServices(), ensureValidRegistry(), executeLifecycle(), fileExists() (+65 more)

### Community 1 - "Runtime Readiness"
Cohesion: 0.07
Nodes (49): Config, Listener, Manager, Service, T, T, Context, Duration (+41 more)

### Community 2 - "Runtime Terminal Models"
Cohesion: 0.10
Nodes (33): CancelFunc, Cmd, Context, Reader, Style, Writer, lifecycleRow, Model (+25 more)

### Community 3 - "Update Maintenance"
Cohesion: 0.10
Nodes (22): asset, installation, Client, Context, Service, Client, Context, Service (+14 more)

### Community 4 - "Configuration Validation"
Cohesion: 0.11
Nodes (38): configDecodeError(), DefaultRuntime(), InferRole(), Load(), readConfigFile(), replaceFile(), rollbackWrites(), safeEnvironmentName() (+30 more)

### Community 5 - "Uninstall Maintenance"
Cohesion: 0.11
Nodes (20): Context, Service, Reader, Result, Store, Writer, Manager, artifactScanLimits (+12 more)

### Community 6 - "CLI Command Tests"
Cohesion: 0.12
Nodes (30): repeatedValue, Context, Service, T, T, commandResponse, fakeCommands, commandKey() (+22 more)

### Community 7 - "Maintenance CLI Tests"
Cohesion: 0.15
Nodes (31): exitCode(), isHelp(), runWithEnvironment(), canonicalCLITestPath(), environmentForTest(), newCLITestStore(), sameStringSlices(), TestDirectoriesAddDoesNotPromptForInitialSetup() (+23 more)

### Community 8 - "Terminal Presentation"
Cohesion: 0.14
Nodes (14): ColorMode, formatProjectRows(), fprint(), fprintf(), fprintln(), pad(), ParseColorMode(), stepStyle() (+6 more)

### Community 9 - "Port Registry"
Cohesion: 0.13
Nodes (31): Config, T, fakeLookup, Listener, ListenerLookup, Problem, ProjectConfig, FirstFree() (+23 more)

### Community 10 - "Presentation Tests"
Cohesion: 0.15
Nodes (33): T, DividerLine(), New(), columnOf(), displayColumn(), stripANSI(), TestHelpRendersBorderlessGroupsWithGloballyAlignedDescriptions(), TestLifecycleModelAlignsHeadersWithDataColumns() (+25 more)

### Community 11 - "Configuration Tests"
Cohesion: 0.16
Nodes (32): assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort(), loadPortFixtureConfig(), runWithConfiguredRoots(), TestCLIRuntimeHelper(), TestExitCodeMapsInterruptedOperationsTo130() (+24 more)

### Community 12 - "Runtime State Storage"
Cohesion: 0.11
Nodes (19): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, File, T, Context, T, T, Mutex (+11 more)

### Community 13 - "Runtime Process Manager"
Cohesion: 0.17
Nodes (13): main(), mustGetwd(), Client, Config, Context, processState, ProgressObserver, Manager (+5 more)

### Community 14 - "Settings Store"
Cohesion: 0.21
Nodes (9): TestHelpListsProjectLifecycleAndPortCommandsWithoutWorker(), T, Settings, canonicalExistingDirectory(), canonicalExistingPath(), Contains(), Store.Save, Settings Safety Tests (+1 more)

### Community 15 - "Security And Release Docs"
Cohesion: 0.15
Nodes (21): Declarative grat.config Configuration, grat Documentation, Local Development Stack Orchestration, Installation-Aware Secure Update, Managed Process-Group Lifecycle, Minimized Command Environment, Release Provenance Verification, Role-Specific Readiness (+13 more)

### Community 16 - "Settings Tests"
Cohesion: 0.31
Nodes (17): Store, T, canonicalPath(), equalStrings(), newTestStore(), TestContainsAcceptsRegularFileBelowRoot(), TestContainsRejectsPathsOutsideRootAndThroughSymlinks(), TestStoreAddCanonicalizesAndDeduplicatesDirectories() (+9 more)

### Community 17 - "Interactive Setup"
Cohesion: 0.28
Nodes (15): collectInitInterview(), parseServiceDefinition(), promptDefault(), promptRequired(), promptServiceName(), TestCollectInitInterviewAcceptsEditsAndAdditionalServices(), TestCollectInitInterviewAcceptsSuppliedProjectName(), TestCollectInitInterviewRequiresAtLeastOneService() (+7 more)

### Community 18 - "Uninstall Tests"
Cohesion: 0.34
Nodes (16): Service, Store, T, fakeUninstallService(), newUninstallStore(), TestDiscoverUninstallArtifactsRejectsScanLimitOverrun(), TestUninstallAbortsBeforePromptsForActiveManagedService(), TestUninstallDefaultYesRemovesOnlyRegisteredProjectArtifacts() (+8 more)

### Community 19 - "Project Documentation"
Cohesion: 0.15
Nodes (13): Code of conduct, Configuration compatibility, Contributing to grat, Cross-Platform Compatibility, Development setup, Focused Pull Requests, Local Quality Gate, Pull requests (+5 more)

### Community 20 - "Help Rendering"
Cohesion: 0.25
Nodes (7): Listener, Renderer, Style, systemListener(), Command, CommandGroup, helpUsageWidth()

### Community 21 - "Linux Listener Inspection"
Cohesion: 0.25
Nodes (8): Listener, T, systemListener(), linuxListeningSocketInodes(), linuxSocketOwnerPIDs(), socketInode(), TestLinuxListenerParsingFindsOnlyTCPListenInodesForTheRequestedPort(), TestLinuxSocketOwnersFindsPIDsFromProcFileDescriptors()

### Community 22 - "Operation Locking"
Cohesion: 0.39
Nodes (7): Context, T, TestLockHonorsCanceledContextWhileContended(), TestLockSerializesCallbacks(), TestLockUsesRestrictivePermissions(), WithLock(), withLockIn()

### Community 23 - "Runtime Progress Reporting"
Cohesion: 0.33
Nodes (6): Manager, Service, ProgressEvent, ProgressObserver, ProgressObserverFunc, ProgressStage

### Community 24 - "Unicode Text Safety"
Cohesion: 0.39
Nodes (6): T, ContainsUnsafe(), Sanitize(), TestSanitizeReplacesEveryUnsafeRune(), TestUnsafeRuneRejectsControlsAndUnicodeFormatCharacters(), UnsafeRune()

### Community 25 - "Project Root Discovery"
Cohesion: 0.47
Nodes (4): T, FindRoot(), TestFindRootReturnsNotFoundOutsideProject(), TestFindRootUsesNearestConfig()

### Community 26 - "Bottle Packaging Tests"
Cohesion: 0.60
Nodes (5): assert_archive_contains(), assert_binary(), assert_file(), assert_mode(), test-homebrew-bottles.sh script

### Community 27 - "Bottle Packaging Script"
Cohesion: 0.70
Nodes (4): package(), usage(), write_formula(), build-homebrew-bottles.sh script

### Community 28 - "Bottle Verification"
Cohesion: 0.60
Nodes (3): usage(), verify_bottle(), verify-homebrew-bottles.sh script

### Community 29 - "Project Agent Rules"
Cohesion: 0.67
Nodes (3): Graphify Before Push, Project Agent Rules, Versioned Graphify Artifacts

### Community 30 - "CLI Test Utilities"
Cohesion: 0.67
Nodes (3): TestListenerOwnerLabelHandlesUnknownPID(), TestLogFollowUsesTrustedExecutable(), T

### Community 31 - "README Contract Script"
Cohesion: 0.83
Nodes (3): require(), require_in(), check-readme.sh script

### Community 32 - "Support Documentation"
Cohesion: 0.50
Nodes (3): Diagnostic Support Request, Sensitive Data Redaction, Support

## Knowledge Gaps
- **76 isolated node(s):** `T`, `T`, `Listener`, `T`, `Listener` (+71 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `New()` connect `Presentation Tests` to `CLI Orchestration`, `Update Maintenance`, `Uninstall Maintenance`, `CLI Command Tests`, `Maintenance CLI Tests`, `Terminal Presentation`, `Configuration Tests`, `Uninstall Tests`?**
  _High betweenness centrality (0.136) - this node is a cross-community bridge._
- **Why does `processAlive()` connect `Runtime Readiness` to `Update Maintenance`, `CLI Command Tests`?**
  _High betweenness centrality (0.050) - this node is a cross-community bridge._
- **Why does `runWithEnvironment()` connect `Maintenance CLI Tests` to `CLI Orchestration`, `Presentation Tests`, `Configuration Tests`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **Are the 27 inferred relationships involving `New()` (e.g. with `configuredRoots()` and `runDirectories()`) actually correct?**
  _`New()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **Are the 14 inferred relationships involving `runWithEnvironment()` (e.g. with `New()` and `Current()`) actually correct?**
  _`runWithEnvironment()` has 14 INFERRED edges - model-reasoned connections that need verification._
- **What connects `T`, `T`, `Listener` to the rest of the system?**
  _78 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `CLI Orchestration` be split into smaller, more focused modules?**
  _Cohesion score 0.06630630630630631 - nodes in this community are weakly interconnected._
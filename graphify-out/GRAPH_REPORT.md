# Graph Report - .  (2026-07-15)

## Corpus Check
- 11 files · ~46,815 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 852 nodes · 1971 edges · 44 communities (34 shown, 10 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 223 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_CLI Recovery|CLI Recovery]]
- [[_COMMUNITY_Presentation Runtime|Presentation Runtime]]
- [[_COMMUNITY_Runtime Manager Tests|Runtime Manager Tests]]
- [[_COMMUNITY_CLI Test Suite|CLI Test Suite]]
- [[_COMMUNITY_Configuration Validation|Configuration Validation]]
- [[_COMMUNITY_Maintenance System|Maintenance System]]
- [[_COMMUNITY_CLI Values And Uninstall Tests|CLI Values And Uninstall Tests]]
- [[_COMMUNITY_Terminal Presentation|Terminal Presentation]]
- [[_COMMUNITY_Port Registry|Port Registry]]
- [[_COMMUNITY_Presentation Tests|Presentation Tests]]
- [[_COMMUNITY_CLI Main And Runtime|CLI Main And Runtime]]
- [[_COMMUNITY_CLI Log Tests|CLI Log Tests]]
- [[_COMMUNITY_Uninstall Maintenance|Uninstall Maintenance]]
- [[_COMMUNITY_Help And Settings|Help And Settings]]
- [[_COMMUNITY_README Documentation|README Documentation]]
- [[_COMMUNITY_Settings Tests|Settings Tests]]
- [[_COMMUNITY_Interactive Setup|Interactive Setup]]
- [[_COMMUNITY_Runtime State Storage|Runtime State Storage]]
- [[_COMMUNITY_Recovery CLI Tests|Recovery CLI Tests]]
- [[_COMMUNITY_Contributing Documentation|Contributing Documentation]]
- [[_COMMUNITY_Darwin Ports And Help|Darwin Ports And Help]]
- [[_COMMUNITY_Linux Listener Inspection|Linux Listener Inspection]]
- [[_COMMUNITY_Runtime Progress|Runtime Progress]]
- [[_COMMUNITY_Unicode Text Safety|Unicode Text Safety]]
- [[_COMMUNITY_Project Root Discovery|Project Root Discovery]]
- [[_COMMUNITY_Bottle Packaging Tests|Bottle Packaging Tests]]
- [[_COMMUNITY_Bottle Packaging|Bottle Packaging]]
- [[_COMMUNITY_Bottle Verification|Bottle Verification]]
- [[_COMMUNITY_Project Agent Rules|Project Agent Rules]]
- [[_COMMUNITY_CLI Test Utilities|CLI Test Utilities]]
- [[_COMMUNITY_Linux Identity Tests|Linux Identity Tests]]
- [[_COMMUNITY_README Contract Script|README Contract Script]]
- [[_COMMUNITY_Support Documentation|Support Documentation]]
- [[_COMMUNITY_Darwin Listener Tests|Darwin Listener Tests]]
- [[_COMMUNITY_Readiness Tests|Readiness Tests]]
- [[_COMMUNITY_Listener Lookup|Listener Lookup]]
- [[_COMMUNITY_Linux Process Identity|Linux Process Identity]]
- [[_COMMUNITY_Code Of Conduct|Code Of Conduct]]
- [[_COMMUNITY_Homebrew Tooling|Homebrew Tooling]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
- [[_COMMUNITY_Verification Tests|Verification Tests]]
- [[_COMMUNITY_Help Contract Tests|Help Contract Tests]]
- [[_COMMUNITY_Settings Store|Settings Store]]

## God Nodes (most connected - your core abstractions)
1. `runWithEnvironment()` - 33 edges
2. `New()` - 32 edges
3. `Service` - 27 edges
4. `T` - 27 edges
5. `T` - 25 edges
6. `Manager` - 23 edges
7. `Renderer` - 22 edges
8. `runPortAssignLocked()` - 21 edges
9. `LifecycleModel` - 19 edges
10. `Renderer` - 19 edges

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

## Communities (44 total, 10 thin omitted)

### Community 0 - "CLI Recovery"
Cohesion: 0.06
Nodes (85): assignReassignedPorts(), configuredRoots(), confirmRecovery(), copyReservations(), defaultEnvironment(), detectServices(), ensureValidRegistry(), executeLifecycle() (+77 more)

### Community 1 - "Presentation Runtime"
Cohesion: 0.08
Nodes (48): CancelFunc, Cmd, Context, Reader, Style, Writer, Manager, Service (+40 more)

### Community 2 - "Runtime Manager Tests"
Cohesion: 0.07
Nodes (51): Duration, Config, Listener, Manager, Service, T, T, Context (+43 more)

### Community 3 - "CLI Test Suite"
Cohesion: 0.09
Nodes (55): assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort(), loadPortFixtureConfig(), runWithConfiguredRoots(), TestCLIRuntimeHelper(), TestExitCodeMapsInterruptedOperationsTo130() (+47 more)

### Community 4 - "Configuration Validation"
Cohesion: 0.09
Nodes (45): configDecodeError(), DefaultRuntime(), InferRole(), Load(), readConfigFile(), replaceFile(), rollbackWrites(), safeEnvironmentName() (+37 more)

### Community 5 - "Maintenance System"
Cohesion: 0.10
Nodes (22): asset, installation, Client, Context, Service, Client, Context, Service (+14 more)

### Community 6 - "CLI Values And Uninstall Tests"
Cohesion: 0.11
Nodes (41): repeatedValue, Service, Store, T, Context, Service, T, commandResponse (+33 more)

### Community 7 - "Terminal Presentation"
Cohesion: 0.14
Nodes (14): ColorMode, formatProjectRows(), fprint(), fprintf(), fprintln(), pad(), ParseColorMode(), stepStyle() (+6 more)

### Community 8 - "Port Registry"
Cohesion: 0.13
Nodes (31): Config, T, fakeLookup, Listener, ListenerLookup, Problem, ProjectConfig, FirstFree() (+23 more)

### Community 9 - "Presentation Tests"
Cohesion: 0.15
Nodes (33): T, DividerLine(), New(), columnOf(), displayColumn(), stripANSI(), TestHelpRendersBorderlessGroupsWithGloballyAlignedDescriptions(), TestLifecycleModelAlignsHeadersWithDataColumns() (+25 more)

### Community 10 - "CLI Main And Runtime"
Cohesion: 0.15
Nodes (16): Client, main(), mustGetwd(), Config, Context, loadedState, processState, ProgressObserver (+8 more)

### Community 11 - "CLI Log Tests"
Cohesion: 0.11
Nodes (19): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, File, T, Context, T, T, Mutex (+11 more)

### Community 12 - "Uninstall Maintenance"
Cohesion: 0.18
Nodes (16): Context, Service, Reader, Result, Store, Writer, artifactScanLimits, installation (+8 more)

### Community 13 - "Help And Settings"
Cohesion: 0.21
Nodes (9): TestHelpListsProjectLifecycleAndPortCommandsWithoutWorker(), T, Settings, canonicalExistingDirectory(), canonicalExistingPath(), Contains(), Store.Save, Settings Safety Tests (+1 more)

### Community 14 - "README Documentation"
Cohesion: 0.13
Nodes (23): README, Declarative grat.config Configuration, grat Documentation, Local Development Stack Orchestration, Installation-Aware Secure Update, Managed Process-Group Lifecycle, Minimized Command Environment, Legacy Process Recovery (+15 more)

### Community 15 - "Settings Tests"
Cohesion: 0.31
Nodes (17): Store, T, canonicalPath(), equalStrings(), newTestStore(), TestContainsAcceptsRegularFileBelowRoot(), TestContainsRejectsPathsOutsideRootAndThroughSymlinks(), TestStoreAddCanonicalizesAndDeduplicatesDirectories() (+9 more)

### Community 16 - "Interactive Setup"
Cohesion: 0.28
Nodes (15): collectInitInterview(), parseServiceDefinition(), promptDefault(), promptRequired(), promptServiceName(), TestCollectInitInterviewAcceptsEditsAndAdditionalServices(), TestCollectInitInterviewAcceptsSuppliedProjectName(), TestCollectInitInterviewRequiresAtLeastOneService() (+7 more)

### Community 17 - "Runtime State Storage"
Cohesion: 0.23
Nodes (6): Manager, Service, loadedState, processState, RecoveryCandidate, Time

### Community 18 - "Recovery CLI Tests"
Cohesion: 0.46
Nodes (14): assertCLIRecoveryState(), assertRecoveryPreview(), cliProcessAlive(), legacyCLIStartIdentity(), recoveryEnvironment(), stopCLIRecoveryGroup(), TestRecoverDeclinedConfirmationLeavesLegacyProcessAndState(), TestRecoverInteractiveConfirmationStopsLegacyProcessAndRemovesState() (+6 more)

### Community 19 - "Contributing Documentation"
Cohesion: 0.15
Nodes (13): Code of conduct, Configuration compatibility, Contributing to grat, Cross-Platform Compatibility, Development setup, Focused Pull Requests, Local Quality Gate, Pull requests (+5 more)

### Community 20 - "Darwin Ports And Help"
Cohesion: 0.25
Nodes (7): Listener, Renderer, Style, systemListener(), Command, CommandGroup, helpUsageWidth()

### Community 21 - "Linux Listener Inspection"
Cohesion: 0.25
Nodes (8): Listener, T, systemListener(), linuxListeningSocketInodes(), linuxSocketOwnerPIDs(), socketInode(), TestLinuxListenerParsingFindsOnlyTCPListenInodesForTheRequestedPort(), TestLinuxSocketOwnersFindsPIDsFromProcFileDescriptors()

### Community 22 - "Runtime Progress"
Cohesion: 0.33
Nodes (6): Manager, Service, ProgressEvent, ProgressObserver, ProgressObserverFunc, ProgressStage

### Community 23 - "Unicode Text Safety"
Cohesion: 0.39
Nodes (6): T, ContainsUnsafe(), Sanitize(), TestSanitizeReplacesEveryUnsafeRune(), TestUnsafeRuneRejectsControlsAndUnicodeFormatCharacters(), UnsafeRune()

### Community 24 - "Project Root Discovery"
Cohesion: 0.47
Nodes (4): T, FindRoot(), TestFindRootReturnsNotFoundOutsideProject(), TestFindRootUsesNearestConfig()

### Community 25 - "Bottle Packaging Tests"
Cohesion: 0.60
Nodes (5): assert_archive_contains(), assert_binary(), assert_file(), assert_mode(), test-homebrew-bottles.sh script

### Community 26 - "Bottle Packaging"
Cohesion: 0.70
Nodes (4): package(), usage(), write_formula(), build-homebrew-bottles.sh script

### Community 27 - "Bottle Verification"
Cohesion: 0.60
Nodes (3): usage(), verify_bottle(), verify-homebrew-bottles.sh script

### Community 28 - "Project Agent Rules"
Cohesion: 0.67
Nodes (3): Graphify Before Push, Project Agent Rules, Versioned Graphify Artifacts

### Community 29 - "CLI Test Utilities"
Cohesion: 0.67
Nodes (3): TestListenerOwnerLabelHandlesUnknownPID(), TestLogFollowUsesTrustedExecutable(), T

### Community 30 - "Linux Identity Tests"
Cohesion: 0.67
Nodes (3): T, TestLinuxProcessStartTicksHandlesClosingParenthesisInCommand(), TestLinuxProcessStartTicksRejectsMismatchedPID()

### Community 31 - "README Contract Script"
Cohesion: 0.83
Nodes (3): require(), require_in(), check-readme.sh script

### Community 32 - "Support Documentation"
Cohesion: 0.50
Nodes (3): Diagnostic Support Request, Sensitive Data Redaction, Support

## Knowledge Gaps
- **79 isolated node(s):** `T`, `T`, `Listener`, `T`, `Listener` (+74 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **10 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `legacyProcessIdentity()` connect `Presentation Runtime` to `Runtime Manager Tests`, `CLI Values And Uninstall Tests`?**
  _High betweenness centrality (0.184) - this node is a cross-community bridge._
- **Why does `New()` connect `Presentation Tests` to `CLI Recovery`, `CLI Test Suite`, `Maintenance System`, `CLI Values And Uninstall Tests`, `Terminal Presentation`, `Uninstall Maintenance`?**
  _High betweenness centrality (0.179) - this node is a cross-community bridge._
- **Are the 16 inferred relationships involving `runWithEnvironment()` (e.g. with `New()` and `Current()`) actually correct?**
  _`runWithEnvironment()` has 16 INFERRED edges - model-reasoned connections that need verification._
- **Are the 27 inferred relationships involving `New()` (e.g. with `configuredRoots()` and `runDirectories()`) actually correct?**
  _`New()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **What connects `T`, `T`, `Listener` to the rest of the system?**
  _82 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `CLI Recovery` be split into smaller, more focused modules?**
  _Cohesion score 0.05898876404494382 - nodes in this community are weakly interconnected._
- **Should `Presentation Runtime` be split into smaller, more focused modules?**
  _Cohesion score 0.07550482879719052 - nodes in this community are weakly interconnected._
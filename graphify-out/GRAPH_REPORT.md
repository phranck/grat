# Graph Report - .  (2026-09-03)

## Corpus Check
- 25 files · ~92,645 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1497 nodes · 3527 edges · 77 communities (62 shown, 15 thin omitted)
- Extraction: 79% EXTRACTED · 21% INFERRED · 0% AMBIGUOUS · INFERRED: 753 edges (avg confidence: 0.8)
- Token cost: 102,007 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Maintenance Test Doubles|Maintenance Test Doubles]]
- [[_COMMUNITY_CLI Command Wiring|CLI Command Wiring]]
- [[_COMMUNITY_Uninstall Tests|Uninstall Tests]]
- [[_COMMUNITY_Lifecycle Terminal View|Lifecycle Terminal View]]
- [[_COMMUNITY_Port Assignment Commands|Port Assignment Commands]]
- [[_COMMUNITY_Configuration Loading|Configuration Loading]]
- [[_COMMUNITY_Lifecycle Commands|Lifecycle Commands]]
- [[_COMMUNITY_Update Command|Update Command]]
- [[_COMMUNITY_Update and Uninstall Service|Update and Uninstall Service]]
- [[_COMMUNITY_Manual Page Model|Manual Page Model]]
- [[_COMMUNITY_CLI Integration Tests|CLI Integration Tests]]
- [[_COMMUNITY_Website and Landing Page|Website and Landing Page]]
- [[_COMMUNITY_Manual Rendering|Manual Rendering]]
- [[_COMMUNITY_Detector Tests|Detector Tests]]
- [[_COMMUNITY_Runtime Manager|Runtime Manager]]
- [[_COMMUNITY_Maintenance Seams|Maintenance Seams]]
- [[_COMMUNITY_Discovery Interview|Discovery Interview]]
- [[_COMMUNITY_Project Discovery|Project Discovery]]
- [[_COMMUNITY_Manual Document Builder|Manual Document Builder]]
- [[_COMMUNITY_Runtime Manager Types|Runtime Manager Types]]
- [[_COMMUNITY_Selection List View|Selection List View]]
- [[_COMMUNITY_Bounded Project Walk|Bounded Project Walk]]
- [[_COMMUNITY_Gates and CI|Gates and CI]]
- [[_COMMUNITY_Backend URL Injection|Backend URL Injection]]
- [[_COMMUNITY_Settings Store Tests|Settings Store Tests]]
- [[_COMMUNITY_Site Badge and Release Order|Site Badge and Release Order]]
- [[_COMMUNITY_Process Identity|Process Identity]]
- [[_COMMUNITY_Managed State Paths|Managed State Paths]]
- [[_COMMUNITY_Go Module Detection|Go Module Detection]]
- [[_COMMUNITY_Release Process|Release Process]]
- [[_COMMUNITY_Python and Django Detection|Python and Django Detection]]
- [[_COMMUNITY_Lifecycle Model|Lifecycle Model]]
- [[_COMMUNITY_Readiness Probing|Readiness Probing]]
- [[_COMMUNITY_Generated Documentation Gate|Generated Documentation Gate]]
- [[_COMMUNITY_Detector Dispatch|Detector Dispatch]]
- [[_COMMUNITY_Installation and Verification|Installation and Verification]]
- [[_COMMUNITY_Service Log Files|Service Log Files]]
- [[_COMMUNITY_Linux Listener Lookup|Linux Listener Lookup]]
- [[_COMMUNITY_Laravel Detection|Laravel Detection]]
- [[_COMMUNITY_Homebrew Tap Release|Homebrew Tap Release]]
- [[_COMMUNITY_Site Badge Rules|Site Badge Rules]]
- [[_COMMUNITY_JavaScript Framework Detection|JavaScript Framework Detection]]
- [[_COMMUNITY_Package Manifest Reading|Package Manifest Reading]]
- [[_COMMUNITY_Node Service Detection|Node Service Detection]]
- [[_COMMUNITY_Node Server Port Reading|Node Server Port Reading]]
- [[_COMMUNITY_Legacy Process Recovery|Legacy Process Recovery]]
- [[_COMMUNITY_Linux Process Identity|Linux Process Identity]]
- [[_COMMUNITY_Homebrew Bottle Packaging|Homebrew Bottle Packaging]]
- [[_COMMUNITY_Bottle Packaging Tests|Bottle Packaging Tests]]
- [[_COMMUNITY_Contributing Guide|Contributing Guide]]
- [[_COMMUNITY_Log Streaming Tests|Log Streaming Tests]]
- [[_COMMUNITY_Port Range Tests|Port Range Tests]]
- [[_COMMUNITY_Bottle Verification|Bottle Verification]]
- [[_COMMUNITY_Agent Rules|Agent Rules]]
- [[_COMMUNITY_CLI Helper Tests|CLI Helper Tests]]
- [[_COMMUNITY_Manual Reference Tests|Manual Reference Tests]]
- [[_COMMUNITY_Gate Script|Gate Script]]
- [[_COMMUNITY_Support Policy|Support Policy]]
- [[_COMMUNITY_Port Ownership Detection|Port Ownership Detection]]
- [[_COMMUNITY_Program Entry Point|Program Entry Point]]
- [[_COMMUNITY_macOS Listener Lookup|macOS Listener Lookup]]
- [[_COMMUNITY_macOS Listener Tests|macOS Listener Tests]]
- [[_COMMUNITY_Process Inspection Tests|Process Inspection Tests]]
- [[_COMMUNITY_System Listener Lookup|System Listener Lookup]]
- [[_COMMUNITY_Favicon Build|Favicon Build]]
- [[_COMMUNITY_Code of Conduct|Code of Conduct]]
- [[_COMMUNITY_Social Card Build|Social Card Build]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
- [[_COMMUNITY_Release Notes Script|Release Notes Script]]
- [[_COMMUNITY_Bottle Verification Tests|Bottle Verification Tests]]
- [[_COMMUNITY_Configuration Compatibility|Configuration Compatibility]]
- [[_COMMUNITY_Colour Options|Colour Options]]
- [[_COMMUNITY_Configuration Manual Page|Configuration Manual Page]]

## God Nodes (most connected - your core abstractions)
1. `join()` - 134 edges
2. `Contains()` - 104 edges
3. `New()` - 49 edges
4. `runWithEnvironment()` - 40 edges
5. `Directory()` - 40 edges
6. `writeProject()` - 32 edges
7. `T` - 30 edges
8. `T` - 30 edges
9. `Service` - 29 edges
10. `T` - 26 edges

## Surprising Connections (you probably didn't know these)
- `SoftwareApplication JSON-LD block` --semantically_similar_to--> `Release Publish Job`  [INFERRED] [semantically similar]
  docs/index.html → .github/workflows/release.yml
- `Artifact attestation verification with gh` --semantically_similar_to--> `gosec Security Scan`  [INFERRED] [semantically similar]
  Documentation.md → .github/workflows/ci.yml
- `README heading and phrase checks` --conceptually_related_to--> `Install section`  [INFERRED]
  scripts/check-docs.sh → docs/index.html
- `The Nine-Step Release Order Across Two Repositories` --implements--> `check-docs.sh, the documentation contract gate`  [INFERRED]
  .claude/skills/release/SKILL.md → scripts/check-docs.sh
- `Local Development Gate` --semantically_similar_to--> `Verify Job`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → .github/workflows/ci.yml

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Readiness is three conditions judged together** — documentation_readiness, documentation_port_variable, documentation_owned_listener_rationale, documentation_service_states, documentation_foreground_requirement, docs_index_process_tree_boundary [EXTRACTED 1.00]
- **Installing and updating grat is gated on provenance** — documentation_installation, documentation_attestation_verification, documentation_grat_update, security_fail_closed_provenance, security_path_resolved_helpers, readme_verify_before_install [EXTRACTED 1.00]
- **What a managed command receives and nothing more** — documentation_environment_baseline, documentation_inherit_env, documentation_port_variable, documentation_backend_url, security_secret_propagation, docs_index_small_environment [EXTRACTED 1.00]

## Communities (77 total, 15 thin omitted)

### Community 0 - "Maintenance Test Doubles"
Cohesion: 0.07
Nodes (57): Buffer, Service, Store, T, T, T, runningProject(), TestAFailingStopIsReportedRatherThanIgnored() (+49 more)

### Community 1 - "CLI Command Wiring"
Cohesion: 0.07
Nodes (61): defaultEnvironment(), exitCode(), Run(), runWithEnvironment(), configuredRoots(), runDirectories(), canonicalCLITestPath(), environmentForTest() (+53 more)

### Community 2 - "Uninstall Tests"
Cohesion: 0.08
Nodes (61): join(), Service, Store, T, Server, T, Context, Service (+53 more)

### Community 3 - "Lifecycle Terminal View"
Cohesion: 0.07
Nodes (47): CancelFunc, Cmd, Context, Reader, Style, Writer, Cmd, Context (+39 more)

### Community 4 - "Port Assignment Commands"
Cohesion: 0.07
Nodes (59): portReassignment, assignReassignedPorts(), copyReservations(), ensureValidRegistry(), hasConfiguredCollision(), listenerOwnerLabel(), newPortReassignLifecycleOperation(), portReassignRowKey() (+51 more)

### Community 5 - "Configuration Loading"
Cohesion: 0.07
Nodes (53): Config, configDecodeError(), DefaultRuntime(), Load(), readConfigFile(), replaceFile(), rollbackWrites(), safeEnvironmentName() (+45 more)

### Community 6 - "Lifecycle Commands"
Cohesion: 0.05
Nodes (51): loadConfig(), loadManager(), executeLifecycle(), lifecycleTitle(), lifecycleTUIStage(), newLifecycleOperation(), progressPresentation(), runLifecycle() (+43 more)

### Community 7 - "Update Command"
Cohesion: 0.08
Nodes (28): runUpdate(), Context, Renderer, updateService, Renderer, Style, Context, Writer (+20 more)

### Community 8 - "Update and Uninstall Service"
Cohesion: 0.09
Nodes (26): asset, installation, T, Client, Context, Service, Client, Context (+18 more)

### Community 9 - "Manual Page Model"
Cohesion: 0.08
Nodes (50): commandDetail, Roles(), Block, Document, T, Block, CommandGroup, Document (+42 more)

### Community 10 - "CLI Integration Tests"
Cohesion: 0.12
Nodes (47): assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort(), loadPortFixtureConfig(), runWithConfiguredRoots(), TestCLIRuntimeHelper(), TestDiscoverAllocatesPortsForExplicitServices() (+39 more)

### Community 11 - "Website and Landing Page"
Cohesion: 0.06
Nodes (49): Written commands carry host and port, because binding every interface is accidental exposure, A Docker-held port belongs to Docker, so it never becomes ready, Hero section with a grat start terminal transcript, Ports section with the role lane scale, The listener must trace back to the command grat started, Reduced-motion respecting in-page scrolling, Services section: start, stop, logs cards, Ten inherited variables plus inherit_env (+41 more)

### Community 12 - "Manual Rendering"
Cohesion: 0.07
Nodes (38): commandDocument(), runManual(), plainManual(), TestBothManualPagesAreReachable(), TestEveryCommandOfTheReferenceHasAManualEntry(), TestTheManualCarriesEveryCommandOfTheReference(), TestTheManualSaysWhyAProjectCanBeRefused(), writeMarkdownManual() (+30 more)

### Community 13 - "Detector Tests"
Cohesion: 0.16
Nodes (45): Directory(), commandOf(), TestAFrameworkWinsOverTheBuildToolUnderIt(), TestAGoLibraryIsReportedRatherThanInvented(), TestAGoProgramThatIgnoresThePortIsReportedRatherThanOffered(), TestAMentionOfThePortIsNotAReadOfIt(), TestAnEmptyDirectoryIsNotAProject(), TestAnExpressServerNeedsThePortInItsSource() (+37 more)

### Community 14 - "Runtime Manager"
Cohesion: 0.11
Nodes (37): Config, Listener, Manager, Service, T, T, Manager, Service (+29 more)

### Community 15 - "Maintenance Seams"
Cohesion: 0.13
Nodes (24): Context, Reader, Result, Store, Writer, Context, T, activeProject (+16 more)

### Community 16 - "Discovery Interview"
Cohesion: 0.10
Nodes (31): detectServices(), discoverHere(), serviceSuggestions(), collectProjectInterview(), parseServiceDefinition(), promptDefault(), promptRequired(), promptServiceName() (+23 more)

### Community 17 - "Project Discovery"
Cohesion: 0.15
Nodes (33): candidate, candidate, absoluteUnder(), allocateServices(), candidateDetail(), chooseCandidates(), discoverBelow(), discoverCandidates() (+25 more)

### Community 18 - "Manual Document Builder"
Cohesion: 0.15
Nodes (14): Item, Document, Item, builder, flatten(), paragraphsOf(), anchor(), Markdown() (+6 more)

### Community 19 - "Runtime Manager Types"
Cohesion: 0.18
Nodes (12): Client, Config, Context, ProgressObserver, Manager, Service, ListenerLookup, readiness (+4 more)

### Community 20 - "Selection List View"
Cohesion: 0.14
Nodes (11): T, NewSelectionModel(), rows(), TestALongListSaysWhatIsOutOfSight(), TestARowThatCannotBeChosenNeverIs(), TestCancellingChoosesNothing(), TestOnlyMarkedRowsComeBack(), TestTheCursorStopsAtBothEndsRatherThanWrapping() (+3 more)

### Community 21 - "Bounded Project Walk"
Cohesion: 0.18
Nodes (18): ErrTooManyEntries, DirEntry, T, T, ErrTooManyEntries, DeeperThanScan(), SkipsScanning(), TestTheScanSkipsWhatCannotHoldAProject() (+10 more)

### Community 22 - "Gates and CI"
Cohesion: 0.12
Nodes (21): GitHub Actions Pinned by Commit Hash, Dependabot Weekly Updates, Local Development Gate, macOS and Linux Platform Parity Requirement, Build gate (cmd/grat), Build-tag coverage of static analysis, CI Workflow, Documentation gate (check-docs.sh) (+13 more)

### Community 23 - "Backend URL Injection"
Cohesion: 0.20
Nodes (16): Context, Duration, loadedState, processState, T, TestProcessIdentitySeparatesRapidProcessStarts(), TestSignalManagedGroupRejectsChangedIdentity(), TestValidateLegacyManagedStateAcceptsDetachedLegacyProcess() (+8 more)

### Community 24 - "Settings Store Tests"
Cohesion: 0.31
Nodes (17): Store, T, canonicalPath(), equalStrings(), newTestStore(), TestContainsAcceptsRegularFileBelowRoot(), TestContainsRejectsPathsOutsideRootAndThroughSymlinks(), TestStoreAddCanonicalizesAndDeduplicatesDirectories() (+9 more)

### Community 25 - "Site Badge and Release Order"
Cohesion: 0.16
Nodes (18): The window where the site names a release nobody can download yet, Why the Badge Is Raised Before the Tag Exists, Step one: commit the raised site badge on a release branch, GITHUB_TOKEN Raises No Workflow Run, The site needs no separate deployment; the badge merge is the deployment, The Site Needs No Separate Deployment, Put the Released Version in the Badge, Pages Concurrency Group Without Cancellation (+10 more)

### Community 26 - "Process Identity"
Cohesion: 0.38
Nodes (16): Cmd, Manager, Service, T, processAlive(), newLegacyRecoveryFixture(), stopFixtureGroup(), TestRecoverRejectsChangedNativeSnapshotIdentityWithoutSignaling() (+8 more)

### Community 27 - "Managed State Paths"
Cohesion: 0.23
Nodes (6): Manager, Service, loadedState, processState, RecoveryCandidate, Time

### Community 28 - "Go Module Detection"
Cohesion: 0.19
Nodes (14): entries(), detectGo(), goPrograms(), goSourceFile(), holdsMainPackage(), isPortEnvironmentCall(), readsPortFromGoSource(), TestTheGoScanStaysInsideTheBoundsEveryOtherScanHas() (+6 more)

### Community 29 - "Release Process"
Cohesion: 0.19
Nodes (15): Release-note Commit Trailer, Releasing grat (release skill), grat Release Process, Version Bump Decision from Commit Prefixes, Package Homebrew bottles, fetch-depth 0 so every tag is present for note generation, GitHub Release Publication, Manual Pages Generated by the Built Binary (+7 more)

### Community 30 - "Python and Django Detection"
Cohesion: 0.16
Nodes (11): readBounded(), detectDjango(), applicationModules(), detectPython(), detectVapor(), Service, Unresolved, Service (+3 more)

### Community 31 - "Lifecycle Model"
Cohesion: 0.19
Nodes (9): Cmd, Context, Model, Msg, Reader, Writer, RunSelection(), selectionTeaModel (+1 more)

### Community 32 - "Readiness Probing"
Cohesion: 0.22
Nodes (11): Context, processState, Manager, Service, legacyProcessIdentity(), parentProcessID(), processAlive(), psField() (+3 more)

### Community 33 - "Generated Documentation Gate"
Cohesion: 0.19
Nodes (13): Manual pages and Documentation.md are generated, not written, Why every documented claim is checked here, Documentation.md diffed against the manual the binary renders, Why a generated document replaces phrase-by-phrase assertions, Go 1.25.13 pin checked in go.mod, README and CONTRIBUTING, Both manual pages must render as valid roff (mandoc lint), README must open with self-updating shields.io badges, require() (+5 more)

### Community 34 - "Detector Dispatch"
Cohesion: 0.18
Nodes (10): InferRole(), fileExists(), detector, Finding, detectRails(), Service, Unresolved, Role (+2 more)

### Community 35 - "Installation and Verification"
Cohesion: 0.22
Nodes (13): Artifact attestation verification with gh, Installation routes (Homebrew, release binary, go install), README installation instructions, Verifying the binary before installing it, not after, CI and release workflow literals asserted from the docs gate, Fail-closed provenance checks for update and direct install, brew and gh resolved through PATH, GitHub Artifact Attestation (+5 more)

### Community 36 - "Service Log Files"
Cohesion: 0.26
Nodes (6): File, T, Manager, Service, newServiceLogFile(), TestNewServiceLogFileTruncatesPreviousOutput()

### Community 37 - "Linux Listener Lookup"
Cohesion: 0.25
Nodes (8): T, Listener, systemListener(), ListeningSocketInodes(), socketInode(), SocketOwnerPIDs(), TestListeningSocketInodesFindsOnlyTCPListenInodesForTheRequestedPort(), TestSocketOwnerPIDsFindsPIDsFromProcFileDescriptors()

### Community 38 - "Laravel Detection"
Cohesion: 0.42
Nodes (8): configuredQueueFallback(), detectLaravel(), environmentAssignment(), environmentValue(), laravelQueueConnection(), laravelQueueWorker(), Service, Unresolved

### Community 39 - "Homebrew Tap Release"
Cohesion: 0.31
Nodes (9): Install section, Homebrew Tap Formula Update, Update the Homebrew tap formula with seven values, The Nine-Step Release Order Across Two Repositories, Why the tap must not be merged before its checks finish, Why the Tap Pull Request Waits for All Three Checks, README heading and phrase checks, Create checksums.txt (+1 more)

### Community 40 - "Site Badge Rules"
Cohesion: 0.29
Nodes (8): Badge and Tag Cannot Both Be First, A Pages deployment from a tag is accepted, marked active and never served, Why a newer badge is allowed and an older one is not, Site badge compared against the newest tag, refusing an older one, Why the badge is whatever main holds rather than written at publish time, Push to main as the only trigger, plus manual dispatch, Tag-Triggered Release, Release Workflow

### Community 41 - "JavaScript Framework Detection"
Cohesion: 0.43
Nodes (6): detectJavaScriptFramework(), frameworkEvidence(), jsFramework, manifest, Service, Unresolved

### Community 43 - "Node Service Detection"
Cohesion: 0.57
Nodes (6): detectNode(), detectSingleService(), namedServices(), manifest, Service, Unresolved

### Community 44 - "Node Server Port Reading"
Cohesion: 0.43
Nodes (6): detectNodeServer(), readsPortFromEnvironment(), startScript(), manifest, Service, Unresolved

### Community 45 - "Legacy Process Recovery"
Cohesion: 0.52
Nodes (4): loadedState, processState, RecoveryCandidate, validateRecoveryCandidate()

### Community 46 - "Linux Process Identity"
Cohesion: 0.43
Nodes (5): T, linuxProcessStartTicks(), processIdentity(), TestLinuxProcessStartTicksHandlesClosingParenthesisInCommand(), TestLinuxProcessStartTicksRejectsMismatchedPID()

### Community 47 - "Homebrew Bottle Packaging"
Cohesion: 0.48
Nodes (5): package(), usage(), write_formula(), write_manual(), build-homebrew-bottles.sh script

### Community 48 - "Bottle Packaging Tests"
Cohesion: 0.52
Nodes (6): assert_archive_contains(), assert_binary(), assert_file(), assert_manual(), assert_mode(), test-homebrew-bottles.sh script

### Community 49 - "Contributing Guide"
Cohesion: 0.33
Nodes (5): Code of conduct, Configuration compatibility, Contributing to grat, Development setup, Pull requests

### Community 50 - "Log Streaming Tests"
Cohesion: 0.50
Nodes (3): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, T

### Community 51 - "Port Range Tests"
Cohesion: 0.60
Nodes (4): TestAWorkerOwnsNoRange(), TestEveryRoleOwnsARangeOfTheStatedWidth(), TestRolesWithDifferentRangesDoNotOverlap(), T

### Community 52 - "Bottle Verification"
Cohesion: 0.60
Nodes (3): usage(), verify_bottle(), verify-homebrew-bottles.sh script

### Community 53 - "Agent Rules"
Cohesion: 0.67
Nodes (3): Graphify Before Push, Project Agent Rules, Versioned Graphify Artifacts

### Community 54 - "CLI Helper Tests"
Cohesion: 0.67
Nodes (3): TestListenerOwnerLabelHandlesUnknownPID(), TestLogFollowUsesTrustedExecutable(), T

### Community 55 - "Manual Reference Tests"
Cohesion: 0.67
Nodes (3): T, TestAnEntryCarriesMoreThanTheOneLineReference(), TestEveryFlagOfAnEntryIsExplained()

### Community 56 - "Gate Script"
Cohesion: 0.67
Nodes (3): GOTOOLCHAIN, step(), gates.sh script

### Community 57 - "Support Policy"
Cohesion: 0.50
Nodes (3): Diagnostic Support Request, Sensitive Data Redaction, Support

### Community 58 - "Port Ownership Detection"
Cohesion: 0.67
Nodes (3): Who decides the port, Reading process.env.PORT and os.Getenv("PORT") from project source, Recognised frameworks on frontend and backend

## Knowledge Gaps
- **188 isolated node(s):** `Reader`, `Store`, `Manager`, `Config`, `StepKind` (+183 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **15 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `join()` connect `Uninstall Tests` to `Maintenance Test Doubles`, `CLI Command Wiring`, `Lifecycle Terminal View`, `Port Assignment Commands`, `Configuration Loading`, `Lifecycle Commands`, `Update Command`, `Update and Uninstall Service`, `Manual Page Model`, `CLI Integration Tests`, `Manual Rendering`, `Detector Tests`, `Runtime Manager`, `Maintenance Seams`, `Discovery Interview`, `Project Discovery`, `Manual Document Builder`, `Runtime Manager Types`, `Selection List View`, `Bounded Project Walk`, `Go Module Detection`, `Python and Django Detection`, `Detector Dispatch`, `Service Log Files`, `Linux Listener Lookup`, `Laravel Detection`, `JavaScript Framework Detection`, `Package Manifest Reading`, `Node Server Port Reading`?**
  _High betweenness centrality (0.358) - this node is a cross-community bridge._
- **Why does `Contains()` connect `Maintenance Test Doubles` to `CLI Command Wiring`, `Uninstall Tests`, `Configuration Loading`, `Manual Page Model`, `CLI Integration Tests`, `Manual Rendering`, `Detector Tests`, `Runtime Manager`, `Maintenance Seams`, `Discovery Interview`, `Project Discovery`, `Selection List View`, `Backend URL Injection`, `Settings Store Tests`, `Python and Django Detection`?**
  _High betweenness centrality (0.146) - this node is a cross-community bridge._
- **Why does `New()` connect `Maintenance Test Doubles` to `CLI Command Wiring`, `Uninstall Tests`, `Lifecycle Commands`, `Update Command`, `Update and Uninstall Service`, `CLI Integration Tests`, `Manual Rendering`, `Maintenance Seams`, `Discovery Interview`, `Runtime Manager Types`, `Settings Store Tests`?**
  _High betweenness centrality (0.093) - this node is a cross-community bridge._
- **Are the 133 inferred relationships involving `join()` (e.g. with `loadConfig()` and `loadPortFixtureConfig()`) actually correct?**
  _`join()` has 133 INFERRED edges - model-reasoned connections that need verification._
- **Are the 100 inferred relationships involving `Contains()` (e.g. with `TestInitRejectsDeprecatedAppFlag()` and `TestInitRejectsInvalidGlobalRegistry()`) actually correct?**
  _`Contains()` has 100 INFERRED edges - model-reasoned connections that need verification._
- **Are the 44 inferred relationships involving `New()` (e.g. with `runWithEnvironment()` and `TestRenderPortReassignSummaryGroupsAssignmentsByProject()`) actually correct?**
  _`New()` has 44 INFERRED edges - model-reasoned connections that need verification._
- **Are the 34 inferred relationships involving `runWithEnvironment()` (e.g. with `configuredRoots()` and `runDirectories()`) actually correct?**
  _`runWithEnvironment()` has 34 INFERRED edges - model-reasoned connections that need verification._
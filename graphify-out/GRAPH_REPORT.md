# Graph Report - .  (2026-09-03)

## Corpus Check
- 2 files · ~92,797 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1530 nodes · 3573 edges · 89 communities (74 shown, 15 thin omitted)
- Extraction: 79% EXTRACTED · 21% INFERRED · 0% AMBIGUOUS · INFERRED: 756 edges (avg confidence: 0.8)
- Token cost: 78,243 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Manual Page Model|Manual Page Model]]
- [[_COMMUNITY_CLI Command Wiring|CLI Command Wiring]]
- [[_COMMUNITY_Uninstall Tests|Uninstall Tests]]
- [[_COMMUNITY_Maintenance Test Doubles|Maintenance Test Doubles]]
- [[_COMMUNITY_Lifecycle Terminal View|Lifecycle Terminal View]]
- [[_COMMUNITY_Port Assignment Commands|Port Assignment Commands]]
- [[_COMMUNITY_Configuration Loading|Configuration Loading]]
- [[_COMMUNITY_Lifecycle Commands|Lifecycle Commands]]
- [[_COMMUNITY_Update Command|Update Command]]
- [[_COMMUNITY_Update and Uninstall Service|Update and Uninstall Service]]
- [[_COMMUNITY_CLI Integration Tests|CLI Integration Tests]]
- [[_COMMUNITY_Manual Rendering|Manual Rendering]]
- [[_COMMUNITY_Detector Tests|Detector Tests]]
- [[_COMMUNITY_Discovery Interview|Discovery Interview]]
- [[_COMMUNITY_Project Discovery|Project Discovery]]
- [[_COMMUNITY_Runtime Manager Types|Runtime Manager Types]]
- [[_COMMUNITY_Maintenance Seams|Maintenance Seams]]
- [[_COMMUNITY_Operation Locking|Operation Locking]]
- [[_COMMUNITY_Selection List View|Selection List View]]
- [[_COMMUNITY_Runtime Manager Tests|Runtime Manager Tests]]
- [[_COMMUNITY_Bounded Project Walk|Bounded Project Walk]]
- [[_COMMUNITY_Gates and CI|Gates and CI]]
- [[_COMMUNITY_Settings Store|Settings Store]]
- [[_COMMUNITY_Website Install Section|Website Install Section]]
- [[_COMMUNITY_Backend URL Injection|Backend URL Injection]]
- [[_COMMUNITY_Site Badge and Release Order|Site Badge and Release Order]]
- [[_COMMUNITY_Process Identity|Process Identity]]
- [[_COMMUNITY_Managed State Paths|Managed State Paths]]
- [[_COMMUNITY_Go Module Detection|Go Module Detection]]
- [[_COMMUNITY_Release Process|Release Process]]
- [[_COMMUNITY_Lifecycle Model|Lifecycle Model]]
- [[_COMMUNITY_Process State and Identity|Process State and Identity]]
- [[_COMMUNITY_Generated Documentation Gate|Generated Documentation Gate]]
- [[_COMMUNITY_Installation and Verification|Installation and Verification]]
- [[_COMMUNITY_Readiness and Exit Status|Readiness and Exit Status]]
- [[_COMMUNITY_Service Log Files|Service Log Files]]
- [[_COMMUNITY_Ports and Allocation|Ports and Allocation]]
- [[_COMMUNITY_Linux Listener Lookup|Linux Listener Lookup]]
- [[_COMMUNITY_Laravel Detection|Laravel Detection]]
- [[_COMMUNITY_Release Note Grouping|Release Note Grouping]]
- [[_COMMUNITY_Command Environment Tests|Command Environment Tests]]
- [[_COMMUNITY_Homebrew Tap Release|Homebrew Tap Release]]
- [[_COMMUNITY_Configuration Contract|Configuration Contract]]
- [[_COMMUNITY_Progress Observer|Progress Observer]]
- [[_COMMUNITY_Site Badge Rules|Site Badge Rules]]
- [[_COMMUNITY_Detector Dispatch|Detector Dispatch]]
- [[_COMMUNITY_JavaScript Framework Detection|JavaScript Framework Detection]]
- [[_COMMUNITY_Package Manifest Reading|Package Manifest Reading]]
- [[_COMMUNITY_Node Service Detection|Node Service Detection]]
- [[_COMMUNITY_Node Server Port Reading|Node Server Port Reading]]
- [[_COMMUNITY_Linux Process Identity|Linux Process Identity]]
- [[_COMMUNITY_Homebrew Bottle Packaging|Homebrew Bottle Packaging]]
- [[_COMMUNITY_Bottle Packaging Tests|Bottle Packaging Tests]]
- [[_COMMUNITY_Contributing Guide|Contributing Guide]]
- [[_COMMUNITY_Environment Isolation|Environment Isolation]]
- [[_COMMUNITY_README and Project Documents|README and Project Documents]]
- [[_COMMUNITY_Release Notes Script|Release Notes Script]]
- [[_COMMUNITY_Log Streaming Tests|Log Streaming Tests]]
- [[_COMMUNITY_Port Range Tests|Port Range Tests]]
- [[_COMMUNITY_Rails Detection|Rails Detection]]
- [[_COMMUNITY_Python Detection|Python Detection]]
- [[_COMMUNITY_Website and Landing Page|Website and Landing Page]]
- [[_COMMUNITY_Bottle Verification|Bottle Verification]]
- [[_COMMUNITY_Agent Rules|Agent Rules]]
- [[_COMMUNITY_CLI Helper Tests|CLI Helper Tests]]
- [[_COMMUNITY_Django Detection|Django Detection]]
- [[_COMMUNITY_Vapor Detection|Vapor Detection]]
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
- [[_COMMUNITY_Release Notes Script Shell|Release Notes Script Shell]]
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
- `Hero with version badge` --semantically_similar_to--> `release-notes.sh, release note generator`  [INFERRED] [semantically similar]
  docs/index.html → scripts/release-notes.sh
- `A release note is read without the diff` --semantically_similar_to--> `Discover names the file rather than guessing a port`  [INFERRED] [semantically similar]
  scripts/release-notes.sh → docs/index.html
- `SoftwareApplication JSON-LD block` --semantically_similar_to--> `Release Publish Job`  [INFERRED] [semantically similar]
  docs/index.html → .github/workflows/release.yml
- `Artifact attestation verification with gh` --semantically_similar_to--> `gosec Security Scan`  [INFERRED] [semantically similar]
  Documentation.md → .github/workflows/ci.yml
- `README heading and phrase checks` --conceptually_related_to--> `Install section`  [INFERRED]
  scripts/check-docs.sh → docs/index.html

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **What reaches a release note** — scripts_release_notes_prefix_mapping, scripts_release_notes_dropped_prefixes, scripts_release_notes_trailer, scripts_release_notes_trailer_overrules_prefixes, scripts_release_notes_reader_perspective [EXTRACTED 1.00]
- **grat's readiness and ownership model** — docs_index_service_contract, docs_index_ownership_boundary, docs_index_readiness, docs_index_stop_identity_check [EXTRACTED 1.00]
- **How a grat.config comes to exist** — docs_index_discover, docs_index_stack_table, docs_index_discover_asks_rather_than_guesses, docs_index_small_environment [EXTRACTED 1.00]

## Communities (89 total, 15 thin omitted)

### Community 0 - "Manual Page Model"
Cohesion: 0.06
Nodes (53): Block, Document, Block, CommandGroup, Document, CommandGroup, T, Item (+45 more)

### Community 1 - "CLI Command Wiring"
Cohesion: 0.07
Nodes (61): defaultEnvironment(), exitCode(), Run(), runWithEnvironment(), configuredRoots(), runDirectories(), canonicalCLITestPath(), environmentForTest() (+53 more)

### Community 2 - "Uninstall Tests"
Cohesion: 0.08
Nodes (61): join(), Service, Store, T, Server, T, Context, Service (+53 more)

### Community 3 - "Maintenance Test Doubles"
Cohesion: 0.08
Nodes (60): Buffer, Service, Store, T, T, T, T, runningProject() (+52 more)

### Community 4 - "Lifecycle Terminal View"
Cohesion: 0.07
Nodes (47): CancelFunc, Cmd, Context, Reader, Style, Writer, Cmd, Context (+39 more)

### Community 5 - "Port Assignment Commands"
Cohesion: 0.07
Nodes (60): portReassignment, assignReassignedPorts(), copyReservations(), ensureValidRegistry(), hasConfiguredCollision(), listenerOwnerLabel(), newPortReassignLifecycleOperation(), portReassignRowKey() (+52 more)

### Community 6 - "Configuration Loading"
Cohesion: 0.07
Nodes (53): Config, configDecodeError(), DefaultRuntime(), Load(), readConfigFile(), replaceFile(), Roles(), rollbackWrites() (+45 more)

### Community 7 - "Lifecycle Commands"
Cohesion: 0.05
Nodes (51): loadConfig(), loadManager(), executeLifecycle(), lifecycleTitle(), lifecycleTUIStage(), newLifecycleOperation(), progressPresentation(), runLifecycle() (+43 more)

### Community 8 - "Update Command"
Cohesion: 0.08
Nodes (28): runUpdate(), Context, Renderer, updateService, Renderer, Style, Context, Writer (+20 more)

### Community 9 - "Update and Uninstall Service"
Cohesion: 0.09
Nodes (26): asset, installation, T, Client, Context, Service, Client, Context (+18 more)

### Community 10 - "CLI Integration Tests"
Cohesion: 0.12
Nodes (47): assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort(), loadPortFixtureConfig(), runWithConfiguredRoots(), TestCLIRuntimeHelper(), TestDiscoverAllocatesPortsForExplicitServices() (+39 more)

### Community 11 - "Manual Rendering"
Cohesion: 0.06
Nodes (40): commandDocument(), runManual(), plainManual(), TestBothManualPagesAreReachable(), TestEveryCommandOfTheReferenceHasAManualEntry(), TestTheManualCarriesEveryCommandOfTheReference(), TestTheManualSaysWhyAProjectCanBeRefused(), writeMarkdownManual() (+32 more)

### Community 12 - "Detector Tests"
Cohesion: 0.16
Nodes (45): Directory(), commandOf(), TestAFrameworkWinsOverTheBuildToolUnderIt(), TestAGoLibraryIsReportedRatherThanInvented(), TestAGoProgramThatIgnoresThePortIsReportedRatherThanOffered(), TestAMentionOfThePortIsNotAReadOfIt(), TestAnEmptyDirectoryIsNotAProject(), TestAnExpressServerNeedsThePortInItsSource() (+37 more)

### Community 13 - "Discovery Interview"
Cohesion: 0.10
Nodes (32): detectServices(), discoverHere(), serviceSuggestions(), collectProjectInterview(), parseServiceDefinition(), promptDefault(), promptRequired(), promptServiceName() (+24 more)

### Community 14 - "Project Discovery"
Cohesion: 0.15
Nodes (33): candidate, candidate, absoluteUnder(), allocateServices(), candidateDetail(), chooseCandidates(), discoverBelow(), discoverCandidates() (+25 more)

### Community 15 - "Runtime Manager Types"
Cohesion: 0.15
Nodes (16): Client, Config, Context, loadedState, processState, ProgressObserver, RecoveryCandidate, Manager (+8 more)

### Community 16 - "Maintenance Seams"
Cohesion: 0.18
Nodes (17): Context, Reader, Result, Store, Writer, activeProject, artifactScanLimits, installation (+9 more)

### Community 17 - "Operation Locking"
Cohesion: 0.17
Nodes (25): FileMode, Context, T, Store, T, TestLockHonorsCanceledContextWhileContended(), TestLockSerializesCallbacks(), TestLockUsesRestrictivePermissions() (+17 more)

### Community 18 - "Selection List View"
Cohesion: 0.14
Nodes (11): T, NewSelectionModel(), rows(), TestALongListSaysWhatIsOutOfSight(), TestARowThatCannotBeChosenNeverIs(), TestCancellingChoosesNothing(), TestOnlyMarkedRowsComeBack(), TestTheCursorStopsAtBothEndsRatherThanWrapping() (+3 more)

### Community 19 - "Runtime Manager Tests"
Cohesion: 0.22
Nodes (23): Config, Listener, Manager, Service, T, fixtureConfig(), fixtureService(), freeTCPPort() (+15 more)

### Community 20 - "Bounded Project Walk"
Cohesion: 0.18
Nodes (18): ErrTooManyEntries, DirEntry, T, T, ErrTooManyEntries, DeeperThanScan(), SkipsScanning(), TestTheScanSkipsWhatCannotHoldAProject() (+10 more)

### Community 21 - "Gates and CI"
Cohesion: 0.12
Nodes (21): GitHub Actions Pinned by Commit Hash, Dependabot Weekly Updates, Local Development Gate, macOS and Linux Platform Parity Requirement, Build gate (cmd/grat), Build-tag coverage of static analysis, CI Workflow, Documentation gate (check-docs.sh) (+13 more)

### Community 22 - "Settings Store"
Cohesion: 0.27
Nodes (4): Settings, canonicalExistingDirectory(), canonicalExistingPath(), Store

### Community 23 - "Website Install Section"
Cohesion: 0.15
Nodes (19): Copy to clipboard buttons, grat discover writes the config, 01. Install, grat landing page, Logs in .grat/log/, Man pages, grat and grat.config, Page metadata and structured data, Section navigation (+11 more)

### Community 24 - "Backend URL Injection"
Cohesion: 0.16
Nodes (15): Context, Duration, loadedState, Context, processState, Manager, Service, waitForExit() (+7 more)

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

### Community 30 - "Lifecycle Model"
Cohesion: 0.19
Nodes (9): Cmd, Context, Model, Msg, Reader, Writer, RunSelection(), selectionTeaModel (+1 more)

### Community 31 - "Process State and Identity"
Cohesion: 0.29
Nodes (12): processState, T, TestProcessIdentitySeparatesRapidProcessStarts(), TestSignalManagedGroupRejectsChangedIdentity(), TestValidateLegacyManagedStateAcceptsDetachedLegacyProcess(), TestValidateManagedStateRejectsLegacyCoarseIdentity(), signalGroup(), signalManagedGroup() (+4 more)

### Community 32 - "Generated Documentation Gate"
Cohesion: 0.19
Nodes (13): Manual pages and Documentation.md are generated, not written, Why every documented claim is checked here, Documentation.md diffed against the manual the binary renders, Why a generated document replaces phrase-by-phrase assertions, Go 1.25.13 pin checked in go.mod, README and CONTRIBUTING, Both manual pages must render as valid roff (mandoc lint), README must open with self-updating shields.io badges, require() (+5 more)

### Community 33 - "Installation and Verification"
Cohesion: 0.22
Nodes (13): Artifact attestation verification with gh, Installation routes (Homebrew, release binary, go install), README installation instructions, Verifying the binary before installing it, not after, CI and release workflow literals asserted from the docs gate, Fail-closed provenance checks for update and direct install, brew and gh resolved through PATH, GitHub Artifact Attestation (+5 more)

### Community 34 - "Readiness and Exit Status"
Cohesion: 0.23
Nodes (12): A Docker-held port belongs to Docker, so it never becomes ready, The listener must trace back to the command grat started, Services section: start, stop, logs cards, Exit status 0, 1, 2 and 130, A managed command has to stay in the foreground, .grat/pid and .grat/log managed state, Why the listener must belong to grat's own process tree, Readiness: alive process, owned listener, 2xx health path (+4 more)

### Community 35 - "Service Log Files"
Cohesion: 0.26
Nodes (6): File, T, Manager, Service, newServiceLogFile(), TestNewServiceLogFileTruncatesPreviousOutput()

### Community 36 - "Ports and Allocation"
Cohesion: 0.22
Nodes (11): Ports section with the role lane scale, Laravel queue worker detection via QUEUE_CONNECTION, Port allocation across registered directories, grat ports assign, grat ports audit, Roles and port ranges, Scan directories and the six-level depth limit, grat/settings.toml and update-check in the user config directory (+3 more)

### Community 37 - "Linux Listener Lookup"
Cohesion: 0.25
Nodes (8): T, Listener, systemListener(), ListeningSocketInodes(), socketInode(), SocketOwnerPIDs(), TestListeningSocketInodesFindsOnlyTCPListenInodesForTheRequestedPort(), TestSocketOwnerPIDsFindsPIDsFromProcFileDescriptors()

### Community 38 - "Laravel Detection"
Cohesion: 0.40
Nodes (9): readBounded(), configuredQueueFallback(), detectLaravel(), environmentAssignment(), environmentValue(), laravelQueueConnection(), laravelQueueWorker(), Service (+1 more)

### Community 39 - "Release Note Grouping"
Cohesion: 0.24
Nodes (10): Discover names the file rather than guessing a port, awk grouping program, Chore, Test and Refactor are dropped, Prefix to heading mapping, A release note is read without the diff, Subject is capitalised into a sentence, The version bump commit is skipped, Release-note trailer (+2 more)

### Community 40 - "Command Environment Tests"
Cohesion: 0.44
Nodes (9): T, containsEnvironmentName(), TestCommandEnvironmentDerivesBackendURLForConsumer(), TestCommandEnvironmentExcludesUnapprovedParentVariables(), TestCommandEnvironmentFallsBackWhenApprovedBackendURLIsAbsent(), TestCommandEnvironmentOmitsBackendURLForProviderAndAmbiguousTopology(), TestCommandEnvironmentPreservesApprovedBackendURLOverride(), TestLaunchDoesNotSourceLoginProfile() (+1 more)

### Community 41 - "Homebrew Tap Release"
Cohesion: 0.31
Nodes (9): Install section, Homebrew Tap Formula Update, Update the Homebrew tap formula with seven values, The Nine-Step Release Order Across Two Repositories, Why the tap must not be merged before its checks finish, Why the Tap Pull Request Waits for All Three Checks, README heading and phrase checks, Create checksums.txt (+1 more)

### Community 42 - "Configuration Contract"
Cohesion: 0.25
Nodes (8): Written commands carry host and port, because binding every interface is accidental exposure, The config file is read as data, only the command is executed, grat.config schema: version, project, runtime, services, runtime table timing overrides, Service states: stopped, running, unhealthy, A [[services]] table, Platform helpers invoked through fixed absolute paths, /bin/sh execution is a trust boundary, not a sandbox

### Community 43 - "Progress Observer"
Cohesion: 0.36
Nodes (5): Manager, Service, ProgressEvent, ProgressObserver, ProgressStage

### Community 44 - "Site Badge Rules"
Cohesion: 0.29
Nodes (8): Badge and Tag Cannot Both Be First, A Pages deployment from a tag is accepted, marked active and never served, Why a newer badge is allowed and an older one is not, Site badge compared against the newest tag, refusing an older one, Why the badge is whatever main holds rather than written at publish time, Push to main as the only trigger, plus manual dispatch, Tag-Triggered Release, Release Workflow

### Community 45 - "Detector Dispatch"
Cohesion: 0.38
Nodes (5): detector, Finding, Service, Unresolved, Role

### Community 46 - "JavaScript Framework Detection"
Cohesion: 0.43
Nodes (6): detectJavaScriptFramework(), frameworkEvidence(), jsFramework, manifest, Service, Unresolved

### Community 48 - "Node Service Detection"
Cohesion: 0.57
Nodes (6): detectNode(), detectSingleService(), namedServices(), manifest, Service, Unresolved

### Community 49 - "Node Server Port Reading"
Cohesion: 0.43
Nodes (6): detectNodeServer(), readsPortFromEnvironment(), startScript(), manifest, Service, Unresolved

### Community 50 - "Linux Process Identity"
Cohesion: 0.43
Nodes (5): T, linuxProcessStartTicks(), processIdentity(), TestLinuxProcessStartTicksHandlesClosingParenthesisInCommand(), TestLinuxProcessStartTicksRejectsMismatchedPID()

### Community 51 - "Homebrew Bottle Packaging"
Cohesion: 0.48
Nodes (5): package(), usage(), write_formula(), write_manual(), build-homebrew-bottles.sh script

### Community 52 - "Bottle Packaging Tests"
Cohesion: 0.52
Nodes (6): assert_archive_contains(), assert_binary(), assert_file(), assert_manual(), assert_mode(), test-homebrew-bottles.sh script

### Community 53 - "Contributing Guide"
Cohesion: 0.33
Nodes (5): Code of conduct, Configuration compatibility, Contributing to grat, Development setup, Pull requests

### Community 54 - "Environment Isolation"
Cohesion: 0.53
Nodes (6): The environment is small on purpose, BACKEND_URL injection from the single backend role, Non-secret environment baseline for a service command, inherit_env, PORT environment variable owned by grat, Limiting accidental secret propagation into services

### Community 55 - "README and Project Documents"
Cohesion: 0.40
Nodes (6): CONTRIBUTING, SECURITY, CODE_OF_CONDUCT and SUPPORT, grat README, MIT License, Quick start with discover, start and status, grat security policy, Private vulnerability reporting to security@layered.work

### Community 56 - "Release Notes Script"
Cohesion: 0.47
Nodes (6): Full changelog compare link, Commit records from git log, release-notes.sh, release note generator, Unit and record separators keep a subject from splitting its record, Tag range resolution, Why not GitHub's generated notes

### Community 57 - "Log Streaming Tests"
Cohesion: 0.50
Nodes (3): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, T

### Community 58 - "Port Range Tests"
Cohesion: 0.60
Nodes (4): TestAWorkerOwnsNoRange(), TestEveryRoleOwnsARangeOfTheStatedWidth(), TestRolesWithDifferentRangesDoNotOverlap(), T

### Community 59 - "Rails Detection"
Cohesion: 0.40
Nodes (4): fileExists(), detectRails(), Service, Unresolved

### Community 60 - "Python Detection"
Cohesion: 0.50
Nodes (4): applicationModules(), detectPython(), Service, Unresolved

### Community 61 - "Website and Landing Page"
Cohesion: 0.40
Nodes (5): Hero with version badge, SoftwareApplication JSON-LD block, grat.layered.work landing page, robots.txt Crawl Policy, Open Graph Share Card

### Community 62 - "Bottle Verification"
Cohesion: 0.60
Nodes (3): usage(), verify_bottle(), verify-homebrew-bottles.sh script

### Community 63 - "Agent Rules"
Cohesion: 0.67
Nodes (3): Graphify Before Push, Project Agent Rules, Versioned Graphify Artifacts

### Community 64 - "CLI Helper Tests"
Cohesion: 0.67
Nodes (3): TestListenerOwnerLabelHandlesUnknownPID(), TestLogFollowUsesTrustedExecutable(), T

### Community 65 - "Django Detection"
Cohesion: 0.50
Nodes (3): detectDjango(), Service, Unresolved

### Community 66 - "Vapor Detection"
Cohesion: 0.50
Nodes (3): detectVapor(), Service, Unresolved

### Community 67 - "Manual Reference Tests"
Cohesion: 0.67
Nodes (3): T, TestAnEntryCarriesMoreThanTheOneLineReference(), TestEveryFlagOfAnEntryIsExplained()

### Community 68 - "Gate Script"
Cohesion: 0.67
Nodes (3): GOTOOLCHAIN, step(), gates.sh script

### Community 69 - "Support Policy"
Cohesion: 0.50
Nodes (3): Diagnostic Support Request, Sensitive Data Redaction, Support

### Community 70 - "Port Ownership Detection"
Cohesion: 0.67
Nodes (3): Who decides the port, Reading process.env.PORT and os.Getenv("PORT") from project source, Recognised frameworks on frontend and backend

## Knowledge Gaps
- **193 isolated node(s):** `Reader`, `Store`, `Manager`, `Config`, `StepKind` (+188 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **15 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `join()` connect `Uninstall Tests` to `Manual Page Model`, `CLI Command Wiring`, `Maintenance Test Doubles`, `Lifecycle Terminal View`, `Port Assignment Commands`, `Configuration Loading`, `Lifecycle Commands`, `Update Command`, `Update and Uninstall Service`, `CLI Integration Tests`, `Manual Rendering`, `Detector Tests`, `Discovery Interview`, `Project Discovery`, `Runtime Manager Types`, `Maintenance Seams`, `Operation Locking`, `Selection List View`, `Runtime Manager Tests`, `Bounded Project Walk`, `Settings Store`, `Go Module Detection`, `Service Log Files`, `Linux Listener Lookup`, `Laravel Detection`, `Detector Dispatch`, `JavaScript Framework Detection`, `Package Manifest Reading`, `Node Server Port Reading`, `Rails Detection`, `Python Detection`, `Django Detection`, `Vapor Detection`?**
  _High betweenness centrality (0.343) - this node is a cross-community bridge._
- **Why does `Contains()` connect `Maintenance Test Doubles` to `Manual Page Model`, `CLI Command Wiring`, `Vapor Detection`, `Uninstall Tests`, `Configuration Loading`, `CLI Integration Tests`, `Manual Rendering`, `Detector Tests`, `Discovery Interview`, `Project Discovery`, `Maintenance Seams`, `Operation Locking`, `Selection List View`, `Runtime Manager Tests`, `Settings Store`, `Python Detection`, `Process State and Identity`?**
  _High betweenness centrality (0.142) - this node is a cross-community bridge._
- **Why does `New()` connect `Maintenance Test Doubles` to `CLI Command Wiring`, `Uninstall Tests`, `Lifecycle Commands`, `Update Command`, `Update and Uninstall Service`, `CLI Integration Tests`, `Manual Rendering`, `Discovery Interview`, `Runtime Manager Types`, `Maintenance Seams`, `Operation Locking`, `Settings Store`?**
  _High betweenness centrality (0.086) - this node is a cross-community bridge._
- **Are the 133 inferred relationships involving `join()` (e.g. with `loadConfig()` and `loadPortFixtureConfig()`) actually correct?**
  _`join()` has 133 INFERRED edges - model-reasoned connections that need verification._
- **Are the 100 inferred relationships involving `Contains()` (e.g. with `TestInitRejectsDeprecatedAppFlag()` and `TestInitRejectsInvalidGlobalRegistry()`) actually correct?**
  _`Contains()` has 100 INFERRED edges - model-reasoned connections that need verification._
- **Are the 44 inferred relationships involving `New()` (e.g. with `runWithEnvironment()` and `TestRenderPortReassignSummaryGroupsAssignmentsByProject()`) actually correct?**
  _`New()` has 44 INFERRED edges - model-reasoned connections that need verification._
- **Are the 34 inferred relationships involving `runWithEnvironment()` (e.g. with `configuredRoots()` and `runDirectories()`) actually correct?**
  _`runWithEnvironment()` has 34 INFERRED edges - model-reasoned connections that need verification._
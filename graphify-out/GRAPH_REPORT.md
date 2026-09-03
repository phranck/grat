# Graph Report - .  (2026-09-03)

## Corpus Check
- 20 files · ~106,061 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1727 nodes · 4024 edges · 87 communities (74 shown, 13 thin omitted)
- Extraction: 79% EXTRACTED · 21% INFERRED · 0% AMBIGUOUS · INFERRED: 864 edges (avg confidence: 0.8)
- Token cost: 92,000 input · 7,379 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Detector Test Suite|Detector Test Suite]]
- [[_COMMUNITY_Service Runtime Manager|Service Runtime Manager]]
- [[_COMMUNITY_Process and Store Types|Process and Store Types]]
- [[_COMMUNITY_Manual Document Model|Manual Document Model]]
- [[_COMMUNITY_CLI Entry and Commands|CLI Entry and Commands]]
- [[_COMMUNITY_Runtime Test Buffers|Runtime Test Buffers]]
- [[_COMMUNITY_Port Assignment and Audit|Port Assignment and Audit]]
- [[_COMMUNITY_Release Asset Client|Release Asset Client]]
- [[_COMMUNITY_Main and Loaded State|Main and Loaded State]]
- [[_COMMUNITY_Lifecycle Operations|Lifecycle Operations]]
- [[_COMMUNITY_CLI Integration Tests|CLI Integration Tests]]
- [[_COMMUNITY_Discover Candidate Selection|Discover Candidate Selection]]
- [[_COMMUNITY_Manual Command Rendering|Manual Command Rendering]]
- [[_COMMUNITY_Config Reference Sections|Config Reference Sections]]
- [[_COMMUNITY_Store Reader and Writer|Store Reader and Writer]]
- [[_COMMUNITY_Discover Interview Flow|Discover Interview Flow]]
- [[_COMMUNITY_Output Presentation and Colour|Output Presentation and Colour]]
- [[_COMMUNITY_Landing Page Hero|Landing Page Hero]]
- [[_COMMUNITY_Release and Homebrew Steps|Release and Homebrew Steps]]
- [[_COMMUNITY_Terminal UI Model|Terminal UI Model]]
- [[_COMMUNITY_Config Loading and Roles|Config Loading and Roles]]
- [[_COMMUNITY_CI Workflow and Gates|CI Workflow and Gates]]
- [[_COMMUNITY_Selection List Model|Selection List Model]]
- [[_COMMUNITY_Listener Ownership Boundary|Listener Ownership Boundary]]
- [[_COMMUNITY_Lifecycle Rendering|Lifecycle Rendering]]
- [[_COMMUNITY_Django Phoenix PHP Detectors|Django Phoenix PHP Detectors]]
- [[_COMMUNITY_Directory Scan Limits|Directory Scan Limits]]
- [[_COMMUNITY_Documentation Gate|Documentation Gate]]
- [[_COMMUNITY_Scan Directory Settings|Scan Directory Settings]]
- [[_COMMUNITY_Config Rejection Tests|Config Rejection Tests]]
- [[_COMMUNITY_Config Ownership Checks|Config Ownership Checks]]
- [[_COMMUNITY_Settings Store Tests|Settings Store Tests]]
- [[_COMMUNITY_Safe Identifier Reading|Safe Identifier Reading]]
- [[_COMMUNITY_Release Badge Ordering|Release Badge Ordering]]
- [[_COMMUNITY_Bun Detector|Bun Detector]]
- [[_COMMUNITY_Node Manifest Reading|Node Manifest Reading]]
- [[_COMMUNITY_Go Program Detector|Go Program Detector]]
- [[_COMMUNITY_Release Note Generation|Release Note Generation]]
- [[_COMMUNITY_Lifecycle Text Layout|Lifecycle Text Layout]]
- [[_COMMUNITY_What Grat Recognises|What Grat Recognises]]
- [[_COMMUNITY_Dotnet and Rust Detectors|Dotnet and Rust Detectors]]
- [[_COMMUNITY_Config Safety Rules|Config Safety Rules]]
- [[_COMMUNITY_Service Log Files|Service Log Files]]
- [[_COMMUNITY_Self Update Command|Self Update Command]]
- [[_COMMUNITY_Static Site Detectors|Static Site Detectors]]
- [[_COMMUNITY_Help Screen Rendering|Help Screen Rendering]]
- [[_COMMUNITY_Install and Attestation|Install and Attestation]]
- [[_COMMUNITY_Command Environment Tests|Command Environment Tests]]
- [[_COMMUNITY_Detector Registry Types|Detector Registry Types]]
- [[_COMMUNITY_Text Safety Filter|Text Safety Filter]]
- [[_COMMUNITY_Deno Detector|Deno Detector]]
- [[_COMMUNITY_JavaScript Framework Table|JavaScript Framework Table]]
- [[_COMMUNITY_Node Detector|Node Detector]]
- [[_COMMUNITY_Linux Process Identity|Linux Process Identity]]
- [[_COMMUNITY_Homebrew Bottle Build|Homebrew Bottle Build]]
- [[_COMMUNITY_Homebrew Bottle Tests|Homebrew Bottle Tests]]
- [[_COMMUNITY_Contributing Guide|Contributing Guide]]
- [[_COMMUNITY_Flask Detector|Flask Detector]]
- [[_COMMUNITY_Identifier Safety Tests|Identifier Safety Tests]]
- [[_COMMUNITY_Health Probe Redirects|Health Probe Redirects]]
- [[_COMMUNITY_Log Streaming Tests|Log Streaming Tests]]
- [[_COMMUNITY_Port Range Tests|Port Range Tests]]
- [[_COMMUNITY_Homebrew Bottle Verification|Homebrew Bottle Verification]]
- [[_COMMUNITY_Agent Rules|Agent Rules]]
- [[_COMMUNITY_CLI Helper Tests|CLI Helper Tests]]
- [[_COMMUNITY_Port Range Reporting Tests|Port Range Reporting Tests]]
- [[_COMMUNITY_Manual Entry Tests|Manual Entry Tests]]
- [[_COMMUNITY_Local Gate Script|Local Gate Script]]
- [[_COMMUNITY_Trailer Check Tests|Trailer Check Tests]]
- [[_COMMUNITY_Support Policy|Support Policy]]
- [[_COMMUNITY_macOS Listener Lookup|macOS Listener Lookup]]
- [[_COMMUNITY_macOS Listener Tests|macOS Listener Tests]]
- [[_COMMUNITY_Readiness Inspection Tests|Readiness Inspection Tests]]
- [[_COMMUNITY_Listener Lookup Interface|Listener Lookup Interface]]
- [[_COMMUNITY_Favicon Build|Favicon Build]]
- [[_COMMUNITY_Code of Conduct|Code of Conduct]]
- [[_COMMUNITY_Social Card Build|Social Card Build]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
- [[_COMMUNITY_Trailer Check Script|Trailer Check Script]]
- [[_COMMUNITY_Release Notes Script|Release Notes Script]]
- [[_COMMUNITY_Bottle Verification Tests|Bottle Verification Tests]]
- [[_COMMUNITY_Config Compatibility Promise|Config Compatibility Promise]]
- [[_COMMUNITY_Colour Options|Colour Options]]

## God Nodes (most connected - your core abstractions)
1. `join()` - 151 edges
2. `Contains()` - 104 edges
3. `Directory()` - 59 edges
4. `New()` - 49 edges
5. `runWithEnvironment()` - 40 edges
6. `writeProject()` - 32 edges
7. `T` - 30 edges
8. `T` - 30 edges
9. `Service` - 29 edges
10. `T` - 26 edges

## Surprising Connections (you probably didn't know these)
- `Hero Section` --semantically_similar_to--> `release-notes.sh, release note generator`  [INFERRED] [semantically similar]
  docs/index.html → scripts/release-notes.sh
- `A release note is read without the diff` --semantically_similar_to--> `Discover names the file rather than guessing a port`  [INFERRED] [semantically similar]
  scripts/release-notes.sh → docs/index.html
- `README heading and phrase checks` --conceptually_related_to--> `Install Section`  [INFERRED]
  scripts/check-docs.sh → docs/index.html
- `SoftwareApplication Structured Data` --semantically_similar_to--> `Release Publish Job`  [INFERRED] [semantically similar]
  docs/index.html → .github/workflows/release.yml
- `Artifact attestation verification with gh` --semantically_similar_to--> `Security Scan (gosec, both GOOS)`  [INFERRED] [semantically similar]
  Documentation.md → .github/workflows/ci.yml

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Service lifecycle command set** — documentation_start, documentation_stop, documentation_restart, documentation_status, documentation_logs, documentation_recover [EXTRACTED 1.00]
- **Machine-wide port allocation flow** — documentation_roles_and_ports, documentation_scan_directories, documentation_ports_audit, documentation_ports_assign, documentation_ports_reassign [EXTRACTED 1.00]
- **Detection refuses rather than guesses the port** — documentation_discover, documentation_port_ownership_question, documentation_source_port_detection, documentation_eleventy_refusal, documentation_laravel_queue_worker_detection [INFERRED 0.85]

## Communities (87 total, 13 thin omitted)

### Community 0 - "Detector Test Suite"
Cohesion: 0.07
Nodes (82): Directory(), commandOf(), TestAFrameworkWinsOverTheBuildToolUnderIt(), TestAGoLibraryIsReportedRatherThanInvented(), TestAGoProgramThatIgnoresThePortIsReportedRatherThanOffered(), TestAMentionOfThePortIsNotAReadOfIt(), TestAnEmptyDirectoryIsNotAProject(), TestAnExpressServerNeedsThePortInItsSource() (+74 more)

### Community 1 - "Service Runtime Manager"
Cohesion: 0.06
Nodes (69): Config, Listener, Manager, Service, T, Context, Duration, loadedState (+61 more)

### Community 2 - "Process and Store Types"
Cohesion: 0.07
Nodes (69): join(), Service, Store, T, Server, T, Context, Service (+61 more)

### Community 3 - "Manual Document Model"
Cohesion: 0.06
Nodes (53): Block, Document, Block, CommandGroup, Document, CommandGroup, T, Item (+45 more)

### Community 4 - "CLI Entry and Commands"
Cohesion: 0.07
Nodes (63): defaultEnvironment(), exitCode(), loadConfig(), Run(), runWithEnvironment(), configuredRoots(), runDirectories(), canonicalCLITestPath() (+55 more)

### Community 5 - "Runtime Test Buffers"
Cohesion: 0.08
Nodes (60): Buffer, Service, Store, T, T, T, T, runningProject() (+52 more)

### Community 6 - "Port Assignment and Audit"
Cohesion: 0.07
Nodes (59): portReassignment, assignReassignedPorts(), copyReservations(), ensureValidRegistry(), hasConfiguredCollision(), listenerOwnerLabel(), newPortReassignLifecycleOperation(), portReassignRowKey() (+51 more)

### Community 7 - "Release Asset Client"
Cohesion: 0.09
Nodes (26): asset, installation, T, Client, Context, Service, Client, Context (+18 more)

### Community 8 - "Main and Loaded State"
Cohesion: 0.08
Nodes (28): Client, main(), mustGetwd(), Config, Context, loadedState, processState, ProgressObserver (+20 more)

### Community 9 - "Lifecycle Operations"
Cohesion: 0.06
Nodes (45): loadManager(), executeLifecycle(), lifecycleTitle(), lifecycleTUIStage(), newLifecycleOperation(), progressPresentation(), runLifecycle(), runLifecycleLocked() (+37 more)

### Community 10 - "CLI Integration Tests"
Cohesion: 0.12
Nodes (47): assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort(), loadPortFixtureConfig(), runWithConfiguredRoots(), TestCLIRuntimeHelper(), TestDiscoverAllocatesPortsForExplicitServices() (+39 more)

### Community 11 - "Discover Candidate Selection"
Cohesion: 0.09
Nodes (42): candidate, candidate, absoluteUnder(), allocateServices(), candidateDetail(), chooseCandidates(), discoverBelow(), discoverCandidates() (+34 more)

### Community 12 - "Manual Command Rendering"
Cohesion: 0.06
Nodes (40): commandDocument(), runManual(), plainManual(), TestBothManualPagesAreReachable(), TestEveryCommandOfTheReferenceHasAManualEntry(), TestTheManualCarriesEveryCommandOfTheReference(), TestTheManualSaysWhyAProjectCanBeRefused(), writeMarkdownManual() (+32 more)

### Community 13 - "Config Reference Sections"
Cohesion: 0.09
Nodes (39): Ports Section: Every Role Owns a Lane, The Environment Is Small on Purpose, BACKEND_URL for a Single Backend, Command Reference, Roles and Port Ranges in grat.config, The runtime Table, grat.config Top Level Keys, Description: What grat Does (+31 more)

### Community 14 - "Store Reader and Writer"
Cohesion: 0.13
Nodes (24): Context, Reader, Result, Store, Writer, Context, T, activeProject (+16 more)

### Community 15 - "Discover Interview Flow"
Cohesion: 0.10
Nodes (31): detectServices(), discoverHere(), serviceSuggestions(), collectProjectInterview(), parseServiceDefinition(), promptDefault(), promptRequired(), promptServiceName() (+23 more)

### Community 16 - "Output Presentation and Colour"
Cohesion: 0.15
Nodes (14): ColorMode, formatProjectRows(), fprint(), fprintf(), fprintln(), isTerminal(), pad(), ParseColorMode() (+6 more)

### Community 17 - "Landing Page Hero"
Cohesion: 0.09
Nodes (34): Recursive Acronym: grat runs approved tasks, Copy to clipboard buttons, grat discover writes the config, Hero Section, 01. Install, Opening Sentence: Replaces Terminal Tabs, Logs in .grat/log/, Man pages, grat and grat.config (+26 more)

### Community 18 - "Release and Homebrew Steps"
Cohesion: 0.10
Nodes (33): Copy command to clipboard button, Install Section, Homebrew Tap Formula Update, Update the Homebrew tap formula with seven values, A Pages deployment from a tag is accepted, marked active and never served, Release-note Commit Trailer, The Nine-Step Release Order Across Two Repositories, Releasing grat (release skill) (+25 more)

### Community 19 - "Terminal UI Model"
Cohesion: 0.14
Nodes (24): CancelFunc, Cmd, Cmd, Context, LifecycleOperation, Model, Msg, Reader (+16 more)

### Community 20 - "Config Loading and Roles"
Cohesion: 0.16
Nodes (23): Config, configDecodeError(), InferRole(), Load(), readConfigFile(), replaceFile(), Roles(), rollbackWrites() (+15 more)

### Community 21 - "CI Workflow and Gates"
Cohesion: 0.12
Nodes (26): GitHub Actions Pinned by Commit Hash, Dependabot Weekly Updates, Local Development Gate, macOS and Linux Platform Parity Requirement, Build Step (cmd/grat), Build-tag coverage of static analysis, Documentation gate (check-docs.sh), Documentation Gate (scripts/check-docs.sh) (+18 more)

### Community 22 - "Selection List Model"
Cohesion: 0.14
Nodes (11): T, NewSelectionModel(), rows(), TestALongListSaysWhatIsOutOfSight(), TestARowThatCannotBeChosenNeverIs(), TestCancellingChoosesNothing(), TestOnlyMarkedRowsComeBack(), TestTheCursorStopsAtBothEndsRatherThanWrapping() (+3 more)

### Community 23 - "Listener Ownership Boundary"
Cohesion: 0.13
Nodes (23): Where the Boundary Sits: Listener Ownership, Docker listener falls outside the started tree, A Docker-held port belongs to Docker, so it never becomes ready, The listener must trace back to the command grat started, Ready means ready, not started, Services Section: Ready Means Ready, Exit Status Codes, A managed command has to stay in the foreground (+15 more)

### Community 24 - "Lifecycle Rendering"
Cohesion: 0.18
Nodes (20): Context, Reader, Style, Writer, LifecycleEvent, lifecycleRow, lifecycleRows(), lifecycleStateStyle() (+12 more)

### Community 25 - "Django Phoenix PHP Detectors"
Cohesion: 0.13
Nodes (18): readBounded(), detectDjango(), detectPhoenix(), configuredQueueFallback(), detectLaravel(), environmentAssignment(), environmentValue(), laravelQueueConnection() (+10 more)

### Community 26 - "Directory Scan Limits"
Cohesion: 0.18
Nodes (18): ErrTooManyEntries, DirEntry, T, T, ErrTooManyEntries, DeeperThanScan(), SkipsScanning(), TestTheScanSkipsWhatCannotHoldAProject() (+10 more)

### Community 27 - "Documentation Gate"
Cohesion: 0.13
Nodes (19): The window where the site names a release nobody can download yet, Badge and Tag Cannot Both Be First, Manual pages and Documentation.md are generated, not written, Why a newer badge is allowed and an older one is not, Why every documented claim is checked here, Documentation.md diffed against the manual the binary renders, Why a generated document replaces phrase-by-phrase assertions, Go 1.25.13 pin checked in go.mod, README and CONTRIBUTING (+11 more)

### Community 28 - "Scan Directory Settings"
Cohesion: 0.27
Nodes (4): Settings, canonicalExistingDirectory(), canonicalExistingPath(), Store

### Community 29 - "Config Rejection Tests"
Cohesion: 0.23
Nodes (18): DefaultRuntime(), TestLoadRejectsDeprecatedAppsTable(), TestLoadRejectsLegacyShellConfig(), TestLoadRejectsOversizedConfigBeforeParsing(), TestLoadRejectsRemovedGitHubWorkerConfiguration(), TestLoadRejectsUnknownFieldsWithStrictDecoder(), TestValidateRejectsBoundedCollectionAndStringInputs(), TestValidateRejectsControlCharactersInProjectName() (+10 more)

### Community 30 - "Config Ownership Checks"
Cohesion: 0.20
Nodes (14): FileInfo, T, T, ownedBy(), OwnedByCurrentUser(), RefuseUnsafeConfig(), asAnotherAccount(), TestAConfigurationOfAnotherAccountIsRefused() (+6 more)

### Community 31 - "Settings Store Tests"
Cohesion: 0.31
Nodes (17): Store, T, canonicalPath(), equalStrings(), newTestStore(), TestContainsAcceptsRegularFileBelowRoot(), TestContainsRejectsPathsOutsideRootAndThroughSymlinks(), TestStoreAddCanonicalizesAndDeduplicatesDirectories() (+9 more)

### Community 32 - "Safe Identifier Reading"
Cohesion: 0.15
Nodes (14): describeRune(), quoteForReason(), safeIdentifier(), unresolvedIdentifier(), applicationModules(), detectPython(), rejectedName, detectVapor() (+6 more)

### Community 33 - "Release Badge Ordering"
Cohesion: 0.17
Nodes (17): Why the Badge Is Raised Before the Tag Exists, Step one: commit the raised site badge on a release branch, GITHUB_TOKEN Raises No Workflow Run, The site needs no separate deployment; the badge merge is the deployment, The Site Needs No Separate Deployment, Put the Released Version in the Badge, Pages Concurrency Group Without Cancellation, Deployment Policy for v* Tags (+9 more)

### Community 34 - "Bun Detector"
Cohesion: 0.20
Nodes (14): bunLockfile(), bunStartScript(), detectBun(), hasExtension(), servesWithBun(), sourceMatches(), detectNodeServer(), readsPortFromEnvironment() (+6 more)

### Community 35 - "Node Manifest Reading"
Cohesion: 0.14
Nodes (9): fileExists(), readManifest(), detectRails(), detectSpringBoot(), Unresolved, Service, Unresolved, Service (+1 more)

### Community 36 - "Go Program Detector"
Cohesion: 0.20
Nodes (14): detectGo(), goPrograms(), goSourceFile(), holdsMainPackage(), isPortEnvironmentCall(), readsPortFromGoSource(), TestTheGoScanStaysInsideTheBoundsEveryOtherScanHas(), goProgram (+6 more)

### Community 37 - "Release Note Generation"
Cohesion: 0.16
Nodes (16): Discover names the file rather than guessing a port, awk grouping program, Chore, Test and Refactor are dropped, Full changelog compare link, Commit records from git log, Prefix to heading mapping, A release note is read without the diff, release-notes.sh, release note generator (+8 more)

### Community 38 - "Lifecycle Text Layout"
Cohesion: 0.33
Nodes (3): truncate(), truncateStyled(), LifecycleModel

### Community 39 - "What Grat Recognises"
Cohesion: 0.14
Nodes (15): grat.config Example, grat.config ownership and permission refusal, How grat Decides What a Project Runs, grat discover, Eleventy recognised in order to be refused, grat.config(7) declarative service description, Laravel Queue Worker Detection, Who decides the port (+7 more)

### Community 40 - "Dotnet and Rust Detectors"
Cohesion: 0.21
Nodes (11): entries(), detectDotnet(), ambiguousRustBinary(), dependsOnRustServer(), detectRust(), readsPortFromRustSource(), DirEntry, Service (+3 more)

### Community 41 - "Config Safety Rules"
Cohesion: 0.17
Nodes (13): Written commands carry host and port, because binding every interface is accidental exposure, grat.config Ownership and Writability Check, The File Is Read as Data, grat.config schema: version, project, runtime, services, Read Names Restricted to a Safe Character Set, A Repository Decides What It Brings With It, Safety, Service states: stopped, running, unhealthy (+5 more)

### Community 42 - "Service Log Files"
Cohesion: 0.21
Nodes (10): File, T, T, isSymlinkRefusal(), newServiceLogFile(), TestNewServiceLogFileTruncatesPreviousOutput(), TestALinkedLogIsRefusedRatherThanFollowed(), TestALinkedStateDirectoryKeepsItsPermissions() (+2 more)

### Community 43 - "Self Update Command"
Cohesion: 0.24
Nodes (9): runUpdate(), Context, Renderer, updateService, Context, Writer, renderSpinnerFrame(), RunSpinner() (+1 more)

### Community 44 - "Static Site Detectors"
Cohesion: 0.38
Nodes (10): detectEleventy(), detectHugo(), detectJekyll(), directoryExists(), eleventyMarker(), firstExisting(), hugoLegacyMarker(), manifest (+2 more)

### Community 45 - "Help Screen Rendering"
Cohesion: 0.36
Nodes (5): Renderer, Style, Command, CommandGroup, helpUsageWidth()

### Community 46 - "Install and Attestation"
Cohesion: 0.29
Nodes (10): Artifact attestation verification with gh, Installation and Attestation Verification, Maintenance and Update Check, grat update, README installation instructions, Verifying the binary before installing it, not after, Fail-closed provenance checks for update and direct install, brew and gh resolved through PATH (+2 more)

### Community 47 - "Command Environment Tests"
Cohesion: 0.44
Nodes (9): T, containsEnvironmentName(), TestCommandEnvironmentDerivesBackendURLForConsumer(), TestCommandEnvironmentExcludesUnapprovedParentVariables(), TestCommandEnvironmentFallsBackWhenApprovedBackendURLIsAbsent(), TestCommandEnvironmentOmitsBackendURLForProviderAndAmbiguousTopology(), TestCommandEnvironmentPreservesApprovedBackendURLOverride(), TestLaunchDoesNotSourceLoginProfile() (+1 more)

### Community 48 - "Detector Registry Types"
Cohesion: 0.32
Nodes (6): detector, Finding, Service, Unresolved, Role, Role

### Community 49 - "Text Safety Filter"
Cohesion: 0.39
Nodes (6): T, ContainsUnsafe(), Sanitize(), TestSanitizeReplacesEveryUnsafeRune(), TestUnsafeRuneRejectsControlsAndUnicodeFormatCharacters(), UnsafeRune()

### Community 50 - "Deno Detector"
Cohesion: 0.43
Nodes (6): chooseDenoTask(), denoTasks(), detectDeno(), readsPortFromDenoSource(), Service, Unresolved

### Community 51 - "JavaScript Framework Table"
Cohesion: 0.43
Nodes (6): detectJavaScriptFramework(), frameworkEvidence(), jsFramework, manifest, Service, Unresolved

### Community 52 - "Node Detector"
Cohesion: 0.57
Nodes (6): detectNode(), detectSingleService(), namedServices(), manifest, Service, Unresolved

### Community 53 - "Linux Process Identity"
Cohesion: 0.43
Nodes (5): T, linuxProcessStartTicks(), processIdentity(), TestLinuxProcessStartTicksHandlesClosingParenthesisInCommand(), TestLinuxProcessStartTicksRejectsMismatchedPID()

### Community 54 - "Homebrew Bottle Build"
Cohesion: 0.48
Nodes (5): package(), usage(), write_formula(), write_manual(), build-homebrew-bottles.sh script

### Community 55 - "Homebrew Bottle Tests"
Cohesion: 0.52
Nodes (6): assert_archive_contains(), assert_binary(), assert_file(), assert_manual(), assert_mode(), test-homebrew-bottles.sh script

### Community 56 - "Contributing Guide"
Cohesion: 0.33
Nodes (5): Code of conduct, Configuration compatibility, Contributing to grat, Development setup, Pull requests

### Community 57 - "Flask Detector"
Cohesion: 0.40
Nodes (5): detectFlask(), flaskModules(), Service, Unresolved, rejectedName

### Community 58 - "Identifier Safety Tests"
Cohesion: 0.53
Nodes (5): TestADetectedNameNeverBecomesASecondCommand(), TestAnOrdinaryNameStillYieldsACommand(), TestAnUnsafeNameSaysWhichCharacter(), TestTheCharacterSetIsWhatItSays(), T

### Community 59 - "Health Probe Redirects"
Cohesion: 0.60
Nodes (5): Client, T, probeClient(), TestAHealthProbeDoesNotFollowARedirectOffTheService(), TestAHealthProbeFollowsARedirectOnTheSameService()

### Community 60 - "Log Streaming Tests"
Cohesion: 0.50
Nodes (3): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, T

### Community 61 - "Port Range Tests"
Cohesion: 0.60
Nodes (4): TestAWorkerOwnsNoRange(), TestEveryRoleOwnsARangeOfTheStatedWidth(), TestRolesWithDifferentRangesDoNotOverlap(), T

### Community 62 - "Homebrew Bottle Verification"
Cohesion: 0.60
Nodes (3): usage(), verify_bottle(), verify-homebrew-bottles.sh script

### Community 63 - "Agent Rules"
Cohesion: 0.67
Nodes (3): Graphify Before Push, Project Agent Rules, Versioned Graphify Artifacts

### Community 64 - "CLI Helper Tests"
Cohesion: 0.67
Nodes (3): TestListenerOwnerLabelHandlesUnknownPID(), TestLogFollowUsesTrustedExecutable(), T

### Community 65 - "Port Range Reporting Tests"
Cohesion: 0.67
Nodes (3): TestAPortInsideItsRangeIsNotReported(), TestAPortOutsideItsRangeDoesNotBlockLoading(), T

### Community 66 - "Manual Entry Tests"
Cohesion: 0.67
Nodes (3): T, TestAnEntryCarriesMoreThanTheOneLineReference(), TestEveryFlagOfAnEntryIsExplained()

### Community 67 - "Local Gate Script"
Cohesion: 0.67
Nodes (3): GOTOOLCHAIN, step(), gates.sh script

### Community 68 - "Trailer Check Tests"
Cohesion: 0.83
Nodes (3): commit(), expect(), test-check-trailers.sh script

### Community 69 - "Support Policy"
Cohesion: 0.50
Nodes (3): Diagnostic Support Request, Sensitive Data Redaction, Support

## Knowledge Gaps
- **214 isolated node(s):** `Reader`, `Store`, `Manager`, `Config`, `StepKind` (+209 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **13 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `join()` connect `Process and Store Types` to `Detector Test Suite`, `Service Runtime Manager`, `Manual Document Model`, `CLI Entry and Commands`, `Runtime Test Buffers`, `Port Assignment and Audit`, `Release Asset Client`, `Main and Loaded State`, `Lifecycle Operations`, `CLI Integration Tests`, `Discover Candidate Selection`, `Manual Command Rendering`, `Store Reader and Writer`, `Discover Interview Flow`, `Output Presentation and Colour`, `Config Loading and Roles`, `Selection List Model`, `Django Phoenix PHP Detectors`, `Directory Scan Limits`, `Scan Directory Settings`, `Config Ownership Checks`, `Safe Identifier Reading`, `Bun Detector`, `Node Manifest Reading`, `Go Program Detector`, `Lifecycle Text Layout`, `Dotnet and Rust Detectors`, `Static Site Detectors`, `Detector Registry Types`, `Deno Detector`, `JavaScript Framework Table`, `Flask Detector`?**
  _High betweenness centrality (0.342) - this node is a cross-community bridge._
- **Why does `Contains()` connect `Runtime Test Buffers` to `Detector Test Suite`, `Safe Identifier Reading`, `Process and Store Types`, `Manual Document Model`, `CLI Entry and Commands`, `Service Runtime Manager`, `CLI Integration Tests`, `Discover Candidate Selection`, `Manual Command Rendering`, `Store Reader and Writer`, `Discover Interview Flow`, `Selection List Model`, `Health Probe Redirects`, `Scan Directory Settings`, `Config Rejection Tests`, `Settings Store Tests`?**
  _High betweenness centrality (0.122) - this node is a cross-community bridge._
- **Why does `New()` connect `Runtime Test Buffers` to `Process and Store Types`, `CLI Entry and Commands`, `Release Asset Client`, `Main and Loaded State`, `Lifecycle Operations`, `CLI Integration Tests`, `Manual Command Rendering`, `Store Reader and Writer`, `Discover Interview Flow`, `Output Presentation and Colour`, `Scan Directory Settings`, `Settings Store Tests`?**
  _High betweenness centrality (0.074) - this node is a cross-community bridge._
- **Are the 150 inferred relationships involving `join()` (e.g. with `loadConfig()` and `loadPortFixtureConfig()`) actually correct?**
  _`join()` has 150 INFERRED edges - model-reasoned connections that need verification._
- **Are the 100 inferred relationships involving `Contains()` (e.g. with `TestInitRejectsDeprecatedAppFlag()` and `TestInitRejectsInvalidGlobalRegistry()`) actually correct?**
  _`Contains()` has 100 INFERRED edges - model-reasoned connections that need verification._
- **Are the 57 inferred relationships involving `Directory()` (e.g. with `discoverCandidates()` and `detectServices()`) actually correct?**
  _`Directory()` has 57 INFERRED edges - model-reasoned connections that need verification._
- **Are the 44 inferred relationships involving `New()` (e.g. with `runWithEnvironment()` and `TestRenderPortReassignSummaryGroupsAssignmentsByProject()`) actually correct?**
  _`New()` has 44 INFERRED edges - model-reasoned connections that need verification._
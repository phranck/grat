# Graph Report - .  (2026-09-03)

## Corpus Check
- 20 files · ~97,659 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1623 nodes · 3770 edges · 81 communities (67 shown, 14 thin omitted)
- Extraction: 79% EXTRACTED · 21% INFERRED · 0% AMBIGUOUS · INFERRED: 785 edges (avg confidence: 0.8)
- Token cost: 104,292 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Maintenance Test Doubles|Maintenance Test Doubles]]
- [[_COMMUNITY_Runtime Manager|Runtime Manager]]
- [[_COMMUNITY_CLI Command Wiring|CLI Command Wiring]]
- [[_COMMUNITY_Manual Page Model|Manual Page Model]]
- [[_COMMUNITY_Configuration Loading|Configuration Loading]]
- [[_COMMUNITY_Lifecycle Terminal View|Lifecycle Terminal View]]
- [[_COMMUNITY_Port Assignment Commands|Port Assignment Commands]]
- [[_COMMUNITY_Manual Rendering|Manual Rendering]]
- [[_COMMUNITY_Update Command|Update Command]]
- [[_COMMUNITY_Program Entry and Types|Program Entry and Types]]
- [[_COMMUNITY_Directories and Config Loading|Directories and Config Loading]]
- [[_COMMUNITY_Update and Uninstall Service|Update and Uninstall Service]]
- [[_COMMUNITY_Maintenance Seams|Maintenance Seams]]
- [[_COMMUNITY_Project Discovery|Project Discovery]]
- [[_COMMUNITY_Detector Tests|Detector Tests]]
- [[_COMMUNITY_Discovery Interview|Discovery Interview]]
- [[_COMMUNITY_Selection List View|Selection List View]]
- [[_COMMUNITY_Website and Landing Page|Website and Landing Page]]
- [[_COMMUNITY_Release Process|Release Process]]
- [[_COMMUNITY_Gates and CI|Gates and CI]]
- [[_COMMUNITY_Configuration Contract|Configuration Contract]]
- [[_COMMUNITY_Bounded Project Walk|Bounded Project Walk]]
- [[_COMMUNITY_Readiness and Ownership Rules|Readiness and Ownership Rules]]
- [[_COMMUNITY_Settings Store|Settings Store]]
- [[_COMMUNITY_Go Module Detection|Go Module Detection]]
- [[_COMMUNITY_Configuration Ownership|Configuration Ownership]]
- [[_COMMUNITY_Uninstall Tests|Uninstall Tests]]
- [[_COMMUNITY_Site Badge and Release Order|Site Badge and Release Order]]
- [[_COMMUNITY_Detected Name Safety|Detected Name Safety]]
- [[_COMMUNITY_Command Reference on the Site|Command Reference on the Site]]
- [[_COMMUNITY_Release Notes Grouping|Release Notes Grouping]]
- [[_COMMUNITY_Recovery Tests|Recovery Tests]]
- [[_COMMUNITY_Files grat Reads and Writes|Files grat Reads and Writes]]
- [[_COMMUNITY_Generated Documentation Gate|Generated Documentation Gate]]
- [[_COMMUNITY_Installation and Attestation|Installation and Attestation]]
- [[_COMMUNITY_Service Log Files|Service Log Files]]
- [[_COMMUNITY_Detector Dispatch|Detector Dispatch]]
- [[_COMMUNITY_What grat Runs|What grat Runs]]
- [[_COMMUNITY_Linux Listener Lookup|Linux Listener Lookup]]
- [[_COMMUNITY_Laravel Detection|Laravel Detection]]
- [[_COMMUNITY_Missing Directory Tests|Missing Directory Tests]]
- [[_COMMUNITY_JavaScript Framework Detection|JavaScript Framework Detection]]
- [[_COMMUNITY_Package Manifest Reading|Package Manifest Reading]]
- [[_COMMUNITY_Node Service Detection|Node Service Detection]]
- [[_COMMUNITY_Node Server Port Reading|Node Server Port Reading]]
- [[_COMMUNITY_Linux Process Identity|Linux Process Identity]]
- [[_COMMUNITY_Homebrew Bottle Packaging|Homebrew Bottle Packaging]]
- [[_COMMUNITY_Bottle Packaging Tests|Bottle Packaging Tests]]
- [[_COMMUNITY_Contributing Guide|Contributing Guide]]
- [[_COMMUNITY_Detected Name Tests|Detected Name Tests]]
- [[_COMMUNITY_Service Contract on the Site|Service Contract on the Site]]
- [[_COMMUNITY_Site Badge Rules|Site Badge Rules]]
- [[_COMMUNITY_Log Streaming Tests|Log Streaming Tests]]
- [[_COMMUNITY_Port Range Tests|Port Range Tests]]
- [[_COMMUNITY_Bottle Verification|Bottle Verification]]
- [[_COMMUNITY_Agent Rules|Agent Rules]]
- [[_COMMUNITY_CLI Helper Tests|CLI Helper Tests]]
- [[_COMMUNITY_Django Detection|Django Detection]]
- [[_COMMUNITY_Manual Reference Tests|Manual Reference Tests]]
- [[_COMMUNITY_Gate Script|Gate Script]]
- [[_COMMUNITY_Trailer Check Tests|Trailer Check Tests]]
- [[_COMMUNITY_Support Policy|Support Policy]]
- [[_COMMUNITY_Release Note Trailer Check|Release Note Trailer Check]]
- [[_COMMUNITY_Port Ownership Detection|Port Ownership Detection]]
- [[_COMMUNITY_macOS Listener Lookup|macOS Listener Lookup]]
- [[_COMMUNITY_macOS Listener Tests|macOS Listener Tests]]
- [[_COMMUNITY_Process Inspection Tests|Process Inspection Tests]]
- [[_COMMUNITY_System Listener Lookup|System Listener Lookup]]
- [[_COMMUNITY_Favicon Build|Favicon Build]]
- [[_COMMUNITY_Code of Conduct|Code of Conduct]]
- [[_COMMUNITY_Social Card Build|Social Card Build]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
- [[_COMMUNITY_Trailer Check Script|Trailer Check Script]]
- [[_COMMUNITY_Release Notes Script|Release Notes Script]]
- [[_COMMUNITY_Bottle Verification Tests|Bottle Verification Tests]]
- [[_COMMUNITY_Configuration Compatibility|Configuration Compatibility]]
- [[_COMMUNITY_Colour Options|Colour Options]]

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
- **grat only acts on what its owner approved** — documentation_safety, documentation_config_ownership, documentation_upward_search_stop, documentation_symlink_refusal, documentation_name_character_set, documentation_start_identity [EXTRACTED 1.00]
- **Every gate runs before a merge reaches main** — workflows_ci_formatting, workflows_ci_vet, workflows_ci_lint, workflows_ci_tests, workflows_ci_build, workflows_ci_vulnerability_scan, workflows_ci_security_scan, workflows_ci_documentation_step, workflows_ci_release_notes_step [EXTRACTED 1.00]
- **One wording across manual, site and share card** — documentation_description, docs_index_lead_sentence, docs_index_acronym, og_card_acronym, docs_index_structured_data [INFERRED 0.85]

## Communities (81 total, 14 thin omitted)

### Community 0 - "Maintenance Test Doubles"
Cohesion: 0.06
Nodes (86): Buffer, Service, Store, T, Server, T, Context, Service (+78 more)

### Community 1 - "Runtime Manager"
Cohesion: 0.05
Nodes (78): Config, Listener, Manager, Service, T, T, Context, Duration (+70 more)

### Community 2 - "CLI Command Wiring"
Cohesion: 0.06
Nodes (81): defaultEnvironment(), exitCode(), Run(), runWithEnvironment(), assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort() (+73 more)

### Community 3 - "Manual Page Model"
Cohesion: 0.06
Nodes (55): commandDetail, Block, Document, Block, CommandGroup, Document, CommandGroup, T (+47 more)

### Community 4 - "Configuration Loading"
Cohesion: 0.06
Nodes (61): Config, configDecodeError(), DefaultRuntime(), InferRole(), Load(), readConfigFile(), replaceFile(), Roles() (+53 more)

### Community 5 - "Lifecycle Terminal View"
Cohesion: 0.07
Nodes (47): CancelFunc, Cmd, Context, Reader, Style, Writer, Cmd, Context (+39 more)

### Community 6 - "Port Assignment Commands"
Cohesion: 0.07
Nodes (60): portReassignment, assignReassignedPorts(), copyReservations(), ensureValidRegistry(), hasConfiguredCollision(), listenerOwnerLabel(), newPortReassignLifecycleOperation(), portReassignRowKey() (+52 more)

### Community 7 - "Manual Rendering"
Cohesion: 0.06
Nodes (50): commandDocument(), runManual(), plainManual(), TestBothManualPagesAreReachable(), TestEveryCommandOfTheReferenceHasAManualEntry(), TestTheManualCarriesEveryCommandOfTheReference(), TestTheManualSaysWhyAProjectCanBeRefused(), writeMarkdownManual() (+42 more)

### Community 8 - "Update Command"
Cohesion: 0.08
Nodes (28): runUpdate(), Context, Renderer, updateService, Renderer, Style, Context, Writer (+20 more)

### Community 9 - "Program Entry and Types"
Cohesion: 0.08
Nodes (28): Client, main(), mustGetwd(), Config, Context, loadedState, processState, ProgressObserver (+20 more)

### Community 10 - "Directories and Config Loading"
Cohesion: 0.05
Nodes (51): loadConfig(), loadManager(), configuredRoots(), runDirectories(), executeLifecycle(), lifecycleTitle(), lifecycleTUIStage(), newLifecycleOperation() (+43 more)

### Community 11 - "Update and Uninstall Service"
Cohesion: 0.09
Nodes (26): asset, installation, T, Client, Context, Service, Client, Context (+18 more)

### Community 12 - "Maintenance Seams"
Cohesion: 0.09
Nodes (42): FileMode, Context, Reader, Result, Store, Writer, Context, T (+34 more)

### Community 13 - "Project Discovery"
Cohesion: 0.09
Nodes (42): candidate, candidate, absoluteUnder(), allocateServices(), candidateDetail(), chooseCandidates(), discoverBelow(), discoverCandidates() (+34 more)

### Community 14 - "Detector Tests"
Cohesion: 0.16
Nodes (45): Directory(), commandOf(), TestAFrameworkWinsOverTheBuildToolUnderIt(), TestAGoLibraryIsReportedRatherThanInvented(), TestAGoProgramThatIgnoresThePortIsReportedRatherThanOffered(), TestAMentionOfThePortIsNotAReadOfIt(), TestAnEmptyDirectoryIsNotAProject(), TestAnExpressServerNeedsThePortInItsSource() (+37 more)

### Community 15 - "Discovery Interview"
Cohesion: 0.10
Nodes (31): detectServices(), discoverHere(), serviceSuggestions(), collectProjectInterview(), parseServiceDefinition(), promptDefault(), promptRequired(), promptServiceName() (+23 more)

### Community 16 - "Selection List View"
Cohesion: 0.14
Nodes (11): T, NewSelectionModel(), rows(), TestALongListSaysWhatIsOutOfSight(), TestARowThatCannotBeChosenNeverIs(), TestCancellingChoosesNothing(), TestOnlyMarkedRowsComeBack(), TestTheCursorStopsAtBothEndsRatherThanWrapping() (+3 more)

### Community 17 - "Website and Landing Page"
Cohesion: 0.13
Nodes (24): Recursive Acronym: grat runs approved tasks, Copy to clipboard buttons, Hero Section, 01. Install, Opening Sentence: Replaces Terminal Tabs, Man pages, grat and grat.config, Page Metadata and Share Tags, Section navigation (+16 more)

### Community 18 - "Release Process"
Cohesion: 0.13
Nodes (24): Homebrew Tap Formula Update, Update the Homebrew tap formula with seven values, Release-note Commit Trailer, The Nine-Step Release Order Across Two Repositories, Releasing grat (release skill), grat Release Process, Why the tap must not be merged before its checks finish, Why the Tap Pull Request Waits for All Three Checks (+16 more)

### Community 19 - "Gates and CI"
Cohesion: 0.16
Nodes (22): GitHub Actions Pinned by Commit Hash, Dependabot Weekly Updates, Local Development Gate, macOS and Linux Platform Parity Requirement, Build Step (cmd/grat), Build-tag coverage of static analysis, Documentation gate (check-docs.sh), Documentation Gate (scripts/check-docs.sh) (+14 more)

### Community 20 - "Configuration Contract"
Cohesion: 0.14
Nodes (22): The Environment Is Small on Purpose, BACKEND_URL for a Single Backend, grat.config Example, Roles and Port Ranges in grat.config, The runtime Table, grat.config Top Level Keys, Description: What grat Does, The Environment a Command Receives (+14 more)

### Community 21 - "Bounded Project Walk"
Cohesion: 0.18
Nodes (18): ErrTooManyEntries, DirEntry, T, T, ErrTooManyEntries, DeeperThanScan(), SkipsScanning(), TestTheScanSkipsWhatCannotHoldAProject() (+10 more)

### Community 22 - "Readiness and Ownership Rules"
Cohesion: 0.14
Nodes (20): A Docker-held port belongs to Docker, so it never becomes ready, The listener must trace back to the command grat started, Services Section: Ready Means Ready, grat.config schema: version, project, runtime, services, Exit Status Codes, A managed command has to stay in the foreground, .grat/pid and .grat/log managed state, Why the listener must belong to grat's own process tree (+12 more)

### Community 23 - "Settings Store"
Cohesion: 0.29
Nodes (4): Settings, canonicalExistingDirectory(), canonicalExistingPath(), Store

### Community 24 - "Go Module Detection"
Cohesion: 0.18
Nodes (16): entries(), detectGo(), goPrograms(), goSourceFile(), holdsMainPackage(), isPortEnvironmentCall(), readsPortFromGoSource(), TestTheGoScanStaysInsideTheBoundsEveryOtherScanHas() (+8 more)

### Community 25 - "Configuration Ownership"
Cohesion: 0.20
Nodes (14): FileInfo, T, T, ownedBy(), OwnedByCurrentUser(), RefuseUnsafeConfig(), asAnotherAccount(), TestAConfigurationOfAnotherAccountIsRefused() (+6 more)

### Community 26 - "Uninstall Tests"
Cohesion: 0.37
Nodes (17): Service, Store, T, fakeUninstallService(), newUninstallStore(), TestDiscoverUninstallArtifactsRejectsScanLimitOverrun(), TestUninstallAbortsBeforePromptsForActiveManagedService(), TestUninstallDefaultYesRemovesOnlyRegisteredProjectArtifacts() (+9 more)

### Community 27 - "Site Badge and Release Order"
Cohesion: 0.16
Nodes (18): The window where the site names a release nobody can download yet, Why the Badge Is Raised Before the Tag Exists, Step one: commit the raised site badge on a release branch, GITHUB_TOKEN Raises No Workflow Run, The site needs no separate deployment; the badge merge is the deployment, The Site Needs No Separate Deployment, Put the Released Version in the Badge, Pages Concurrency Group Without Cancellation (+10 more)

### Community 28 - "Detected Name Safety"
Cohesion: 0.15
Nodes (14): describeRune(), quoteForReason(), safeIdentifier(), unresolvedIdentifier(), applicationModules(), detectPython(), rejectedName, detectVapor() (+6 more)

### Community 29 - "Command Reference on the Site"
Cohesion: 0.15
Nodes (17): Install Section, Command Reference, grat directories add, remove, list, grat discover, Version, Colour and Help Options, Maintenance and Update Check, Port Allocation, grat ports audit (+9 more)

### Community 30 - "Release Notes Grouping"
Cohesion: 0.16
Nodes (16): Discover names the file rather than guessing a port, awk grouping program, Chore, Test and Refactor are dropped, Full changelog compare link, Commit records from git log, Prefix to heading mapping, A release note is read without the diff, release-notes.sh, release note generator (+8 more)

### Community 31 - "Recovery Tests"
Cohesion: 0.46
Nodes (14): assertCLIRecoveryState(), assertRecoveryPreview(), cliProcessAlive(), legacyCLIStartIdentity(), recoveryEnvironment(), stopCLIRecoveryGroup(), TestRecoverDeclinedConfirmationLeavesLegacyProcessAndState(), TestRecoverInteractiveConfirmationStopsLegacyProcessAndRemovesState() (+6 more)

### Community 32 - "Files grat Reads and Writes"
Cohesion: 0.14
Nodes (15): Written commands carry host and port, because binding every interface is accidental exposure, grat.config Ownership and Writability Check, The File Is Read as Data, Files grat Reads and Writes, grat logs, Read Names Restricted to a Safe Character Set, grat recover, A Repository Decides What It Brings With It (+7 more)

### Community 33 - "Generated Documentation Gate"
Cohesion: 0.17
Nodes (14): Manual pages and Documentation.md are generated, not written, Why every documented claim is checked here, Documentation.md diffed against the manual the binary renders, Why a generated document replaces phrase-by-phrase assertions, Go 1.25.13 pin checked in go.mod, README and CONTRIBUTING, Both manual pages must render as valid roff (mandoc lint), README must open with self-updating shields.io badges, README heading and phrase checks (+6 more)

### Community 34 - "Installation and Attestation"
Cohesion: 0.22
Nodes (13): Artifact attestation verification with gh, Installation and Attestation Verification, README installation instructions, Verifying the binary before installing it, not after, CI and release workflow literals asserted from the docs gate, Fail-closed provenance checks for update and direct install, brew and gh resolved through PATH, GitHub Artifact Attestation (+5 more)

### Community 35 - "Service Log Files"
Cohesion: 0.21
Nodes (10): File, T, T, isSymlinkRefusal(), newServiceLogFile(), TestNewServiceLogFileTruncatesPreviousOutput(), TestALinkedLogIsRefusedRatherThanFollowed(), TestALinkedStateDirectoryKeepsItsPermissions() (+2 more)

### Community 36 - "Detector Dispatch"
Cohesion: 0.20
Nodes (9): fileExists(), detector, Finding, detectRails(), Service, Unresolved, Role, Service (+1 more)

### Community 37 - "What grat Runs"
Cohesion: 0.22
Nodes (11): Where the Boundary Sits: Listener Ownership, grat discover writes the config, Stack to Command Table, What It Runs Section, How grat Decides What a Project Runs, Laravel Queue Worker Detection, Owned Listener Condition, Refusing Beats a Guessed Port (+3 more)

### Community 38 - "Linux Listener Lookup"
Cohesion: 0.25
Nodes (8): T, Listener, systemListener(), ListeningSocketInodes(), socketInode(), SocketOwnerPIDs(), TestListeningSocketInodesFindsOnlyTCPListenInodesForTheRequestedPort(), TestSocketOwnerPIDsFindsPIDsFromProcFileDescriptors()

### Community 39 - "Laravel Detection"
Cohesion: 0.40
Nodes (9): readBounded(), configuredQueueFallback(), detectLaravel(), environmentAssignment(), environmentValue(), laravelQueueConnection(), laravelQueueWorker(), Service (+1 more)

### Community 40 - "Missing Directory Tests"
Cohesion: 0.47
Nodes (8): Store, T, canonical(), storeWithMissingDirectory(), TestAMissingDirectoryCanBeRemoved(), TestAnotherDirectoryCanBeAddedWhilstOneIsGone(), TestSettingsStillLoadWhenADirectoryIsGone(), TestShapeIsStillRefused()

### Community 41 - "JavaScript Framework Detection"
Cohesion: 0.43
Nodes (6): detectJavaScriptFramework(), frameworkEvidence(), jsFramework, manifest, Service, Unresolved

### Community 43 - "Node Service Detection"
Cohesion: 0.57
Nodes (6): detectNode(), detectSingleService(), namedServices(), manifest, Service, Unresolved

### Community 44 - "Node Server Port Reading"
Cohesion: 0.43
Nodes (6): detectNodeServer(), readsPortFromEnvironment(), startScript(), manifest, Service, Unresolved

### Community 45 - "Linux Process Identity"
Cohesion: 0.43
Nodes (5): T, linuxProcessStartTicks(), processIdentity(), TestLinuxProcessStartTicksHandlesClosingParenthesisInCommand(), TestLinuxProcessStartTicksRejectsMismatchedPID()

### Community 46 - "Homebrew Bottle Packaging"
Cohesion: 0.48
Nodes (5): package(), usage(), write_formula(), write_manual(), build-homebrew-bottles.sh script

### Community 47 - "Bottle Packaging Tests"
Cohesion: 0.52
Nodes (6): assert_archive_contains(), assert_binary(), assert_file(), assert_manual(), assert_mode(), test-homebrew-bottles.sh script

### Community 48 - "Contributing Guide"
Cohesion: 0.33
Nodes (5): Code of conduct, Configuration compatibility, Contributing to grat, Development setup, Pull requests

### Community 49 - "Detected Name Tests"
Cohesion: 0.53
Nodes (5): TestADetectedNameNeverBecomesASecondCommand(), TestAnOrdinaryNameStillYieldsACommand(), TestAnUnsafeNameSaysWhichCharacter(), TestTheCharacterSetIsWhatItSays(), T

### Community 50 - "Service Contract on the Site"
Cohesion: 0.33
Nodes (6): Logs in .grat/log/, Port ownership follows the process chain, Readiness is a listener plus a health path, What a service has to do, 03. Services, ready means ready, Stop signals only on a matching identity

### Community 51 - "Site Badge Rules"
Cohesion: 0.40
Nodes (6): Badge and Tag Cannot Both Be First, A Pages deployment from a tag is accepted, marked active and never served, Why a newer badge is allowed and an older one is not, Site badge compared against the newest tag, refusing an older one, Why the badge is whatever main holds rather than written at publish time, Push to main as the only trigger, plus manual dispatch

### Community 52 - "Log Streaming Tests"
Cohesion: 0.50
Nodes (3): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, T

### Community 53 - "Port Range Tests"
Cohesion: 0.60
Nodes (4): TestAWorkerOwnsNoRange(), TestEveryRoleOwnsARangeOfTheStatedWidth(), TestRolesWithDifferentRangesDoNotOverlap(), T

### Community 54 - "Bottle Verification"
Cohesion: 0.60
Nodes (3): usage(), verify_bottle(), verify-homebrew-bottles.sh script

### Community 55 - "Agent Rules"
Cohesion: 0.67
Nodes (3): Graphify Before Push, Project Agent Rules, Versioned Graphify Artifacts

### Community 56 - "CLI Helper Tests"
Cohesion: 0.67
Nodes (3): TestListenerOwnerLabelHandlesUnknownPID(), TestLogFollowUsesTrustedExecutable(), T

### Community 57 - "Django Detection"
Cohesion: 0.50
Nodes (3): detectDjango(), Service, Unresolved

### Community 58 - "Manual Reference Tests"
Cohesion: 0.67
Nodes (3): T, TestAnEntryCarriesMoreThanTheOneLineReference(), TestEveryFlagOfAnEntryIsExplained()

### Community 59 - "Gate Script"
Cohesion: 0.67
Nodes (3): GOTOOLCHAIN, step(), gates.sh script

### Community 60 - "Trailer Check Tests"
Cohesion: 0.83
Nodes (3): commit(), expect(), test-check-trailers.sh script

### Community 61 - "Support Policy"
Cohesion: 0.50
Nodes (3): Diagnostic Support Request, Sensitive Data Redaction, Support

### Community 62 - "Release Note Trailer Check"
Cohesion: 0.50
Nodes (4): Release-note Commit Trailer, Release Notes Step (trailer checks), Shallow Checkout Needs a Fetch of main, Why a Broken Trailer Fails Silently

### Community 63 - "Port Ownership Detection"
Cohesion: 0.67
Nodes (3): Who decides the port, Reading process.env.PORT and os.Getenv("PORT") from project source, Recognised frameworks on frontend and backend

## Knowledge Gaps
- **196 isolated node(s):** `Reader`, `Store`, `Manager`, `Config`, `StepKind` (+191 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **14 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `join()` connect `CLI Command Wiring` to `Maintenance Test Doubles`, `Runtime Manager`, `Manual Page Model`, `Configuration Loading`, `Lifecycle Terminal View`, `Port Assignment Commands`, `Manual Rendering`, `Update Command`, `Program Entry and Types`, `Directories and Config Loading`, `Update and Uninstall Service`, `Maintenance Seams`, `Project Discovery`, `Detector Tests`, `Discovery Interview`, `Selection List View`, `Bounded Project Walk`, `Settings Store`, `Go Module Detection`, `Configuration Ownership`, `Uninstall Tests`, `Detected Name Safety`, `Detector Dispatch`, `Linux Listener Lookup`, `Laravel Detection`, `Missing Directory Tests`, `JavaScript Framework Detection`, `Package Manifest Reading`, `Node Server Port Reading`, `Django Detection`?**
  _High betweenness centrality (0.339) - this node is a cross-community bridge._
- **Why does `Contains()` connect `Maintenance Test Doubles` to `Runtime Manager`, `CLI Command Wiring`, `Manual Page Model`, `Configuration Loading`, `Manual Rendering`, `Maintenance Seams`, `Project Discovery`, `Detector Tests`, `Discovery Interview`, `Selection List View`, `Settings Store`, `Uninstall Tests`, `Detected Name Safety`, `Recovery Tests`?**
  _High betweenness centrality (0.131) - this node is a cross-community bridge._
- **Why does `New()` connect `Maintenance Test Doubles` to `CLI Command Wiring`, `Manual Rendering`, `Update Command`, `Program Entry and Types`, `Directories and Config Loading`, `Update and Uninstall Service`, `Maintenance Seams`, `Discovery Interview`, `Settings Store`, `Uninstall Tests`?**
  _High betweenness centrality (0.083) - this node is a cross-community bridge._
- **Are the 133 inferred relationships involving `join()` (e.g. with `loadConfig()` and `loadPortFixtureConfig()`) actually correct?**
  _`join()` has 133 INFERRED edges - model-reasoned connections that need verification._
- **Are the 100 inferred relationships involving `Contains()` (e.g. with `TestInitRejectsDeprecatedAppFlag()` and `TestInitRejectsInvalidGlobalRegistry()`) actually correct?**
  _`Contains()` has 100 INFERRED edges - model-reasoned connections that need verification._
- **Are the 44 inferred relationships involving `New()` (e.g. with `runWithEnvironment()` and `TestRenderPortReassignSummaryGroupsAssignmentsByProject()`) actually correct?**
  _`New()` has 44 INFERRED edges - model-reasoned connections that need verification._
- **Are the 34 inferred relationships involving `runWithEnvironment()` (e.g. with `configuredRoots()` and `runDirectories()`) actually correct?**
  _`runWithEnvironment()` has 34 INFERRED edges - model-reasoned connections that need verification._
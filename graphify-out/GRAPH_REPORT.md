# Graph Report - .  (2026-09-02)

## Corpus Check
- 4 files · ~123,948 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2079 nodes · 5028 edges · 103 communities (85 shown, 18 thin omitted)
- Extraction: 78% EXTRACTED · 22% INFERRED · 0% AMBIGUOUS · INFERRED: 1085 edges (avg confidence: 0.8)
- Token cost: 84,171 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_CLI Dispatch and Test Harness|CLI Dispatch and Test Harness]]
- [[_COMMUNITY_Expose and Funnel Tests|Expose and Funnel Tests]]
- [[_COMMUNITY_Log Streaming|Log Streaming]]
- [[_COMMUNITY_Manual Page Model|Manual Page Model]]
- [[_COMMUNITY_Configuration Loading|Configuration Loading]]
- [[_COMMUNITY_Maintenance Test Doubles|Maintenance Test Doubles]]
- [[_COMMUNITY_Manual Rendering|Manual Rendering]]
- [[_COMMUNITY_Expose and Hide Commands|Expose and Hide Commands]]
- [[_COMMUNITY_Update and Uninstall Service|Update and Uninstall Service]]
- [[_COMMUNITY_CLI Command Wiring|CLI Command Wiring]]
- [[_COMMUNITY_CLI Integration Tests|CLI Integration Tests]]
- [[_COMMUNITY_Project Discovery|Project Discovery]]
- [[_COMMUNITY_Detector Tests|Detector Tests]]
- [[_COMMUNITY_Lifecycle Rendering|Lifecycle Rendering]]
- [[_COMMUNITY_Discovery Interview|Discovery Interview]]
- [[_COMMUNITY_README and Project Overview|README and Project Overview]]
- [[_COMMUNITY_Directories and Prompts|Directories and Prompts]]
- [[_COMMUNITY_Landing Page Behaviour and Rationale|Landing Page Behaviour and Rationale]]
- [[_COMMUNITY_Runtime Manager Tests|Runtime Manager Tests]]
- [[_COMMUNITY_Port Assignment Commands|Port Assignment Commands]]
- [[_COMMUNITY_Program Entry and Runtime Types|Program Entry and Runtime Types]]
- [[_COMMUNITY_Port Registry|Port Registry]]
- [[_COMMUNITY_Terminal Rendering|Terminal Rendering]]
- [[_COMMUNITY_Tailscale Setup Seams|Tailscale Setup Seams]]
- [[_COMMUNITY_Lifecycle Terminal View|Lifecycle Terminal View]]
- [[_COMMUNITY_Selection List View|Selection List View]]
- [[_COMMUNITY_Configuration Contract|Configuration Contract]]
- [[_COMMUNITY_Gates and CI|Gates and CI]]
- [[_COMMUNITY_Lifecycle Rows|Lifecycle Rows]]
- [[_COMMUNITY_Tailscale Command Client|Tailscale Command Client]]
- [[_COMMUNITY_Website and Landing Page|Website and Landing Page]]
- [[_COMMUNITY_Release Process|Release Process]]
- [[_COMMUNITY_Bounded Project Walk|Bounded Project Walk]]
- [[_COMMUNITY_Public Access Commands|Public Access Commands]]
- [[_COMMUNITY_Backend URL Injection|Backend URL Injection]]
- [[_COMMUNITY_Settings Store|Settings Store]]
- [[_COMMUNITY_Linux Listener Lookup|Linux Listener Lookup]]
- [[_COMMUNITY_Setup Approval Tests|Setup Approval Tests]]
- [[_COMMUNITY_Process Identity|Process Identity]]
- [[_COMMUNITY_Settings Store Tests|Settings Store Tests]]
- [[_COMMUNITY_Managed State Paths|Managed State Paths]]
- [[_COMMUNITY_Go Module Detection|Go Module Detection]]
- [[_COMMUNITY_Lifecycle Layout|Lifecycle Layout]]
- [[_COMMUNITY_Python and Django Detection|Python and Django Detection]]
- [[_COMMUNITY_Readiness Probing|Readiness Probing]]
- [[_COMMUNITY_Detector Dispatch|Detector Dispatch]]
- [[_COMMUNITY_Update Routes and Provenance|Update Routes and Provenance]]
- [[_COMMUNITY_Tailscale Stage Tests|Tailscale Stage Tests]]
- [[_COMMUNITY_Funnel Parsing Tests|Funnel Parsing Tests]]
- [[_COMMUNITY_Site Navigation Script|Site Navigation Script]]
- [[_COMMUNITY_Ports and Directories Commands|Ports and Directories Commands]]
- [[_COMMUNITY_Service Log Files|Service Log Files]]
- [[_COMMUNITY_Update Command|Update Command]]
- [[_COMMUNITY_Registry Lock Tests|Registry Lock Tests]]
- [[_COMMUNITY_Help Output|Help Output]]
- [[_COMMUNITY_Funnel Operations|Funnel Operations]]
- [[_COMMUNITY_Tailscale Status Parsing|Tailscale Status Parsing]]
- [[_COMMUNITY_Tailscale Errors and Funnel Identity|Tailscale Errors and Funnel Identity]]
- [[_COMMUNITY_Laravel Detection|Laravel Detection]]
- [[_COMMUNITY_Security Policy|Security Policy]]
- [[_COMMUNITY_Logs Command|Logs Command]]
- [[_COMMUNITY_JavaScript Framework Detection|JavaScript Framework Detection]]
- [[_COMMUNITY_Package Manifest Reading|Package Manifest Reading]]
- [[_COMMUNITY_Node Service Detection|Node Service Detection]]
- [[_COMMUNITY_Node Server Port Reading|Node Server Port Reading]]
- [[_COMMUNITY_Tailscale Client Interface|Tailscale Client Interface]]
- [[_COMMUNITY_Project Examples|Project Examples]]
- [[_COMMUNITY_Homebrew Bottle Packaging|Homebrew Bottle Packaging]]
- [[_COMMUNITY_Bottle Packaging Tests|Bottle Packaging Tests]]
- [[_COMMUNITY_Contributing Guide|Contributing Guide]]
- [[_COMMUNITY_Port Reassignment View|Port Reassignment View]]
- [[_COMMUNITY_Port Range Tests|Port Range Tests]]
- [[_COMMUNITY_Quiet Command Tests|Quiet Command Tests]]
- [[_COMMUNITY_Documentation Gate|Documentation Gate]]
- [[_COMMUNITY_Bottle Verification|Bottle Verification]]
- [[_COMMUNITY_Agent Rules|Agent Rules]]
- [[_COMMUNITY_CLI Helper Tests|CLI Helper Tests]]
- [[_COMMUNITY_Manual Reference Tests|Manual Reference Tests]]
- [[_COMMUNITY_Funnel Identity Tests|Funnel Identity Tests]]
- [[_COMMUNITY_Environment Isolation|Environment Isolation]]
- [[_COMMUNITY_README Gate|README Gate]]
- [[_COMMUNITY_Gate Script|Gate Script]]
- [[_COMMUNITY_macOS Listener Tests|macOS Listener Tests]]
- [[_COMMUNITY_Process Inspection Tests|Process Inspection Tests]]
- [[_COMMUNITY_System Listener Lookup|System Listener Lookup]]
- [[_COMMUNITY_Favicon Build|Favicon Build]]
- [[_COMMUNITY_Code of Conduct|Code of Conduct]]
- [[_COMMUNITY_Social Card Build|Social Card Build]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
- [[_COMMUNITY_Release Notes Script|Release Notes Script]]
- [[_COMMUNITY_Bottle Verification Tests|Bottle Verification Tests]]
- [[_COMMUNITY_Port Allocation|Port Allocation]]
- [[_COMMUNITY_Readiness Checks|Readiness Checks]]
- [[_COMMUNITY_Service Management|Service Management]]
- [[_COMMUNITY_Sigstore Attestation|Sigstore Attestation]]
- [[_COMMUNITY_Trusted Configuration Boundary|Trusted Configuration Boundary]]
- [[_COMMUNITY_Fixed Helper Paths|Fixed Helper Paths]]
- [[_COMMUNITY_Release Provenance|Release Provenance]]
- [[_COMMUNITY_Settings Storage|Settings Storage]]

## God Nodes (most connected - your core abstractions)
1. `join()` - 151 edges
2. `Contains()` - 147 edges
3. `runWithEnvironment()` - 81 edges
4. `New()` - 67 edges
5. `newCLITestStore()` - 49 edges
6. `grat` - 45 edges
7. `exposeEnvironment()` - 40 edges
8. `Directory()` - 40 edges
9. `exposeProject()` - 35 edges
10. `runPortAssignLocked()` - 33 edges

## Surprising Connections (you probably didn't know these)
- `Checksum and attestation verification` --semantically_similar_to--> `gosec Security Scan`  [INFERRED] [semantically similar]
  Documentation.md → .github/workflows/ci.yml
- `Does grat Fit Your Project` --semantically_similar_to--> `Readiness: Process, Owned Listener, Health Path`  [INFERRED] [semantically similar]
  README.md → Documentation.md
- `Port Lane Scale Visualisation` --semantically_similar_to--> `Roles and Port Ranges`  [INFERRED] [semantically similar]
  docs/index.html → Documentation.md
- `Start identity and process-group ownership` --semantically_similar_to--> `Listener must belong to the started process tree`  [INFERRED] [semantically similar]
  Documentation.md → docs/index.html
- `Local Development Gate` --semantically_similar_to--> `Verify Job`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → .github/workflows/ci.yml

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **One Version Named in Three Places** — docs_index_version_badge, workflows_pages_badge_rewrite_step, release_skill_badge_before_tag_rationale, workflows_release_site_job [INFERRED 0.85]
- **The GITHUB_TOKEN Event Gap and Its Three Answers** — release_skill_github_token_raises_no_event, workflows_pages_workflow_call_trigger, workflows_pages_release_published_trigger, workflows_release_site_job, workflows_pages_deployment_policy_v_tags [EXTRACTED 1.00]
- **What One Release Publishes** — workflows_release_build_job, workflows_release_manual_pages_step, workflows_release_homebrew_bottles, workflows_release_checksums, workflows_release_notes_script [EXTRACTED 1.00]

## Communities (103 total, 18 thin omitted)

### Community 0 - "CLI Dispatch and Test Harness"
Cohesion: 0.06
Nodes (90): parseGlobalOptions(), runWithEnvironment(), canonicalCLITestPath(), environmentForTest(), newCLITestStore(), sameStringSlices(), TestDirectoriesAddDoesNotPromptForInitialSetup(), TestDirectoriesCommandsPersistAndListConfiguredRoots() (+82 more)

### Community 1 - "Expose and Funnel Tests"
Cohesion: 0.05
Nodes (89): Buffer, TestOnlyGratSpeaksDuringExpose(), configuredServices(), settle(), stoppedWithOpenFunnel(), TestAPortChangeClosesTheFunnel(), TestStartNamesAFunnelThatIsAlreadyOpen(), TestStopClosesItWithoutATerminalToo() (+81 more)

### Community 2 - "Log Streaming"
Cohesion: 0.06
Nodes (73): TestOutputLogStreamsBeforeInputReachesEOF(), notifyingWriter, join(), T, Service, Store, T, Server (+65 more)

### Community 3 - "Manual Page Model"
Cohesion: 0.05
Nodes (60): Block, Document, Block, CommandGroup, Document, Role, Section, CommandGroup (+52 more)

### Community 4 - "Configuration Loading"
Cohesion: 0.05
Nodes (65): Config, configDecodeError(), DefaultRuntime(), FunnelPublicPorts(), Load(), publicPortList(), readConfigFile(), replaceFile() (+57 more)

### Community 5 - "Maintenance Test Doubles"
Cohesion: 0.07
Nodes (43): T, Context, Service, Reader, Result, Store, Writer, Context (+35 more)

### Community 6 - "Manual Rendering"
Cohesion: 0.06
Nodes (54): commandDocument(), runManual(), plainManual(), TestBothManualPagesAreReachable(), TestEveryCommandOfTheReferenceHasAManualEntry(), TestTheManualCarriesEveryCommandOfTheReference(), TestTheManualSaysWhyAProjectCanBeRefused(), writeMarkdownManual() (+46 more)

### Community 7 - "Expose and Hide Commands"
Cohesion: 0.08
Nodes (53): loadConfig(), configuredPath(), containsName(), forgetExposePaths(), funnelsToClose(), parseExposeArguments(), reachableAt(), readyTailscale() (+45 more)

### Community 8 - "Update and Uninstall Service"
Cohesion: 0.09
Nodes (26): asset, installation, T, Client, Context, Service, Client, Context (+18 more)

### Community 9 - "CLI Command Wiring"
Cohesion: 0.11
Nodes (50): configuredRoots(), confirmRecovery(), copyReservations(), defaultEnvironment(), detectServices(), executeLifecycle(), fileExists(), hasConfiguredCollision() (+42 more)

### Community 10 - "CLI Integration Tests"
Cohesion: 0.11
Nodes (48): exitCode(), assertGloballyUniqueRolePorts(), cliHelperCommand(), containsArgument(), freeCLITCPPort(), loadPortFixtureConfig(), runWithConfiguredRoots(), TestCLIRuntimeHelper() (+40 more)

### Community 11 - "Project Discovery"
Cohesion: 0.09
Nodes (42): candidate, candidate, absoluteUnder(), allocateServices(), candidateDetail(), chooseCandidates(), discoverBelow(), discoverCandidates() (+34 more)

### Community 12 - "Detector Tests"
Cohesion: 0.16
Nodes (45): Directory(), commandOf(), TestAFrameworkWinsOverTheBuildToolUnderIt(), TestAGoLibraryIsReportedRatherThanInvented(), TestAGoProgramThatIgnoresThePortIsReportedRatherThanOffered(), TestAMentionOfThePortIsNotAReadOfIt(), TestAnEmptyDirectoryIsNotAProject(), TestAnExpressServerNeedsThePortInItsSource() (+37 more)

### Community 13 - "Lifecycle Rendering"
Cohesion: 0.09
Nodes (41): lifecycleTUIStage(), progressPresentation(), funnelWithdrawalCollector, funnelWithdrawalReporter, executeLifecycle(), lifecycleTitle(), lifecycleTUIStage(), newLifecycleOperation() (+33 more)

### Community 14 - "Discovery Interview"
Cohesion: 0.10
Nodes (39): detectServices(), discoverHere(), serviceSuggestions(), collectProjectInterview(), parseServiceDefinition(), promptDefault(), promptRequired(), promptServiceName() (+31 more)

### Community 15 - "README and Project Overview"
Cohesion: 0.05
Nodes (44): GitHub Artifact Attestation, BACKEND_URL Discovery, Code of Conduct, Command contract, Commands, Configuration reference, Contents, Contributing and support (+36 more)

### Community 16 - "Directories and Prompts"
Cohesion: 0.07
Nodes (36): configuredRoots(), runDirectories(), readPromptLine(), confirmRecovery(), hasLiveRecoveryCandidate(), renderRecoveryPreview(), runRecover(), publicAddresses() (+28 more)

### Community 17 - "Landing Page Behaviour and Rationale"
Cohesion: 0.08
Nodes (43): Why grat Asks Before Installing Tailscale, Copy-to-Clipboard Command Button, Copy Command Button Handler, Allocation Across Every Known Project, grat discover Writes the Config, grat discover Framework Detection, Why a Docker-Held Port Never Reports Ready, grat expose Terminal Transcript (+35 more)

### Community 18 - "Runtime Manager Tests"
Cohesion: 0.11
Nodes (37): Config, Listener, Manager, Service, T, T, Manager, Service (+29 more)

### Community 19 - "Port Assignment Commands"
Cohesion: 0.12
Nodes (39): assignReassignedPorts(), ensureValidRegistry(), newPortReassignLifecycleOperation(), portReassignRowKey(), renderPortReassignSummary(), runPortAssignLocked(), stopReassignProjects(), validatePortReassignReport() (+31 more)

### Community 20 - "Program Entry and Runtime Types"
Cohesion: 0.13
Nodes (18): main(), mustGetwd(), Client, Config, Context, loadedState, processState, ProgressObserver (+10 more)

### Community 21 - "Port Registry"
Cohesion: 0.13
Nodes (31): Config, T, fakeLookup, Listener, ListenerLookup, Problem, ProjectConfig, FirstFree() (+23 more)

### Community 22 - "Terminal Rendering"
Cohesion: 0.17
Nodes (10): formatProjectRows(), fprint(), fprintln(), pad(), stepStyle(), terminalSafe(), ProjectGroup, ProjectRowsOptions (+2 more)

### Community 23 - "Tailscale Setup Seams"
Cohesion: 0.17
Nodes (22): InstallCommand, Client, CommandClient, Context, Reader, Writer, Approval, Approver (+14 more)

### Community 24 - "Lifecycle Terminal View"
Cohesion: 0.14
Nodes (24): CancelFunc, Cmd, Cmd, Context, LifecycleOperation, Model, Msg, Reader (+16 more)

### Community 25 - "Selection List View"
Cohesion: 0.14
Nodes (11): T, NewSelectionModel(), rows(), TestALongListSaysWhatIsOutOfSight(), TestARowThatCannotBeChosenNeverIs(), TestCancellingChoosesNothing(), TestOnlyMarkedRowsComeBack(), TestTheCursorStopsAtBothEndsRatherThanWrapping() (+3 more)

### Community 26 - "Configuration Contract"
Cohesion: 0.12
Nodes (24): grat.config Compatibility Promise, Listener must belong to the started process tree, Stack to command table, What it runs section, BACKEND_URL Topology Discovery, Non-secret Environment Baseline, grat expose status, grat.config(7) Configuration Schema (+16 more)

### Community 27 - "Gates and CI"
Cohesion: 0.11
Nodes (23): GitHub Actions Pinned by Commit Hash, Dependabot Weekly Updates, Local Development Gate, grat(1), grat manual, Badge and Tag Cannot Both Be First, Build gate (cmd/grat), Build-tag coverage of static analysis (+15 more)

### Community 28 - "Lifecycle Rows"
Cohesion: 0.18
Nodes (20): Context, Reader, Style, Writer, LifecycleEvent, lifecycleRow, lifecycleRows(), lifecycleStateStyle() (+12 more)

### Community 29 - "Tailscale Command Client"
Cohesion: 0.15
Nodes (18): CommandClient, Context, Duration, Reader, Writer, Locate(), ErrNoInstallPath, InstallCommand (+10 more)

### Community 30 - "Website and Landing Page"
Cohesion: 0.13
Nodes (22): grat.layered.work Landing Page, Collapsing Navigation Panel, Duration-based Scroll Respecting prefers-reduced-motion, SoftwareApplication structured data, Site Version Badge, robots.txt Crawl Policy, Open Graph Share Card, Why the Badge Is Raised Before the Tag Exists (+14 more)

### Community 31 - "Release Process"
Cohesion: 0.14
Nodes (22): Man Pages grat(1) and grat.config(7), Homebrew Tap Formula Update, Release-note Commit Trailer, The Nine-Step Release Order Across Two Repositories, grat Release Process, Why the Tap Pull Request Waits for All Three Checks, Version Bump Decision from Commit Prefixes, Attest Release Binary (+14 more)

### Community 32 - "Bounded Project Walk"
Cohesion: 0.18
Nodes (18): ErrTooManyEntries, DirEntry, T, T, ErrTooManyEntries, DeeperThanScan(), SkipsScanning(), TestTheScanSkipsWhatCannotHoldAProject() (+10 more)

### Community 33 - "Public Access Commands"
Cohesion: 0.17
Nodes (20): macOS and Linux Platform Parity Requirement, Public access section, --always flag, Exit Status Codes, grat expose, [services.expose] table, A Funnel Outlives the Service Behind It, grat hide (+12 more)

### Community 34 - "Backend URL Injection"
Cohesion: 0.18
Nodes (17): Context, Duration, loadedState, processState, T, TestProcessIdentitySeparatesRapidProcessStarts(), TestSignalManagedGroupRejectsChangedIdentity(), TestValidateLegacyManagedStateAcceptsDetachedLegacyProcess() (+9 more)

### Community 35 - "Settings Store"
Cohesion: 0.27
Nodes (4): Settings, canonicalExistingDirectory(), canonicalExistingPath(), Store

### Community 36 - "Linux Listener Lookup"
Cohesion: 0.15
Nodes (14): T, T, Listener, systemListener(), linuxListeningSocketInodes(), linuxSocketOwnerPIDs(), socketInode(), TestLinuxListenerParsingFindsOnlyTCPListenInodesForTheRequestedPort() (+6 more)

### Community 37 - "Setup Approval Tests"
Cohesion: 0.17
Nodes (16): Approval, Context, MachineChange, SetupEvent, T, answeringApprover, recordingObserver, missingTailscale() (+8 more)

### Community 38 - "Process Identity"
Cohesion: 0.33
Nodes (17): Cmd, Manager, Service, T, legacyProcessIdentity(), processAlive(), newLegacyRecoveryFixture(), stopFixtureGroup() (+9 more)

### Community 39 - "Settings Store Tests"
Cohesion: 0.31
Nodes (17): Store, T, canonicalPath(), equalStrings(), newTestStore(), TestContainsAcceptsRegularFileBelowRoot(), TestContainsRejectsPathsOutsideRootAndThroughSymlinks(), TestStoreAddCanonicalizesAndDeduplicatesDirectories() (+9 more)

### Community 40 - "Managed State Paths"
Cohesion: 0.23
Nodes (6): Manager, Service, loadedState, processState, RecoveryCandidate, Time

### Community 41 - "Go Module Detection"
Cohesion: 0.19
Nodes (14): entries(), detectGo(), goPrograms(), goSourceFile(), holdsMainPackage(), isPortEnvironmentCall(), readsPortFromGoSource(), TestTheGoScanStaysInsideTheBoundsEveryOtherScanHas() (+6 more)

### Community 42 - "Lifecycle Layout"
Cohesion: 0.33
Nodes (3): truncate(), truncateStyled(), LifecycleModel

### Community 43 - "Python and Django Detection"
Cohesion: 0.16
Nodes (11): readBounded(), detectDjango(), applicationModules(), detectPython(), detectVapor(), Service, Unresolved, Service (+3 more)

### Community 44 - "Readiness Probing"
Cohesion: 0.22
Nodes (11): Context, processState, Manager, Service, legacyProcessIdentity(), parentProcessID(), processAlive(), psField() (+3 more)

### Community 45 - "Detector Dispatch"
Cohesion: 0.18
Nodes (10): InferRole(), fileExists(), detector, Finding, detectRails(), Service, Unresolved, Role (+2 more)

### Community 46 - "Update Routes and Provenance"
Cohesion: 0.18
Nodes (13): Checksum and attestation verification, Documentation Installation Reference, grat update by Installation Route, settings.toml, grat update, Background Update Check, Installation, Verify Before Install (+5 more)

### Community 47 - "Tailscale Stage Tests"
Cohesion: 0.33
Nodes (12): T, stageFor(), TestABrowserCommandExistsForEverySupportedSystem(), TestAFailedStatusCallMeansTheServiceIsNotRunning(), TestARunningBackendIsReady(), TestAStartingBackendIsNotYetReady(), TestAStatusWithoutATailnetMeansSignedOut(), TestEveryStateThatSigningInFixesReportsTheSameStage() (+4 more)

### Community 48 - "Funnel Parsing Tests"
Cohesion: 0.31
Nodes (12): T, parseFunnels(), TestAMissingInstallationNamesWhereItLooked(), TestClosingAFunnelRepeatsEveryFlagAndDropsTheTarget(), TestOpeningAFunnelUsesTheDocumentedArguments(), TestParsingAnEmptyServeConfigReportsNothing(), TestParsingRejectsAHostPortWithoutAUsablePort(), TestParsingTheServeConfigReportsEveryPathOfOnePublicPort() (+4 more)

### Community 49 - "Site Navigation Script"
Cohesion: 0.21
Nodes (12): The Button States Itself Through aria-expanded, Escape key handler, Escape Keydown Handler, In-Page Script Block, Nav Panel Link Click Handler, Nav link click handler, Nav toggle click handler, Duration Over Distance, and None When Motion Is Refused (+4 more)

### Community 50 - "Ports and Directories Commands"
Cohesion: 0.20
Nodes (12): Port Lane Scale Visualisation, grat directories add, grat directories list, grat directories remove, grat discover, Funnel path collision refusal, Laravel queue worker detection, Port Allocation and Audit (+4 more)

### Community 51 - "Service Log Files"
Cohesion: 0.26
Nodes (6): File, T, Manager, Service, newServiceLogFile(), TestNewServiceLogFileTruncatesPreviousOutput()

### Community 52 - "Update Command"
Cohesion: 0.24
Nodes (9): runUpdate(), Context, Renderer, updateService, Context, Writer, renderSpinnerFrame(), RunSpinner() (+1 more)

### Community 53 - "Registry Lock Tests"
Cohesion: 0.33
Nodes (9): Context, T, Registry Lock Tests, TestRegistryLockHonorsContextWhileContended(), TestRegistryLockReleasesAfterCallbackPanic(), TestRegistryLockSerializesCallbacks(), TestRegistryLockUsesProvidedGratConfigurationDirectory(), WithRegistryLock() (+1 more)

### Community 54 - "Help Output"
Cohesion: 0.36
Nodes (5): Renderer, Style, Command, CommandGroup, helpUsageWidth()

### Community 55 - "Funnel Operations"
Cohesion: 0.44
Nodes (4): Context, Funnel, CommandClient, funnelArguments()

### Community 56 - "Tailscale Status Parsing"
Cohesion: 0.24
Nodes (9): Writer, boundedWriter, NewCommandClient(), parseStatus(), portOf(), serveConfig, status, webHandler (+1 more)

### Community 57 - "Tailscale Errors and Funnel Identity"
Cohesion: 0.24
Nodes (4): Client, ErrCommandFailed, ErrNotInstalled, Funnel

### Community 58 - "Laravel Detection"
Cohesion: 0.42
Nodes (8): configuredQueueFallback(), detectLaravel(), environmentAssignment(), environmentValue(), laravelQueueConnection(), laravelQueueWorker(), Service, Unresolved

### Community 59 - "Security Policy"
Cohesion: 0.25
Nodes (6): Reporting a vulnerability, Security policy, Supported versions, Diagnostic Support Request, Sensitive Data Redaction, Support

### Community 60 - "Logs Command"
Cohesion: 0.43
Nodes (6): outputLog(), runLogs(), writerAdapter, Context, Renderer, Writer

### Community 61 - "JavaScript Framework Detection"
Cohesion: 0.43
Nodes (6): detectJavaScriptFramework(), frameworkEvidence(), jsFramework, manifest, Service, Unresolved

### Community 63 - "Node Service Detection"
Cohesion: 0.57
Nodes (6): detectNode(), detectSingleService(), namedServices(), manifest, Service, Unresolved

### Community 64 - "Node Server Port Reading"
Cohesion: 0.43
Nodes (6): detectNodeServer(), readsPortFromEnvironment(), startScript(), manifest, Service, Unresolved

### Community 65 - "Tailscale Client Interface"
Cohesion: 0.57
Nodes (3): Context, Funnel, Client

### Community 66 - "Project Examples"
Cohesion: 0.29
Nodes (7): Go HTTP API, Laravel, Project examples, Python with FastAPI, React, Laravel, and a queue worker, React with Vite, Swift with Vapor

### Community 67 - "Homebrew Bottle Packaging"
Cohesion: 0.48
Nodes (5): package(), usage(), write_formula(), write_manual(), build-homebrew-bottles.sh script

### Community 68 - "Bottle Packaging Tests"
Cohesion: 0.52
Nodes (6): assert_archive_contains(), assert_binary(), assert_file(), assert_manual(), assert_mode(), test-homebrew-bottles.sh script

### Community 69 - "Contributing Guide"
Cohesion: 0.33
Nodes (5): Code of conduct, Configuration compatibility, Contributing to grat, Development setup, Pull requests

### Community 70 - "Port Reassignment View"
Cohesion: 0.40
Nodes (5): newLifecycleOperation(), selectPortServices(), Config, Service, LifecycleOperation

### Community 71 - "Port Range Tests"
Cohesion: 0.60
Nodes (4): TestAWorkerOwnsNoRange(), TestEveryRoleOwnsARangeOfTheStatedWidth(), TestRolesWithDifferentRangesDoNotOverlap(), T

### Community 72 - "Quiet Command Tests"
Cohesion: 0.60
Nodes (4): T, TestRunQuietlyKeepsOutputOffTheTerminal(), TestRunQuietlyPutsTheOutputInTheError(), TestTailLinesKeepsTheEnd()

### Community 73 - "Documentation Gate"
Cohesion: 0.70
Nodes (4): require(), require_doc(), require_in(), check-docs.sh script

### Community 74 - "Bottle Verification"
Cohesion: 0.60
Nodes (3): usage(), verify_bottle(), verify-homebrew-bottles.sh script

### Community 75 - "Agent Rules"
Cohesion: 0.67
Nodes (3): Graphify Before Push, Project Agent Rules, Versioned Graphify Artifacts

### Community 76 - "CLI Helper Tests"
Cohesion: 0.67
Nodes (3): TestListenerOwnerLabelHandlesUnknownPID(), TestLogFollowUsesTrustedExecutable(), T

### Community 77 - "Manual Reference Tests"
Cohesion: 0.67
Nodes (3): T, TestAnEntryCarriesMoreThanTheOneLineReference(), TestEveryFlagOfAnEntryIsExplained()

### Community 78 - "Funnel Identity Tests"
Cohesion: 0.67
Nodes (3): T, TestAFunnelIsItselfAndNothingElse(), TestOneFunnelBelongsToOneService()

### Community 79 - "Environment Isolation"
Cohesion: 0.50
Nodes (4): BACKEND_URL, inherit_env, Non-Secret Environment Baseline, Trusted Local Configurations

### Community 80 - "README Gate"
Cohesion: 0.83
Nodes (3): require(), require_in(), check-readme.sh script

### Community 81 - "Gate Script"
Cohesion: 0.67
Nodes (3): GOTOOLCHAIN, step(), gates.sh script

## Ambiguous Edges - Review These
- `Site Version Badge` → `SoftwareApplication JSON-LD`  [AMBIGUOUS]
  docs/index.html · relation: shares_data_with

## Knowledge Gaps
- **263 isolated node(s):** `Store`, `ProgressObserver`, `LifecycleEvent`, `ColorMode`, `Service` (+258 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **18 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `Site Version Badge` and `SoftwareApplication JSON-LD`?**
  _Edge tagged AMBIGUOUS (relation: shares_data_with) - confidence is low._
- **Why does `join()` connect `Log Streaming` to `CLI Dispatch and Test Harness`, `Expose and Funnel Tests`, `Manual Page Model`, `Configuration Loading`, `Maintenance Test Doubles`, `Manual Rendering`, `Expose and Hide Commands`, `Update and Uninstall Service`, `CLI Integration Tests`, `Project Discovery`, `Detector Tests`, `Discovery Interview`, `Runtime Manager Tests`, `Port Assignment Commands`, `Program Entry and Runtime Types`, `Port Registry`, `Terminal Rendering`, `Selection List View`, `Tailscale Command Client`, `Bounded Project Walk`, `Settings Store`, `Linux Listener Lookup`, `Go Module Detection`, `Lifecycle Layout`, `Python and Django Detection`, `Detector Dispatch`, `Service Log Files`, `Registry Lock Tests`, `Tailscale Errors and Funnel Identity`, `Laravel Detection`, `Logs Command`, `JavaScript Framework Detection`, `Package Manifest Reading`, `Node Server Port Reading`?**
  _High betweenness centrality (0.306) - this node is a cross-community bridge._
- **Why does `Contains()` connect `Expose and Funnel Tests` to `CLI Dispatch and Test Harness`, `Log Streaming`, `Manual Page Model`, `Configuration Loading`, `Maintenance Test Doubles`, `Manual Rendering`, `Expose and Hide Commands`, `CLI Integration Tests`, `Project Discovery`, `Detector Tests`, `Discovery Interview`, `Runtime Manager Tests`, `Selection List View`, `Backend URL Injection`, `Settings Store`, `Setup Approval Tests`, `Settings Store Tests`, `Python and Django Detection`, `Tailscale Stage Tests`, `Funnel Parsing Tests`, `Funnel Operations`, `Quiet Command Tests`?**
  _High betweenness centrality (0.177) - this node is a cross-community bridge._
- **Why does `New()` connect `Expose and Funnel Tests` to `CLI Dispatch and Test Harness`, `Log Streaming`, `Settings Store`, `Maintenance Test Doubles`, `Manual Rendering`, `Expose and Hide Commands`, `Update and Uninstall Service`, `CLI Command Wiring`, `CLI Integration Tests`, `Setup Approval Tests`, `Settings Store Tests`, `Discovery Interview`, `Tailscale Stage Tests`, `Directories and Prompts`, `Program Entry and Runtime Types`, `Funnel Operations`, `Terminal Rendering`, `Tailscale Setup Seams`?**
  _High betweenness centrality (0.098) - this node is a cross-community bridge._
- **Are the 150 inferred relationships involving `join()` (e.g. with `loadConfig()` and `loadPortFixtureConfig()`) actually correct?**
  _`join()` has 150 INFERRED edges - model-reasoned connections that need verification._
- **Are the 143 inferred relationships involving `Contains()` (e.g. with `TestInitRejectsDeprecatedAppFlag()` and `TestInitRejectsInvalidGlobalRegistry()`) actually correct?**
  _`Contains()` has 143 INFERRED edges - model-reasoned connections that need verification._
- **Are the 62 inferred relationships involving `runWithEnvironment()` (e.g. with `configuredRoots()` and `runDirectories()`) actually correct?**
  _`runWithEnvironment()` has 62 INFERRED edges - model-reasoned connections that need verification._
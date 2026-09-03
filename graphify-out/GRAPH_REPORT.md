# Graph Report - .  (2026-09-03)

## Corpus Check
- 10 files · ~108,144 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1749 nodes · 4091 edges · 99 communities (84 shown, 15 thin omitted)
- Extraction: 79% EXTRACTED · 21% INFERRED · 0% AMBIGUOUS · INFERRED: 877 edges (avg confidence: 0.8)
- Token cost: 96,000 input · 7,783 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Program Entry Point|Program Entry Point]]
- [[_COMMUNITY_CLI Entry and Dispatch|CLI Entry and Dispatch]]
- [[_COMMUNITY_Service Manager Loading|Service Manager Loading]]
- [[_COMMUNITY_Ports Command and Reassignment|Ports Command and Reassignment]]
- [[_COMMUNITY_CLI Integration Tests|CLI Integration Tests]]
- [[_COMMUNITY_CLI Helper Tests|CLI Helper Tests]]
- [[_COMMUNITY_Log Streaming Tests|Log Streaming Tests]]
- [[_COMMUNITY_Config Types and Roles|Config Types and Roles]]
- [[_COMMUNITY_Config Rejection Tests|Config Rejection Tests]]
- [[_COMMUNITY_Update Service and Installation|Update Service and Installation]]
- [[_COMMUNITY_Uninstall Artifact Scan|Uninstall Artifact Scan]]
- [[_COMMUNITY_Uninstall Command|Uninstall Command]]
- [[_COMMUNITY_Listener Lookup Interface|Listener Lookup Interface]]
- [[_COMMUNITY_macOS Listener Lookup|macOS Listener Lookup]]
- [[_COMMUNITY_macOS Listener Tests|macOS Listener Tests]]
- [[_COMMUNITY_Registry Lock|Registry Lock]]
- [[_COMMUNITY_Project Registry Scan|Project Registry Scan]]
- [[_COMMUNITY_Registry Scan Tests|Registry Scan Tests]]
- [[_COMMUNITY_Detector Result Types|Detector Result Types]]
- [[_COMMUNITY_Help Screen Rendering|Help Screen Rendering]]
- [[_COMMUNITY_Lifecycle Text Layout|Lifecycle Text Layout]]
- [[_COMMUNITY_Lifecycle Event Types|Lifecycle Event Types]]
- [[_COMMUNITY_Presentation Renderer|Presentation Renderer]]
- [[_COMMUNITY_Lifecycle Tea Model|Lifecycle Tea Model]]
- [[_COMMUNITY_Presentation Types|Presentation Types]]
- [[_COMMUNITY_Spinner Runner|Spinner Runner]]
- [[_COMMUNITY_Project Root Lookup|Project Root Lookup]]
- [[_COMMUNITY_Service Log Files|Service Log Files]]
- [[_COMMUNITY_Runtime Manager Interface|Runtime Manager Interface]]
- [[_COMMUNITY_Manager Lifecycle Tests|Manager Lifecycle Tests]]
- [[_COMMUNITY_Process Group Signalling|Process Group Signalling]]
- [[_COMMUNITY_Service Launch Environment|Service Launch Environment]]
- [[_COMMUNITY_Process Stop and Readiness|Process Stop and Readiness]]
- [[_COMMUNITY_Command Environment Tests|Command Environment Tests]]
- [[_COMMUNITY_Linux Process Identity|Linux Process Identity]]
- [[_COMMUNITY_Progress Observation|Progress Observation]]
- [[_COMMUNITY_Recovery of Legacy State|Recovery of Legacy State]]
- [[_COMMUNITY_Readiness Inspection Tests|Readiness Inspection Tests]]
- [[_COMMUNITY_Managed State on Disk|Managed State on Disk]]
- [[_COMMUNITY_Settings Store|Settings Store]]
- [[_COMMUNITY_Settings Store Tests|Settings Store Tests]]
- [[_COMMUNITY_Text Safety Filter|Text Safety Filter]]
- [[_COMMUNITY_Version Reporting|Version Reporting]]
- [[_COMMUNITY_Homebrew Bottle Build|Homebrew Bottle Build]]
- [[_COMMUNITY_Release Build Script|Release Build Script]]
- [[_COMMUNITY_Homebrew Bottle Tests|Homebrew Bottle Tests]]
- [[_COMMUNITY_Bottle Verification Tests|Bottle Verification Tests]]
- [[_COMMUNITY_Homebrew Bottle Verification|Homebrew Bottle Verification]]
- [[_COMMUNITY_Agent Rules|Agent Rules]]
- [[_COMMUNITY_Code of Conduct|Code of Conduct]]
- [[_COMMUNITY_Contributing Guide|Contributing Guide]]
- [[_COMMUNITY_Build and Attestation Workflow|Build and Attestation Workflow]]
- [[_COMMUNITY_README Introduction|README Introduction]]
- [[_COMMUNITY_Support Policy|Support Policy]]
- [[_COMMUNITY_Security and Verification Workflow|Security and Verification Workflow]]
- [[_COMMUNITY_Release Publication Workflow|Release Publication Workflow]]
- [[_COMMUNITY_Discover Command|Discover Command]]
- [[_COMMUNITY_Discover in Place|Discover in Place]]
- [[_COMMUNITY_Discover Interview|Discover Interview]]
- [[_COMMUNITY_Port Range Reporting Tests|Port Range Reporting Tests]]
- [[_COMMUNITY_Port Range Tests|Port Range Tests]]
- [[_COMMUNITY_Detector Registry and Tests|Detector Registry and Tests]]
- [[_COMMUNITY_Node Manifest Reading|Node Manifest Reading]]
- [[_COMMUNITY_Dotnet Detector|Dotnet Detector]]
- [[_COMMUNITY_Django Detector|Django Detector]]
- [[_COMMUNITY_Go Program Detector|Go Program Detector]]
- [[_COMMUNITY_JavaScript Framework Table|JavaScript Framework Table]]
- [[_COMMUNITY_Node Detector|Node Detector]]
- [[_COMMUNITY_Node Server Port Reading|Node Server Port Reading]]
- [[_COMMUNITY_Laravel and Queue Detector|Laravel and Queue Detector]]
- [[_COMMUNITY_Python and Swift Detectors|Python and Swift Detectors]]
- [[_COMMUNITY_Manual Entry Tests|Manual Entry Tests]]
- [[_COMMUNITY_Config Manual Page|Config Manual Page]]
- [[_COMMUNITY_Selection List Model|Selection List Model]]
- [[_COMMUNITY_Bounded Directory Walk|Bounded Directory Walk]]
- [[_COMMUNITY_Health Probe Redirects|Health Probe Redirects]]
- [[_COMMUNITY_Favicon Build|Favicon Build]]
- [[_COMMUNITY_Social Card Build|Social Card Build]]
- [[_COMMUNITY_Local Gate Script|Local Gate Script]]
- [[_COMMUNITY_Release Notes Script|Release Notes Script]]
- [[_COMMUNITY_Config Compatibility Promise|Config Compatibility Promise]]
- [[_COMMUNITY_Ports and Settings Documented|Ports and Settings Documented]]
- [[_COMMUNITY_Command Environment Documented|Command Environment Documented]]
- [[_COMMUNITY_Runtime Guarantees Documented|Runtime Guarantees Documented]]
- [[_COMMUNITY_Config Schema Documented|Config Schema Documented]]
- [[_COMMUNITY_Landing Page Structure|Landing Page Structure]]
- [[_COMMUNITY_Colour Options|Colour Options]]
- [[_COMMUNITY_Release Note Generator|Release Note Generator]]
- [[_COMMUNITY_Readiness Documented|Readiness Documented]]
- [[_COMMUNITY_Identifier Safety Tests|Identifier Safety Tests]]
- [[_COMMUNITY_Trailer Check Script|Trailer Check Script]]
- [[_COMMUNITY_Trailer Check Tests|Trailer Check Tests]]
- [[_COMMUNITY_Deno Detector|Deno Detector]]
- [[_COMMUNITY_Phoenix Detector|Phoenix Detector]]
- [[_COMMUNITY_Symfony Detector|Symfony Detector]]

## God Nodes (most connected - your core abstractions)
1. `join()` - 152 edges
2. `Contains()` - 104 edges
3. `Directory()` - 63 edges
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
- **The three-condition readiness contract** — documentation_readiness, documentation_start, documentation_status, documentation_port_variable, docs_index_process_tree_boundary [INFERRED 0.85]
- **Machine-wide port registry and allocation** — documentation_roles_and_ports, documentation_port_allocation, documentation_scan_directories, documentation_ports_audit, documentation_ports_assign, documentation_ports_reassign [INFERRED 0.85]
- **Proving a process is grat's own before signalling it** — documentation_start_identity, documentation_shutdown_sequence, documentation_stop, documentation_recover, documentation_status [INFERRED 0.85]

## Communities (99 total, 15 thin omitted)

### Community 3 - "CLI Entry and Dispatch"
Cohesion: 0.07
Nodes (61): Run(), Context, Writer, environment, Reader, Store, updateService, uninstallService (+53 more)

### Community 7 - "Service Manager Loading"
Cohesion: 0.06
Nodes (47): Manager, loadManager(), loadConfig(), Config, lifecycleProgressRenderer, ProgressEvent, StepKind, LifecycleEvent (+39 more)

### Community 14 - "Ports Command and Reassignment"
Cohesion: 0.14
Nodes (29): Report, portReassignment, ProjectConfig, Reservation, runPorts(), Context, environment, Renderer (+21 more)

### Community 8 - "CLI Integration Tests"
Cohesion: 0.12
Nodes (47): TestVersionCommandsRenderTheToolVersion(), T, TestExitCodeMapsInterruptedOperationsTo130(), TestInitAllocatesPortsForExplicitServices(), TestInitRejectsInvalidGlobalRegistry(), TestInitRejectsDeprecatedAppFlag(), TestRunRejectsUnknownCommand(), TestRunRejectsRemovedWorkerCommand() (+39 more)

### Community 72 - "CLI Helper Tests"
Cohesion: 0.67
Nodes (3): TestLogFollowUsesTrustedExecutable(), T, TestListenerOwnerLabelHandlesUnknownPID()

### Community 68 - "Log Streaming Tests"
Cohesion: 0.50
Nodes (3): notifyingWriter, TestOutputLogStreamsBeforeInputReachesEOF(), T

### Community 17 - "Config Types and Roles"
Cohesion: 0.16
Nodes (23): Config, Role, PortRange, Project, Runtime, Durations, Service, Load() (+15 more)

### Community 28 - "Config Rejection Tests"
Cohesion: 0.23
Nodes (18): DefaultRuntime(), TestLoadRejectsOversizedConfigBeforeParsing(), T, TestLoadRejectsUnknownFieldsWithStrictDecoder(), TestLoadRejectsLegacyShellConfig(), TestLoadRejectsDeprecatedAppsTable(), TestLoadRejectsRemovedGitHubWorkerConfiguration(), TestValidateRequiresAbsoluteHealthPathForPort() (+10 more)

### Community 6 - "Update Service and Installation"
Cohesion: 0.09
Nodes (26): Service, Context, Client, installation, Result, DefaultService(), runCommand(), runningBuildInfo() (+18 more)

### Community 11 - "Uninstall Artifact Scan"
Cohesion: 0.13
Nodes (24): installationKind, installation, uninstallArtifacts, artifactScanLimits, Context, Store, Reader, Writer (+16 more)

### Community 1 - "Uninstall Command"
Cohesion: 0.07
Nodes (69): TestUninstallDefaultYesRemovesOnlyRegisteredProjectArtifacts(), T, TestUninstallKeepsDeclinedArtifactClass(), TestUninstallRejectsNonInteractiveCleanup(), TestUninstallUsesOperationLockBeforePreflight(), TestUninstallAbortsBeforePromptsForActiveManagedService(), TestUninstallSkipsSymlinkedDirectoriesOutsideRegisteredRoots(), TestDiscoverUninstallArtifactsRejectsScanLimitOverrun() (+61 more)

### Community 51 - "Registry Lock"
Cohesion: 0.33
Nodes (9): WithRegistryLock(), Context, withRegistryLockIn(), TestRegistryLockUsesProvidedGratConfigurationDirectory(), T, TestRegistryLockHonorsContextWhileContended(), TestRegistryLockSerializesCallbacks(), TestRegistryLockReleasesAfterCallbackPanic() (+1 more)

### Community 33 - "Project Registry Scan"
Cohesion: 0.23
Nodes (16): scanLimits, scanCounters, Source, Reservation, ProjectConfig, Config, Problem, Report (+8 more)

### Community 42 - "Registry Scan Tests"
Cohesion: 0.26
Nodes (14): Scan(), fakeLookup, TestScanDoesNotExecuteConfig(), T, TestFirstFreeSkipsConfiguredAndLivePorts(), TestFirstFreeSkipsLivePortWhenOwnerPIDIsUnknown(), TestFirstFreeTreatsVisibleOwnerPIDAsOccupied(), TestAddListenersRecordsUnknownOwner() (+6 more)

### Community 56 - "Detector Result Types"
Cohesion: 0.32
Nodes (6): Role, Service, Role, Unresolved, Finding, detector

### Community 52 - "Help Screen Rendering"
Cohesion: 0.36
Nodes (5): Command, CommandGroup, Renderer, helpUsageWidth(), Style

### Community 43 - "Lifecycle Text Layout"
Cohesion: 0.33
Nodes (3): truncate(), LifecycleModel, truncateStyled()

### Community 23 - "Lifecycle Event Types"
Cohesion: 0.18
Nodes (20): LifecycleStage, LifecycleService, LifecycleGroup, LifecycleOperation, LifecycleEvent, LifecycleRunner, lifecycleRow, sanitizeLifecycleOperation() (+12 more)

### Community 4 - "Presentation Renderer"
Cohesion: 0.08
Nodes (60): NewLifecycleModel(), New(), DividerLine(), TestRendererUsesPlainTextForNonTerminalOutput(), T, TestRendererUsesSemanticColorWhenForced(), TestRunSpinnerRendersBeforeRunnerStarts(), TestRunSpinnerReturnsRunnerError() (+52 more)

### Community 16 - "Lifecycle Tea Model"
Cohesion: 0.14
Nodes (24): Msg, lifecycleTeaModel, CancelFunc, lifecycleCompleteMessage, lifecycleSpinnerMessage, Cmd, Model, View (+16 more)

### Community 12 - "Presentation Types"
Cohesion: 0.15
Nodes (14): ColorMode, StepKind, ProjectGroup, ProjectRowsOptions, ParseColorMode(), formatProjectRows(), fprintf(), fprint() (+6 more)

### Community 49 - "Spinner Runner"
Cohesion: 0.24
Nodes (9): SpinnerRunner, RunSpinner(), Context, Writer, renderSpinnerFrame(), runUpdate(), Context, updateService (+1 more)

### Community 32 - "Project Root Lookup"
Cohesion: 0.20
Nodes (14): FindRoot(), TestFindRootUsesNearestConfig(), T, TestFindRootReturnsNotFoundOutsideProject(), OwnedByCurrentUser(), ownedBy(), FileInfo, RefuseUnsafeConfig() (+6 more)

### Community 47 - "Service Log Files"
Cohesion: 0.21
Nodes (10): newServiceLogFile(), File, TestNewServiceLogFileTruncatesPreviousOutput(), T, isSymlinkRefusal(), TestALinkedLogIsRefusedRatherThanFollowed(), T, TestALinkedStateDirectoryKeepsItsPermissions() (+2 more)

### Community 13 - "Runtime Manager Interface"
Cohesion: 0.16
Nodes (15): Manager, Status, Service, Config, ListenerLookup, ProgressObserver, RecoveryCandidate, Context (+7 more)

### Community 22 - "Manager Lifecycle Tests"
Cohesion: 0.23
Nodes (22): TestFixtureManagerUsesRuntimeBackendPortRange(), T, TestStartAndStopRequiresOwnedHealthyListener(), TestStartLaunchesSelectedBackendBeforeSelectedConsumers(), TestRestartEmitsOrderedLifecycleEvents(), TestStartRejectsUnhealthyHTTPResponse(), TestStartGracefullyStopsPreviouslyStartedServicesWhenCancelled(), TestStatusIgnoresLegacyPIDFiles() (+14 more)

### Community 46 - "Process Group Signalling"
Cohesion: 0.27
Nodes (13): TestSignalGroupStopsIsolatedSession(), processState, validateManagedState(), validateLegacyManagedState(), signalManagedGroup(), Signal, signalGroup(), TestValidateManagedStateRejectsLegacyCoarseIdentity() (+5 more)

### Community 29 - "Process Stop and Readiness"
Cohesion: 0.16
Nodes (15): Context, loadedState, waitForExit(), Duration, readiness, Manager, Context, Service (+7 more)

### Community 54 - "Command Environment Tests"
Cohesion: 0.44
Nodes (9): TestLaunchDoesNotSourceLoginProfile(), T, TestCommandEnvironmentExcludesUnapprovedParentVariables(), TestCommandEnvironmentDerivesBackendURLForConsumer(), TestCommandEnvironmentPreservesApprovedBackendURLOverride(), TestCommandEnvironmentFallsBackWhenApprovedBackendURLIsAbsent(), TestCommandEnvironmentOmitsBackendURLForProviderAndAmbiguousTopology(), TestLaunchKeepsLogDestinationAvailableAfterManagerExit() (+1 more)

### Community 60 - "Linux Process Identity"
Cohesion: 0.43
Nodes (5): processIdentity(), linuxProcessStartTicks(), TestLinuxProcessStartTicksHandlesClosingParenthesisInCommand(), T, TestLinuxProcessStartTicksRejectsMismatchedPID()

### Community 57 - "Progress Observation"
Cohesion: 0.36
Nodes (5): ProgressStage, ProgressEvent, Service, ProgressObserver, Manager

### Community 38 - "Recovery of Legacy State"
Cohesion: 0.38
Nodes (16): processAlive(), TestRecoverStopsValidatedLegacyProcess(), T, TestRecoverRemovesStaleLegacyStateWithoutSignaling(), TestRecoveryCandidatesRejectsChangedLegacyIdentityWithoutSignaling(), TestRecoveryCandidatesRejectsChangedProcessGroupWithoutSignaling(), TestRecoveryCandidatesRejectsVersionTwoStateWithoutSignaling(), TestRecoverRejectsStateChangedAfterPreviewWithoutSignaling() (+8 more)

### Community 34 - "Managed State on Disk"
Cohesion: 0.22
Nodes (7): processState, Time, loadedState, RecoveryCandidate, Service, Manager, refuseLinkedDirectory()

### Community 27 - "Settings Store"
Cohesion: 0.27
Nodes (4): Settings, Store, canonicalExistingPath(), canonicalExistingDirectory()

### Community 35 - "Settings Store Tests"
Cohesion: 0.31
Nodes (17): TestStoreLoadReportsMissingSettings(), T, TestStoreAddCanonicalizesAndDeduplicatesDirectories(), TestStoreAddResolvesRelativeDirectoriesAgainstWorkingDirectory(), TestStoreAddRejectsMissingAndNonDirectoryPaths(), TestStoreRemovePersistsRemainingDirectories(), TestStoreRejectsInvalidSettingsDocuments(), TestStoreSaveUsesRestrictivePermissions() (+9 more)

### Community 58 - "Text Safety Filter"
Cohesion: 0.39
Nodes (6): UnsafeRune(), ContainsUnsafe(), Sanitize(), TestUnsafeRuneRejectsControlsAndUnicodeFormatCharacters(), T, TestSanitizeReplacesEveryUnsafeRune()

### Community 10 - "Version Reporting"
Cohesion: 0.06
Nodes (40): Current(), TestCurrentReturnsSourceVersion(), T, TestCurrentPrefixesLinkerOverrideWithV(), commandDocument(), Document, runManual(), Writer (+32 more)

### Community 61 - "Homebrew Bottle Build"
Cohesion: 0.48
Nodes (5): build-homebrew-bottles.sh script, usage(), write_formula(), package(), write_manual()

### Community 62 - "Homebrew Bottle Tests"
Cohesion: 0.52
Nodes (6): test-homebrew-bottles.sh script, assert_file(), assert_archive_contains(), assert_binary(), assert_mode(), assert_manual()

### Community 70 - "Homebrew Bottle Verification"
Cohesion: 0.60
Nodes (3): verify-homebrew-bottles.sh script, usage(), verify_bottle()

### Community 71 - "Agent Rules"
Cohesion: 0.67
Nodes (3): Project Agent Rules, Graphify Before Push, Versioned Graphify Artifacts

### Community 63 - "Contributing Guide"
Cohesion: 0.33
Nodes (5): Contributing to grat, Development setup, Pull requests, Configuration compatibility, Code of conduct

### Community 21 - "Build and Attestation Workflow"
Cohesion: 0.12
Nodes (24): grat README, README installation instructions, Cross-Platform Build Job, GitHub Artifact Attestation, MIT License, CGO_ENABLED=0 for Uniform Static Binaries, Verifying the binary before installing it, not after, Installation and Attestation Verification (+16 more)

### Community 18 - "README Introduction"
Cohesion: 0.11
Nodes (27): Quick start with discover, start and status, Does grat fit your project?, Readiness and Status, Who decides the port, Reading process.env.PORT and os.Getenv("PORT") from project source, Refusing Beats a Guessed Port, Worker role without a port, Why the listener must belong to grat's own process tree (+19 more)

### Community 80 - "Support Policy"
Cohesion: 0.50
Nodes (3): Support, Diagnostic Support Request, Sensitive Data Redaction

### Community 19 - "Security and Verification Workflow"
Cohesion: 0.12
Nodes (26): Verify Job, Platform Verification Matrix, Vulnerability Scan (govulncheck), Dependabot Weekly Updates, GitHub Actions Pinned by Commit Hash, Security Scan (gosec, both GOOS), Documentation Gate (scripts/check-docs.sh), Weekly Scheduled Security Scan (+18 more)

### Community 5 - "Release Publication Workflow"
Cohesion: 0.05
Nodes (62): Tag-Triggered Release, Release Publish Job, SHA-256 Release Checksums, GitHub Release Publication, check-docs.sh script, require(), require_in(), require_doc() (+54 more)

### Community 9 - "Discover Command"
Cohesion: 0.09
Nodes (42): candidate, Service, runDiscover(), Context, Reader, environment, Renderer, discoverBelow() (+34 more)

### Community 48 - "Discover in Place"
Cohesion: 0.24
Nodes (9): discoverHere(), Context, Reader, Renderer, serviceSuggestions(), Unresolved, detectServices(), TestInitReportsWhatItCouldNotResolve() (+1 more)

### Community 44 - "Discover Interview"
Cohesion: 0.30
Nodes (13): collectProjectInterview(), Reader, Writer, promptRequired(), promptServiceName(), promptDefault(), parseServiceDefinition(), validateServiceDefinitions() (+5 more)

### Community 73 - "Port Range Reporting Tests"
Cohesion: 0.67
Nodes (3): TestAPortOutsideItsRangeDoesNotBlockLoading(), T, TestAPortInsideItsRangeIsNotReported()

### Community 69 - "Port Range Tests"
Cohesion: 0.60
Nodes (4): TestEveryRoleOwnsARangeOfTheStatedWidth(), T, TestRolesWithDifferentRangesDoNotOverlap(), TestAWorkerOwnsNoRange()

### Community 0 - "Detector Registry and Tests"
Cohesion: 0.06
Nodes (90): Directory(), writeProject(), T, commandOf(), Finding, TestNodeYieldsOneServicePerConventionalScript(), TestNodePrefersDevFrontendOverDev(), TestNodeUsesPnpmWhenTheManifestSaysSo() (+82 more)

### Community 45 - "Node Manifest Reading"
Cohesion: 0.15
Nodes (7): fileExists(), detectRails(), Service, Unresolved, detectSpringBoot(), Service, Unresolved

### Community 26 - "Dotnet Detector"
Cohesion: 0.15
Nodes (17): readBounded(), entries(), DirEntry, detectDotnet(), Service, Unresolved, detectFlask(), Service (+9 more)

### Community 74 - "Django Detector"
Cohesion: 0.50
Nodes (3): detectDjango(), Service, Unresolved

### Community 40 - "Go Program Detector"
Cohesion: 0.20
Nodes (14): goProgram, detectGo(), Service, Unresolved, goPrograms(), holdsMainPackage(), readsPortFromGoSource(), isPortEnvironmentCall() (+6 more)

### Community 31 - "JavaScript Framework Table"
Cohesion: 0.20
Nodes (16): jsFramework, detectJavaScriptFramework(), manifest, Service, Unresolved, frameworkEvidence(), detectHugo(), Service (+8 more)

### Community 30 - "Node Detector"
Cohesion: 0.19
Nodes (16): readManifest(), Unresolved, detectNode(), Service, Unresolved, namedServices(), manifest, detectSingleService() (+8 more)

### Community 39 - "Node Server Port Reading"
Cohesion: 0.20
Nodes (14): detectNodeServer(), manifest, Service, Unresolved, readsPortFromEnvironment(), servesWithBun(), bunLockfile(), detectBun() (+6 more)

### Community 55 - "Laravel and Queue Detector"
Cohesion: 0.42
Nodes (8): detectLaravel(), Service, Unresolved, laravelQueueWorker(), laravelQueueConnection(), configuredQueueFallback(), environmentValue(), environmentAssignment()

### Community 36 - "Python and Swift Detectors"
Cohesion: 0.15
Nodes (14): detectPython(), Service, Unresolved, applicationModules(), detectVapor(), Service, Unresolved, safeIdentifier() (+6 more)

### Community 77 - "Manual Entry Tests"
Cohesion: 0.67
Nodes (3): TestAnEntryCarriesMoreThanTheOneLineReference(), T, TestEveryFlagOfAnEntryIsExplained()

### Community 2 - "Config Manual Page"
Cohesion: 0.06
Nodes (53): ConfigPage(), ConfigDocument(), Document, field, fieldList(), Block, runtimeFields(), roleRanges() (+45 more)

### Community 20 - "Selection List Model"
Cohesion: 0.14
Nodes (11): SelectionItem, SelectionModel, NewSelectionModel(), rows(), SelectionItem, TestTheCursorStopsAtBothEndsRatherThanWrapping(), T, TestOnlyMarkedRowsComeBack() (+3 more)

### Community 24 - "Bounded Directory Walk"
Cohesion: 0.18
Nodes (18): SkipsScanning(), DeeperThanScan(), Walk(), DirEntry, ErrTooManyEntries, TestTheScanSkipsWhatCannotHoldAProject(), T, TestTheScanStopsBelowItsDepth() (+10 more)

### Community 67 - "Health Probe Redirects"
Cohesion: 0.60
Nodes (5): TestAHealthProbeDoesNotFollowARedirectOffTheService(), T, TestAHealthProbeFollowsARedirectOnTheSameService(), probeClient(), Client

### Community 78 - "Local Gate Script"
Cohesion: 0.67
Nodes (3): gates.sh script, GOTOOLCHAIN, step()

### Community 53 - "Ports and Settings Documented"
Cohesion: 0.29
Nodes (10): Roles and Port Ranges, Scan Directories, Background update check and GRAT_NO_UPDATE_CHECK, grat ports audit, grat ports assign, grat/settings.toml and update-check in the user config directory, Port allocation across registered directories, grat ports audit, assign, reassign (+2 more)

### Community 50 - "Command Environment Documented"
Cohesion: 0.38
Nodes (11): Non-secret environment baseline for a service command, BACKEND_URL for a Single Backend, inherit_env, The Environment Is Small on Purpose, PORT environment variable owned by grat, A [[services]] table, Limiting accidental secret propagation into services, What It Runs Section (+3 more)

### Community 25 - "Runtime Guarantees Documented"
Cohesion: 0.18
Nodes (21): Exit Status Codes, Start Identity Before Any Signal, Shutdown sequence: SIGTERM to the process group, then SIGKILL, Services Section: Ready Means Ready, .grat/pid and .grat/log managed state, Command Reference, grat start, grat stop (+13 more)

### Community 37 - "Config Schema Documented"
Cohesion: 0.12
Nodes (17): /bin/sh execution is a trust boundary, not a sandbox, runtime table timing overrides, grat.config(7) Manual, Service states: stopped, running, unhealthy, grat.config schema: version, project, runtime, services, The File Is Read as Data, Platform helpers invoked through fixed absolute paths, Written commands carry host and port, because binding every interface is accidental exposure (+9 more)

### Community 15 - "Landing Page Structure"
Cohesion: 0.11
Nodes (28): SoftwareApplication Structured Data, grat Share Card, robots.txt Crawl Policy, Install Section, Ports Section: Every Role Owns a Lane, grat.layered.work landing page, Hero Section, Reduced motion gets no animation (+20 more)

### Community 41 - "Release Note Generator"
Cohesion: 0.16
Nodes (16): release-notes.sh, release note generator, Tag range resolution, Commit records from git log, awk grouping program, Prefix to heading mapping, Chore, Test and Refactor are dropped, Release-note trailer, A trailer overrules a dropped prefix (+8 more)

### Community 65 - "Readiness Documented"
Cohesion: 0.33
Nodes (6): What a service has to do, Port ownership follows the process chain, 03. Services, ready means ready, Readiness is a listener plus a health path, Stop signals only on a matching identity, Logs in .grat/log/

### Community 64 - "Identifier Safety Tests"
Cohesion: 0.53
Nodes (5): TestADetectedNameNeverBecomesASecondCommand(), T, TestAnUnsafeNameSaysWhichCharacter(), TestAnOrdinaryNameStillYieldsACommand(), TestTheCharacterSetIsWhatItSays()

### Community 79 - "Trailer Check Tests"
Cohesion: 0.83
Nodes (3): test-check-trailers.sh script, commit(), expect()

### Community 59 - "Deno Detector"
Cohesion: 0.43
Nodes (6): detectDeno(), Service, Unresolved, denoTasks(), chooseDenoTask(), readsPortFromDenoSource()

### Community 75 - "Phoenix Detector"
Cohesion: 0.50
Nodes (3): detectPhoenix(), Service, Unresolved

### Community 76 - "Symfony Detector"
Cohesion: 0.50
Nodes (3): detectSymfony(), Service, Unresolved

## Knowledge Gaps
- **211 isolated node(s):** `Reader`, `Store`, `Manager`, `Config`, `StepKind` (+206 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **15 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `join()` connect `Uninstall Command` to `Detector Registry and Tests`, `Config Manual Page`, `CLI Entry and Dispatch`, `Presentation Renderer`, `Update Service and Installation`, `Service Manager Loading`, `CLI Integration Tests`, `Discover Command`, `Version Reporting`, `Uninstall Artifact Scan`, `Presentation Types`, `Runtime Manager Interface`, `Ports Command and Reassignment`, `Config Types and Roles`, `Selection List Model`, `Manager Lifecycle Tests`, `Bounded Directory Walk`, `Dotnet Detector`, `Settings Store`, `Node Detector`, `JavaScript Framework Table`, `Project Root Lookup`, `Project Registry Scan`, `Python and Swift Detectors`, `Node Server Port Reading`, `Go Program Detector`, `Lifecycle Text Layout`, `Node Manifest Reading`, `Discover in Place`, `Registry Lock`, `Laravel and Queue Detector`, `Detector Result Types`, `Deno Detector`, `Service Launch Environment`, `Django Detector`, `Phoenix Detector`, `Symfony Detector`?**
  _High betweenness centrality (0.347) - this node is a cross-community bridge._
- **Why does `Contains()` connect `Presentation Renderer` to `Detector Registry and Tests`, `Uninstall Command`, `Config Manual Page`, `CLI Entry and Dispatch`, `Python and Swift Detectors`, `Health Probe Redirects`, `Settings Store Tests`, `CLI Integration Tests`, `Discover Command`, `Version Reporting`, `Uninstall Artifact Scan`, `Discover Interview`, `Process Group Signalling`, `Discover in Place`, `Selection List Model`, `Manager Lifecycle Tests`, `Settings Store`, `Config Rejection Tests`?**
  _High betweenness centrality (0.122) - this node is a cross-community bridge._
- **Why does `New()` connect `Presentation Renderer` to `Uninstall Command`, `CLI Entry and Dispatch`, `Settings Store Tests`, `Update Service and Installation`, `Service Manager Loading`, `CLI Integration Tests`, `Version Reporting`, `Uninstall Artifact Scan`, `Discover Interview`, `Presentation Types`, `Runtime Manager Interface`, `Settings Store`?**
  _High betweenness centrality (0.073) - this node is a cross-community bridge._
- **Are the 151 inferred relationships involving `join()` (e.g. with `loadConfig()` and `loadPortFixtureConfig()`) actually correct?**
  _`join()` has 151 INFERRED edges - model-reasoned connections that need verification._
- **Are the 100 inferred relationships involving `Contains()` (e.g. with `TestInitRejectsDeprecatedAppFlag()` and `TestInitRejectsInvalidGlobalRegistry()`) actually correct?**
  _`Contains()` has 100 INFERRED edges - model-reasoned connections that need verification._
- **Are the 61 inferred relationships involving `Directory()` (e.g. with `discoverCandidates()` and `detectServices()`) actually correct?**
  _`Directory()` has 61 INFERRED edges - model-reasoned connections that need verification._
- **Are the 44 inferred relationships involving `New()` (e.g. with `runWithEnvironment()` and `TestRenderPortReassignSummaryGroupsAssignmentsByProject()`) actually correct?**
  _`New()` has 44 INFERRED edges - model-reasoned connections that need verification._
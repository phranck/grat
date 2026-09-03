package detect

import (
	"strings"
)

const (
	// springBootGradlePlugin is how both Gradle manifests name the plugin.
	springBootGradlePlugin = "org.springframework.boot"

	// springBootMavenParent and springBootMavenPlugin are the two ways a Maven
	// project declares itself a Spring Boot one. A project can use either.
	springBootMavenParent = "spring-boot-starter-parent"
	springBootMavenPlugin = "spring-boot-maven-plugin"
)

const (
	// springBootMavenCommand starts the application through Maven, which forks
	// a child JVM to run it. That child is below the mvn process grat started,
	// so the listener traces back to grat and readiness arrives.
	springBootMavenCommand = "SERVER_PORT=$PORT mvn spring-boot:run"

	// springBootGradleCommand starts it through Gradle, and --no-daemon is a
	// condition rather than a preference.
	//
	// bootRun is a JavaExec, so the application is a child JVM either way. What
	// the flag decides is whose child it is. Gradle runs its build in a
	// long-lived daemon by default, and that daemon is a process grat did not
	// start, so the application would hang below it and the listener would
	// never trace back to the command grat ran. Without the flag the port is
	// held, the application works, and readiness never arrives.
	springBootGradleCommand = "SERVER_PORT=$PORT ./gradlew bootRun --no-daemon"
)

// detectSpringBoot recognises a Spring Boot application and which build tool
// starts it, because the answer changes the command.
//
// The port arrives as SERVER_PORT rather than PORT. Spring Boot's relaxed
// binding maps an environment variable onto a property by lower-casing it and
// turning underscores into dots, so SERVER_PORT is what reaches server.port,
// and PORT reaches nothing at all.
func detectSpringBoot(root string) ([]Service, []Unresolved) {
	for _, manifest := range []string{"build.gradle.kts", "build.gradle"} {
		data, ok := readBounded(join(root, manifest))
		if !ok {
			continue
		}
		if !strings.Contains(string(data), springBootGradlePlugin) {
			continue
		}
		if !fileExists(join(root, "gradlew")) {
			return nil, []Unresolved{{
				Marker: manifest,
				Reason: "the project uses the Spring Boot Gradle plugin but carries no gradlew wrapper, and grat proposes no command that depends on a Gradle installed elsewhere",
			}}
		}
		return []Service{service("backend", springBootGradleCommand, manifest)}, nil
	}

	data, ok := readBounded(join(root, "pom.xml"))
	if !ok {
		return nil, nil
	}
	pom := string(data)
	if !strings.Contains(pom, springBootMavenParent) && !strings.Contains(pom, springBootMavenPlugin) {
		return nil, nil
	}
	return []Service{service("backend", springBootMavenCommand, "pom.xml")}, nil
}

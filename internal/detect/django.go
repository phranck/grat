package detect

import "regexp"

// djangoManagementPattern matches the call that every Django manage.py makes.
// The file name alone proves nothing, because manage.py is a common name for a
// project's own script; this line is what Django's own project template writes.
var djangoManagementPattern = regexp.MustCompile(`execute_from_command_line`)

// djangoSettingsPattern matches manage.py naming the settings module, which is
// the second half of what the template writes and what tells the command which
// project it is running.
var djangoSettingsPattern = regexp.MustCompile(`DJANGO_SETTINGS_MODULE`)

// detectDjango recognises a Django project.
//
// Django is the one Python framework where nothing has to be guessed. The marker
// is exact, the command needs no module name because manage.py sets the settings
// module itself, and the address is one positional argument. Source: the
// django-admin reference at docs.djangoproject.com and the project template
// django/conf/project_template/manage.py-tpl.
func detectDjango(root string) ([]Service, []Unresolved) {
	data, ok := readBounded(join(root, "manage.py"))
	if !ok {
		return nil, nil
	}
	if !djangoManagementPattern.Match(data) {
		// A script that merely shares the name runs something else entirely.
		return nil, nil
	}
	if !djangoSettingsPattern.Match(data) {
		return nil, []Unresolved{{
			Marker: "manage.py",
			Reason: "manage.py runs Django commands but names no settings module, so which project it serves is not decidable",
		}}
	}

	return []Service{service("backend", "python manage.py runserver 127.0.0.1:$PORT", "manage.py")}, nil
}

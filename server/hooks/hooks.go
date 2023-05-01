package hooks

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/soft-serve/server/config"
	"github.com/charmbracelet/soft-serve/server/utils"
)

// The names of git server-side hooks.
const (
	PreReceiveHook  = "pre-receive"
	UpdateHook      = "update"
	PostReceiveHook = "post-receive"
	PostUpdateHook  = "post-update"
)

// GenerateHooks generates git server-side hooks for a repository. Currently, it supports the following hooks:
// - pre-receive
// - update
// - post-receive
// - post-update
//
// This function should be called by the backend when a repository is created.
// TODO: support context
func GenerateHooks(ctx context.Context, cfg *config.Config, repo string) error {
	repo = utils.SanitizeRepo(repo) + ".git"
	hooksPath := filepath.Join(cfg.DataPath, "repos", repo, "hooks")
	if err := os.MkdirAll(hooksPath, os.ModePerm); err != nil {
		return err
	}

	ex, err := os.Executable()
	if err != nil {
		return err
	}

	dp, err := filepath.Abs(cfg.DataPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for data path: %w", err)
	}

	cp := filepath.Join(dp, "config.yaml")
	// Add extra environment variables to the hooks here.
	envs := []string{}

	for _, hook := range []string{
		PreReceiveHook,
		UpdateHook,
		PostReceiveHook,
		PostUpdateHook,
	} {
		var data bytes.Buffer
		var args string

		// Hooks script/directory path
		hp := filepath.Join(hooksPath, hook)

		// Write the hooks primary script
		if err := os.WriteFile(hp, []byte(hookTemplate), os.ModePerm); err != nil {
			return err
		}

		// Create ${hook}.d directory.
		hp += ".d"
		if err := os.MkdirAll(hp, os.ModePerm); err != nil {
			return err
		}

		switch hook {
		case UpdateHook:
			args = "$1 $2 $3"
		case PostUpdateHook:
			args = "$@"
		}

		if err := hooksTmpl.Execute(&data, struct {
			Executable string
			Config     string
			Envs       []string
			Hook       string
			Args       string
		}{
			Executable: ex,
			Config:     cp,
			Envs:       envs,
			Hook:       hook,
			Args:       args,
		}); err != nil {
			log.WithPrefix("backend.hooks").Error("failed to execute hook template", "err", err)
			continue
		}

		// Write the soft-serve hook inside ${hook}.d directory.
		hp = filepath.Join(hp, "soft-serve")
		err = os.WriteFile(hp, data.Bytes(), os.ModePerm) //nolint:gosec
		if err != nil {
			log.WithPrefix("backend.hooks").Error("failed to write hook", "err", err)
			continue
		}
	}

	return nil
}

const (
	// hookTemplate allows us to run multiple hooks from a directory. It should
	// support every type of git hook, as it proxies both stdin and arguments.
	hookTemplate = `#!/usr/bin/env bash
# AUTO GENERATED BY SOFT SERVE, DO NOT MODIFY
data=$(cat)
exitcodes=""
hookname=$(basename $0)
GIT_DIR=${GIT_DIR:-$(dirname $0)/..}
for hook in ${GIT_DIR}/hooks/${hookname}.d/*; do
  # Avoid running non-executable hooks
  test -x "${hook}" && test -f "${hook}" || continue

  # Run the actual hook
  echo "${data}" | "${hook}" "$@"

  # Store the exit code for later use
  exitcodes="${exitcodes} $?"
done

# Exit on the first non-zero exit code.
for i in ${exitcodes}; do
  [ ${i} -eq 0 ] || exit ${i}
done
`
)

var (
	// hooksTmpl is the soft-serve hook that will be run by the git hooks
	// inside the hooks directory.
	hooksTmpl = template.Must(template.New("hooks").Parse(`#!/usr/bin/env bash
# AUTO GENERATED BY SOFT SERVE, DO NOT MODIFY
if [ -z "$SOFT_SERVE_REPO_NAME" ]; then
	echo "Warning: SOFT_SERVE_REPO_NAME not defined. Skipping hooks."
	exit 0
fi
{{ range $_, $env := .Envs }}
{{ $env }} \{{ end }}
{{ .Executable }} hook --config "{{ .Config }}" {{ .Hook }} {{ .Args }}
`))
)

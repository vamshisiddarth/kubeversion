package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/yourusername/kubeversion/internal/download"
)

type Release struct {
	TagName string `json:"tag_name"`
}

func fetchAvailableVersions() ([]string, error) {
	resp, err := http.Get("https://api.github.com/repos/kubernetes/kubernetes/releases")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch versions: HTTP %d", resp.StatusCode)
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to decode versions: %w", err)
	}

	versions := make([]string, 0, len(releases))
	for _, release := range releases {
		if strings.HasPrefix(release.TagName, "v") {
			versions = append(versions, release.TagName)
		}
	}

	// Sort versions in descending order
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))
	return versions, nil
}

func getInstalledVersions() (map[string]bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	kubectlDir := filepath.Join(homeDir, ".kubeversion", "versions")
	entries, err := os.ReadDir(kubectlDir)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]bool), nil
		}
		return nil, fmt.Errorf("failed to read versions directory: %w", err)
	}

	installed := make(map[string]bool)
	for _, entry := range entries {
		version := strings.TrimPrefix(entry.Name(), "kubectl-")
		installed[version] = true
	}

	return installed, nil
}

func ListVersions() error {
	availableVersions, err := fetchAvailableVersions()
	if err != nil {
		return fmt.Errorf("failed to get available versions: %w", err)
	}

	installedVersions, err := getInstalledVersions()
	if err != nil {
		return fmt.Errorf("failed to get installed versions: %w", err)
	}

	if !isKubeversionInPath() {
		fmt.Printf("\nNOTE: kubeversion's bin directory is not in your PATH.\n")
		fmt.Printf("To use managed kubectl versions, add this to your shell config:\n")
		fmt.Printf("export PATH=\"$HOME/.kubeversion/bin:$PATH\"\n\n")
	}

	// Create templates for the prompt
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F449 {{ . | cyan }}",
		Inactive: "  {{ . }}",
		Selected: "\U0001F44D {{ . | green }}",
		Details: `
{{ "Version:" | faint }}  {{ . }}
{{ if .Installed }}{{ "Status:" | faint }}  Already Installed{{ else }}{{ "Status:" | faint }}  Not Installed{{ end }}`,
	}

	// Create version items with installation status
	type versionItem struct {
		Version    string
		Installed  bool
		displayStr string
	}

	var items []versionItem
	for _, v := range availableVersions {
		item := versionItem{
			Version:   v,
			Installed: installedVersions[v],
		}
		if item.Installed {
			item.displayStr = fmt.Sprintf("%s (installed)", v)
		} else {
			item.displayStr = v
		}
		items = append(items, item)
	}

	prompt := promptui.Select{
		Label:     "Select Kubectl Version",
		Items:     items,
		Templates: templates,
		Size:      10, // Show 10 items at a time
		Searcher: func(input string, index int) bool {
			item := items[index]
			return strings.Contains(strings.ToLower(item.Version), strings.ToLower(input))
		},
		Stdout: &bellSkipper{}, // Skip the bell sound
	}

	index, _, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			return nil // User cancelled
		}
		return fmt.Errorf("prompt failed: %w", err)
	}

	selectedVersion := items[index].Version
	if items[index].Installed {
		return SwitchVersion(selectedVersion)
	}

	// Install and then switch to the selected version
	if err := download.InstallVersion(selectedVersion); err != nil {
		return fmt.Errorf("failed to install version %s: %w", selectedVersion, err)
	}

	return SwitchVersion(selectedVersion)
}

// bellSkipper implements an io.Writer that skips the terminal bell character.
type bellSkipper struct{}

func (bs *bellSkipper) Write(b []byte) (int, error) {
	const charBell = 7 // Bell character
	if len(b) == 1 && b[0] == charBell {
		return 0, nil
	}
	return os.Stdout.Write(b)
}

// Add Close method to implement io.WriteCloser
func (bs *bellSkipper) Close() error {
	return nil
}

func SwitchVersion(version string) error {
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	kubectlDir := filepath.Join(homeDir, ".kubeversion", "versions")
	versionPath := filepath.Join(kubectlDir, fmt.Sprintf("kubectl-%s", version))

	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		return fmt.Errorf("version %s is not installed. Use 'kubeversion install %s' first", version, version)
	}

	binDir := filepath.Join(homeDir, ".kubeversion", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	kubectlPath := filepath.Join(binDir, "kubectl")
	if err := os.RemoveAll(kubectlPath); err != nil {
		return fmt.Errorf("failed to remove existing kubectl: %w", err)
	}

	if err := os.Symlink(versionPath, kubectlPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	// Check if PATH is properly configured
	if !isKubeversionInPath() {
		fmt.Printf("\nIMPORTANT: To use kubeversion's kubectl, add this to your shell config (~/.bashrc, ~/.zshrc, etc.):\n")
		fmt.Printf("export PATH=\"$HOME/.kubeversion/bin:$PATH\"\n")
		fmt.Printf("Then restart your shell or run: source ~/.bashrc (or ~/.zshrc)\n\n")
	}

	fmt.Printf("Successfully switched to kubectl %s\n", version)
	return nil
}

// isKubeversionInPath checks if kubeversion's bin directory is in PATH
func isKubeversionInPath() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	kubeversionBin := filepath.Join(homeDir, ".kubeversion", "bin")
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	// Check if our bin directory is in PATH
	for _, path := range paths {
		if path == kubeversionBin {
			return true
		}
	}
	return false
}

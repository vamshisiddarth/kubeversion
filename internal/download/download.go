package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/schollz/progressbar/v3"
)

const (
	downloadURLFormat = "https://dl.k8s.io/release/%s/bin/%s/%s/kubectl"
	versionPrefix     = "v"
)

func InstallVersion(version string) error {
	if !strings.HasPrefix(version, versionPrefix) {
		version = versionPrefix + version
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	kubectlDir := filepath.Join(homeDir, ".kubeversion", "versions")
	if err := os.MkdirAll(kubectlDir, 0755); err != nil {
		return fmt.Errorf("failed to create kubectl directory: %w", err)
	}

	url := fmt.Sprintf(downloadURLFormat, version, runtime.GOOS, runtime.GOARCH)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download kubectl: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download kubectl: HTTP %d", resp.StatusCode)
	}

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading kubectl",
	)

	destPath := filepath.Join(kubectlDir, fmt.Sprintf("kubectl-%s", version))
	out, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("failed to create kubectl file: %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(io.MultiWriter(out, bar), resp.Body); err != nil {
		return fmt.Errorf("failed to save kubectl: %w", err)
	}

	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// ImageSize represents the size of the image for a specific architecture.
type ImageSize struct {
	Architecture string `json:"architecture"`
	Size         int64  `json:"size"`
}

// Response represents the response to be returned to the user.
type Response struct {
	Image string           `json:"image"`
	Sizes map[string]int64 `json:"sizes"`
}

func main() {
	http.HandleFunc("/get-uncompressed-size", getUncompressedSizeHandler)
	fmt.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", nil)
}

// getUncompressedSizeHandler handles the /get-uncompressed-size route.
func getUncompressedSizeHandler(w http.ResponseWriter, r *http.Request) {
	image := r.URL.Query().Get("image")
	if image == "" {
		http.Error(w, "Missing 'image' parameter", http.StatusBadRequest)
		return
	}

	architectures, err := getSupportedArchitectures(image)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting architectures: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Determine the download directory. /tmp is used by default.
	downloadDir := os.Getenv("DOWNLOAD_DIR")
	if downloadDir == "" {
		downloadDir = "/tmp"
	}

	tempDir, err := os.MkdirTemp(downloadDir, "skopeo_downloads-*")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating temp directory: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	defer os.RemoveAll(tempDir)

	imageName := sanitizeImageName(image)

	var wg sync.WaitGroup
	results := make(chan ImageSize, len(architectures))

	for _, arch := range architectures {
		wg.Add(1)
		go func(architecture string) {
			defer wg.Done()
			filePath := filepath.Join(tempDir, fmt.Sprintf("%s_%s.tar", imageName, architecture))
			size, err := downloadImageAndGetSize(image, architecture, filePath)
			if err != nil {
				fmt.Printf("Error downloading image for architecture %s: %s\n", architecture, err.Error())
				return
			}
			results <- ImageSize{Architecture: architecture, Size: size}
		}(arch)
	}

	wg.Wait()
	close(results)

	sizes := make(map[string]int64)
	for result := range results {
		sizes[result.Architecture] = result.Size
	}

	response := Response{Image: image, Sizes: sizes}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// sanitizeImageName replaces slashes and colons in the image name with underscores.
func sanitizeImageName(image string) string {
	return strings.NewReplacer("/", "_", ":", "_").Replace(image)
}

// getSupportedArchitectures gets the list of supported architectures for a Docker image.
func getSupportedArchitectures(image string) ([]string, error) {
	cmdArgs := GenerateSkopeoInspectCmdArgs(image)
	cmd := exec.Command("skopeo", cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var manifest struct {
		Manifests []struct {
			Platform struct {
				Architecture string `json:"architecture"`
			} `json:"platform"`
		} `json:"manifests"`
	}
	if err := json.Unmarshal(output, &manifest); err != nil {
		return nil, err
	}

	var architectures []string
	for _, m := range manifest.Manifests {
		architectures = append(architectures, m.Platform.Architecture)
	}

	// Ensure at least "amd64" is included if no architectures were found
	if len(architectures) == 0 {
		architectures = []string{"amd64"}
	}

	return architectures, nil
}

// downloadImageAndGetSize downloads the Docker image for a specific architecture and returns its uncompressed size.
func downloadImageAndGetSize(image, architecture, filePath string) (int64, error) {
	cmdArgs := GenerateSkopeoCmdArgs(image, filePath, architecture)

	fmt.Printf("Executing skopeo with arguments: %v\n", cmdArgs)
	cmd := exec.Command("skopeo", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("skopeo output: %s, error: %s", string(output), err.Error())
	}

	fmt.Printf("skopeo output for architecture %s: %s\n", architecture, string(output))

	// Ensure the file was created
	if _, err := os.Stat(filePath); err != nil {
		return 0, fmt.Errorf("error verifying file creation: %s", err.Error())
	}

	size, err := getFileSize(filePath)
	if err != nil {
		return 0, err
	}

	return size, nil
}

// GenerateSkopeoCmdArgs generates the command line arguments for the skopeo copy command based on environment variables and input parameters.
func GenerateSkopeoCmdArgs(imageName, targetFilename, architecture string) []string {
	cmdArgs := []string{"copy", "--remove-signatures"}

	// Check and add registry credentials if they are set
	registryUsername, usernameSet := os.LookupEnv("REGISTRY_USERNAME")
	registryPassword, passwordSet := os.LookupEnv("REGISTRY_PASSWORD")

	if usernameSet && passwordSet {
		cmdArgs = append(cmdArgs, "--src-username", registryUsername, "--src-password", registryPassword)
	}

	// Add architecture override if specified
	if architecture != "" {
		cmdArgs = append(cmdArgs, "--override-arch", architecture)
	}

	// Add the rest of the command source image and destination tar
	cmdArgs = append(cmdArgs, fmt.Sprintf("docker://%s", imageName), fmt.Sprintf("docker-archive://%s", targetFilename))

	return cmdArgs
}

// GenerateSkopeoInspectCmdArgs generates the command line arguments for the skopeo inspect to fetch all supported architectures
func GenerateSkopeoInspectCmdArgs(imageName string) []string {
	return []string{"inspect", "--raw", fmt.Sprintf("docker://%s", imageName)}
}

// getFileSize returns the size of the file at the given path in bytes.
func getFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

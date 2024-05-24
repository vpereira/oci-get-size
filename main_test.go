package main

import (
	"os"
	"testing"
)

func TestSanitizeImageName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"registry.suse.com/bci/bci-busybox:latest", "registry.suse.com_bci_bci-busybox_latest"},
		{"docker.io/library/nginx:1.19.6", "docker.io_library_nginx_1.19.6"},
		{"quay.io/coreos/etcd:v3.4.13", "quay.io_coreos_etcd_v3.4.13"},
		{"gcr.io/google-containers/pause:3.2", "gcr.io_google-containers_pause_3.2"},
		{"myregistry.local:5000/test/image:v1", "myregistry.local_5000_test_image_v1"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := sanitizeImageName(test.input)
			if result != test.expected {
				t.Errorf("sanitizeImageName(%q) = %q; want %q", test.input, result, test.expected)
			}
		})
	}
}

func TestGetFileSize(t *testing.T) {
	// Test case 1: Test for an existing empty file
	emptyFileName := "empty_test_file"
	f, err := os.Create(emptyFileName)
	if err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
	}
	f.Close()
	defer os.Remove(emptyFileName) // Clean up

	size, err := getFileSize(emptyFileName)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if size != 0 {
		t.Errorf("Expected size 0, got %d", size)
	}

	// Test case 2: Test for the existing main.go file
	size, err = getFileSize("main.go")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if size == 0 {
		t.Errorf("Expected size greater than 0, got %d", size)
	}
}

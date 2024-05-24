package main

import (
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

package apt_test

import (
	"context"
	"time"
)

// func TestInstall(t *testing.T) {
// 	tests := []testingmodels.TestInfo{
// 		//testingmodels.NewTestInfo("", ""),
// 		testingmodels.NewTestInfo("nginx", ""),
// 		testingmodels.NewTestInfo("nginx", "1.18.0-6ubuntu14.3"),
// 		//testingmodels.NewTestInfo("nginx", "1.18.0-6ubuntu14.3=abc"),
// 	}

// 	for _, tc := range tests {
// 		tc := tc

// 		t.Run(tc.String(), func(t *testing.T) {
// 			installer := apt.NewAptInstaller()
// 			context, cancel := CreateTestContext()
// 			err := installer.Install(context, tc.Input)
// 			assert.NilError(t, err)
// 			defer cancel()
// 		})
// 	}
// }

func CreateTestContext() (context.Context, context.CancelFunc) {
	const testTimeout = time.Minute * 1
	ctx := context.Background()
	return context.WithTimeout(ctx, testTimeout)
}

// func TestExtractVersion(t *testing.T) {
// 	t.Parallel()

// 	tests := []struct {
// 		input    string
// 		expected string
// 	}{
// 		{
// 			input:    "",
// 			expected: "",
// 		},
// 		{
// 			input: `
// Package: nginx
// Status: install ok installed
// Priority: optional
// Section: httpd
// Installed-Size: 49
// Maintainer: Ubuntu Developers <ubuntu-devel-discuss@lists.ubuntu.com>
// Architecture: arm64
// Version: 1.18.0-6ubuntu14.3
// Depends: nginx-core (<< 1.18.0-6ubuntu14.3.1~) | nginx-full (<< 1.18.0-6ubuntu14.3.1~) | nginx-light (<< 1.18.0-6ubuntu14.3.1~) | nginx-extras (<< 1.18.0-6ubuntu14.3.1~), nginx-core (>= 1.18.0-6ubuntu14.3) | nginx-full (>= 1.18.0-6ubuntu14.3) | nginx-light (>= 1.18.0-6ubuntu14.3) | nginx-extras (>= 1.18.0-6ubuntu14.3)
// Breaks: libnginx-mod-http-lua (<< 1.18.0-6ubuntu5)
// Description: small, powerful, scalable web/proxy server
//  Nginx ("engine X") is a high-performance web and reverse proxy server
//  created by Igor Sysoev. It can be used both as a standalone web server
//  and as a proxy to reduce the load on back-end HTTP or mail servers.
//  .
//  This is a dependency package to install either nginx-core (by default),
//  nginx-full, nginx-light or nginx-extras.
// Homepage: https://nginx.net
// Original-Maintainer: Debian Nginx Maintainers <pkg-nginx-maintainers@alioth-lists.debian.net>
// `,
// 			expected: "1.18.0-6ubuntu14.3",
// 		},
// 	}

// 	for _, tc := range tests {
// 		tc := tc

// 		t.Run(tc.input, func(t *testing.T) {
// 			t.Parallel()

// 			actual := apt.ExtractVersion(tc.input)
// 			assert.Equal(t, actual, tc.expected)
// 		})
// 	}
// }

package missing_header_fields

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/scroll-tech/go-ethereum/common"
	"github.com/scroll-tech/go-ethereum/log"
)

func TestManagerDownload(t *testing.T) {
	t.Skip("skipping test due to long runtime/downloading file")
	log.Root().SetHandler(log.StdoutHandler)

	// TODO: replace with actual sha256 hash and downloadURL
	sha256 := [32]byte(common.FromHex("0x250c097758924bc21d072e8dc57f4a2357ffaafb20e85eacea5c18dfe70e62b4"))
	downloadURL := "https://ftp.halifax.rwth-aachen.de/ubuntu-releases/robots.txt"
	filePath := filepath.Join(t.TempDir(), "test_file_path")
	manager := NewManager(context.Background(), filePath, downloadURL, sha256)

	_, _, err := manager.GetMissingHeaderFields(0)
	require.NoError(t, err)

	// Check if the file was downloaded and tmp file was removed
	_, err = os.Stat(filePath)
	require.NoError(t, err)
	_, err = os.Stat(filePath + ".tmp")
	require.Error(t, err)
}

func TestManagerChecksum(t *testing.T) {
	downloadURL := "" // since the file exists we don't need to download it
	filePath := filepath.Join("testdata", "missing-headers.bin")

	// Checksum doesn't match
	{
		sha256 := [32]byte(common.FromHex("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))

		manager := NewManager(context.Background(), filePath, downloadURL, sha256)

		_, _, err := manager.GetMissingHeaderFields(0)
		require.ErrorContains(t, err, "expectedChecksum mismatch")
	}

	// Checksum matches
	{
		sha256 := [32]byte(common.FromHex("0xfa5d9de3dfdae76a9abd03f7c28274ab223d74a90ed1735e74be1f8fc9c2a435"))
		manager := NewManager(context.Background(), filePath, downloadURL, sha256)

		difficulty, extra, err := manager.GetMissingHeaderFields(0)
		require.NoError(t, err)
		require.Equal(t, expectedMissingHeaders1[0].difficulty, difficulty)
		require.Equal(t, expectedMissingHeaders1[0].extra, extra)
	}
}

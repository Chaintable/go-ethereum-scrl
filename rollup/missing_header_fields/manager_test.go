package missing_header_fields

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/scroll-tech/go-ethereum/common"
	"github.com/scroll-tech/go-ethereum/log"
	"github.com/scroll-tech/go-ethereum/params"
)

func TestManagerDownload(t *testing.T) {
	t.Skip("skipping test due to long runtime/downloading file")
	log.Root().SetHandler(log.StdoutHandler)

	sha256 := *params.ScrollSepoliaChainConfig.Scroll.MissingHeaderFieldsSHA256
	downloadURL := "https://scroll-block-missing-metadata.s3.us-west-2.amazonaws.com/" + params.ScrollSepoliaChainConfig.ChainID.String() + ".bin"
	filePath := filepath.Join(t.TempDir(), "test_file_path")
	manager := NewManager(context.Background(), filePath, downloadURL, sha256)

	_, _, _, err := manager.GetMissingHeaderFields(0)
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

		_, _, _, err := manager.GetMissingHeaderFields(0)
		require.ErrorContains(t, err, "expectedChecksum mismatch")
	}

	// Checksum matches
	{
		sha256 := [32]byte(common.FromHex("e5a1e71338cd899e46ff28a9ae81b8f2579e429e18cec463104fb246a6e23502"))
		manager := NewManager(context.Background(), filePath, downloadURL, sha256)

		difficulty, stateRoot, extra, err := manager.GetMissingHeaderFields(0)
		require.NoError(t, err)
		require.Equal(t, expectedMissingHeaders[0].difficulty, difficulty)
		require.Equal(t, expectedMissingHeaders[0].stateRoot, stateRoot)
		require.Equal(t, expectedMissingHeaders[0].extra, extra)
	}
}

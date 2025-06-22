package blob_client

import (
	"context"
	"testing"

	"github.com/scroll-tech/go-ethereum/common"
)

func TestGetBlobByVersionedHashAndBlockTime(t *testing.T) {
	apiEndpoint := "https://scroll-sepolia-blob-data.s3.us-west-2.amazonaws.com/"
	awsS3Client := NewAwsS3Client(apiEndpoint)

	ctx := context.Background()
	versionedHash := common.HexToHash("0x01e7f0962458d4a4ff61bad08437c3972d7cb443d7ccfdd41b32904e8a5fe24b") 
	_, err := awsS3Client.GetBlobByVersionedHashAndBlockTime(ctx, versionedHash, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
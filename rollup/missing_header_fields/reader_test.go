package missing_header_fields

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/scroll-tech/go-ethereum/common"
)

type header struct {
	number     uint64
	difficulty uint64
	extra      []byte
}

var expectedMissingHeaders1 = []header{
	{0, 1, common.FromHex("000000000000000000000000000000000000000000000000000000000000000048c3f81f3d998b6652900e1c3183736c238fe4290000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")},
	{1, 2, common.FromHex("d88304031d846765746888676f312e31392e31856c696e7578000000000000001982b5c754257988f9486b158a33709645735e8e965912c508aee9b0513cc2f22fe13f0835ce1e11abe666c9dba6a1259612b812783cc457e5b34b025980635501")},
	{2, 2, common.FromHex("d88304031d846765746888676f312e31392e31856c696e757800000000000000237c933578bf062f86a30cdc71b0e946f0f685711e0e9cceeb1c953ed816d2694347e1e59625545c4040f2604b75448ccb5360fdcb378741331c1d4c0d342a7e01")},
	{3, 2, common.FromHex("d88304031d846765746888676f312e31392e31856c696e75780000000000000012388a2df0f522f96e67564d38be64b5ca7fb37ef9b3f88de875d08653871407584b180917a47dc4abec60bf8da462c617328b9d2da8c4bb9978e018b44ec07401")},
	{4, 2, common.FromHex("d88304031d846765746888676f312e31392e31856c696e75780000000000000091e57e01b8ed1b433b2bd04e272f9eaf986f3fa728c8fc2b4112352101d24ba76ff1ee64e9a1f8a47c4c49e362e318b2b4767088514f72a7ba9bb7a45b4b447700")},
	{5, 2, common.FromHex("d88304031d846765746888676f312e31392e31856c696e7578000000000000007ab6b6bd8d52c9beffe935e1bc805d9d4ad62d54485104e80943537d380d6f425ad055ad510c498d1e6efc2aa7e7cc7e1b6166f8421e94b13e291196cba1934a00")},
	{6, 2, common.FromHex("d88304031d846765746888676f312e31392e31856c696e7578000000000000008f7175ed80593d395069afed2d970e505b076e15642e1f6e3bcb2f589a2a47fa5841868b80fcc70f18cd52de027bdb5ab881fdefc5ebf0d8034cb35e926e89f300")},
	{7, 2, common.FromHex("d88304031d846765746888676f312e31392e31856c696e75780000000000000006224d2201ba60083743844ce0c2ec4b0ab3e79b69f64eae9a10055fe704380c6410c4f5119cb834f43705c1a785758170a868a38e536432e3a5a5805c83b13801")},
	{8, 2, common.FromHex("d88304031d846765746888676f312e31392e31856c696e757800000000000000b5c1f5c8aa79582f4b9a66ae8a59561d6c357deaefc7353ea6ace017e76e4b367118e0fd55b9f4cd0235ee1f14222e9b558156b6253e84f71d8048e13643af4801")},
	{9, 2, common.FromHex("d88304031d846765746888676f312e31392e31856c696e757800000000000000dff93471464bf856b2f633ac16b54c3ff88219a8c83067495ff9f16035c91fb56de8f6914eb4cbd8fbe8a54854e32d697a81408e20cdaa52fed9689d21f7ad0201")},
	{10, 2, common.FromHex("d88304031d846765746888676f312e31392e31856c696e7578000000000000003c55c63554686f48d9e6dd78b8a7849152e33f169f843e623aa604f2c777eec9539d301fe5f1bea84e5cb7c40e74b723e7700b95eab08bc441d65092b40b548d01")},
	{11, 2, common.FromHex("d88304031d846765746888676f312e31392e31856c696e75780000000000000013908919960db7ad7459bbeb172049e0dc71fb3a56c7e89bd3699b9ab2cd64d77acc2ba79bb4b9dc98baba0a8e32fd1ca791065b7bd535225f55e254359338d401")},
	{12, 2, common.FromHex("d88304031d846765746888676f312e31392e31856c696e7578000000000000009b6d0158e0da8bb62c00ff406321b323caf26ffe8d0ea7f181080ea08e236afd63116dea58c752206ff4939240b80aaf1467bb7bb7ad2d8f7d8583f95eec725e01")},
}

func TestReader_Read(t *testing.T) {
	expectedVanities := map[int][32]byte{
		0: [32]byte(common.FromHex("0000000000000000000000000000000000000000000000000000000000000000")),
		1: [32]byte(common.FromHex("0xd88304031d846765746888676f312e31392e31856c696e757800000000000000")),
	}

	reader, err := NewReader("testdata/missing-headers.bin")
	require.NoError(t, err)

	require.Len(t, reader.sortedVanities, len(expectedVanities))
	for i, expectedVanity := range expectedVanities {
		require.Equal(t, expectedVanity, reader.sortedVanities[i])
	}

	readAndAssertHeader(t, reader, expectedMissingHeaders1, 0)
	readAndAssertHeader(t, reader, expectedMissingHeaders1, 0)
	readAndAssertHeader(t, reader, expectedMissingHeaders1, 1)
	readAndAssertHeader(t, reader, expectedMissingHeaders1, 6)

	// we don't allow reading previous headers
	_, _, err = reader.Read(5)
	require.Error(t, err)

	readAndAssertHeader(t, reader, expectedMissingHeaders1, 8)
	readAndAssertHeader(t, reader, expectedMissingHeaders1, 8)

	// we don't allow reading previous headers
	_, _, err = reader.Read(5)
	require.Error(t, err)

	// we don't allow reading previous headers
	_, _, err = reader.Read(6)
	require.Error(t, err)

	readAndAssertHeader(t, reader, expectedMissingHeaders1, 9)
	readAndAssertHeader(t, reader, expectedMissingHeaders1, 10)
	readAndAssertHeader(t, reader, expectedMissingHeaders1, 11)
	readAndAssertHeader(t, reader, expectedMissingHeaders1, 12)

	// no data anymore
	_, _, err = reader.Read(13)
	require.Error(t, err)
}

func readAndAssertHeader(t *testing.T, reader *Reader, expectedHeaders []header, headerNum uint64) {
	difficulty, extra, err := reader.Read(headerNum)
	require.NoError(t, err)
	require.Equalf(t, expectedHeaders[headerNum].difficulty, difficulty, "expected difficulty %d, got %d", expectedHeaders[headerNum].difficulty, difficulty)
	require.Equal(t, expectedHeaders[headerNum].extra, extra)
}

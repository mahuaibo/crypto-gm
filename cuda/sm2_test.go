//+build cuda

package cuda

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
to test: LD_LIBRARY_PATH=/usr/local/lib go test  -tags cuda
*/

func TestVerifySignatureGPU(t *testing.T) {
	info := [10]struct {
		es, ss, sks string
	}{
		{"e197e49321287c3718f5f2510ade487350fb43a742729d1c40a32da76f239222", "3044022079a9e01462fb3776a83a6850aba15dcf5f4456e31425bf5aeeb4677b9e770c1702201fbb4c38d38fb20f9235adb340a3f27a072d1917538ebdeef06ed57bd6167365", "80fb020fbe465a3800bd023c69480d6e29f1a86394683fe0bc8c10a1f83fb475"},
		{"355db09e7194d20f958ad5ebd682ce3ca9e745e0ffc536ad8548a9ea29ccc100", "3046022100fe4b3c11721c031d45528bbe1ae614c60161c14250d98e1cb923872c1b1d8cc302210098b75798cc71716431124eec265a88f7ed95a1c973b21741c629300a45166345", "f4e6cff146d0b591bf53889c057a8ae97d6f89dae6a1e562e337ab6249bc402e"},
		{"57f82104e6260ed2548ce195cdec89cbddc47673b94e954d0429f0d5f91e301a", "3046022100ce56ce7274d905a9c954d38868160bfced2aa7ad11cac57cdf4e55b133ce8898022100fffbf46a1dc4f09e59794378ee8c592a361cc39b2294bfe18afa9f36954892a8", "cc1530a92c406b0e9cd11a845716e5f1efb4612e5dd2e995f6d901bfabdea8ee"},
		{"85fd4e6127156c738093b5a9fcc1b8afc32ea60adb25f5c438161b094b86e6cd", "3046022100dfd5257556a4b2d1e2e058cf06edb235f96a0d4074f1714a4b32ac0d2e15b84b022100e950e4542d2cb3b658be072af881271713159dfe49e9611bffd7a54afde18ef6", "3ac54f142bb3579658d67a3829b2a633dec52dd6299b408a3de76961d7a7ba47"},
		{"22db90cb6c8c15b8805e3d463b31bd32a8c2e4f2d8f028563a939afcd3d3bffe", "30440220361b19d2a96a212e26ae5a102d5b9505bbeeb6707dc81b843a36b16a988bae8b022002dc2e6294057b1669f61ec2c1fe1b940742e267d5816476e33da041286dfdce", "fcbfb0a4a07997dda8810b835a4ae14c6959e506dafcd1c53ec4719bf1ef4d7d"},
		{"5b4748c97e0115154be4cf89c7e9743c16527119547d5eae19acf204686a1839", "3046022100bc81908f9ad16fa58be78b36c70965671a10b4090bb9795d5a21f79a460c4c1c022100fe1faef889bff2882d37d6ee612fe3fac6099fd2aff0bc93aa593e91eda84aa2", "aa159de3ffb09fbfa70fd3f135ceb46409e0e866858c1ca1297ad6cf69d4f410"},
		{"1e070f08b972c5c5b13cdab1f42f74e346bdd9cf4693d85058e1134da346f68e", "304502210096cf73d2781cd004abdd6b00fc8b8b41df4ac6655f055b94ac0d9ec45ee4e14f0220316f3bf230b1779be123ebcbd6428e887c3c4b7d1e9c8907efd3e31f5fe5887c", "ece2749e0b091f17f8bddbd6707a4401770c0c8fe887f7c2f921beea25b6999d"},
		{"cc7d1e4ba3c0c35855638dcb55199edec4f565a90e965abaf907211475a9a8d1", "3045022020b2412563a6e86520b06c017e71825e7477a7db37ea5d0e845c80348a235e0a022100ba23b70ece87c12b321eb38ffcd68353ccc5b1447403b51ce6d537881bdce86c", "997979f904db67fc6f5c8572a7776eb16ef70a7994812dfee49062ff7f27b954"},
		{"b27ff60a3286ccf2accf4f3fa183bf55212c5121d77d119de65f0fcad1879ae8", "304502200dec31254ac5bc158b8cedbe4348abe8fa7348cee61c235cd6a6f77bd15a274b022100feb9323e0c6513b88c6dda9e01da9f508a17f3103e58aedc7340075dbb7314fc", "0dc49da4205a21b61bf18e5d757f97092a719eb58f6ff4a663b64243b0965201"},
		{"b0bcfd241ac426c13d535ca7dba3763d18494cc21abf474ae4f7493a510bba4a", "3046022100e5c0293e9149bdd905efa26111b2dfd770611235b64d986cabc236c400bee4fc022100d9828bd20c90cc384b66524c82872b438c5f54df84e84eb6108dc0437754bb6a", "acc4af8ffed59594371db6e229b9c3fc1d0392243dc5901cfbf9af2105fcb061"},
	}
	es, ss, sks := make([][]byte, 10), make([][]byte, 10), make([][]byte, 10)
	for i := range info {
		es[i], _ = hex.DecodeString(info[i].es)
		ss[i], _ = hex.DecodeString(info[i].ss)
		sks[i], _ = hex.DecodeString(info[i].sks)
	}
	es[2][1] = 0
	ss[6][22] = 10 //5?
	r, _ := VerifySignatureGPUM(ss, es, sks)
	for i, result := range r {
		if i == 2 || i == 3 || i == 6 || i == 7 { //第i笔的签名错误会影响到临近的一笔验签
			assert.False(t, result == 0)
			continue
		}
		assert.True(t, result == 0)
	}
}

func TestVerifySignatureGPUSingle(t *testing.T) {
	s, _ := hex.DecodeString("30450221009b357468b832499ec086bddecbfeac1c48c84014721027d635cb5b0c4a876f1d022038677b9cb1d9dbe684d8659ead78a358c0fc5c2e419cbac54cff6adb8ce1ea43")
	e, _ := hex.DecodeString("1943e5cce3a05cbb544abdc293b22ae34026cdd71cc0e1b098be10844e29e438")
	pk, _ := hex.DecodeString("68cdc06c1c12ac407d3a4b4556230c974ed01a399b8e26a22ee5e7f412178ee3")
	r, err := VerifySignatureGPUM([][]byte{s}, [][]byte{e}, [][]byte{pk})
	assert.Nil(t, err[0])
	assert.True(t, r[0] == 0)

	t.Run("签名错误的情况", func(t *testing.T) {
		ss := make([]byte, len(s))
		copy(ss, s)
		ss[11] = 0
		r, err := VerifySignatureGPUM([][]byte{ss}, [][]byte{e}, [][]byte{pk})
		assert.Nil(t, err[0])
		assert.True(t, r[0] != 0)
	})

	t.Run("消息错误的情况", func(t *testing.T) {
		ee := make([]byte, len(e))
		copy(ee, s)
		ee[10] = 0
		r, err := VerifySignatureGPUM([][]byte{s}, [][]byte{ee}, [][]byte{pk})
		assert.Nil(t, err[0])
		assert.True(t, r[0] != 0)
	})

	t.Run("公钥错误的情况", func(t *testing.T) {
		pks := make([]byte, len(pk))
		copy(pks, pk)
		pks[10] = 0
		r, err := VerifySignatureGPUM([][]byte{s}, [][]byte{e}, [][]byte{pks})
		assert.Nil(t, err[0])
		assert.True(t, r[0] != 0)
	})

}

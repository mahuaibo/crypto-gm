package gm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateKGC(t *testing.T) {
	kgc := GenerateKGC()
	assert.NotNil(t, kgc)
	key := kgc.GenerateKey([]byte("hyperchain"))
	assert.NotNil(t, key)
	assert.Equal(t, key.KGCPubKey, kgc.KGCPubKey.Pub)
}

func BenchmarkKGC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		kgc := GenerateKGC()
		assert.NotNil(b, kgc)
	}
} //BenchmarkKGC-4   	     550	   2217265 ns/op

func BenchmarkKGC2(b *testing.B) {
	kgc := GenerateKGC()
	assert.NotNil(b, kgc)
	for i := 0; i < b.N; i++ {
		key := kgc.GenerateKey([]byte("hyperchain"))
		assert.NotNil(b, key)
		assert.Equal(b, key.KGCPubKey, kgc.KGCPubKey.Pub)
	}
}

func TestSM9Key_Bytes(t *testing.T) {
	kgc := GenerateKGC()
	assert.NotNil(t, kgc)
	key := kgc.GenerateKey([]byte("hyperchain"))
	bs, err := key.Bytes()
	assert.Nil(t, err)
	newKey := new(SM9Key).FromBytes(bs, nil)
	assert.Equal(t, key.K, newKey.K)
	assert.Equal(t, key.ID, newKey.ID)
	assert.Equal(t, key.KGCPubKey, newKey.KGCPubKey)
}

func TestSM9_Sign(t *testing.T) {
	kgc := GenerateKGC()
	assert.NotNil(t, kgc)
	key := kgc.GenerateKey([]byte("hyperchain"))
	sm9 := new(SM9)
	pri, err := key.Bytes()
	assert.Nil(t, err)
	s, err := sm9.Sign(pri, []byte(msg))
	assert.Nil(t, err)
	ID, err := key.PublicKey()
	assert.Nil(t, err)
	pub, err := ID.Bytes()
	assert.Nil(t, err)
	ok, err := sm9.Verify(pub, s, []byte(msg))
	assert.Nil(t, err)
	assert.True(t, ok)
}

func BenchmarkSM9_Sign(b *testing.B) {
	kgc := GenerateKGC()
	assert.NotNil(b, kgc)
	key := kgc.GenerateKey([]byte("hyperchain"))
	sm9 := new(SM9)
	pri, err := key.Bytes()
	assert.Nil(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sm9.Sign(pri, []byte(msg))
		assert.Nil(b, err)
	}
} //BenchmarkSM9_Sign-4   	      50	  26396904 ns/op

func BenchmarkSM9_Verify(t *testing.B) {
	kgc := GenerateKGC()
	assert.NotNil(t, kgc)
	key := kgc.GenerateKey([]byte("hyperchain"))
	sm9 := new(SM9)
	pri, err := key.Bytes()
	assert.Nil(t, err)
	s, err := sm9.Sign(pri, []byte(msg))
	assert.Nil(t, err)
	ID, err := key.PublicKey()
	assert.Nil(t, err)
	pub, err := ID.Bytes()
	assert.Nil(t, err)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		_, err := sm9.Verify(pub, s, []byte(msg))
		assert.Nil(t, err)
	}
} //BenchmarkSM9_Verify-4   	      30	  45641334 ns/op

func TestSM9Key_Sign(t *testing.T) {
	kgc := GenerateKGC()
	assert.NotNil(t, kgc)
	key := kgc.GenerateKey([]byte("hyperchain"))
	s, err := key.Sign(nil, []byte(msg))
	assert.Nil(t, err)
	ID, err := key.PublicKey()
	assert.Nil(t, err)
	b, err := ID.Verify(nil, s, []byte(msg))
	assert.Nil(t, err)
	assert.True(t, b)
}

func TestID(t *testing.T) {
	kgc := GenerateKGC()
	assert.NotNil(t, kgc)
	key := kgc.GenerateKey([]byte("hyperchain"))
	oldID, err := key.PublicKey()
	assert.Nil(t, err)
	b, err := oldID.Bytes()
	assert.Nil(t, err)
	id := new(ID)
	newID := id.FromBytes(b, nil)
	assert.Equal(t, oldID, newID)
	assert.False(t, oldID.Symmetric())
	assert.False(t, oldID.Private())
	pubKey, err := oldID.PublicKey()
	assert.Nil(t, err)
	assert.Equal(t, pubKey, oldID)
}

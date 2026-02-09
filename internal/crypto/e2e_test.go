package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"

	"golang.org/x/crypto/nacl/box"
)

// testKeyPair derives a deterministic keypair for tests using PIN "123456"
// and a 16-byte zero salt. This avoids repeated boilerplate in edge-case tests.
func testKeyPair(t *testing.T) *KeyPair {
	t.Helper()
	kp, err := DeriveKeyPair("123456", make([]byte, 16))
	if err != nil {
		t.Fatalf("testKeyPair: %v", err)
	}
	return kp
}

// testEncrypt encrypts plaintext with SealAnonymous for the given keypair.
func testEncrypt(t *testing.T, plaintext []byte, kp *KeyPair) []byte {
	t.Helper()
	ciphertext, err := box.SealAnonymous(nil, plaintext, &kp.PublicKey, rand.Reader)
	if err != nil {
		t.Fatalf("testEncrypt: %v", err)
	}
	return ciphertext
}

func TestIsEncrypted(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"encrypted value", "e2e::abc123==", true},
		{"plain text", "Hello world", false},
		{"empty string", "", false},
		{"prefix only", "e2e::", true},
		{"partial prefix", "e2e:", false},
		{"embedded prefix", "some e2e::data", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEncrypted(tt.value); got != tt.want {
				t.Errorf("IsEncrypted(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestDeriveKeyPair_Deterministic(t *testing.T) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		t.Fatalf("generating salt: %v", err)
	}

	pin := "123456"

	kp1, err := DeriveKeyPair(pin, salt)
	if err != nil {
		t.Fatalf("DeriveKeyPair (1st call): %v", err)
	}

	kp2, err := DeriveKeyPair(pin, salt)
	if err != nil {
		t.Fatalf("DeriveKeyPair (2nd call): %v", err)
	}

	if kp1.PublicKey != kp2.PublicKey {
		t.Error("public keys differ for same PIN + salt")
	}
	if kp1.PrivateKey != kp2.PrivateKey {
		t.Error("private keys differ for same PIN + salt")
	}
}

func TestDeriveKeyPair_DifferentPINs(t *testing.T) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		t.Fatalf("generating salt: %v", err)
	}

	kp1, err := DeriveKeyPair("123456", salt)
	if err != nil {
		t.Fatalf("DeriveKeyPair (PIN 123456): %v", err)
	}

	kp2, err := DeriveKeyPair("654321", salt)
	if err != nil {
		t.Fatalf("DeriveKeyPair (PIN 654321): %v", err)
	}

	if kp1.PublicKey == kp2.PublicKey {
		t.Error("different PINs produced the same public key")
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		t.Fatalf("generating salt: %v", err)
	}

	kp, err := DeriveKeyPair("999999", salt)
	if err != nil {
		t.Fatalf("DeriveKeyPair: %v", err)
	}

	plaintext := "Your verification code is 847291"

	// Encrypt with SealAnonymous (simulates what the server does)
	ciphertext, err := box.SealAnonymous(nil, []byte(plaintext), &kp.PublicKey, rand.Reader)
	if err != nil {
		t.Fatalf("SealAnonymous: %v", err)
	}

	// Decrypt
	got, err := Decrypt(ciphertext, kp)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}

	if string(got) != plaintext {
		t.Errorf("Decrypt = %q, want %q", string(got), plaintext)
	}
}

func TestDecryptField(t *testing.T) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		t.Fatalf("generating salt: %v", err)
	}

	kp, err := DeriveKeyPair("111111", salt)
	if err != nil {
		t.Fatalf("DeriveKeyPair: %v", err)
	}

	plaintext := "Code: 123456"

	// Build the "e2e::<base64>" value
	ciphertext, err := box.SealAnonymous(nil, []byte(plaintext), &kp.PublicKey, rand.Reader)
	if err != nil {
		t.Fatalf("SealAnonymous: %v", err)
	}
	field := EncryptedPrefix + base64.StdEncoding.EncodeToString(ciphertext)

	// DecryptField should return the original plaintext
	got, err := DecryptField(field, kp)
	if err != nil {
		t.Fatalf("DecryptField: %v", err)
	}
	if got != plaintext {
		t.Errorf("DecryptField = %q, want %q", got, plaintext)
	}

	// Non-encrypted fields should pass through unchanged
	plain := "not encrypted"
	got, err = DecryptField(plain, kp)
	if err != nil {
		t.Fatalf("DecryptField (plain): %v", err)
	}
	if got != plain {
		t.Errorf("DecryptField (plain) = %q, want %q", got, plain)
	}
}

func TestDecryptField_EmptyString(t *testing.T) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		t.Fatalf("generating salt: %v", err)
	}

	kp, err := DeriveKeyPair("000000", salt)
	if err != nil {
		t.Fatalf("DeriveKeyPair: %v", err)
	}

	got, err := DecryptField("", kp)
	if err != nil {
		t.Fatalf("DecryptField: %v", err)
	}
	if got != "" {
		t.Errorf("DecryptField(\"\") = %q, want empty string", got)
	}
}

func TestVerify(t *testing.T) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		t.Fatalf("generating salt: %v", err)
	}

	kp, err := DeriveKeyPair("222222", salt)
	if err != nil {
		t.Fatalf("DeriveKeyPair: %v", err)
	}

	// Create a verifier using CreateVerifier
	verifierB64, err := CreateVerifier(kp)
	if err != nil {
		t.Fatalf("CreateVerifier: %v", err)
	}

	// Verify should succeed with the correct keypair
	if !Verify(kp, verifierB64) {
		t.Error("Verify returned false for correct keypair")
	}

	// Verify should fail with a different keypair
	otherSalt := make([]byte, 16)
	if _, err := rand.Read(otherSalt); err != nil {
		t.Fatalf("generating other salt: %v", err)
	}
	otherKP, err := DeriveKeyPair("333333", otherSalt)
	if err != nil {
		t.Fatalf("DeriveKeyPair (other): %v", err)
	}
	if Verify(otherKP, verifierB64) {
		t.Error("Verify returned true for wrong keypair")
	}
}

func TestCreateVerifier(t *testing.T) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		t.Fatalf("generating salt: %v", err)
	}

	kp, err := DeriveKeyPair("444444", salt)
	if err != nil {
		t.Fatalf("DeriveKeyPair: %v", err)
	}

	v1, err := CreateVerifier(kp)
	if err != nil {
		t.Fatalf("CreateVerifier (1st): %v", err)
	}
	v2, err := CreateVerifier(kp)
	if err != nil {
		t.Fatalf("CreateVerifier (2nd): %v", err)
	}

	// SealAnonymous uses ephemeral keys, so two verifiers should differ
	if v1 == v2 {
		t.Error("two CreateVerifier calls produced identical ciphertexts (expected ephemeral randomness)")
	}

	// Both should still verify correctly
	if !Verify(kp, v1) {
		t.Error("Verify failed for first verifier")
	}
	if !Verify(kp, v2) {
		t.Error("Verify failed for second verifier")
	}
}

func TestVerify_InvalidBase64(t *testing.T) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		t.Fatalf("generating salt: %v", err)
	}

	kp, err := DeriveKeyPair("555555", salt)
	if err != nil {
		t.Fatalf("DeriveKeyPair: %v", err)
	}

	if Verify(kp, "not-valid-base64!!!") {
		t.Error("Verify returned true for invalid base64")
	}
}

func TestDecryptField_InvalidBase64(t *testing.T) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		t.Fatalf("generating salt: %v", err)
	}

	kp, err := DeriveKeyPair("666666", salt)
	if err != nil {
		t.Fatalf("DeriveKeyPair: %v", err)
	}

	_, err = DecryptField("e2e::not-valid-base64!!!", kp)
	if err == nil {
		t.Error("DecryptField should return error for invalid base64")
	}
}

// ---------------------------------------------------------------------------
// Decrypt error paths
// ---------------------------------------------------------------------------

func TestDecrypt_EmptyCiphertext(t *testing.T) {
	kp := testKeyPair(t)
	_, err := Decrypt([]byte{}, kp)
	if err == nil {
		t.Error("Decrypt should return error for empty ciphertext")
	}
}

func TestDecrypt_TruncatedCiphertext(t *testing.T) {
	kp := testKeyPair(t)

	// SealedBox overhead is 48 bytes (32 ephemeral pubkey + 16 MAC).
	// Any ciphertext shorter than that is invalid.
	truncated := make([]byte, 47)
	if _, err := rand.Read(truncated); err != nil {
		t.Fatalf("generating random bytes: %v", err)
	}

	_, err := Decrypt(truncated, kp)
	if err == nil {
		t.Error("Decrypt should return error for ciphertext shorter than SealedBox overhead")
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	kp := testKeyPair(t)
	ciphertext := testEncrypt(t, []byte("secret message"), kp)

	// Flip the last byte to corrupt the MAC
	ciphertext[len(ciphertext)-1] ^= 0xff

	_, err := Decrypt(ciphertext, kp)
	if err == nil {
		t.Error("Decrypt should return error for tampered ciphertext")
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	kp := testKeyPair(t)
	ciphertext := testEncrypt(t, []byte("for keypair A only"), kp)

	// Derive a different keypair
	otherKP, err := DeriveKeyPair("654321", make([]byte, 16))
	if err != nil {
		t.Fatalf("DeriveKeyPair (other): %v", err)
	}

	_, err = Decrypt(ciphertext, otherKP)
	if err == nil {
		t.Error("Decrypt should return error when using the wrong keypair")
	}
}

// ---------------------------------------------------------------------------
// DecryptField edge cases
// ---------------------------------------------------------------------------

func TestDecryptField_PrefixOnlyNoData(t *testing.T) {
	kp := testKeyPair(t)

	// "e2e::" with no base64 payload — base64 decodes "" to empty []byte,
	// which is too short for SealedBox and should fail during decryption.
	_, err := DecryptField("e2e::", kp)
	if err == nil {
		t.Error("DecryptField should return error for prefix with no ciphertext data")
	}
}

func TestDecryptField_NonE2EPrefix(t *testing.T) {
	kp := testKeyPair(t)
	value := "other::abc123"

	got, err := DecryptField(value, kp)
	if err != nil {
		t.Fatalf("DecryptField: %v", err)
	}
	if got != value {
		t.Errorf("DecryptField(%q) = %q, want value returned unchanged", value, got)
	}
}

func TestDecryptField_PlainText(t *testing.T) {
	kp := testKeyPair(t)
	value := "Hello world"

	got, err := DecryptField(value, kp)
	if err != nil {
		t.Fatalf("DecryptField: %v", err)
	}
	if got != value {
		t.Errorf("DecryptField(%q) = %q, want value returned unchanged", value, got)
	}
}

func TestDecryptField_UnicodeRoundTrip(t *testing.T) {
	kp := testKeyPair(t)
	plaintext := "Votre code: \U0001f511 123456"

	ciphertext := testEncrypt(t, []byte(plaintext), kp)
	field := EncryptedPrefix + base64.StdEncoding.EncodeToString(ciphertext)

	got, err := DecryptField(field, kp)
	if err != nil {
		t.Fatalf("DecryptField: %v", err)
	}
	if got != plaintext {
		t.Errorf("DecryptField = %q, want %q", got, plaintext)
	}
}

// ---------------------------------------------------------------------------
// DeriveKeyPair edge cases
// ---------------------------------------------------------------------------

func TestDeriveKeyPair_EmptySalt(t *testing.T) {
	// Argon2id accepts an empty salt without error.
	kp, err := DeriveKeyPair("123456", []byte{})
	if err != nil {
		t.Fatalf("DeriveKeyPair with empty salt: %v", err)
	}
	if kp == nil {
		t.Fatal("DeriveKeyPair returned nil keypair for empty salt")
	}
}

func TestDeriveKeyPair_EmptyPIN(t *testing.T) {
	// Argon2id accepts an empty password without error.
	kp, err := DeriveKeyPair("", make([]byte, 16))
	if err != nil {
		t.Fatalf("DeriveKeyPair with empty PIN: %v", err)
	}
	if kp == nil {
		t.Fatal("DeriveKeyPair returned nil keypair for empty PIN")
	}
}

func TestDeriveKeyPair_LongPIN(t *testing.T) {
	longPIN := strings.Repeat("A", 1000)
	kp, err := DeriveKeyPair(longPIN, make([]byte, 16))
	if err != nil {
		t.Fatalf("DeriveKeyPair with 1000-char PIN: %v", err)
	}
	if kp == nil {
		t.Fatal("DeriveKeyPair returned nil keypair for long PIN")
	}
}

func TestDeriveKeyPair_DifferentSalts(t *testing.T) {
	pin := "123456"
	salt1 := make([]byte, 16) // all zeros
	salt2 := make([]byte, 16)
	salt2[0] = 1 // differs in first byte

	kp1, err := DeriveKeyPair(pin, salt1)
	if err != nil {
		t.Fatalf("DeriveKeyPair (salt1): %v", err)
	}

	kp2, err := DeriveKeyPair(pin, salt2)
	if err != nil {
		t.Fatalf("DeriveKeyPair (salt2): %v", err)
	}

	if kp1.PublicKey == kp2.PublicKey {
		t.Error("same PIN with different salts produced the same public key")
	}
	if kp1.PrivateKey == kp2.PrivateKey {
		t.Error("same PIN with different salts produced the same private key")
	}
}

// ---------------------------------------------------------------------------
// Verify edge cases
// ---------------------------------------------------------------------------

func TestVerify_EmptyVerifier(t *testing.T) {
	kp := testKeyPair(t)
	if Verify(kp, "") {
		t.Error("Verify should return false for empty verifier string")
	}
}

func TestVerify_TruncatedVerifier(t *testing.T) {
	kp := testKeyPair(t)

	// 32 bytes of random data is valid base64 but too short for a SealedBox
	// (needs at least 48 bytes overhead + plaintext).
	short := make([]byte, 32)
	if _, err := rand.Read(short); err != nil {
		t.Fatalf("generating random bytes: %v", err)
	}
	verifier := base64.StdEncoding.EncodeToString(short)

	if Verify(kp, verifier) {
		t.Error("Verify should return false for truncated verifier")
	}
}

func TestVerify_CorruptedVerifier(t *testing.T) {
	kp := testKeyPair(t)

	verifierB64, err := CreateVerifier(kp)
	if err != nil {
		t.Fatalf("CreateVerifier: %v", err)
	}

	// Decode, corrupt a byte in the ciphertext body, re-encode
	raw, err := base64.StdEncoding.DecodeString(verifierB64)
	if err != nil {
		t.Fatalf("decoding verifier: %v", err)
	}
	raw[len(raw)-1] ^= 0xff
	corrupted := base64.StdEncoding.EncodeToString(raw)

	if Verify(kp, corrupted) {
		t.Error("Verify should return false for corrupted verifier")
	}
}

// ---------------------------------------------------------------------------
// CreateVerifier
// ---------------------------------------------------------------------------

func TestCreateVerifier_Roundtrip(t *testing.T) {
	kp := testKeyPair(t)

	verifier, err := CreateVerifier(kp)
	if err != nil {
		t.Fatalf("CreateVerifier: %v", err)
	}

	if !Verify(kp, verifier) {
		t.Error("Verify returned false for verifier created with same keypair")
	}
}

func TestCreateVerifier_WrongKey(t *testing.T) {
	kpA := testKeyPair(t)

	verifier, err := CreateVerifier(kpA)
	if err != nil {
		t.Fatalf("CreateVerifier: %v", err)
	}

	// Derive a different keypair
	kpB, err := DeriveKeyPair("654321", make([]byte, 16))
	if err != nil {
		t.Fatalf("DeriveKeyPair (B): %v", err)
	}

	if Verify(kpB, verifier) {
		t.Error("Verify should return false when verifying with a different keypair")
	}
}

// ---------------------------------------------------------------------------
// ClearCachedKeyPair
// ---------------------------------------------------------------------------

func TestClearCachedKeyPair(t *testing.T) {
	// Set the package-level cache to a known value
	cachedKeyPair = testKeyPair(t)
	if cachedKeyPair == nil {
		t.Fatal("cachedKeyPair should be non-nil after assignment")
	}

	ClearCachedKeyPair()

	if cachedKeyPair != nil {
		t.Error("cachedKeyPair should be nil after ClearCachedKeyPair")
	}
}

// ---------------------------------------------------------------------------
// Cross-compatibility / regression test
// ---------------------------------------------------------------------------

func TestDeriveKeyPair_KnownTestVector(t *testing.T) {
	// Deterministic inputs: PIN "123456", 16 zero-byte salt.
	// The expected public key was generated once and hardcoded here to detect
	// any accidental change to the Argon2id parameters or key derivation logic.
	pin := "123456"
	salt := make([]byte, 16)
	expectedPubKeyHex := "584d44d2c4b7ebf636a528b1c18c42499d9590b4a9480759888c7d99a64c005d"

	kp, err := DeriveKeyPair(pin, salt)
	if err != nil {
		t.Fatalf("DeriveKeyPair: %v", err)
	}

	gotHex := hex.EncodeToString(kp.PublicKey[:])
	if gotHex != expectedPubKeyHex {
		t.Errorf("public key mismatch\n  got:  %s\n  want: %s", gotHex, expectedPubKeyHex)
	}
}

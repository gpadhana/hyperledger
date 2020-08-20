// +build pkcs11

/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pkcs11

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/sw"
	"github.com/hyperledger/fabric/bccsp/utils"
	"github.com/miekg/pkcs11"
	"github.com/stretchr/testify/require"
)

func defaultOptions() PKCS11Opts {
	lib, pin, label := FindPKCS11Lib()
	return PKCS11Opts{
		Library:                 lib,
		Label:                   label,
		Pin:                     pin,
		Hash:                    "SHA2",
		Security:                256,
		SoftwareVerify:          false,
		createSessionRetryDelay: time.Millisecond,
	}
}

func newKeyStore(t *testing.T) (bccsp.KeyStore, func()) {
	tempDir, err := ioutil.TempDir("", "pkcs11_ks")
	require.NoError(t, err)
	ks, err := sw.NewFileBasedKeyStore(nil, tempDir, false)
	require.NoError(t, err)

	return ks, func() { os.RemoveAll(tempDir) }
}

func newImpl(t *testing.T, opts PKCS11Opts) (*impl, func()) {
	ks, ksCleanup := newKeyStore(t)
	provider, err := New(opts, ks)
	require.NoError(t, err)

	pi := provider.(*impl)
	cleanup := func() {
		pi.ctx.Destroy()
		ksCleanup()
	}
	return pi, cleanup
}

func TestNew(t *testing.T) {
	ks, cleanup := newKeyStore(t)
	defer cleanup()

	t.Run("DefaultConfig", func(t *testing.T) {
		opts := defaultOptions()
		opts.createSessionRetryDelay = 0

		provider, err := New(opts, ks)
		require.NoError(t, err)
		require.NotNil(t, provider)
		pi := provider.(*impl)
		defer func() { pi.ctx.Destroy() }()

		curve, err := curveForSecurityLevel(opts.Security)
		require.NoError(t, err)

		require.NotNil(t, pi.BCCSP)
		require.Equal(t, opts.Pin, pi.pin)
		require.NotNil(t, pi.ctx)
		require.True(t, curve.Equal(pi.curve))
		require.Equal(t, opts.SoftwareVerify, pi.softVerify)
		require.Equal(t, opts.Immutable, pi.immutable)
		require.Equal(t, defaultCreateSessionRetries, pi.createSessionRetries)
		require.Equal(t, defaultCreateSessionRetryDelay, pi.createSessionRetryDelay)
		require.Equal(t, defaultSessionCacheSize, cap(pi.sessPool))
	})
	t.Run("ConditionalOverride", func(t *testing.T) {
		opts := defaultOptions()
		opts.createSessionRetries = 3
		opts.createSessionRetryDelay = time.Second
		opts.sessionCacheSize = -1

		provider, err := New(opts, ks)
		require.NoError(t, err)
		require.NotNil(t, provider)
		pi := provider.(*impl)
		defer func() { pi.ctx.Destroy() }()

		require.Equal(t, 3, pi.createSessionRetries)
		require.Equal(t, time.Second, pi.createSessionRetryDelay)
		require.Nil(t, pi.sessPool)
	})
}

func TestInvalidNewParameter(t *testing.T) {
	ks, cleanup := newKeyStore(t)
	defer cleanup()

	t.Run("BadSecurityLevel", func(t *testing.T) {
		opts := defaultOptions()
		opts.Security = 0

		_, err := New(opts, ks)
		require.EqualError(t, err, "Failed initializing configuration: Security level not supported [0]")
	})

	t.Run("BadHashFamily", func(t *testing.T) {
		opts := defaultOptions()
		opts.Hash = "SHA8"

		_, err := New(opts, ks)
		require.EqualError(t, err, "Failed initializing fallback SW BCCSP: Failed initializing configuration at [256,SHA8]: Hash Family not supported [SHA8]")
	})

	t.Run("BadKeyStore", func(t *testing.T) {
		_, err := New(defaultOptions(), nil)
		require.EqualError(t, err, "Failed initializing fallback SW BCCSP: Invalid bccsp.KeyStore instance. It must be different from nil.")
	})

	t.Run("MissingLibrary", func(t *testing.T) {
		opts := defaultOptions()
		opts.Library = ""

		_, err := New(opts, ks)
		require.Error(t, err)
		require.Contains(t, err.Error(), "pkcs11: library path not provided")
	})
}

func TestFindPKCS11LibEnvVars(t *testing.T) {
	const (
		dummy_PKCS11_LIB   = "/usr/lib/pkcs11"
		dummy_PKCS11_PIN   = "23456789"
		dummy_PKCS11_LABEL = "testing"
	)

	// Set environment variables used for test and preserve
	// original values for restoration after test completion
	orig_PKCS11_LIB := os.Getenv("PKCS11_LIB")
	orig_PKCS11_PIN := os.Getenv("PKCS11_PIN")
	orig_PKCS11_LABEL := os.Getenv("PKCS11_LABEL")

	t.Run("ExplicitEnvironment", func(t *testing.T) {
		os.Setenv("PKCS11_LIB", dummy_PKCS11_LIB)
		os.Setenv("PKCS11_PIN", dummy_PKCS11_PIN)
		os.Setenv("PKCS11_LABEL", dummy_PKCS11_LABEL)

		lib, pin, label := FindPKCS11Lib()
		require.EqualValues(t, dummy_PKCS11_LIB, lib, "FindPKCS11Lib did not return expected library")
		require.EqualValues(t, dummy_PKCS11_PIN, pin, "FindPKCS11Lib did not return expected pin")
		require.EqualValues(t, dummy_PKCS11_LABEL, label, "FindPKCS11Lib did not return expected label")
	})

	t.Run("MissingEnvironment", func(t *testing.T) {
		os.Unsetenv("PKCS11_LIB")
		os.Unsetenv("PKCS11_PIN")
		os.Unsetenv("PKCS11_LABEL")

		_, pin, label := FindPKCS11Lib()
		require.EqualValues(t, "98765432", pin, "FindPKCS11Lib did not return expected pin")
		require.EqualValues(t, "ForFabric", label, "FindPKCS11Lib did not return expected label")
	})

	os.Setenv("PKCS11_LIB", orig_PKCS11_LIB)
	os.Setenv("PKCS11_PIN", orig_PKCS11_PIN)
	os.Setenv("PKCS11_LABEL", orig_PKCS11_LABEL)
}

func TestInvalidSKI(t *testing.T) {
	pi, cleanup := newImpl(t, defaultOptions())
	defer cleanup()

	_, err := pi.GetKey(nil)
	require.EqualError(t, err, "Failed getting key for SKI [[]]: invalid SKI. Cannot be of zero length")

	_, err = pi.GetKey([]byte{0, 1, 2, 3, 4, 5, 6})
	require.Error(t, err)
	require.True(t, strings.HasPrefix(err.Error(), "Failed getting key for SKI [[0 1 2 3 4 5 6]]: "))
}

func TestKeyGenECDSAOpts(t *testing.T) {
	tests := map[string]struct {
		curve     elliptic.Curve
		immutable bool
		opts      bccsp.KeyGenOpts
	}{
		"Default":             {elliptic.P256(), false, &bccsp.ECDSAKeyGenOpts{Temporary: false}},
		"P256":                {elliptic.P256(), false, &bccsp.ECDSAP256KeyGenOpts{Temporary: false}},
		"P384":                {elliptic.P384(), false, &bccsp.ECDSAP384KeyGenOpts{Temporary: false}},
		"Immutable":           {elliptic.P384(), true, &bccsp.ECDSAP384KeyGenOpts{Temporary: false}},
		"Ephemeral/Default":   {elliptic.P256(), false, &bccsp.ECDSAKeyGenOpts{Temporary: true}},
		"Ephemeral/P256":      {elliptic.P256(), false, &bccsp.ECDSAP256KeyGenOpts{Temporary: true}},
		"Ephemeral/P384":      {elliptic.P384(), false, &bccsp.ECDSAP384KeyGenOpts{Temporary: true}},
		"Ephemeral/Immutable": {elliptic.P384(), true, &bccsp.ECDSAP384KeyGenOpts{Temporary: false}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			opts := defaultOptions()
			opts.Immutable = tt.immutable
			pi, cleanup := newImpl(t, opts)
			defer cleanup()

			k, err := pi.KeyGen(tt.opts)
			require.NoError(t, err)
			require.True(t, k.Private(), "key should be private")
			require.False(t, k.Symmetric(), "key should be asymmetric")

			ecdsaKey := k.(*ecdsaPrivateKey).pub
			require.Equal(t, tt.curve, ecdsaKey.pub.Curve, "wrong curve")

			raw, err := k.Bytes()
			require.EqualError(t, err, "Not supported.")
			require.Empty(t, raw, "result should be empty")

			pk, err := k.PublicKey()
			require.NoError(t, err)
			require.NotNil(t, pk)

			sess, err := pi.getSession()
			require.NoError(t, err)
			defer pi.returnSession(sess)

			for _, kt := range []keyType{publicKeyType, privateKeyType} {
				handle, err := pi.findKeyPairFromSKI(sess, k.SKI(), kt)
				require.NoError(t, err)

				attr, err := pi.ctx.GetAttributeValue(sess, handle, []*pkcs11.Attribute{{Type: pkcs11.CKA_TOKEN}})
				require.NoError(t, err)
				require.Len(t, attr, 1)

				if tt.opts.Ephemeral() {
					require.Equal(t, []byte{0}, attr[0].Value)
				} else {
					require.Equal(t, []byte{1}, attr[0].Value)
				}

				attr, err = pi.ctx.GetAttributeValue(sess, handle, []*pkcs11.Attribute{{Type: pkcs11.CKA_MODIFIABLE}})
				require.NoError(t, err)
				require.Len(t, attr, 1)

				if tt.immutable {
					require.Equal(t, []byte{0}, attr[0].Value)
				} else {
					require.Equal(t, []byte{1}, attr[0].Value)
				}
			}
		})
	}
}

func TestKeyGenMissingOpts(t *testing.T) {
	pi, cleanup := newImpl(t, defaultOptions())
	defer cleanup()

	_, err := pi.KeyGen(bccsp.KeyGenOpts(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), "Invalid Opts parameter. It must not be nil")
}

func TestECDSAGetKeyBySKI(t *testing.T) {
	pi, cleanup := newImpl(t, defaultOptions())
	defer cleanup()

	k, err := pi.KeyGen(&bccsp.ECDSAKeyGenOpts{Temporary: false})
	require.NoError(t, err)

	k2, err := pi.GetKey(k.SKI())
	require.NoError(t, err)

	require.True(t, k2.Private(), "key should be private")
	require.False(t, k2.Symmetric(), "key should be asymmetric")
	require.Equalf(t, k.SKI(), k2.SKI(), "expected %x got %x", k.SKI(), k2.SKI())
}

func TestECDSAPublicKeyFromPrivateKey(t *testing.T) {
	pi, cleanup := newImpl(t, defaultOptions())
	defer cleanup()

	k, err := pi.KeyGen(&bccsp.ECDSAKeyGenOpts{Temporary: false})
	require.NoError(t, err)

	pk, err := k.PublicKey()
	require.NoError(t, err)
	require.False(t, pk.Private(), "key should be public")
	require.False(t, pk.Symmetric(), "key should be asymmetric")
	require.Equal(t, k.SKI(), pk.SKI(), "SKI should be the same")

	raw, err := pk.Bytes()
	require.NoError(t, err)
	require.NotEmpty(t, raw, "marshaled ECDSA public key must not be empty")
}

func TestECDSASign(t *testing.T) {
	pi, cleanup := newImpl(t, defaultOptions())
	defer cleanup()

	k, err := pi.KeyGen(&bccsp.ECDSAKeyGenOpts{Temporary: false})
	require.NoError(t, err)

	digest, err := pi.Hash([]byte("Hello World"), &bccsp.SHAOpts{})
	require.NoError(t, err)

	signature, err := pi.Sign(k, digest, nil)
	require.NoError(t, err)
	require.NotEmpty(t, signature, "signature must not be empty")

	t.Run("NoKey", func(t *testing.T) {
		_, err := pi.Sign(nil, digest, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Invalid Key. It must not be nil")
	})

	t.Run("BadSKI", func(t *testing.T) {
		_, err := pi.Sign(&ecdsaPrivateKey{ski: []byte("bad-ski")}, digest, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Private key not found")
	})

	t.Run("MissingDigest", func(t *testing.T) {
		_, err = pi.Sign(k, nil, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Invalid digest. Cannot be empty")
	})
}

func TestECDSAVerify(t *testing.T) {
	pi, cleanup := newImpl(t, defaultOptions())
	defer cleanup()

	k, err := pi.KeyGen(&bccsp.ECDSAKeyGenOpts{Temporary: false})
	require.NoError(t, err)
	pk, err := k.PublicKey()
	require.NoError(t, err)

	digest, err := pi.Hash([]byte("Hello, World."), &bccsp.SHAOpts{})
	require.NoError(t, err)
	otherDigest, err := pi.Hash([]byte("Bye, World."), &bccsp.SHAOpts{})
	require.NoError(t, err)

	signature, err := pi.Sign(k, digest, nil)
	require.NoError(t, err)

	tests := map[string]bool{
		"WithSoftVerify":    true,
		"WithoutSoftVerify": false,
	}
	for name, softVerify := range tests {
		t.Run(name, func(t *testing.T) {
			opts := defaultOptions()
			opts.SoftwareVerify = softVerify
			pi, cleanup := newImpl(t, opts)
			defer cleanup()

			valid, err := pi.Verify(k, signature, digest, nil)
			require.NoError(t, err)
			require.True(t, valid, "signature should be valid from private key")

			valid, err = pi.Verify(pk, signature, digest, nil)
			require.NoError(t, err)
			require.True(t, valid, "signature should be valid from public key")

			valid, err = pi.Verify(k, signature, otherDigest, nil)
			require.NoError(t, err)
			require.False(t, valid, "signature should be valid from private key")

			valid, err = pi.Verify(pk, signature, otherDigest, nil)
			require.NoError(t, err)
			require.False(t, valid, "signature should not be valid from public key")
		})
	}

	t.Run("MissingKey", func(t *testing.T) {
		_, err := pi.Verify(nil, signature, digest, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Invalid Key. It must not be nil")
	})

	t.Run("MissingSignature", func(t *testing.T) {
		_, err := pi.Verify(pk, nil, digest, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Invalid signature. Cannot be empty")
	})

	t.Run("MissingDigest", func(t *testing.T) {
		_, err = pi.Verify(pk, signature, nil, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Invalid digest. Cannot be empty")
	})
}

func TestECDSALowS(t *testing.T) {
	pi, cleanup := newImpl(t, defaultOptions())
	defer cleanup()

	k, err := pi.KeyGen(&bccsp.ECDSAKeyGenOpts{Temporary: false})
	require.NoError(t, err)

	digest, err := pi.Hash([]byte("Hello World"), &bccsp.SHAOpts{})
	require.NoError(t, err)

	// Ensure that signature with low-S are generated
	t.Run("GeneratesLowS", func(t *testing.T) {
		signature, err := pi.Sign(k, digest, nil)
		require.NoError(t, err)

		_, S, err := utils.UnmarshalECDSASignature(signature)
		require.NoError(t, err)

		if S.Cmp(utils.GetCurveHalfOrdersAt(k.(*ecdsaPrivateKey).pub.pub.Curve)) >= 0 {
			t.Fatal("Invalid signature. It must have low-S")
		}

		valid, err := pi.Verify(k, signature, digest, nil)
		require.NoError(t, err)
		require.True(t, valid, "signature should be valid")
	})

	// Ensure that signature with high-S are rejected.
	t.Run("RejectsHighS", func(t *testing.T) {
		for {
			R, S, err := pi.signP11ECDSA(k.SKI(), digest)
			require.NoError(t, err)
			if S.Cmp(utils.GetCurveHalfOrdersAt(k.(*ecdsaPrivateKey).pub.pub.Curve)) > 0 {
				sig, err := utils.MarshalECDSASignature(R, S)
				require.NoError(t, err)

				valid, err := pi.Verify(k, sig, digest, nil)
				require.Error(t, err, "verification must fail for a signature with high-S")
				require.False(t, valid, "signature must not be valid with high-S")
				return
			}
		}
	})
}

func TestInitialize(t *testing.T) {
	// Setup PKCS11 library and provide initial set of values
	lib, pin, label := FindPKCS11Lib()

	t.Run("MissingLibrary", func(t *testing.T) {
		_, err := (&impl{}).initialize(PKCS11Opts{Library: "", Pin: pin, Label: label})
		require.Error(t, err)
		require.Contains(t, err.Error(), "pkcs11: library path not provided")
	})

	t.Run("BadLibraryPath", func(t *testing.T) {
		_, err := (&impl{}).initialize(PKCS11Opts{Library: "badLib", Pin: pin, Label: label})
		require.Error(t, err)
		require.Contains(t, err.Error(), "pkcs11: instantiation failed for badLib")
	})

	t.Run("BadLabel", func(t *testing.T) {
		_, err := (&impl{}).initialize(PKCS11Opts{Library: lib, Pin: pin, Label: "badLabel"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "could not find token with label")
	})

	t.Run("MissingPin", func(t *testing.T) {
		_, err := (&impl{}).initialize(PKCS11Opts{Library: lib, Pin: "", Label: label})
		require.Error(t, err)
		require.Contains(t, err.Error(), "Login failed: pkcs11")
	})
}

func TestNamedCurveFromOID(t *testing.T) {
	tests := map[string]struct {
		oid   asn1.ObjectIdentifier
		curve elliptic.Curve
	}{
		"P224":    {oidNamedCurveP224, elliptic.P224()},
		"P256":    {oidNamedCurveP256, elliptic.P256()},
		"P384":    {oidNamedCurveP384, elliptic.P384()},
		"P521":    {oidNamedCurveP521, elliptic.P521()},
		"unknown": {asn1.ObjectIdentifier{4, 9, 15, 1}, nil},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.curve, namedCurveFromOID(tt.oid))
		})
	}
}

func TestCurveForSecurityLevel(t *testing.T) {
	tests := map[int]struct {
		expectedErr string
		curve       asn1.ObjectIdentifier
	}{
		256: {curve: oidNamedCurveP256},
		384: {curve: oidNamedCurveP384},
		512: {expectedErr: "Security level not supported [512]"},
	}

	for level, tt := range tests {
		t.Run(strconv.Itoa(level), func(t *testing.T) {
			curve, err := curveForSecurityLevel(level)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.curve, curve)
		})
	}
}

func TestPKCS11GetSession(t *testing.T) {
	opts := defaultOptions()
	opts.sessionCacheSize = 5
	pi, cleanup := newImpl(t, opts)
	defer cleanup()

	sessionCacheSize := opts.sessionCacheSize
	var sessions []pkcs11.SessionHandle
	for i := 0; i < 3*sessionCacheSize; i++ {
		session, err := pi.getSession()
		require.NoError(t, err)
		sessions = append(sessions, session)
	}

	// Return all sessions, should leave sessionCacheSize cached
	for _, session := range sessions {
		pi.returnSession(session)
	}

	// Lets break OpenSession, so non-cached session cannot be opened
	oldSlot := pi.slot
	pi.slot = ^uint(0)

	// Should be able to get sessionCacheSize cached sessions
	sessions = nil
	for i := 0; i < sessionCacheSize; i++ {
		session, err := pi.getSession()
		require.NoError(t, err)
		sessions = append(sessions, session)
	}

	_, err := pi.getSession()
	require.EqualError(t, err, "OpenSession failed: pkcs11: 0x3: CKR_SLOT_ID_INVALID")

	// Load cache with bad sessions
	for i := 0; i < sessionCacheSize; i++ {
		pi.returnSession(pkcs11.SessionHandle(^uint(0)))
	}

	// Fix OpenSession so non-cached sessions can be opened
	pi.slot = oldSlot

	// Request a session, return, and re-acquire. The pool should be emptied
	// before creating a new session so when returned, it should be the only
	// session in the cache.
	sess, err := pi.getSession()
	require.NoError(t, err)
	pi.returnSession(sess)
	sess2, err := pi.getSession()
	require.NoError(t, err)
	require.Equal(t, sess, sess2, "expected to get back the same session")

	// Cleanup
	for _, session := range sessions {
		pi.returnSession(session)
	}
}

func TestSessionHandleCaching(t *testing.T) {
	verifyHandleCache := func(t *testing.T, pi *impl, sess pkcs11.SessionHandle, k bccsp.Key) {
		pubHandle, err := pi.findKeyPairFromSKI(sess, k.SKI(), publicKeyType)
		require.NoError(t, err)
		h, ok := pi.cachedHandle(publicKeyType, k.SKI())
		require.True(t, ok)
		require.Equal(t, h, pubHandle)

		privHandle, err := pi.findKeyPairFromSKI(sess, k.SKI(), privateKeyType)
		require.NoError(t, err)
		h, ok = pi.cachedHandle(privateKeyType, k.SKI())
		require.True(t, ok)
		require.Equal(t, h, privHandle)
	}

	t.Run("SessionCacheDisabled", func(t *testing.T) {
		opts := defaultOptions()
		opts.sessionCacheSize = -1

		pi, cleanup := newImpl(t, opts)
		defer cleanup()

		require.Nil(t, pi.sessPool, "sessPool channel should be nil")
		require.Empty(t, pi.sessions, "sessions set should be empty")
		require.Empty(t, pi.handleCache, "handleCache should be empty")

		sess1, err := pi.getSession()
		require.NoError(t, err)
		require.Len(t, pi.sessions, 1, "expected one open session")

		sess2, err := pi.getSession()
		require.NoError(t, err)
		require.Len(t, pi.sessions, 2, "expected two open sessions")

		// Generate a key
		k, err := pi.KeyGen(&bccsp.ECDSAP256KeyGenOpts{Temporary: false})
		require.NoError(t, err)
		verifyHandleCache(t, pi, sess1, k)
		require.Len(t, pi.handleCache, 2, "expected two handles in handle cache")

		pi.returnSession(sess1)
		require.Len(t, pi.sessions, 1, "expected one open session")
		verifyHandleCache(t, pi, sess1, k)
		require.Len(t, pi.handleCache, 2, "expected two handles in handle cache")

		pi.returnSession(sess2)
		require.Empty(t, pi.sessions, "expected sessions to be empty")
		require.Empty(t, pi.handleCache, "expected handles to be cleared")

		pi.slot = ^uint(0) // break OpenSession
		_, err = pi.getSession()
		require.EqualError(t, err, "OpenSession failed: pkcs11: 0x3: CKR_SLOT_ID_INVALID")
		require.Empty(t, pi.sessions, "expected sessions to be empty")
	})

	t.Run("SessionCacheEnabled", func(t *testing.T) {
		opts := defaultOptions()
		opts.sessionCacheSize = 1

		pi, cleanup := newImpl(t, opts)
		defer cleanup()

		require.NotNil(t, pi.sessPool, "sessPool channel should not be nil")
		require.Equal(t, 1, cap(pi.sessPool))
		require.Len(t, pi.sessions, 1, "sessions should contain login session")
		require.Len(t, pi.sessPool, 1, "sessionPool should hold login session")
		require.Empty(t, pi.handleCache, "handleCache should be empty")

		sess1, err := pi.getSession()
		require.NoError(t, err)
		require.Len(t, pi.sessions, 1, "expected one open session (sess1 from login)")
		require.Len(t, pi.sessPool, 0, "sessionPool should be empty")

		sess2, err := pi.getSession()
		require.NoError(t, err)
		require.Len(t, pi.sessions, 2, "expected two open sessions (sess1 and sess2)")
		require.Len(t, pi.sessPool, 0, "sessionPool should be empty")

		// Generate a key
		k, err := pi.KeyGen(&bccsp.ECDSAP256KeyGenOpts{Temporary: false})
		require.NoError(t, err)
		verifyHandleCache(t, pi, sess1, k)
		require.Len(t, pi.handleCache, 2, "expected two handles in handle cache")

		pi.returnSession(sess1)
		require.Len(t, pi.sessions, 2, "expected two open sessions (sess2 in-use, sess1 cached)")
		require.Len(t, pi.sessPool, 1, "sessionPool should have one handle (sess1)")
		verifyHandleCache(t, pi, sess1, k)
		require.Len(t, pi.handleCache, 2, "expected two handles in handle cache")

		pi.returnSession(sess2)
		require.Len(t, pi.sessions, 1, "expected one cached session (sess1)")
		require.Len(t, pi.sessPool, 1, "sessionPool should have one handle (sess1)")
		require.Len(t, pi.handleCache, 2, "expected two handles in handle cache")

		sess1, err = pi.getSession()
		require.NoError(t, err)
		require.Len(t, pi.sessions, 1, "expected one open session (sess1)")
		require.Len(t, pi.sessPool, 0, "sessionPool should be empty")
		require.Len(t, pi.handleCache, 2, "expected two handles in handle cache")

		pi.slot = ^uint(0) // break OpenSession
		_, err = pi.getSession()
		require.EqualError(t, err, "OpenSession failed: pkcs11: 0x3: CKR_SLOT_ID_INVALID")
		require.Len(t, pi.sessions, 1, "expected one active session (sess1)")
		require.Len(t, pi.sessPool, 0, "sessionPool should be empty")
		require.Len(t, pi.handleCache, 2, "expected two handles in handle cache")

		// Return a busted session that should be cached
		pi.returnSession(pkcs11.SessionHandle(^uint(0)))
		require.Len(t, pi.sessions, 1, "expected one active session (sess1)")
		require.Len(t, pi.sessPool, 1, "sessionPool should contain busted session")
		require.Len(t, pi.handleCache, 2, "expected two handles in handle cache")

		// Return sess1 that should be discarded
		pi.returnSession(sess1)
		require.Len(t, pi.sessions, 0, "expected sess1 to be removed")
		require.Len(t, pi.sessPool, 1, "sessionPool should contain busted session")
		require.Empty(t, pi.handleCache, "expected handles to be purged on removal of last tracked session")

		// Try to get broken session from cache
		_, err = pi.getSession()
		require.EqualError(t, err, "OpenSession failed: pkcs11: 0x3: CKR_SLOT_ID_INVALID")
		require.Empty(t, pi.sessions, "expected sessions to be empty")
		require.Len(t, pi.sessPool, 0, "sessionPool should be empty")
	})
}

func TestKeyCache(t *testing.T) {
	opts := defaultOptions()
	opts.sessionCacheSize = 1
	pi, cleanup := newImpl(t, opts)
	defer cleanup()

	require.Empty(t, pi.keyCache)

	_, err := pi.GetKey([]byte("nonsense-key"))
	require.Error(t, err) // message comes from software keystore
	require.Empty(t, pi.keyCache)

	k, err := pi.KeyGen(&bccsp.ECDSAP256KeyGenOpts{Temporary: false})
	require.NoError(t, err)
	_, ok := pi.cachedKey(k.SKI())
	require.False(t, ok, "created keys are not (currently) cached")

	key, err := pi.GetKey(k.SKI())
	require.NoError(t, err)
	cached, ok := pi.cachedKey(k.SKI())
	require.True(t, ok, "key should be cached")
	require.Same(t, key, cached, "key from cache should be what was found")

	// Kill all valid cached sessions
	pi.slot = ^uint(0)
	sess, err := pi.getSession()
	require.NoError(t, err)
	require.Len(t, pi.sessions, 1, "should have one active session")
	require.Len(t, pi.sessPool, 0, "sessionPool should be empty")

	pi.returnSession(pkcs11.SessionHandle(^uint(0)))
	require.Len(t, pi.sessions, 1, "should have one active session")
	require.Len(t, pi.sessPool, 1, "sessionPool should be empty")

	_, ok = pi.cachedKey(k.SKI())
	require.True(t, ok, "key should remain in cache due to active sessions")

	// Force caches to be cleared
	pi.returnSession(sess)
	require.Empty(t, pi.sessions, "sessions should be empty")
	require.Empty(t, pi.keyCache, "key cache should be empty")

	_, ok = pi.cachedKey(k.SKI())
	require.False(t, ok, "key should not be in cache")
}

// This helps verify that we're delegating to the software provider.
// This is not intended to test the software provider implementation.
func TestDelegation(t *testing.T) {
	pi, cleanup := newImpl(t, defaultOptions())
	defer cleanup()

	k, err := pi.KeyGen(&bccsp.AES256KeyGenOpts{})
	require.NoError(t, err)

	t.Run("KeyGen", func(t *testing.T) {
		k, err := pi.KeyGen(&bccsp.AES256KeyGenOpts{})
		require.NoError(t, err)
		require.True(t, k.Private())
		require.True(t, k.Symmetric())
	})

	t.Run("KeyDeriv", func(t *testing.T) {
		k, err := pi.KeyDeriv(k, &bccsp.HMACDeriveKeyOpts{Arg: []byte{1}})
		require.NoError(t, err)
		require.True(t, k.Private())
	})

	t.Run("KeyImport", func(t *testing.T) {
		raw := make([]byte, 32)
		_, err := rand.Read(raw)
		require.NoError(t, err)

		k, err := pi.KeyImport(raw, &bccsp.AES256ImportKeyOpts{})
		require.NoError(t, err)
		require.True(t, k.Private())
	})

	t.Run("GetKey", func(t *testing.T) {
		k, err := pi.GetKey(k.SKI())
		require.NoError(t, err)
		require.True(t, k.Private())
	})

	t.Run("Hash", func(t *testing.T) {
		digest, err := pi.Hash([]byte("message"), &bccsp.SHA3_384Opts{})
		require.NoError(t, err)
		require.NotEmpty(t, digest)
	})

	t.Run("GetHash", func(t *testing.T) {
		h, err := pi.GetHash(&bccsp.SHA256Opts{})
		require.NoError(t, err)
		require.Equal(t, sha256.New(), h)
	})

	t.Run("Sign", func(t *testing.T) {
		_, err := pi.Sign(k, []byte("message"), nil)
		require.EqualError(t, err, "Unsupported 'SignKey' provided [*sw.aesPrivateKey]")
	})

	t.Run("Verify", func(t *testing.T) {
		_, err := pi.Verify(k, []byte("signature"), []byte("digest"), nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Unsupported 'VerifyKey' provided")
	})

	t.Run("EncryptDecrypt", func(t *testing.T) {
		msg := []byte("message")
		ct, err := pi.Encrypt(k, msg, &bccsp.AESCBCPKCS7ModeOpts{})
		require.NoError(t, err)

		pt, err := pi.Decrypt(k, ct, &bccsp.AESCBCPKCS7ModeOpts{})
		require.NoError(t, err)
		require.Equal(t, msg, pt)
	})
}

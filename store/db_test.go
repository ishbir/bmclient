// Copyright (c) 2015 Monetas.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package store_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/monetas/bmclient/store"
)

func TestOpenClose(t *testing.T) {
	f, err := ioutil.TempFile("", "tempstore")
	if err != nil {
		panic(err)
	}
	fName := f.Name()
	f.Close()

	pass := []byte("password")
	passNew := []byte("new_password")

	// Create a new database

	s, err := store.Open(fName, pass)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	err = s.Close()
	if err != nil {
		t.Error("Failed to close database:", err)
	}

	// Try opening same database but with incorrect passphrase
	s, err = store.Open(fName, passNew)
	if err != store.ErrDecryptionFailed {
		t.Error("Expected ErrDecryptionFailed got", err)
	}

	// Try re-opening database with correct passphrase, to make sure decryption
	// works.
	s, err = store.Open(fName, pass)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}

	// Change passphrase and close database.
	err = s.ChangePassphrase(passNew)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Close()
	if err != nil {
		t.Error("Failed to close database:", err)
	}

	// Re-open database with new passphrase to see if ChangePassphrase was
	// successful.
	s, err = store.Open(fName, passNew)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	err = s.Close()
	if err != nil {
		t.Error("Failed to close database:", err)
	}

	os.Remove(fName)
}

// Package repository contains helper methods for working with a Git repo.
package repository

import (
	"bytes"
	"strings"

	"github.com/MichaelMure/git-bug/util/git"
	"github.com/MichaelMure/git-bug/util/lamport"
)

// RepoCommon represent the common function the we want all the repo to implement
type RepoCommon interface {
	// GetPath returns the path to the repo.
	GetPath() string

	// GetUserName returns the name the the user has used to configure git
	GetUserName() (string, error)

	// GetUserEmail returns the email address that the user has used to configure git.
	GetUserEmail() (string, error)

	// GetCoreEditor returns the name of the editor that the user has used to configure git.
	GetCoreEditor() (string, error)

	// StoreConfig store a single key/value pair in the config of the repo
	StoreConfig(key string, value string) error

	// ReadConfigs read all key/value pair matching the key prefix
	ReadConfigs(keyPrefix string) (map[string]string, error)

	// RmConfigs remove all key/value pair matching the key prefix
	RmConfigs(keyPrefix string) error
}

// Repo represents a source code repository.
type Repo interface {
	RepoCommon

	// FetchRefs fetch git refs from a remote
	FetchRefs(remote string, refSpec string) (string, error)

	// PushRefs push git refs to a remote
	PushRefs(remote string, refSpec string) (string, error)

	// StoreData will store arbitrary data and return the corresponding hash
	StoreData(data []byte) (git.Hash, error)

	// ReadData will attempt to read arbitrary data from the given hash
	ReadData(hash git.Hash) ([]byte, error)

	// StoreTree will store a mapping key-->Hash as a Git tree
	StoreTree(mapping []TreeEntry) (git.Hash, error)

	// StoreCommit will store a Git commit with the given Git tree
	StoreCommit(treeHash git.Hash) (git.Hash, error)

	// StoreCommit will store a Git commit with the given Git tree
	StoreCommitWithParent(treeHash git.Hash, parent git.Hash) (git.Hash, error)

	// UpdateRef will create or update a Git reference
	UpdateRef(ref string, hash git.Hash) error

	// ListRefs will return a list of Git ref matching the given refspec
	ListRefs(refspec string) ([]string, error)

	// RefExist will check if a reference exist in Git
	RefExist(ref string) (bool, error)

	// CopyRef will create a new reference with the same value as another one
	CopyRef(source string, dest string) error

	// ListCommits will return the list of tree hashes of a ref, in chronological order
	ListCommits(ref string) ([]git.Hash, error)

	// ListEntries will return the list of entries in a Git tree
	ListEntries(hash git.Hash) ([]TreeEntry, error)

	// FindCommonAncestor will return the last common ancestor of two chain of commit
	FindCommonAncestor(hash1 git.Hash, hash2 git.Hash) (git.Hash, error)

	// GetTreeHash return the git tree hash referenced in a commit
	GetTreeHash(commit git.Hash) (git.Hash, error)
}

type ClockedRepo interface {
	Repo

	// LoadClocks read the clocks values from the on-disk repo
	LoadClocks() error

	// WriteClocks write the clocks values into the repo
	WriteClocks() error

	// CreateTimeIncrement increment the creation clock and return the new value.
	CreateTimeIncrement() (lamport.Time, error)

	// EditTimeIncrement increment the edit clock and return the new value.
	EditTimeIncrement() (lamport.Time, error)

	// CreateWitness witness another create time and increment the corresponding
	// clock if needed.
	CreateWitness(time lamport.Time) error

	// EditWitness witness another edition time and increment the corresponding
	// clock if needed.
	EditWitness(time lamport.Time) error
}

func prepareTreeEntries(entries []TreeEntry) bytes.Buffer {
	var buffer bytes.Buffer

	for _, entry := range entries {
		buffer.WriteString(entry.Format())
	}

	return buffer
}

func readTreeEntries(s string) ([]TreeEntry, error) {
	split := strings.Split(s, "\n")

	casted := make([]TreeEntry, len(split))
	for i, line := range split {
		if line == "" {
			continue
		}

		entry, err := ParseTreeEntry(line)

		if err != nil {
			return nil, err
		}

		casted[i] = entry
	}

	return casted, nil
}

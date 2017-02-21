package snapshot

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/containerd"
	"github.com/docker/containerd/testutil"
	"github.com/stretchr/testify/assert"
)

/*
make sure the content of parent appears in the child (#523)
make sure Remove() fails for in-use snapshots
make sure View() returns RO mounts, and Prepare() returns RW ones
Commit() can be called while snapshot is in use.
Commit() called serially results in the original parent from the start of the transaction.
X Create a New Layer on top of base layer with Prepare, Stat on new layer, should return Active layer.
X Commit a New Layer on top of base layer with Prepare & Commit , Stat on new layer, should return Committed layer.
X Creating two layers with Prepare or View with same key must fail.
Multiple Commit on same key should be successful & contents should match with last Commit changes.
Removing intermediate snapshot should always fail.
Deletion of files/folder of base layer in new layer, On Commit, those files should not be visible.
Movement of files/folder from base layer to folder in new layer should be allowed.
Verify Copy at destination and deletion at source.
Modification on same file with multiple commits (new layers). Read on top layer, should result last modification which was committed by previous layer.
*/

// Create a New Layer on top of base layer with Prepare, Stat on new layer, should return Active layer.
func testSnapshotterStatActive(t *testing.T, snapshotter Snapshotter, work string) {
	preparing := filepath.Join(work, "preparing")
	if err := os.MkdirAll(preparing, 0777); err != nil {
		t.Fatal(err)
	}

	mounts, err := snapshotter.Prepare(preparing, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(mounts) < 1 {
		t.Fatal("expected mounts to have entries")
	}

	if err = containerd.MountAll(mounts, preparing); err != nil {
		t.Fatal(err)
	}
	defer testutil.Unmount(t, preparing)

	if err = ioutil.WriteFile(filepath.Join(preparing, "foo"), []byte("foo\n"), 0777); err != nil {
		t.Fatal(err)
	}

	si, err := snapshotter.Stat(preparing)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, si.Name, preparing)
}

// Commit a New Layer on top of base layer with Prepare & Commit , Stat on new layer, should return Committed layer.
func testSnapshotterStatCommited(t *testing.T, snapshotter Snapshotter, work string) {
	preparing := filepath.Join(work, "preparing")
	if err := os.MkdirAll(preparing, 0777); err != nil {
		t.Fatal(err)
	}

	mounts, err := snapshotter.Prepare(preparing, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(mounts) < 1 {
		t.Fatal("expected mounts to have entries")
	}

	if err = containerd.MountAll(mounts, preparing); err != nil {
		t.Fatal(err)
	}
	defer testutil.Unmount(t, preparing)

	if err = ioutil.WriteFile(filepath.Join(preparing, "foo"), []byte("foo\n"), 0777); err != nil {
		t.Fatal(err)
	}

	committed := filepath.Join(work, "committed")
	if err = snapshotter.Commit(committed, preparing); err != nil {
		t.Fatal(err)
	}

	si, err := snapshotter.Stat(committed)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, si.Name, committed)
}

// Creating two layers with Prepare or View with same key must fail.
func testSnapshotterSameKeyPrepareView(t *testing.T, snapshotter Snapshotter, work string) {
	newlayer := filepath.Join(work, "newlayer")
	if err := os.MkdirAll(newlayer, 0777); err != nil {
		t.Fatal(err)
	}

	_, err := snapshotter.Prepare(newlayer, "")
	if err != nil {
		t.Fatal(err)
	}

	//Must Fail.
	_, err = snapshotter.View(newlayer, "")
	if err == nil {
		t.Errorf("expected error but View and Prepare created with same key")
	}

}

// Creating two layers with Prepare or View with same key must fail.
func testSnapshotterMultiCommit(t *testing.T, snapshotter Snapshotter, work string) {
	layer1 := filepath.Join(work, "layer1")
	if err := os.MkdirAll(newlayer, 0777); err != nil {
		t.Fatal(err)
	}

	mounts, err := snapshotter.Prepare(layer1, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(mounts) < 1 {
		t.Fatal("expected mounts to have entries")
	}

	if err = containerd.MountAll(mounts, layer1); err != nil {
		t.Fatal(err)
	}
	defer testutil.Unmount(t, layer1)

	if err = ioutil.WriteFile(filepath.Join(layer1, "foo"), []byte("foo\n"), 0777); err != nil {
		t.Fatal(err)
	}

	committed := filepath.Join(work, "committed")
	if err = snapshotter.Commit(committed, layer1); err != nil {
		t.Fatal(err)
	}

	if err = ioutil.WriteFile(filepath.Join(layer1, "foo"), []byte("foo1 bar\n"), 0777); err != nil {
		t.Fatal(err)
	}

	committed = filepath.Join(work, "committed")
	if err = snapshotter.Commit(committed, layer1); err != nil {
		t.Fatal(err)
	}

	if err = ioutil.WriteFile(filepath.Join(layer1, "foo"), []byte("foo2 bar baz\n"), 0777); err != nil {
		t.Fatal(err)
	}

	committed = filepath.Join(work, "committed")
	if err = snapshotter.Commit(committed, layer1); err != nil {
		t.Fatal(err)
	}

}

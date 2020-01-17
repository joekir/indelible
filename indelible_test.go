package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSetAppendOnly_withtmpfile_setsattr(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "indelible_test")
	if err != nil {
		t.Fatalf("Unable to create file: %v", err)
	}

	defer os.Remove(tmpfile.Name())

	path, err := filepath.Abs(tmpfile.Name())
	if err != nil {
		t.Fatalf("Unable to get Absolute Path of tmpfile: %v", err)
	}

	t.Log(path)
	err = setAppendOnly(path)
	if err != nil {
		t.Fatalf("Set Chattr +a failed: %v", err)
	}

	attr := int32(FS_EXTENT_FL)
	err = ioctl(tmpfile, _FS_IOC_SETFLAGS, &attr)
	if err != nil {
		t.Fatalf("Failed to restore file-attrs so it can be deleted: %v", err)
	}

	if err = tmpfile.Close(); err != nil {
		t.Fatalf("Unable to close file: %v", err)
	}
}

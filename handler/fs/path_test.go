//go:build unit
// +build unit

package fs_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"reactor-crw/handler/fs"
)

func TestNewPathResolver(t *testing.T) {
	paths := []struct {
		name        string
		msg         string
		errExpected bool
	}{
		{" ", "Expected an error during resolving empty path", true},
		{"invalid_path", "Expected an error during resolving non-existent dir", true},
		{os.TempDir(), "Wasn't expected an error", false},
	}

	for _, path := range paths {
		_, err := fs.NewPathResolver(path.name)
		if path.errExpected {
			require.Errorf(t, err, path.msg)
		}

		if !path.errExpected {
			require.NoError(t, err, path.msg)
		}
	}
}

func TestPathResolver_CreateFolder(t *testing.T) {
	p, _ := fs.NewPathResolver(os.TempDir())

	folders := []struct {
		name        string
		errExpected bool
	}{
		{name: "test"},
		{name: "_123_test"},
		{name: "109483___"},
		{name: "%&*$test", errExpected: false},
		{name: "byVPBIQvHVU?wmode=transparent&rel=0", errExpected: false},
	}

	for _, folder := range folders {
		err := p.CreateFolder(folder.name)
		if !folder.errExpected {
			require.NoError(t, err, "Wasn't expected an error during creating %s", folder.name)
		}

		if folder.errExpected {
			require.Error(t, err, "Expected an error during creating %s", folder.name)
		}

		p.Remove(folder.name)
	}

	p = &fs.PathResolver{}
	err := p.CreateFolder("test")
	require.Error(t, err, "Expected an error during creating folder")
}

func TestPathResolver_CreateFile(t *testing.T) {
	p, _ := fs.NewPathResolver(os.TempDir())
	_, err := p.CreateFile("filename")
	require.Error(t, err, "Expected an error during file. Folder wasn't created yet.")

	_ = p.CreateFolder("test")

	//_, err = p.CreateFile(" ")
	//require.Error(t, err, "Expected an error during creating file with an empty name")

	_, err = p.CreateFile("filename")
	require.NoError(t, err, "Wasn't expected an error during creating filename")

	filePath := "./test/filename"
	_, err = os.Stat(filePath)
	require.NoError(t, err, "File %s wasn't created", filePath)

	p.Remove(filePath)
}

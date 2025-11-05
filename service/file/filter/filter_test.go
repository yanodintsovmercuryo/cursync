package filter_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const (
	testDir = "/test"
)

func TestFilter_GetFilePatterns(t *testing.T) {
	t.Run("from flag", func(t *testing.T) {
		t.Parallel()
		f, finish := setUp(t)
		defer finish()

		flagValue := "*.txt,*.md"

		result, err := f.filter.GetFilePatterns(flagValue)
		require.NoError(t, err)

		expected := []string{"*.txt", "*.md"}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("empty patterns", func(t *testing.T) {
		t.Parallel()
		f, finish := setUp(t)
		defer finish()

		flagValue := ""

		result, err := f.filter.GetFilePatterns(flagValue)
		require.NoError(t, err)

		require.Empty(t, result)
	})

	t.Run("with spaces and empty values", func(t *testing.T) {
		t.Parallel()
		f, finish := setUp(t)
		defer finish()

		flagValue := " *.txt , *.md , "

		result, err := f.filter.GetFilePatterns(flagValue)
		require.NoError(t, err)

		expected := []string{"*.txt", "*.md"}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestFilter_FindFilesByPatterns(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		f, finish := setUp(t)
		defer finish()

		dir := testDir
		patterns := []string{"*.txt"}
		allFiles := []string{"/test/file1.txt", "/test/file2.txt", "/test/file3.md"}

		f.fileOpsMock.EXPECT().
			FindAllFiles(dir).
			Return(allFiles, nil).
			Times(1)

		result, err := f.filter.FindFilesByPatterns(dir, patterns)
		require.NoError(t, err)

		expectedCount := 2
		if diff := cmp.Diff(expectedCount, len(result)); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("error finding files", func(t *testing.T) {
		t.Parallel()
		f, finish := setUp(t)
		defer finish()

		dir := testDir
		patterns := []string{"*.txt"}
		expectedErr := errors.New("find error")

		f.fileOpsMock.EXPECT().
			FindAllFiles(dir).
			Return(nil, expectedErr).
			Times(1)

		result, err := f.filter.FindFilesByPatterns(dir, patterns)
		require.ErrorIs(t, err, expectedErr)
		require.Empty(t, result)
	})

	t.Run("empty patterns returns all files", func(t *testing.T) {
		t.Parallel()
		f, finish := setUp(t)
		defer finish()

		dir := testDir
		patterns := []string{}
		allFiles := []string{"/test/file1.txt", "/test/file2.md"}

		f.fileOpsMock.EXPECT().
			FindAllFiles(dir).
			Return(allFiles, nil).
			Times(1)

		result, err := f.filter.FindFilesByPatterns(dir, patterns)
		require.NoError(t, err)

		if diff := cmp.Diff(allFiles, result); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestFilter_CleanupExtraFilesByPatterns(t *testing.T) {
	t.Run("error walking destination directory", func(t *testing.T) {
		t.Parallel()
		f, finish := setUp(t)
		defer finish()

		srcFiles := []string{"/src/file.txt"}
		srcBase := "/src"
		dstBase := "/dst"
		patterns := []string{"*.txt"}
		expectedErr := errors.New("walk error")

		f.fileOpsMock.EXPECT().
			FindAllFiles(dstBase).
			Return(nil, expectedErr).
			Times(1)

		err := f.filter.CleanupExtraFilesByPatterns(srcFiles, srcBase, dstBase, patterns)
		require.Error(t, err)
		require.Contains(t, err.Error(), "error walking destination directory")
	})
}

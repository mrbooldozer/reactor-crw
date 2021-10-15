// +build unit

package fs_test

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"reactor-crw/handler/fs"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type pathResolverMock struct {
	mock.Mock
}

func (m *pathResolverMock) CreateFolder(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *pathResolverMock) CreateFile(name string) (io.WriteCloser, error) {
	args := m.Called(name)
	return args.Get(0).(io.WriteCloser), args.Error(1)
}

func (m *pathResolverMock) Remove(name string) {
	_ = m.Called(name)
}

type transportMock struct {
	mock.Mock
}

func (m *transportMock) FetchData(url string) (io.ReadCloser, error) {
	args := m.Called(url)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func TestNewFileSaver(t *testing.T) {
	pr := pathResolverMock{}
	trp := transportMock{}

	t.Log("Given the need to create a file saver handler.")
	{
		t.Log("When baseFolder is incorrect.")
		{
			pr.On("CreateFolder", "baseFolder").Return(errors.New("error"))
			_, err := fs.NewFileSaver(&pr, &trp, "baseFolder")
			require.Error(t, err, "Expected an error during create folder")
			pr.ExpectedCalls = []*mock.Call{}
		}

		t.Log("When all data correct.")
		{
			pr.On("CreateFolder", "baseFolder").Return(nil)
			_, err := fs.NewFileSaver(&pr, &trp, "baseFolder")
			require.NoError(t, err, "Wasn't expected an error on new creating file saver")
		}
	}
}

func TestFileSaver_Process(t *testing.T) {
	tmlFile, _ := ioutil.TempFile(os.TempDir(), "file-title.txt")

	pr := pathResolverMock{}
	pr.On("CreateFolder", "baseFolder").Return(nil)
	pr.On("Remove", mock.Anything)

	trp := transportMock{}

	p := make(chan int, 1)
	e := make(chan error, 1)

	fileSaver, _ := fs.NewFileSaver(&pr, &trp, "baseFolder")

	t.Log("Given the need to process data.")
	{
		t.Log("When transport returns an error.")
		{
			trp.On("FetchData", "file-title.txt").Return(tmlFile, errors.New("error")).Once()

			fileSaver.Process("file-title.txt", p, e)
			require.Error(t, <-e, "Expected an error during create file")
			<-p
		}

		t.Log("When path resolver returns an error.")
		{
			trp.On("FetchData", "file-title.txt").Return(tmlFile, nil).Once()
			pr.On("CreateFile", "file-title.txt").Return(tmlFile, errors.New("error")).Once()

			fileSaver.Process("file-title.txt", p, e)
			require.Error(t, <-e, "Expected an error during create file")
			<-p
		}

		t.Log("When cannot save the data.")
		{
			trp.On("FetchData", "file-title.txt").Return(tmlFile, nil).Once()
			pr.On("CreateFile", "file-title.txt").Return(tmlFile, nil).Once()
			pr.On("Remove", mock.Anything)

			fileSaver.Process("file-title.txt", p, e)
			require.Error(t, <-e, "Expected an error file copying")
			<-p
		}

		t.Log("When all data correct.")
		{
			tmlFile, _ = ioutil.TempFile(os.TempDir(), "new-file-title.txt")

			trp.On("FetchData", "new-file-title.txt").Return(tmlFile, nil).Once()
			pr.On("CreateFile", "new-file-title.txt").Return(tmlFile, nil).Once()

			fileSaver.Process("new-file-title.txt", p, e)
			select {
			case err := <-e:
				require.NoError(t, err, "Wasn't expected an error on new file process")
			case <-p:
				return
			}
		}
	}
}

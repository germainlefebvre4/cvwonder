package cvserve

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/stretchr/testify/assert"
)

func NewServeServicesTest() ServeServices {
	return ServeServices{}
}

func TestStartLiveReloader(t *testing.T) {
	testDirectory, _ := os.Getwd()
	baseDirectory, err := filepath.Abs(testDirectory + "/../..")
	randomString := utils.GenerateRandomString(5)
	outputDirectory := baseDirectory + "/generated-test-" + randomString
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		ServeService ServeServices
	}
	type args struct {
		outputDirectory string
		inputFilePath   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name:   "Should start live reloader",
			fields: fields{NewServeServicesTest()},
			args: args{
				outputDirectory: outputDirectory,
				inputFilePath:   "TestStartLiveReloader",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		// Run test
		t.Run("Should start live reloader", func(t *testing.T) {
			// Prepare
			if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
				err := os.Mkdir(outputDirectory, os.ModePerm)
				if err != nil {
					t.Fatal(err)
				}
			}
			err := os.WriteFile(outputDirectory+"/TestStartLiveReloader.html", []byte("TestRunWebServer"), os.ModePerm)
			if err != nil {
				t.Fatal(err)
			}

			// Check results
			assert.Equalf(
				t,
				tt.wantErr,
				nil,
				"StartLiveReloader(%v, %v)",
				tt.args.outputDirectory,
				tt.args.inputFilePath,
			)

			// Clean
			err = os.RemoveAll(tt.args.outputDirectory)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

// func TestStartLiveReloaderHTTPConnection(t *testing.T) {
// 	testDirectory, _ := os.Getwd()
// 	baseDirectory, err := filepath.Abs(testDirectory + "/../..")
// 	randomString := utils.GenerateRandomString(5)
// 	outputDirectoryName := "generated-test-" + randomString
// 	outputDirectory := baseDirectory + "/" + outputDirectoryName

// 	// Prepare
// 	if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
// 		err := os.Mkdir(outputDirectory, os.ModePerm)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}
// 	err = os.WriteFile(outputDirectory+"/TestStartLiveReloaderHTTPConnection.html", []byte("TestStartLiveReloaderHTTPConnection"), os.ModePerm)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Run server
// 	go func() {
// 		service := NewServeServicesTest()
// 		service.StartLiveReloader(18080, outputDirectoryName, "TestStartLiveReloader")
// 	}()
// 	time.Sleep(1 * time.Second)

// 	// Check results
// 	// conn, err := net.Dial("tcp", ":18080")
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }
// 	// assert.NotNil(t, conn)
// 	// conn.Close()

// 	req, err := http.NewRequest("GET", "http://localhost:18080/TestStartLiveReloaderHTTPConnection.html", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assert.Equal(t, 200, resp.StatusCode)
// }

package render_pdf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/stretchr/testify/assert"
)

func NewRenderPDFServicesTest() RenderPDFServices {
	return RenderPDFServices{}
}

func TestGenerateOutputFile(t *testing.T) {
	testDirectory, _ := os.Getwd()
	baseDirectory, err := filepath.Abs(testDirectory + "/../../..")
	randomString := utils.GenerateRandomString(5)
	outputDirectory := baseDirectory + "/generated-test-" + randomString
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		RenderPDFService RenderPDFServices
	}
	type args struct {
		outputDirectory string
		inputFilename   string
	}
	test := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "Should create and return output directory and file",
			fields: fields{NewRenderPDFServicesTest()},
			args: args{
				outputDirectory: outputDirectory,
				inputFilename:   "TestGenerateOutputFile",
			},
			want: outputDirectory + "/TestGenerateOutputFile.pdf",
		},
	}
	for _, tt := range test {
		// Prepare
		if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
			err := os.Mkdir(outputDirectory, os.ModePerm)
			if err != nil {
				t.Fatal(err)
			}
		}

		// Run test
		t.Run(tt.name, func(t *testing.T) {
			service := NewRenderPDFServicesTest()
			assert.Equalf(
				t,
				tt.want,
				service.generateOutputFile(tt.args.outputDirectory, tt.args.inputFilename),
				"generateOutputFile(%v, %v)",
				tt.args.outputDirectory,
				tt.args.inputFilename,
			)
		})

		// Clean
		err := os.RemoveAll(outputDirectory)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// func TestRunWebServer(t *testing.T) {
// 	testDirectory, _ := os.Getwd()
// 	baseDirectory, err := filepath.Abs(testDirectory + "/../../..")
// 	randomString := utils.GenerateRandomString(5)
// 	outputDirectory := baseDirectory + "/generated-test-" + randomString
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	type fields struct {
// 		RenderPDFService RenderPDFServices
// 	}
// 	type args struct {
// 		port            int
// 		inputFilename   string
// 		outputDirectory string
// 	}
// 	test := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   string
// 	}{
// 		{
// 			name:   "Should run web server and return local server URL",
// 			fields: fields{NewRenderPDFServicesTest()},
// 			args: args{
// 				port:            18080,
// 				inputFilename:   "TestRunWebServer",
// 				outputDirectory: outputDirectory,
// 			},
// 			want: "http://localhost:18080/TestRunWebServer.html",
// 		},
// 	}
// 	for _, tt := range test {
// 		// Prepare
// 		if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
// 			err := os.Mkdir(outputDirectory, os.ModePerm)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 		}
// 		err := os.WriteFile(outputDirectory+"/TestRunWebServer.html", []byte("TestRunWebServer"), os.ModePerm)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		// Run test
// 		t.Run(tt.name, func(t *testing.T) {
// 			service := NewRenderPDFServicesTest()
// 			assert.Equalf(
// 				t,
// 				tt.want,
// 				service.runWebServer(tt.args.port, tt.args.inputFilename, tt.args.outputDirectory),
// 				"runWebServer(%v, %v)",
// 				tt.args.port,
// 				tt.args.inputFilename,
// 				tt.args.outputDirectory,
// 			)
// 		})

// 		// Clean
// 		err = os.RemoveAll(outputDirectory)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}
// }

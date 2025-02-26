package render_html

import (
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/germainlefebvre4/cvwonder/internal/fixtures"
	"github.com/germainlefebvre4/cvwonder/internal/model"
	theme_config "github.com/germainlefebvre4/cvwonder/internal/themes/config"
	utils "github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/stretchr/testify/assert"
)

func NewRenderHTMLServicesTest() RenderHTMLServices {
	return RenderHTMLServices{}
}

func TestRenderFormatHTML(t *testing.T) {
	testDirectory, _ := os.Getwd()
	baseDirectory, err := filepath.Abs(testDirectory + "/../../..")
	randomString := utils.GenerateRandomString(5)
	outputDirectory := baseDirectory + "/generated-test-" + randomString
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		RenderHTMLService RenderHTMLServices
	}
	type args struct {
		cv              model.CV
		baseDirectory   string
		outputDirectory string
		inputFilename   string
		themeName       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name:   "Should render HTML",
			fields: fields{NewRenderHTMLServicesTest()},
			args: args{
				cv:              fixtures.CvModelGood01,
				baseDirectory:   baseDirectory,
				outputDirectory: outputDirectory,
				inputFilename:   "cv",
				themeName:       "test",
			},
			wantErr: nil,
		},
	}
	// Create the test theme
	if _, err := os.Stat(baseDirectory + "/themes"); os.IsNotExist(err) {
		err := os.Mkdir(baseDirectory+"/themes", os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Stat(baseDirectory + "/themes/test"); os.IsNotExist(err) {
		err := os.Mkdir(baseDirectory+"/themes/test", os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, tt := range tests {
		// Prepare
		// Create the theme template file
		if _, err := os.Stat(baseDirectory + "/themes/test"); os.IsNotExist(err) {
			err := os.Mkdir(baseDirectory+"/themes/test", os.ModePerm)
			if err != nil {
				t.Fatal(err)
			}
		}
		err = os.WriteFile(baseDirectory+"/themes/test/index.html", fixtures.CvHtmlTemplate01, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

		// Run test
		t.Run(tt.name, func(t *testing.T) {
			service := NewRenderHTMLServicesTest()
			assert.Equalf(
				t,
				tt.wantErr,
				service.RenderFormatHTML(tt.args.cv, tt.args.baseDirectory, tt.args.outputDirectory, tt.args.inputFilename, tt.args.themeName),
				"RenderFormatHTML(%v, %v, %v, %v)",
				tt.args.cv,
				tt.args.outputDirectory,
				tt.args.inputFilename,
				tt.args.themeName,
			)
		})

		// Clean
		err = os.RemoveAll(tt.args.outputDirectory)
		if err != nil {
			t.Fatal(err)
		}
	}
	err = os.RemoveAll(baseDirectory + "/themes/test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTemplateFunctions(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		args args
		want template.FuncMap
	}{
		{
			name: "Should return template functions",
			args: args{},
			want: template.FuncMap{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.IsType(
				t,
				tt.want,
				getTemplateFunctions(),
				"getTemplateFunctions()",
			)
		})
	}
}

func TestTemplateFunctions(t *testing.T) {
	testDirectory, _ := os.Getwd()
	baseDirectory, err := filepath.Abs(testDirectory + "/../../..")
	randomString := utils.GenerateRandomString(5)
	outputDirectory := baseDirectory + "/generated-test-" + randomString
	themeDirectory := baseDirectory + "/themes/test-" + randomString
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		themeDirectory      string
		themeIndexFilePath  string
		outputFilePath      string
		cv                  model.CV
		templateFileContent []byte
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected []byte
	}{
		{
			name: "Should increase a number",
			args: args{
				themeDirectory:      themeDirectory,
				themeIndexFilePath:  themeDirectory + "/index.html",
				outputFilePath:      outputDirectory + "/" + utils.GenerateRandomString(5) + ".html",
				cv:                  fixtures.CvModelGood01,
				templateFileContent: fixtures.CvHtmlTemplateFunctionsTemplateInc01,
			},
			wantErr:  false,
			expected: fixtures.CvHtmlTemplateFunctionsGeneratedInc01,
		},
		{
			name: "Should decrease a number",
			args: args{
				themeDirectory:      themeDirectory,
				themeIndexFilePath:  themeDirectory + "/index.html",
				outputFilePath:      outputDirectory + "/" + utils.GenerateRandomString(5) + ".html",
				cv:                  fixtures.CvModelGood01,
				templateFileContent: fixtures.CvHtmlTemplateFunctionsTemplateDec01,
			},
			wantErr:  false,
			expected: fixtures.CvHtmlTemplateFunctionsGeneratedDec01,
		},
		{
			name: "Should list items",
			args: args{
				themeDirectory:      themeDirectory,
				themeIndexFilePath:  themeDirectory + "/index.html",
				outputFilePath:      outputDirectory + "/" + utils.GenerateRandomString(5) + ".html",
				cv:                  fixtures.CvModelGood01,
				templateFileContent: fixtures.CvHtmlTemplateFunctionsTemplateList01,
			},
			wantErr:  false,
			expected: fixtures.CvHtmlTemplateFunctionsGeneratedList01,
		},
		{
			name: "Should join items",
			args: args{
				themeDirectory:      themeDirectory,
				themeIndexFilePath:  themeDirectory + "/index.html",
				outputFilePath:      outputDirectory + "/" + utils.GenerateRandomString(5) + ".html",
				cv:                  fixtures.CvModelGood01,
				templateFileContent: fixtures.CvHtmlTemplateFunctionsTemplateJoin01,
			},
			wantErr:  false,
			expected: fixtures.CvHtmlTemplateFunctionsGeneratedJoin01,
		},
		{
			name: "Should split items",
			args: args{
				themeDirectory:      themeDirectory,
				themeIndexFilePath:  themeDirectory + "/index.html",
				outputFilePath:      outputDirectory + "/" + utils.GenerateRandomString(5) + ".html",
				cv:                  fixtures.CvModelGood01,
				templateFileContent: fixtures.CvHtmlTemplateFunctionsTemplateSplit01,
			},
			wantErr:  false,
			expected: fixtures.CvHtmlTemplateFunctionsGeneratedSplit01,
		},
		{
			name: "Should trim items",
			args: args{
				themeDirectory:      themeDirectory,
				themeIndexFilePath:  themeDirectory + "/index.html",
				outputFilePath:      outputDirectory + "/" + utils.GenerateRandomString(5) + ".html",
				cv:                  fixtures.CvModelGood01,
				templateFileContent: fixtures.CvHtmlTemplateFunctionsTemplateTrim01,
			},
			wantErr:  false,
			expected: fixtures.CvHtmlTemplateFunctionsGeneratedTrim01,
		},
		{
			name: "Should lower items",
			args: args{
				themeDirectory:      themeDirectory,
				themeIndexFilePath:  themeDirectory + "/index.html",
				outputFilePath:      outputDirectory + "/" + utils.GenerateRandomString(5) + ".html",
				cv:                  fixtures.CvModelGood01,
				templateFileContent: fixtures.CvHtmlTemplateFunctionsTemplateLower01,
			},
			wantErr:  false,
			expected: fixtures.CvHtmlTemplateFunctionsGeneratedLower01,
		},
		{
			name: "Should upper items",
			args: args{
				themeDirectory:      themeDirectory,
				themeIndexFilePath:  themeDirectory + "/index.html",
				outputFilePath:      outputDirectory + "/" + utils.GenerateRandomString(5) + ".html",
				cv:                  fixtures.CvModelGood01,
				templateFileContent: fixtures.CvHtmlTemplateFunctionsTemplateUpper01,
			},
			wantErr:  false,
			expected: fixtures.CvHtmlTemplateFunctionsGeneratedUpper01,
		},
		{
			name: "Should replace items",
			args: args{
				themeDirectory:      themeDirectory,
				themeIndexFilePath:  themeDirectory + "/index.html",
				outputFilePath:      outputDirectory + "/" + utils.GenerateRandomString(5) + ".html",
				cv:                  fixtures.CvModelGood01,
				templateFileContent: fixtures.CvHtmlTemplateFunctionsTemplateReplace01,
			},
			wantErr:  false,
			expected: fixtures.CvHtmlTemplateFunctionsGeneratedReplace01,
		},
	}
	// Prepare
	// Create the generated directory
	if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
		err := os.Mkdir(outputDirectory, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	}
	// Create the test theme
	if _, err := os.Stat(baseDirectory + "/themes"); os.IsNotExist(err) {
		err := os.Mkdir(baseDirectory+"/themes", os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, tt := range tests {
		// Prepare the theme
		if _, err := os.Stat(tt.args.themeDirectory); os.IsNotExist(err) {
			err := os.Mkdir(tt.args.themeDirectory, os.ModePerm)
			if err != nil {
				t.Fatal(err)
			}
		}
		err = os.WriteFile(tt.args.themeDirectory+"/theme.yaml", theme_config.ThemeConfigYamlGood01, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

		// Run test
		t.Run(tt.name, func(t *testing.T) {
			// Create output file
			var outputFile *os.File
			outputFile, err = os.Create(tt.args.outputFilePath)
			if err != nil {
			}
			// Create theme template file
			err = os.WriteFile(tt.args.themeIndexFilePath, tt.args.templateFileContent, os.ModePerm)
			if err != nil {
				t.Fatal(err)
			}
			// Generate html template
			templateFunctions := getTemplateFunctions()
			tmpl, err := template.New("index.html").Funcs(templateFunctions).ParseFiles(tt.args.themeIndexFilePath)
			if err != nil {
				t.Fatal(err)
			}
			err = tmpl.ExecuteTemplate(outputFile, "index.html", tt.args.cv)
			if err != nil {
				t.Fatal(err)
			}
			// Read output file
			var outputFileContent []byte
			outputFileContent, err = os.ReadFile(tt.args.outputFilePath)
			if err != nil {
			}
			// Check
			res := assert.Equal(t, tt.expected, outputFileContent)
			// Clean
			if res {
				err = os.Remove(tt.args.outputFilePath)
				if err != nil {
					t.Fatal(err)
				}
			}
		})

	}
	// Clean
	err = os.RemoveAll(outputDirectory)
	if err != nil {
		t.Fatal(err)
	}
	err = os.RemoveAll(themeDirectory)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateTemplateFile(t *testing.T) {
	testDirectory, _ := os.Getwd()
	baseDirectory, err := filepath.Abs(testDirectory + "/../../..")
	randomString := utils.GenerateRandomString(5)
	outputDirectory := baseDirectory + "/generated-test-" + randomString
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		RenderHTMLService RenderHTMLServices
	}
	type args struct {
		themeDirectory    string
		outputDirectory   string
		outputFilePath    string
		outputTmpFilePath string
		cv                model.CV
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Should generate template file",
			fields: fields{NewRenderHTMLServicesTest()},
			args: args{
				themeDirectory:    baseDirectory + "/themes/test",
				outputDirectory:   outputDirectory,
				outputFilePath:    outputDirectory + "/cv.html",
				outputTmpFilePath: outputDirectory + "/cv.html.tmp",
				cv:                fixtures.CvModelGood01,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		// Prepare
		// TODO: Add theme files
		// Create the test theme
		if _, err := os.Stat(baseDirectory + "/themes"); os.IsNotExist(err) {
			err := os.Mkdir(baseDirectory+"/themes", os.ModePerm)
			if err != nil {
				t.Fatal(err)
			}
		}
		if _, err := os.Stat(baseDirectory + "/themes/test"); os.IsNotExist(err) {
			err := os.Mkdir(baseDirectory+"/themes/test", os.ModePerm)
			if err != nil {
				t.Fatal(err)
			}
		}
		err = os.WriteFile(baseDirectory+"/themes/test/index.html", fixtures.CvHtmlTemplate01, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(baseDirectory+"/themes/test/theme.yaml", theme_config.ThemeConfigYamlGood01, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

		// Run test
		t.Run(tt.name, func(t *testing.T) {
			service := NewRenderHTMLServicesTest()
			service.generateTemplateFile(tt.args.themeDirectory, tt.args.outputDirectory, tt.args.outputFilePath, tt.args.outputTmpFilePath, tt.args.cv)
			assert.DirExists(t, tt.args.outputDirectory)
			assert.FileExists(t, tt.args.outputFilePath)
			assert.FileExists(t, tt.args.outputTmpFilePath)
		})

		// Clean
		err := os.RemoveAll(tt.args.outputDirectory)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCopyTemplateFileContent(t *testing.T) {
	testDirectory, _ := os.Getwd()
	baseDirectory, err := filepath.Abs(testDirectory + "/../../..")
	randomString := utils.GenerateRandomString(5)
	outputDirectory := baseDirectory + "/generated-test-" + randomString
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should copy template file content",
			args: args{outputDirectory + "/TestCopyTemplateFileContent.test.tmp", outputDirectory + "/TestCopyTemplateFileContent.test"},
		},
	}
	for _, tt := range tests {
		// Prepare
		if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
			err := os.Mkdir(outputDirectory, os.ModePerm)
			if err != nil {
				t.Fatal(err)
			}
		}

		// Run test
		t.Run(tt.name, func(t *testing.T) {
			// Prepare
			err := os.WriteFile(tt.args.src, []byte("test"), os.ModePerm)
			if err != nil {
				t.Fatal(err)
			}
			f1, err1 := os.ReadFile(tt.args.src)
			if err1 != nil {
				t.Fatal(err1)
			}

			copyTemplateFileContent(tt.args.src, tt.args.dst)
			assert.FileExists(t, tt.args.dst)

			f2, err2 := os.ReadFile(tt.args.dst)
			if err2 != nil {
				t.Fatal(err2)
			}
			assert.Equal(t, f1, f2)

			// Clean
			err = os.Remove(tt.args.dst)
			if err != nil {
				t.Fatal(err)
			}
		})

		// Clean
		err = os.RemoveAll(outputDirectory)
		if err != nil {
			t.Fatal(err)
		}
	}
}

package cvbulk

import "github.com/germainlefebvre4/cvwonder/internal/model"

// BulkGenerateInterface defines the contract for bulk CV generation.
type BulkGenerateInterface interface {
// BulkGenerate processes a slice of InputFile concurrently using a worker
// pool of the given concurrency. Each file is generated under outputDir,
// mirroring the relative path from inputDir. Returns one BulkResult per file.
BulkGenerate(
inputFiles []model.InputFile,
inputDir string,
outputDir string,
themeName string,
format string,
validate bool,
concurrency int,
) []BulkResult
}

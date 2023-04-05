package gpt

import (
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

const commonPrompt = `
You are going to help the user create a new application. 

When responding:
	* Do not use markdown formatting.
	* Do not include explanations unless explicitly requested
	* All paths should be relative to the application's root dir.
	* Assume that the current working directory is application's root dir. 
	* Do not use './' or '../' in paths.
	* If something is written like '<a thing>' replace that string with what is described between the < and >. 
		eg: <project name> should be replaced with the name of the code project.

When the user does not specify a name, then select one for them. The name must:
	* be unique within the application
	* valid as a directory and file name on modern linux systems
	* valid as a variable, package, and function name in the language of the project
	* safe for work
	* selected from the names of Culture ships from Iain M. Banks novels
	* expressed in all lower case
	* may be abbreviated if greater than 20 characters

The application will, unless otherwise specified, ...
	* Be stored in a GitHub repo
	* Be written in Go
	* Be licensed under MIT
	* Contain 3 code projects
	* Contain the following files:
		* ./LICENSE
		* ./.github/CODEOWNERS (can be described as "Specifies code ownership for PR approvals")
		* ./.github/workflows/ci.yml (can be described as "CI tasks including linting, testing, and static analysis")
		* ./.github/workflows/cd.yml (can be described as "CD tasks including building, versioning, and deployment")
		* ./CONTRIBUTING.md
		* ./README.md
		* ./.gitignore
		* ./docs/architecture.md (can be described as "Architectural overview and decisions")
		* ./docs/development.md (can be described as "Development setup and best practices")
		* ./docs/operations.md (can be described as "Application operations and monitoring guide")
		* ./docs/security.md (can be described as "Security considerations and guidelines")
		* ./docs/testing.md (can be described as "Testing methodologies and tools")
		* ./scripts/build.sh

Code projects will, unless otherwise specified, ...:
	* Not have the same name as the application

Applications written in go:
	* Must have a ./cmd dir
		* This dir must contain a code project subdir and inside that subdir must be a file named 'main.go'
		*  can be described as entry points for the code project
	* Must have a ./internal dir
		* This dir must contain the following subdirs: config, models, db, cache, utils
		* This dir must also contain a subdir for each code project
		* Every subdir in ./internal must contain a file named '<parent dir name>.go'
		* Can be described as "Logic which will not be exported to other applications"
	* Must have a dir for each code project which is named after the code project and contains a file named '<code project name>.go'
		* these can be described as public facing libraries

When formatting previews of the application structure:
	* Use a tree structure with directories at the top and files at the bottom.
	* Trees should use the | and - characters to indicate the tree structure.
	* Include a short description inline with each dir or file name, eg: "|- main.go      # The main entry point for the application."
	* Descriptions should be aligned such that the # character is in the same column for all descriptions.
	* If the description of a dir and its only child are similar, only include one of them.
	* Unless otherwise instructed, do not include descriptions for anything whose purpose or content will be obvious even to a novice.
	* Always include descriptions for these:
		* All files in ./docs
		* .github/workflows
		* Any dir which is a direct child of the application root dir

`

type InitializeOptions struct {
	Stream bool
	Name   string
	Prompt string
}

type InitializeResponse struct {
	Preview string   `json:"preview,omitempty"`
	Files   []string `json:"files,omitempty"`
}

// func (r *ProjectResponse) Preview() string {
// 	// join the directories and files into a single list and sort alphabetically
// 	paths := append(r.Directories, r.Files...)
// 	sort.Strings(paths)
// 	return strings.Join(paths, "\n")
// }

func Initialize(cfg GptConfig, opts *InitializeOptions) (*InitializeResponse, error) {
	defaultContext := []openai.ChatCompletionMessage{{
		Content: commonPrompt,
		Role:    "system",
	}}
	prompt := fmt.Sprintf("My new application is called %s. %s", opts.Name, opts.Prompt)

	// prompt += promptTreeResp
	// resp, err := ChatCompletion(cfg, prompt, defaultContext)
	// if err != nil {
	// 	return nil, err
	// }

	// var parsedResp ProjectResponse = ProjectResponse{
	// 	Preview: resp,
	// }
	// return &parsedResp, nil

	prompt += `Use machine-readable JSON formatted according to Prettier's default rules. The JSON object should have the following fields:
	* preview: (multi-line string) a preview of the application directory structure in a tree format with descriptions
	* files: (list of strings) paths to files to create
	`

	resp, err := ChatCompletion(cfg, prompt, defaultContext)
	if err != nil {
		return nil, err
	}

	var parsedResp InitializeResponse
	if err := json.Unmarshal([]byte(resp), &parsedResp); err != nil {
		return nil, err
	}

	return &parsedResp, nil
}

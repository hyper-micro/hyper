package command

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hyper-micro/hyper/tools/config"
)

type InitCommandArgs struct {
	ProjectName string
	Mod         string
}

type InitCommand struct {
	args InitCommandArgs
}

func NewInitCommand(args InitCommandArgs) *InitCommand {
	return &InitCommand{
		args: args,
	}
}

func (cmd *InitCommand) Replace() error {
	if cmd.args.ProjectName == "" {
		return errors.New("project name empty")
	}

	allowedRegex := regexp.MustCompile(`^[0-9a-zA-Z\-_]+$`)
	if !allowedRegex.MatchString(cmd.args.ProjectName) {
		return errors.New("illegal characters")
	}

	nPath, err := os.Getwd()
	if err != nil {
		return err
	}
	projectPath := fmt.Sprintf("%s/%s", nPath, cmd.args.ProjectName)

	pathInfo, err := os.Stat(projectPath)
	if err != nil {
		return err
	}
	if !pathInfo.IsDir() {
		return fmt.Errorf("project name not exists")
	}
	if err = cmd.replaceKeywords(projectPath); err != nil {
		return err
	}

	return nil
}

func (cmd *InitCommand) Run() (err error) {
	projectPath, err := cmd.localPath()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			_ = os.RemoveAll(projectPath)
		}
	}()

	if err = cmd.fetchProjectTemplate(projectPath); err != nil {
		return
	}

	if err = cmd.cleanUpGit(projectPath); err != nil {
		return
	}

	if err = cmd.replaceKeywords(projectPath); err != nil {
		return
	}

	return
}

func (cmd *InitCommand) localPath() (string, error) {
	if cmd.args.ProjectName == "" {
		return "", errors.New("project name empty")
	}

	allowedRegex := regexp.MustCompile(`^[0-9a-zA-Z\-_]+$`)
	if !allowedRegex.MatchString(cmd.args.ProjectName) {
		return "", errors.New("illegal characters")
	}

	nPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	projectPath := fmt.Sprintf("%s/%s", nPath, cmd.args.ProjectName)

	pathInfo, err := os.Stat(projectPath)
	if err == nil {
		if pathInfo.IsDir() {
			return "", fmt.Errorf("the folder '%s' already exists", projectPath)
		}
		return "", fmt.Errorf("the file '%s' already exists", projectPath)
	}
	if !os.IsNotExist(err) {
		return "", err
	}

	return projectPath, nil
}

func (cmd *InitCommand) fetchProjectTemplate(projectPath string) (err error) {

	repo, err := git.PlainClone(projectPath, false, &git.CloneOptions{
		URL:           config.HyperProjectTemplateGithubRepo,
		ReferenceName: plumbing.NewBranchReferenceName(config.HyperProjectTemplateGithubBranch),
		Progress:      io.Discard,
	})
	if err != nil {
		return
	}

	tagRefs, err := repo.Tags()
	if err != nil {
		return
	}

	var latestTag = plumbing.NewBranchReferenceName(config.HyperProjectTemplateGithubBranch)
	err = tagRefs.ForEach(func(t *plumbing.Reference) error {
		latestTag = t.Name()
		return nil
	})
	if err != nil {
		return
	}

	workTree, err := repo.Worktree()
	if err != nil {
		return
	}
	err = workTree.Checkout(&git.CheckoutOptions{
		Branch: latestTag,
	})
	if err != nil {
		return
	}

	return nil
}

func (cmd *InitCommand) cleanUpGit(projectPath string) error {
	gitPath := fmt.Sprintf("%s/%s", projectPath, ".git")
	return os.RemoveAll(gitPath)
}

func (cmd *InitCommand) replaceKeywords(projectPath string) error {
	var src2replace = map[string]string{
		config.HyperProjectTemplateGoModName: cmd.args.Mod,
		"projectTemplateService":             cmd.projectCamelName(),
		"ProjectTemplateService":             cmd.projectCamelNameUpper(),
	}

	return filepath.Walk(projectPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fileData := string(data)
		for src, replace := range src2replace {
			fileData = strings.ReplaceAll(fileData, src, replace)
		}

		if err := os.WriteFile(path, []byte(fileData), info.Mode()); err != nil {
			return err
		}

		return nil
	})
}

func (cmd *InitCommand) projectCamelName() string {
	projectName := strings.ToLower(cmd.args.ProjectName)
	re := regexp.MustCompile(`[-_]([a-zA-Z0-9])`)
	resultBytes := re.ReplaceAllFunc([]byte(projectName), func(bytes []byte) []byte {
		return []byte(strings.ToUpper(string(bytes)))
	})
	result := string(resultBytes)
	result = strings.ReplaceAll(result, "-", "")
	result = strings.ReplaceAll(result, "_", "")
	return result
}

func (cmd *InitCommand) projectCamelNameUpper() string {
	projectName := cmd.projectCamelName()
	if len(projectName) == 0 {
		return ""
	}
	return strings.ToUpper(projectName[:1]) + projectName[1:]
}

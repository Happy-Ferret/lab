package git

import (
	// "bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lighttiger2505/lab/cmd"
)

type Editor struct {
	Program    string
	Topic      string
	File       string
	Message    string
	CS         string
	openEditor func(program, file string) error
}

func NewEditor(filePrefix, topic, message string, editFunc func(program, file string) error) (editor *Editor, err error) {
	messageFile, err := getMessageFile(filePrefix)
	if err != nil {
		return
	}

	program, err := GitEditor()
	if err != nil {
		return
	}

	cs := CommentChar()

	if editFunc == nil {
		editFunc = openTextEditor
	}

	editor = &Editor{
		Program:    program,
		Topic:      topic,
		File:       messageFile,
		Message:    message,
		CS:         cs,
		openEditor: editFunc,
	}

	return
}

func (e *Editor) DeleteFile() error {
	return os.Remove(e.File)
}

func (e *Editor) EditTitleAndDescription() (title, body string, err error) {
	content, err := e.openAndEdit()
	if err != nil {
		return
	}

	content = bytes.TrimSpace(content)
	reader := bytes.NewReader(content)
	title, body, err = readTitleAndBody(reader, e.CS)

	if err != nil || title == "" {
		defer e.DeleteFile()
	}

	return
}

func (e *Editor) openAndEdit() (content []byte, err error) {
	err = e.writeContent()
	if err != nil {
		return
	}

	err = e.openEditor(e.Program, e.File)
	if err != nil {
		err = fmt.Errorf("error using text editor for %s message", e.Topic)
		defer e.DeleteFile()
		return
	}

	content, err = e.readContent()

	return
}

func (e *Editor) writeContent() (err error) {
	// only write message if file doesn't exist
	if !e.isFileExist() && e.Message != "" {
		err = ioutil.WriteFile(e.File, []byte(e.Message), 0644)
		if err != nil {
			return
		}
	}

	return
}

func (e *Editor) isFileExist() bool {
	_, err := os.Stat(e.File)
	return err == nil || !os.IsNotExist(err)
}

func (e *Editor) readContent() (content []byte, err error) {
	return ioutil.ReadFile(e.File)
}

func openTextEditor(program, file string) error {
	editCmd := cmd.NewBasicCmd(program)
	r := regexp.MustCompile("[mg]?vi[m]$")
	if r.MatchString(program) {
		editCmd.WithArg("--cmd")
		editCmd.WithArg("set ft=markdown tw=0 wrap lbr")
	}
	editCmd.WithArg(file)
	// Reattach stdin to the console before opening the editor
	setConsole(editCmd)

	return editCmd.Spawn()
}

func sweepMarkdownComment(text string) string {
	r := regexp.MustCompile("<!--[\\s\\S]*?-->[\\n]*")
	return r.ReplaceAllString(text, "")
}

func parceTitleAndBody(text string) (title, body string) {
	r := regexp.MustCompile("\\S")
	var titleParts, bodyParts []string
	for _, line := range strings.Split(text, "\n") {
		if len(bodyParts) == 0 && r.MatchString(line) {
			titleParts = append(titleParts, line)
		} else {
			bodyParts = append(bodyParts, line)
		}
	}

	title = strings.TrimSpace(strings.Join(titleParts, " "))
	body = strings.TrimSpace(strings.Join(bodyParts, "\n"))
	return
}

func readTitleAndBody(reader io.Reader, cs string) (title, body string, err error) {
	// Reading message file
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	text := string(b)

	// Sweep markdown text
	sweepText := sweepMarkdownComment(text)
	// Parce title and body from message file
	title, body = parceTitleAndBody(sweepText)
	return
}

func getMessageFile(about string) (string, error) {
	gitDir, err := GitDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(gitDir, fmt.Sprintf("%s_EDITMSG", about)), nil
}

func setConsole(cmd *cmd.BasicCmd) {
	stdin, err := os.OpenFile("/dev/tty", os.O_RDONLY, 0660)
	if err == nil {
		cmd.Stdin = stdin
	}
}

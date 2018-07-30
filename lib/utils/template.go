package utils

import (
	"bufio"
	"bytes"
	"text/template"
)

func ExecuteTemplateToBuffer(templatePath string, params interface{}) (*bytes.Buffer, error) {
	template, err := template.ParseFiles(templatePath)

	if err != nil {
		return nil, err
	}

	fileBuffer := new(bytes.Buffer)
	fileWriter := bufio.NewWriter(fileBuffer)

	if err := template.Execute(fileWriter, params); err != nil {
		return nil, err
	}

	if err := fileWriter.Flush(); err != nil {
		return nil, err
	}

	return fileBuffer, nil
}

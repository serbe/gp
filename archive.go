package main

import (
	"archive/zip"
	"errors"
	"io"
	"os"
)

func decompress(archive string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	for _, file := range reader.File {
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()
		targetFile, err := os.OpenFile(file.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()
		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}

func compress(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()
	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return errors.New("Source is Dir")
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Method = zip.Deflate
	writer, err := archive.CreateHeader(header)
	if err != nil {
		return err
	}
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(writer, file)
	return err
}

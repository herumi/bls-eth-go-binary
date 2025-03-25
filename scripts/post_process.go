package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Run the `go mod vendor` command
	cmd := exec.Command("go", "mod", "vendor")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running `go mod vendor`: %v\n", err)
		os.Exit(1)
	}

	// Copy non-Go files back into the vendor directory
	directories := []string{"bls/include", "bls/include/bls", "bls/include/mcl", "bls/lib", "bls/tests"}
	for _, dir := range directories {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				destPath := filepath.Join("vendor", path)
				err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
				if err != nil {
					return err
				}
				err = copyFile(path, destPath)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error copying files from %s: %v\n", dir, err)
			os.Exit(1)
		}
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Close()
}

package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Run(ctx context.Context, backupDir string, serverAddr string) error {
	if err := createDirs(backupDir, serverAddr); err != nil {
		return err
	}

	t := time.NewTicker(time.Hour)

	for ; true; <-t.C {
		serverBackups, err := getData(serverAddr + "/files/")
		if err != nil {
			return err
		}
		clientBackups, err := getCurrentBackups(backupDir)
		if err != nil {
			return err
		}
		backupsToDownload := diff(serverBackups, clientBackups)
		backupsToDelete := diff(clientBackups, serverBackups)

		for _, backup := range backupsToDownload {
			if err = downloadFile(backupDir+backup, serverAddr+"/backups"+backup); err != nil {
				return err
			}
		}
		if err := deleteFiles(backupsToDelete, backupDir); err != nil {
			return err
		}
	}

	return nil
}

func createDirs(backupDir string, serverAddr string) error {
	dirs, err := getData(serverAddr + "/dirs/")
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		_, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("mkdir -p %s%s", backupDir, dir)).Output()
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadFile(filepath string, url string) (err error) {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func getData(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []string{}, fmt.Errorf("bad status: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	slice := strings.Split(string(bodyBytes), "\n")
	if len(slice) > 0 {
		slice = slice[:len(slice)-1]
	}

	return slice, nil
}

func getCurrentBackups(dir string) ([]string, error) {
	cmd := fmt.Sprintf(
		"(cd %s && find . -type f | cut -c 2- | sort)",
		dir,
	)

	stdout, err := exec.Command("/bin/bash", "-c", cmd).Output()
	if err != nil {
		return []string{}, err
	}
	slice := strings.Split(string(stdout), "\n")

	return slice, nil
}

func diff(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func deleteFiles(files []string, basedir string) error {
	filesPath := make([]string, 0, len(files))
	for _, file := range files {
		if file != "" {
			filesPath = append(filesPath, basedir+file)
		}
	}
	if len(filesPath) == 0 {
		return nil
	}

	cmd := fmt.Sprintf("rm %s", strings.Join(filesPath, " "))
	_, err := exec.Command("/bin/bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println(cmd)
		return err
	}

	return nil
}

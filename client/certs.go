package client

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func certPath(pubKeyPath string) string {
	base := filepath.Base(pubKeyPath)
	dir := filepath.Dir(pubKeyPath)
	prefix := strings.Split(base, ".")[0]
	return path.Join(dir, prefix+"-cert.pub")
}

// Returns full path for all the public keys in the directory given
func listPubKeysInDir(dir string) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to enumerate files from %s", dir)
	}
	files := []string{}
	for _, fileInfo := range fileInfos {
		if strings.HasSuffix(fileInfo.Name(), ".pub") && !strings.Contains(fileInfo.Name(), "cert") {
			files = append(files, path.Join(dir, fileInfo.Name()))
		}
	}
	return files, nil
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func updateUsersCertAuthority(filePath string, trustedUserCAs [][]byte) error {
	f, err := os.Create(filePath)
	if err != nil {
		return errors.Wrapf(err, "failed to create file %s", filePath)
	}
	defer f.Close()
	for _, b := range trustedUserCAs {
		f.Write(b)
		if b[len(b)-1] != '\n' {
			f.WriteString("\n")
		}
	}
	return nil
}

func updateKnownHostsCertAuthority(filePath string, trustedHostCAs [][]byte) error {
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		return errors.Wrapf(err, "Failed to read %s", filePath)
	}

	re := regexp.MustCompile(`(?ms:^#accord-trusted-hosts-start(.*)#accord-trusted-hosts-end)`)
	newlines := strings.Split(re.ReplaceAllString(string(input), ""), "\n")
	// the last synt
	newlines = deleteEmpty(newlines)
	newlines = append(newlines, "#accord-trusted-hosts-start")
	for _, b := range trustedHostCAs {

		if b[len(b)-1] == '\n' {
			newlines = append(newlines, "@cert-authority * "+string(b[:len(b)-1]))
		} else {
			newlines = append(newlines, "@cert-authority * "+string(b))
		}

	}
	newlines = append(newlines, "#accord-trusted-hosts-end", "\n")
	backupFile := filePath + ".bak"
	log.Println("Copied old file to " + backupFile)
	err = os.Rename(filePath, backupFile)
	if err != nil {
		return errors.Wrapf(err, "Failed to rename file to %s", backupFile)
	}
	err = ioutil.WriteFile(filePath, []byte(strings.Join(newlines, "\n")), 0644)
	if err != nil {
		return errors.Wrapf(err, "Failed to write to %s", filePath)
	}
	return nil
}

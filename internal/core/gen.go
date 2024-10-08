package core

import (
	"cmp"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/internal/pkged"
	"github.com/artalkjs/artalk/v2/internal/utils"
)

func Gen(genType string, specificPath string, overwrite bool) {
	locale := cmp.Or(os.Getenv("ATK_LOCALE"), "en")

	// check if generate config file
	isGenConf := false
	if genType == "config" || genType == "conf" || genType == "artalk.yml" {
		isGenConf = true
		genType = "artalk.yml"
	}

	// get generation content
	var fileStr string
	if isGenConf {
		fileStr = config.Template(locale)
		// gen random `app_key`
		appKey := utils.RandomStringWithAlphabet(16, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*1234567890")
		fileStr = strings.Replace(fileStr, `app_key: ""`, fmt.Sprintf(`app_key: "%s"`, appKey), 1)
	} else {
		fileStr = getEmbedFile(genType)
	}

	genFullPath := filepath.Base(genType) // generate file in work dir
	if specificPath != "" {
		genFullPath = specificPath
	}

	absPath, err := filepath.Abs(genFullPath)
	if err != nil {
		log.Fatal(err)
	}
	if s, err := os.Stat(absPath); err == nil && s.IsDir() {
		absPath = filepath.Join(absPath, filepath.Base(genType))
	}

	if utils.CheckFileExist(absPath) && !overwrite {
		log.Fatal(i18n.T("{{name}} already exists", map[string]interface{}{"name": i18n.T("File")}) + ": " + absPath)
	}

	if err := utils.EnsureDir(filepath.Dir(absPath)); err != nil {
		log.Fatal("Failed to create target directory: ", err)
	}

	dst, err := os.Create(absPath)
	if err != nil {
		log.Fatal("Failed to create target file: ", err)
	}
	defer dst.Close()

	if _, err = dst.Write([]byte(fileStr)); err != nil {
		log.Fatal("Failed to write target file: ", err)
	}

	log.Info("File Generated: " + absPath)
}

func getEmbedFile(filename string) string {
	file, err := pkged.FS().Open(strings.TrimPrefix(filename, "/"))
	if err != nil {
		log.Fatal("Invalid built-in resource `"+filename+"`: ", err)
	}

	buf, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Read built-in resources `"+filename+"` error: ", err)
	}

	return string(buf)
}

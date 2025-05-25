package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

var (
	LatestVersionFile = "/wire/otas/latest"
	DoNotAcceptFile   = "/wire/otas/dnar"
	FullPath          = "/wire/otas/full"
)

var TargetMap = []string{"dev", "oskr", "whiskey", "orange", "dvt3", "dvt2"}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func shouldNotAccept() bool {
	return fileExists(DoNotAcceptFile)
}

func normVer(v string) string {
	var r []rune
	for _, c := range v {
		if !unicode.IsLetter(c) {
			r = append(r, c)
		}
	}
	return strings.TrimSpace(string(r))
}

func getLatestVersion() string {
	b, err := os.ReadFile(LatestVersionFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func targetToString(t string) string {
	i, err := strconv.Atoi(t)
	if err != nil || i < 0 || i >= len(TargetMap) {
		return ""
	}
	return TargetMap[i]
}

func getOTA(version, target string) ([]byte, error) {
	if shouldNotAccept() {
		return nil, errors.New("server busy")
	}
	version = normVer(version)
	latest := getLatestVersion()
	if version == latest {
		return nil, errors.New("already on latest")
	}
	p := filepath.Join(FullPath, target, latest+".ota")
	if fileExists(p) {
		return os.ReadFile(p)
	}
	return nil, errors.New("ota not found")
}

func otaHandler(w http.ResponseWriter, r *http.Request) {
	if shouldNotAccept() && r.FormValue("isBuildServer") != "true" {
		http.Error(w, "busy uploading", 500)
		return
	}

	path := r.URL.Path
	if strings.HasPrefix(path, "/vic/full") || strings.HasPrefix(path, "/vic/diff") {
		ver := normVer(r.FormValue("victorversion"))
		t := r.FormValue("victortarget")
		fmt.Println(ver, t)
		if t == "" {
			http.Error(w, "no target", 404)
			return
		}
		t = targetToString(t)
		data, err := getOTA(ver, t)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.Write(data)
		return
	}

	if strings.HasPrefix(path, "/vic/latest/") {
		target := strings.Split(strings.Split(path, "/vic/latest/")[1], ".ota")[0]
		latestOta, err := os.ReadFile(filepath.Join(FullPath, target, getLatestVersion()+".ota"))
		if err != nil {
			http.Error(w, "target not found", 404)
			return
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(latestOta)))
		w.Write(latestOta)
		return
	}

	http.NotFound(w, r)
}

func main() {
	http.HandleFunc("/vic/", otaHandler)
	http.Handle("/", http.FileServer(http.Dir(FullPath+"/")))
	fmt.Println("listening on :5901")
	http.ListenAndServe(":5901", nil)
}

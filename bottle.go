package bottle

import (
	"bytes"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/yuin/cidre"
	"html/template"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

func isDir(assetdir func(string) ([]string, error), path string) bool {
	_, err := assetdir(path)
	return err == nil
}

func walkAssetDir(assetdir func(string) ([]string, error), path string, fn func(string, string) error) {
	dirstack := []string{strings.TrimRight(path, "/")}
	curdir := ""
	for len(dirstack) > 0 {
		curdir, dirstack = dirstack[len(dirstack)-1], dirstack[:len(dirstack)-1]
		files, _ := assetdir(curdir)
		for _, filename := range files {
			path := curdir + "/" + filename
			err := fn(curdir, filename)
			if isDir(assetdir, path) && err != filepath.SkipDir {
				dirstack = append(dirstack, path)
			}
		}
	}
}

type HtmlTemplateRenderer struct {
	*cidre.HtmlTemplateRenderer
	Asset    func(string) ([]byte, error)
	AssetDir func(string) ([]string, error)
}

func NewHtmlTemplateRenderer(r cidre.Renderer, asset func(string) ([]byte, error), assetdir func(string) ([]string, error)) *HtmlTemplateRenderer {
	rndr := &HtmlTemplateRenderer{
		HtmlTemplateRenderer: r.(*cidre.HtmlTemplateRenderer),
		Asset:                asset,
		AssetDir:             assetdir,
	}
	return rndr
}

func (rndr *HtmlTemplateRenderer) Compile() {
	if len(rndr.Config.TemplateDirectory) == 0 {
		return
	}

	funcMap := template.FuncMap{
		"include": func(name string, param interface{}) template.HTML {
			var buf bytes.Buffer
			rndr.RenderTemplateFile(&buf, name, param)
			return template.HTML(buf.String())
		},
		"raw": func(h string) template.HTML { return template.HTML(h) },
		// parse time dummy function
		"yield": func() template.HTML { return template.HTML("") },
	}

	extendsReg := regexp.MustCompile(regexp.QuoteMeta(rndr.Config.LeftDelim) + `/\*\s*extends\s*([^\s]+)\s*\*/` + regexp.QuoteMeta(rndr.Config.RightDelim))
	walkAssetDir(rndr.AssetDir, rndr.Config.TemplateDirectory, func(dir, filename string) error {
		path := dir + "/" + filename
		if isDir(rndr.AssetDir, path) {
			return nil
		}
		if !strings.HasSuffix(filename, ".tpl") {
			return nil
		}
		tplname := filename[0 : len(filename)-len(".tpl")]
		bts, err1 := rndr.Asset(path)
		if err1 != nil {
			panic(err1)
		}
		matches := extendsReg.FindAllSubmatch(bts, -1)
		if len(matches) > 0 {
			rndr.SetLayout(tplname, string(matches[0][1]))
		}
		tplobj, err2 := template.New("").Delims(rndr.Config.LeftDelim, rndr.Config.RightDelim).Funcs(rndr.Config.FuncMap).Funcs(funcMap).Parse(string(bts))
		if err2 != nil {
			panic(err2)
		}
		rndr.SetTemplate(tplname, tplobj)
		return nil
	})
}

func Static(mt *cidre.MountPoint, n, p, local string, asset func(string) ([]byte, error), assetdir func(string) ([]string, error), middlewares ...interface{}) {
	path := strings.Trim(p, "/")
	server := http.StripPrefix(mt.Path+path, http.FileServer(&assetfs.AssetFS{Asset: asset, AssetDir: assetdir, Prefix: local}))
	mt.Route(n, path+"/(?P<path>.*)", "GET", true, server.ServeHTTP, middlewares...)
}


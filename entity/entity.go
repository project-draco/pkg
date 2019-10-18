package entity

import (
	"regexp"
	"strings"

	"github.com/project-draco/naming"
)

// Entity represents a source-code entity,
// such as class, interface, method, or field.
type Entity string

var (
	classNameRegexp *regexp.Regexp
	fileNameRegexp  *regexp.Regexp
	pathRegexp      *regexp.Regexp
)

// QueryString returns a entity's corresponding string
// suitable for searching where the file path are removed
func (e Entity) QueryString() string {
	q := strings.TrimSpace(string(e))
	q = strings.Replace(q, "/body", "", -1)
	q = strings.Replace(q, "/parameters", "", -1)
	idx := strings.Index(q, ".java/")
	if idx == -1 {
		return ""
	}
	underscoreidx := strings.LastIndex(q[:idx], "_")
	if underscoreidx == -1 {
		underscoreidx = 0
	}
	q = q[underscoreidx:]
	arr := strings.Split(q, "/[CN]/")
	if len(arr) > 2 {
		q = arr[0] + "/[CN]/" + arr[len(arr)-1]
	}
	q = naming.RemoveGenerics(q)
	parenthesisidx := strings.Index(q, "(")
	if parenthesisidx != -1 {
		arr := strings.Split(q[parenthesisidx+1:len(q)-1], ",")
		for i := range arr {
			arr[i] = lastSubstring(arr[i], ".")
		}
		q = q[:parenthesisidx] + "(" + strings.Join(arr, ",") + ")"
	}
	return q
}

// Classname returns the class name of the entity
func (e Entity) Classname() string {
	if classNameRegexp == nil {
		var err error
		classNameRegexp, err = regexp.Compile(`.+\.java/\[CN\]/([^\[]+)/`)
		if err != nil {
			panic(err)
		}
	}
	return classNameRegexp.FindAllStringSubmatch(e.QueryString(), -1)[0][1]
}

// Filename returns the file name of the entity
func (e Entity) Filename() string {
	return QuerystringFilename(e.QueryString())
}

// QuerystringFilename returns a string corresponding to the filename,
// suitable for querying
func QuerystringFilename(qs string) string {
	if fileNameRegexp == nil {
		fileNameRegexp = regexp.MustCompile(`\_([^\.]+)\.java/\[CN\]/`)
	}
	submatch := fileNameRegexp.FindAllStringSubmatch(qs, -1)
	if len(submatch) == 0 || len(submatch[0]) < 2 {
		return ""
	}
	return submatch[0][1]
}

// Path returns the path of the entity
func (e Entity) Path() string {
	if pathRegexp == nil {
		var err error
		pathRegexp, err = regexp.Compile(`([^\.]+)\.java/\[CN\]/`)
		if err != nil {
			panic(err)
		}
	}
	submatch := pathRegexp.FindAllStringSubmatch(string(e), -1)
	if len(submatch) == 0 || len(submatch[0]) < 2 {
		return ""
	}
	return submatch[0][1]
}

// Name returns the name of the entity
func (e Entity) Name() string {
	qs := e.QueryString()
	qs = qs[strings.LastIndex(qs, "/")+1:]
	parenthesisidx := strings.Index(qs, "(")
	if parenthesisidx == -1 {
		return qs
	}
	return qs[:parenthesisidx]
}

// Parameters return the parameters of the entity, if any
func (e Entity) Parameters() []string {
	qs := e.QueryString()
	parenthesisidx := strings.Index(qs, "(")
	if parenthesisidx == -1 {
		return nil
	}
	qs = qs[parenthesisidx+1 : len(qs)-1]
	if qs == "" {
		return []string{}
	}
	return strings.Split(qs, ",")
}

func lastSubstring(s, sep string) string {
	arr := strings.Split(strings.TrimSpace(s), sep)
	return arr[len(arr)-1]
}

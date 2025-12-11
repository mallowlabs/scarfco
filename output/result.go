package output

type Result struct {
	Files []ResultFile
}

type ResultFile struct {
	Name   string
	Errors []ResultError
}

type ResultError struct {
	Source   string
	Severity string
	Message  string
	Line     int
}

func (r *Result) AddFile(name string) *ResultFile {
	f := ResultFile{Name: name}
	r.Files = append(r.Files, f)
	return &r.Files[len(r.Files)-1]
}

func (f *ResultFile) AddError(source, severity, message string, line int) {
	e := ResultError{Source: source, Severity: severity, Message: message, Line: line}
	f.Errors = append(f.Errors, e)
}

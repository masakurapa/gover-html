package profile

type Packages map[string]*Package

type Package struct {
	Dir        string
	ImportPath string
	Error      *struct {
		Err string
	}
}

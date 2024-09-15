package stringset

// creating shared packages like utils, common and base; these package names don't provide any insight into what the package provides

// instead create an expressive package name for a code group with high cohesion that doesn't really belong anywhere else

// and expose an expressive api
type Set map[string]struct{}

func New(keys ...string) Set {
	// ...
	return nil
}

func (s Set) Sort() []string {
	// ...
	return nil
}

// utility packages often handle common facilities; if we have a client and server package, where should we put the common types? One solution is to combine the client and server into a single package

func main() {

}

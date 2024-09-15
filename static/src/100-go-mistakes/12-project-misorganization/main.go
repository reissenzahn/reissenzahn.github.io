package main

// https://github.com/golang-standards/project-layout

// /cmd - the main source files; /cmd/foo/main.go for a foo application
// /internal - private code that should not be imported by others
// /pkg - public code to be exposed to others
// /test - external tests and test data
// /configs - configuration
// /docs - design and user documentation
// /examples - usage examples
// /api - api contract files
// /web - web app assets (static files, etc.)
// /build - packaging and CI files
// /scripts - scripts for installation, analysis, etc.
// /vendor - application dependencies

// there is no concept of sub-packages; net/http can only see exported net members

// subdirectories help to organize packages by cohesion

// avoid huge packages and tiny packages

// name packages after what they provide; rather than what they contain using short, concise, expressive and, by convention, a single lowercase word

// export as little as possible to reduce the coupling between packages

func main() {

}


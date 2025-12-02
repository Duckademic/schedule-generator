package components

// GeneratorComponent represents any component participating in the generation process.
type GeneratorComponent interface {
	Run()                          // The main improvement of schedule for generator
	GetErrorService() ErrorService // Each component must embed an ErrorService to report generation errors.
}

package command

type Interface interface {
	Process(cmdArgs string) string
	Name() string
	Help() string
}

package central

type Options struct {
	keys []string
	values []string
}

func NewOptions() *Options {
	return &Options{}
}

func (options *Options) AddOption(key string, value string) {
	options.keys = append(options.keys, key)
	options.values = append(options.values, value)
}

func (options *Options) GetKeys() []string {
	return options.keys
}

func (options *Options) GetValues() []string {
	return options.values
}

func (options *Options) HasOptions() bool {
	return len(options.keys) > 0
}

func (options *Options) String() string {
	string := ""

	for i := 0; i < len(options.keys); i++ {
		string += " [" + options.keys[i] + "] " + options.values[i]
	}

	return string
}
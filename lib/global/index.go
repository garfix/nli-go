package global

type index struct {
	Type         string
	Version      string
	Uses         map[string]string
	Read         []string
	Write		 []string
	Solution	 []string
}
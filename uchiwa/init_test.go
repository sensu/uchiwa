package uchiwa

func init() {
	c, _ := LoadConfig("../test/gotest/config_test.json")
	d := New(c)
	Build(d)
}

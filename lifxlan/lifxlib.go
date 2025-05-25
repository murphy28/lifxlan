package lifxlan

func init() {
	if err := InitializeProducts(); err != nil {
		panic(err)
	}
}

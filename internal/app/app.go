package app

func Run() {
	server := NewServer()
	defer server.Close()

	server.Start()
}

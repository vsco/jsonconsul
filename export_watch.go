package jsonconsul

func (c *JsonExport) RunWatcher() {
	for {
		err := c.Run()
		if err != nil {
			log.Println(err)
			break
		}

		log.Println("Waiting", time.Second*c.WatchFrequency)
		<-time.After(time.Second * c.WatchFrequency)
	}
}

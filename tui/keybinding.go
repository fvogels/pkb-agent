package tui

type KeyBinding struct {
	Key         string
	Description string
	Message     Message
}

func Match(messageQueue MessageQueue, message MsgKey, bindings ...KeyBinding) {
	for _, binding := range bindings {
		if message.Key == binding.Key {
			messageQueue.Enqueue(binding.Message)
		}
	}
}

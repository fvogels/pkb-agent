package tui

type KeyBinding struct {
	Key         string
	Description string
	Message     Message
}

func HandleKeyBindings(messageQueue MessageQueue, message MsgKey, bindings ...KeyBinding) bool {
	for _, binding := range bindings {
		if message.Key == binding.Key {
			messageQueue.Enqueue(binding.Message)
			return true
		}
	}

	return false
}

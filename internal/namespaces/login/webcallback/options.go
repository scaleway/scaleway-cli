package webcallback

type Options func(*WebCallback)

func WithPort(port int) Options {
	return func(callback *WebCallback) {
		callback.port = port
	}
}

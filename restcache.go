package discordgo

// RestCache represents a basic REST caching service
type RestCache interface {
	Set(string, []byte) error
	Get(string) ([]byte, bool)
}

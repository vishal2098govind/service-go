package tagdebug

type Conf struct {
	apiHost string `conf:"test=123"`
}
type ConfPublic struct {
	ApiHost string `conf:"test=123"`
}

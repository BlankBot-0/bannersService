package cache

type (
	Cache interface {
		Get(key string) string
		Set(key string, value string)
	}
)

type Deps struct {
	Cache Cache
}

type CacheSystem struct {
	Deps
}

func NewCacheSystem(deps Deps) *CacheSystem {
	return &CacheSystem{
		Deps: deps,
	}
}

func (c *CacheSystem) Set(key string, value string) {

}

func (c *CacheSystem) Get(key string) string {}

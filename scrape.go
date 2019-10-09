package scrape

import (
	"os"

	"github.com/javscrape/go-scrape/net"
)

var debug = false

// IScrape ...
type IScrape interface {
	GrabSample(b bool)
	IsGrabSample() (b bool)
	CacheImage(path string)
	SortOut(path string)
	Find(name string) (msg *[]*Content, e error)
}

type scrapeImpl struct {
	grabs  []IGrab
	sample bool
	cache  string
	out    string
}

// SortOut ...
func (impl *scrapeImpl) SortOut(path string) {
	impl.out = path
}

// CacheImage ...
func (impl *scrapeImpl) CacheImage(path string) {
	impl.cache = path
}

// IsGrabSample ...
func (impl *scrapeImpl) IsGrabSample() bool {
	return impl.sample
}

// GrabSample ...
func (impl *scrapeImpl) GrabSample(b bool) {
	impl.sample = b
	if !impl.sample {
		return
	}
	for _, grab := range impl.grabs {
		grab.Sample(b)
	}
}

// DebugOn ...
func DebugOn() {
	debug = true
}

// NewScrape ...
func NewScrape(grabs ...IGrab) IScrape {
	return &scrapeImpl{grabs: grabs}
}

// Find ...
func (impl *scrapeImpl) Find(name string) (msg *[]*Content, e error) {
	msg = new([]*Content)
	for _, grab := range impl.grabs {
		iGrab, e := grab.Find(name)
		if e != nil {
			log.With("name", grab.Name(), "find", name).Error(e)
			continue
		}
		e = iGrab.Decode(msg)
		if e != nil {
			log.With("name", grab.Name(), "find", name).Error(e)
		}
	}

	if impl.cache == "" {
		c := net.NewCache(impl.cache)
		e := imageCache(c, *msg)
		if e != nil {
			return nil, e
		}
	}
	return
}

func imageCache(cache *net.Cache, msg []*Content) (e error) {
	path := make(chan string)
	go func(path chan<- string) {
		defer close(path)
		for _, m := range msg {
			path <- m.Image
			path <- m.Thumb
			for _, act := range m.Actors {
				path <- act.Image
			}
			for _, s := range m.Sample {
				path <- s.Image
				path <- s.Thumb
			}
		}
	}(path)

	for p := range path {
		if p != "" {
			e = cache.Get(p)
			if e != nil && !os.IsExist(e) {
				log.Error(e)
			}
		}
	}
	return nil
}

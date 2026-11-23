package connect

import "errors"

// PreferPublic returns the first public endpoint, or the first private one if
// none is public. isPublic / isPrivate classify each endpoint.
func PreferPublic[E any](endpoints []E, isPublic, isPrivate func(E) bool) (E, error) {
	var zero E

	for _, e := range endpoints {
		if isPublic(e) {
			return e, nil
		}
	}
	for _, e := range endpoints {
		if isPrivate(e) {
			return e, nil
		}
	}

	return zero, errors.New("no endpoint found")
}

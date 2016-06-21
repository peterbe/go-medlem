package iris

import (
	"sort"
)

type Routes []*Route

//implementing the sort.Interface for type 'Routes'

func (routes Routes) Len() int {
	return len(routes)
}

func (routes Routes) Less(r1, r2 int) bool {
	//sort by longest path parts no  longest fullpath, longest first.
	r1Len := len(routes[r1].pathParts)
	r2Len := len(routes[r2].pathParts)
	if r1Len == r2Len { //if same len then static path is first.
		return routes[r1].isStatic
	}

	return r1Len > r2Len

}

func (routes Routes) Swap(r1, r2 int) {
	routes[r1], routes[r2] = routes[r2], routes[r1]
}

//end

type branch struct {
	prefix string
	routes Routes
}

type tree []*branch

//implementing the sort.Interface for type 'tree'

func (branches tree) Len() int {
	return len(branches)
}

func (branches tree) Less(r1, r2 int) bool {
	//sort by longest path prefix, longest first.
	return len(branches[r1].prefix) > len(branches[r2].prefix)
}

func (branches tree) Swap(r1, r2 int) {
	branches[r1], branches[r2] = branches[r2], branches[r1]
}

//end

type Trees map[string]tree

func (_trees Trees) addRoute(method string, route *Route) {
	if _trees[method] == nil {
		_trees[method] = make([]*branch, 0)
	}
	ok := false
	var _branch *branch
	index := 0
	for index, _branch = range _trees[method] {
		//check if route has parameters or * after the prefix, if yes then add a slash to the end
		routePref := route.PathPrefix

		if _branch.prefix == routePref {
			_trees[method][index].routes = append(_branch.routes, route)
			ok = true
			break
		}
	}
	if !ok {
		_branch = &branch{prefix: route.PathPrefix, routes: make([]*Route, 0)}
		_branch.routes = append(_branch.routes, route)
		_trees[method] = append(_trees[method], _branch)
	}

}

// sort is used one time before the Listen, at the Build state
func (_trees Trees) sort() {
	for method, _tree := range _trees {
		for index, _ := range _tree {
			// sort the routes by the biggest path parts
			sort.Sort(_trees[method][index].routes)
		}
		// sort the branches by the longest prefix
		sort.Sort(_trees[method])
	}
}

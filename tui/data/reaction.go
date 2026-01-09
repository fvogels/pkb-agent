package data

func DefineReaction(reaction func(), dependencies ...Observable) {
	for _, dependency := range dependencies {
		dependency.Observe(reaction)
	}
}

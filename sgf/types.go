package sgf

// Collection represents an sgf `Collection` object
type Collection struct {
	Trees []*GameTree
}

// GameTree represents an SGF `GameTree`
type GameTree struct {
	// Principal is the principal variation
	Principal Sequence
	// Children represent alternate paths
	Children []*GameTree
}

// Sequence is an SGF `Sequence`
type Sequence struct {
	Nodes []Node
}

// Node is an SGF `Node`
type Node struct {
	Props []Property
}

// Property is an SGF `Property`
type Property struct {
	Prop   string
	Values []PropValue
}

// PropValue is an SGF `PropValue`. Property values are stored as raw
// strings for maximum compatibility and flexibility, but convenience
// methods are provided to interpret properties in the standad SGF
// formats.
type PropValue string

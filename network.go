package main

type Network struct {
	Nodes map[int]Radar
	Zipf  Zipf
}

func NewNetwork(numNodes int) Network {

	nodes := make(map[int]Radar)

	for i := range numNodes {
		nodes[i] = newRadar(i + 1)
	}

	return Network{
		Nodes: nodes,
	}
}

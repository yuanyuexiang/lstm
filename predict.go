package lstm

import (
	"context"

	"github.com/owulveryck/lstm/datasetter"
	G "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

// Predict ...
func (m *Model) Predict(ctx context.Context, dataSet datasetter.ReadWriter) error {
	hiddenT := tensor.New(tensor.Of(tensor.Float32), tensor.WithShape(m.hiddenSize))
	cellT := tensor.New(tensor.Of(tensor.Float32), tensor.WithShape(m.hiddenSize))
	lstm := m.newLSTM(hiddenT, cellT)
	// We need an empty memory to start...
	prevHidden := G.NewVector(lstm.g, tensor.Float32, G.WithName("hₜ₋₁"), G.WithShape(m.hiddenSize), G.WithValue(hiddenT))
	prevCell := G.NewVector(lstm.g, tensor.Float32, G.WithName("Cₜ₋₁"), G.WithShape(m.hiddenSize), G.WithValue(cellT))
	// First pass to get update the hidden state and the cell according to the input
	_, _, err := lstm.forwardStep(dataSet, prevHidden, prevCell, 0)
	if err != nil {
		return err
	}
	// TODO Then run this as many times as needed until a stop

	//g := lstm.g.SubgraphRoots(dataSet.GetComputedVectors()...)
	//machine := G.NewTapeMachine(g, G.ExecuteFwdOnly())
	machine := G.NewLispMachine(lstm.g, G.ExecuteFwdOnly())
	err = machine.RunAll()
	return err

}

package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	I   = interface{}
	In  = <-chan I
	Out = In
	Bi  = chan I
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	ch := in
	for _, stage := range stages {
		ch = interStage(stage(ch), done)
	}
	return ch
}

func interStage(in In, done In) Out {
	out := make(Bi)
	go func(out Bi) {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			}
		}
	}(out)
	return out
}

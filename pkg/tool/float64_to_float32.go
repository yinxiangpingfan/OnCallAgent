package tool

func ToFloat32(f []float64) []float32 {
	res := make([]float32, len(f))
	for i, v := range f {
		res[i] = float32(v)
	}
	return res
}

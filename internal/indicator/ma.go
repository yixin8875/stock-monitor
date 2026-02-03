package indicator

// MA 计算移动平均线
func MA(closes []float64, period int) []float64 {
	if len(closes) < period {
		return nil
	}

	result := make([]float64, len(closes))
	for i := range result {
		if i < period-1 {
			result[i] = 0
			continue
		}
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += closes[i-j]
		}
		result[i] = sum / float64(period)
	}
	return result
}

// LastMA 计算最新的MA值
func LastMA(closes []float64, period int) float64 {
	if len(closes) < period {
		return 0
	}
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += closes[len(closes)-1-i]
	}
	return sum / float64(period)
}

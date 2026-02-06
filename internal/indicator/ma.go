package indicator

// MA 计算移动平均线（滑动窗口，O(n)复杂度）
func MA(closes []float64, period int) []float64 {
	if len(closes) < period {
		return nil
	}

	result := make([]float64, len(closes))
	// 计算第一个窗口的和
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += closes[i]
	}
	result[period-1] = sum / float64(period)

	// 滑动窗口：加入新值，移除旧值
	for i := period; i < len(closes); i++ {
		sum += closes[i] - closes[i-period]
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

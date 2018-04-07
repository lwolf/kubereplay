package helpers

func BlueGreenReplicas(n int32, segmentSize int32) (blueReplicas int32, greenReplicas int32) {
	if segmentSize == 100 {
		return n, 0
	}
	if n == 1 {
		return 1, 1
	}
	coef := float32(segmentSize) / 100
	blueReplicas = int32(float32(n) * coef)
	greenReplicas = n - blueReplicas
	return
}

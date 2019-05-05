package metric

const prefix = "readflow_"

func metricName(name string) string {
	return prefix + name
}

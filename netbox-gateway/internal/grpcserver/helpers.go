package grpcserver

func getStringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func getInt64Value(i *int64) int64 {
	if i != nil {
		return *i
	}
	return 0
}

func getTenantID(tenant interface{}) int64 {
	// Implementar conforme estrutura real do NetBox
	// Por enquanto retorna 0
	return 0
}

package l21

type LegacyReaderAdapter struct {
	legacy *LegacyReader
}

func (a *LegacyReaderAdapter) ReadData() string {
	return a.legacy.GetLegacyData()
}

func NewLegacyReaderAdapter(legacy *LegacyReader) *LegacyReaderAdapter {
	return &LegacyReaderAdapter{legacy: legacy}
}

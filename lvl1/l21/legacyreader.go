package l21

type LegacyReader struct{}

func (l *LegacyReader) GetLegacyData() string {
	return "Данные из legacyReader"
}

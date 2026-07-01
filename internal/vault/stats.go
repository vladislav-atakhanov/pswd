package vault

type Stats struct {
	Passwords           int
	Devices             int
	HeaderSize          int
	BodySize            int
	ActiveEncryptedSize int
	DataEnd             int
	OrphanedCount       int
	OrphanedSize        int
}

func (v *Vault) Stats() Stats {
	s := Stats{
		Passwords:  len(v.content),
		Devices:    len(v.devices),
		HeaderSize: v.HeaderLength(),
		DataEnd:    v.dataEnd,
	}
	s.BodySize = v.dataEnd - v.HeaderLength()

	for _, item := range v.content {
		s.ActiveEncryptedSize += item.length
	}

	s.OrphanedCount = len(v.orphanedSpans)
	for _, sp := range v.orphanedSpans {
		s.OrphanedSize += sp.Length
	}

	return s
}

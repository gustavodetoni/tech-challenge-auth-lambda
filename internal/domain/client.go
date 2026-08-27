package domain

type DocumentType string

const (
	DocumentTypeCPF  DocumentType = "CPF"
	DocumentTypeCNPJ DocumentType = "CNPJ"
)

type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
)

type Client struct {
	ID             string
	DocumentType   DocumentType
	DocumentNumber string
	Name           string
	Status         Status
}

func (c Client) IsActive() bool {
	return c.Status == StatusActive
}

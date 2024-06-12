package models

type StorageBucket struct {
	ID        string `bun:"id,pk"`
	Name      string `bun:"name"`
	Slug      string `bun:"slug,unique"`
	Host      string `bun:"host"`
	KeyId     string `bun:"key_id"`
	SecretKey string `bun:"secret_key"`

	UserID string `bun:"user_id"`
	User   *User  `bun:"rel:belongs-to,join:user_id=id"`
}

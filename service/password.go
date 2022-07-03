package service

import (
	"context"
	"crypto/rand"
	"crypto/subtle"

	"github.com/rs/zerolog"
	"github.com/stockwayup/pass/conf"
	"golang.org/x/crypto/argon2"
)

const saltSize = 16

type Password struct {
	cfg *conf.Config
}

func NewPasswordSvc(cfg *conf.Config) *Password {
	return &Password{cfg: cfg}
}

func (s *Password) HashPassword(ctx context.Context, password []byte) (hash []byte, salt []byte, err error) {
	salt = make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("generate salt")

		return hash, salt, err
	}

	hash = argon2.IDKey(
		password,
		salt,
		s.cfg.Password.Time,
		s.cfg.Password.Memory,
		s.cfg.Password.Threads,
		s.cfg.Password.KeyLen,
	)

	return hash, salt, nil
}

func (s *Password) IsValid(in, hash, salt []byte) (bool, error) {
	comparisonHash := argon2.IDKey(
		in,
		salt,
		s.cfg.Password.Time,
		s.cfg.Password.Memory,
		s.cfg.Password.Threads,
		s.cfg.Password.KeyLen,
	)

	return subtle.ConstantTimeCompare(hash, comparisonHash) == 1, nil
}

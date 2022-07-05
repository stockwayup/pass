package service // nolint: testpackage

import (
	"context"
	"testing"

	"github.com/stockwayup/pass/conf"
)

func TestPassword_HashPassword(t *testing.T) {
	t.Parallel()

	type fields struct {
		cfg *conf.Config
	}

	type args struct {
		ctx      context.Context // nolint: containedctx
		password []byte
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantHash []byte
		wantSalt []byte
		wantErr  bool
	}{
		{
			name: "success",
			fields: fields{
				cfg: &conf.Config{
					Env: "test",
					Password: struct {
						Time    uint32 `json:"time"    default:"1"`
						Memory  uint32 `json:"memory"  default:"65536"`
						Threads uint8  `json:"threads" default:"4"`
						KeyLen  uint32 `json:"key_len" default:"32"`
					}{
						Time:    1,
						Memory:  65536,
						Threads: 1,
						KeyLen:  32,
					},
				},
			},
			args: args{
				ctx:      context.Background(),
				password: []byte("password"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &Password{
				cfg: tt.fields.cfg,
			}

			_, _, err := s.HashPassword(tt.args.ctx, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}

func BenchmarkPassword_HashPassword_1_65536(b *testing.B) {
	benchmarkPasswordHashPassword(b, uint32(1), uint32(65536))
}

func BenchmarkPassword_HashPassword_2_65536(b *testing.B) {
	benchmarkPasswordHashPassword(b, uint32(2), uint32(65536))
}

func BenchmarkPassword_HashPassword_3_65536(b *testing.B) {
	benchmarkPasswordHashPassword(b, uint32(3), uint32(65536))
}

func BenchmarkPassword_HashPassword_1_131072(b *testing.B) {
	benchmarkPasswordHashPassword(b, uint32(1), uint32(131072))
}

func BenchmarkPassword_HashPassword_2_131072(b *testing.B) {
	benchmarkPasswordHashPassword(b, uint32(2), uint32(131072))
}

func benchmarkPasswordHashPassword(b *testing.B, time, memory uint32) {
	b.Helper()

	for n := 0; n < b.N; n++ {
		s := &Password{
			cfg: &conf.Config{
				Env: "test",
				Password: struct {
					Time    uint32 `json:"time"    default:"1"`
					Memory  uint32 `json:"memory"  default:"65536"`
					Threads uint8  `json:"threads" default:"4"`
					KeyLen  uint32 `json:"key_len" default:"32"`
				}{
					Time:    time,
					Memory:  memory,
					Threads: 1,
					KeyLen:  32,
				},
			},
		}

		_, _, _ = s.HashPassword(context.Background(), []byte("password"))
	}
}

func BenchmarkPassword_IsValid_1_65536(b *testing.B) {
	benchmarkPasswordIsValid(b, uint32(1), uint32(65536))
}

func BenchmarkPassword_IsValid_2_65536(b *testing.B) {
	benchmarkPasswordIsValid(b, uint32(2), uint32(65536))
}

func BenchmarkPassword_IsValid_3_65536(b *testing.B) {
	benchmarkPasswordIsValid(b, uint32(3), uint32(65536))
}

func BenchmarkPassword_IsValid_1_131072(b *testing.B) {
	benchmarkPasswordIsValid(b, uint32(1), uint32(131072))
}

func BenchmarkPassword_IsValid_2_131072(b *testing.B) {
	benchmarkPasswordIsValid(b, uint32(2), uint32(131072))
}

func benchmarkPasswordIsValid(b *testing.B, time, memory uint32) {
	b.Helper()

	for n := 0; n < b.N; n++ {
		s := &Password{
			cfg: &conf.Config{
				Env: "test",
				Password: struct {
					Time    uint32 `json:"time"    default:"1"`
					Memory  uint32 `json:"memory"  default:"65536"`
					Threads uint8  `json:"threads" default:"4"`
					KeyLen  uint32 `json:"key_len" default:"32"`
				}{
					Time:    time,
					Memory:  memory,
					Threads: 1,
					KeyLen:  32,
				},
			},
		}

		_, _ = s.IsValid(
			[]byte("password"),
			[]byte("0uEhQE4YRUz/LVMJGabSwUIy1hIapH5pvnXr3/Uf5fY"),
			[]byte("2DHUWesLywClWUMBglAoMA"),
		)
	}
}

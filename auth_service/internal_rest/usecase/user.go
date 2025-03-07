package usecase

type AuthUseCase struct {
	repo ISsoRepository
}

func NewAuthUseCase(repo ISsoRepository) *AuthUseCase {
	return &AuthUseCase{repo: repo}
}

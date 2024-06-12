package account

import (
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/account"
	"samm/pkg/utils"
	"time"
)

func CreateAccountBuilder(payload *account.StoreAccountDto) domain.Account {
	accountDomain := domain.Account{}
	accountDomain.Name.Ar = payload.Name.Ar
	accountDomain.Name.En = payload.Name.En
	accountDomain.Email = payload.Email
	password, _ := utils.HashPassword(payload.Password)
	accountDomain.Password = password
	accountDomain.CreatedAt = time.Now()
	accountDomain.UpdatedAt = time.Now()
	accountDomain.Country.Id = payload.Country.Id
	accountDomain.Country.PhonePrefix = payload.Country.PhonePrefix
	accountDomain.Country.Currency = payload.Country.Currency
	accountDomain.Country.Timezone = payload.Country.Timezone
	accountDomain.Country.Name.Ar = payload.Country.Name.Ar
	accountDomain.Country.Name.En = payload.Country.Name.En

	return accountDomain
}

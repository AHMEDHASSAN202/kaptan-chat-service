package brand

import (
	"github.com/jinzhu/copier"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/brand"
	"samm/pkg/utils"
	"time"
)

func convertDtoToCorrespondingDomain(dto *brand.UpdateBrandDto) *domain.Brand {
	brandDoc := domain.Brand{}
	copier.Copy(&brandDoc, dto)
	brandDoc.ID = utils.ConvertStringIdToObjectId(dto.Id)
	brandDoc.UpdatedAt = time.Now()
	return &brandDoc
}

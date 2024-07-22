package location

import "samm/pkg/utils/dto"

type FindLocationMobileDto struct {
	WithCollectionMethod bool `query:"with_collection_method"`
	dto.MobileHeaders
}

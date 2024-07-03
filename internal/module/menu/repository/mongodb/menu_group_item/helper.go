package menu_group_item

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/menu/dto/menu_group"
	"samm/pkg/database/mongodb"
	"samm/pkg/utils"
	"strings"
	"time"
)

func AddAvailabilityQuery(countryId, field string) bson.M {
	currentTime := time.Now().UTC().Format(utils.DefaultTimeFormat)
	currentDay := utils.GetDay(countryId)
	return bson.M{
		"$expr": bson.M{
			"$or": bson.A{
				bson.M{"$ne": bson.A{field, nil}},
				bson.M{"$ne": bson.A{field, make([]interface{}, 0)}},
				bson.D{{"$gt", bson.A{
					bson.D{
						{"$size", bson.D{
							{"$filter", bson.D{
								{"input", field},
								{"as", "hours"},
								{"cond", bson.D{
									{"$and", bson.A{
										bson.D{
											{"$eq", bson.A{
												"$$hours.day", strings.ToLower(currentDay),
											}},
										},
										bson.D{
											{"$lte", bson.A{
												"$$hours.from",
												currentTime,
											}},
										},
										bson.D{
											{"$gte", bson.A{
												"$$hours.to",
												currentTime,
											}},
										},
										bson.D{
											{"$gte", bson.A{
												"$$hours.to",
												"$$hours.from",
											}},
										},
									}},
								}},
							}},
						}},
					},
					0,
				}}},
				bson.D{{"$gt", bson.A{
					bson.D{
						{"$size", bson.D{
							{"$filter", bson.D{
								{"input", field},
								{"as", "hours"},
								{"cond", bson.D{
									{"$and", bson.A{
										bson.D{
											{"$eq", bson.A{
												"$$hours.day", strings.ToLower(currentDay),
											}},
										},
										bson.D{
											{"$gte", bson.A{
												currentTime,
												"$$hours.from",
											}},
										},
										bson.D{
											{"$gte", bson.A{
												"$$hours.from",
												"$$hours.to",
											}},
										},
									}},
								}},
							}},
						}},
					},
					0,
				}}},
				bson.D{{"$gt", bson.A{
					bson.D{
						{"$size", bson.D{
							{"$filter", bson.D{
								{"input", field},
								{"as", "hours"},
								{"cond", bson.D{
									{"$and", bson.A{
										bson.D{
											{"$eq", bson.A{
												"$$hours.day", strings.ToLower(currentDay),
											}},
										},
										bson.D{
											{"$lte", bson.A{
												currentTime,
												"$$hours.to",
											}},
										},
										bson.D{
											{"$gte", bson.A{
												"$$hours.from",
												"$$hours.to",
											}},
										},
									}},
								}},
							}},
						}},
					},
					0,
				}}},
			},
		},
	}
}

func createIndexes(collection *mongo.Collection) {
	mongodb.CreateIndex(collection, false,
		bson.E{"menu_group._id", mongodb.IndexType.Asc},
		bson.E{"category._id", mongodb.IndexType.Asc},
	)

	mongodb.CreateIndex(collection, false,
		bson.E{"menu_group.branch_ids", mongodb.IndexType.Asc},
		bson.E{"menu_group.status", mongodb.IndexType.Asc},
		bson.E{"category.status", mongodb.IndexType.Asc},
		bson.E{"category.sort", mongodb.IndexType.Asc},
		bson.E{"category._id", mongodb.IndexType.Asc},
		bson.E{"status", mongodb.IndexType.Asc},
		bson.E{"sort", mongodb.IndexType.Asc},
		bson.E{"modifier_group_ids", mongodb.IndexType.Asc},
		bson.E{"name.ar", mongodb.IndexType.Text},
		bson.E{"name.en", mongodb.IndexType.Text},
	)
}

func getProductAndModifierId(order *menu_group.FilterMenuGroupItemsForOrder) ([]primitive.ObjectID, []primitive.ObjectID) {
	var modifierIds []primitive.ObjectID
	var productIds []primitive.ObjectID

	for _, item := range order.MenuItems {
		productIds = append(productIds, utils.ConvertStringIdToObjectId(item.Id))
		for _, modifier := range item.ModifierIds {
			modifierIds = append(modifierIds, utils.ConvertStringIdToObjectId(modifier.Id))
		}
	}
	return productIds, modifierIds
}

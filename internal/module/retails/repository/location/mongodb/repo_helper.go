package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"samm/pkg/utils"
	"strings"
	"time"
)

func AddIsOpenFieldPipeline(CountryId string) bson.M {
	currentTime := time.Now().UTC().Format("15:04:05")
	currentDay := utils.GetDayByCountry(CountryId)

	return bson.M{
		"$addFields": bson.D{
			{"is_open", bson.D{
				{"$cond", bson.A{
					bson.D{
						{"$or", bson.A{
							bson.D{
								{"$and", bson.A{
									bson.D{{"$gt", bson.A{
										bson.D{
											{"$size", bson.D{
												{"$filter", bson.D{
													{"input", "$working_hour"},
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
									bson.D{
										{"$eq", bson.A{
											"$open",
											true,
										}}},
								},
								},
							},
							bson.D{
								{"$and", bson.A{
									bson.D{{"$gt", bson.A{
										bson.D{
											{"$size", bson.D{
												{"$filter", bson.D{
													{"input", "$working_hour"},
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
									bson.D{
										{"$eq", bson.A{
											"$open",
											true,
										}}},
								},
								},
							},
							bson.D{
								{"$and", bson.A{
									bson.D{{"$gt", bson.A{
										bson.D{
											{"$size", bson.D{
												{"$filter", bson.D{
													{"input", "$working_hour"},
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
									bson.D{
										{"$eq", bson.A{
											"$open",
											true,
										}}},
								},
								},
							},
							bson.D{
								{"$and", bson.A{
									bson.D{{"$gt", bson.A{
										bson.D{
											{"$size", bson.D{
												{"$filter", bson.D{
													{"input", "$working_hour"},
													{"as", "hours"},
													{"cond", bson.D{
														{"$and", bson.A{
															bson.D{
																{"$eq", bson.A{
																	"$$hours.day", strings.ToLower(currentDay),
																}},
															},
															bson.D{
																{"$eq", bson.A{
																	"$$hours.is_full_day", true,
																}},
															},
														}},
													}},
												}},
											}},
										},
										0,
									}}},
									bson.D{
										{"$eq", bson.A{
											"$open",
											true,
										}}},
								},
								},
							},
						}},
					},
					true,
					false,
				}},
			}},
		},
	}
}

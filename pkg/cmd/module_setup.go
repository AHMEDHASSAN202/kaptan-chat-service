package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// createFolders creates a list of folders within the specified root directory.
func createFolders(root string, folders []string) error {
	for _, folder := range folders {
		path := filepath.Join(root, folder)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
		fmt.Println("Created folder:", path)
	}
	return nil
}

// createMainFile creates a main.go file in the root directory.
func createMainFile(root string) error {
	mainFilePath := filepath.Join(root, "main.go")
	content := `package main

	import "fmt"

		func main() {
			fmt.Println("Hello, World!")
		}
	
	`
	if err := os.WriteFile(mainFilePath, []byte(content), 0644); err != nil {
		return err
	}
	fmt.Println("Created file:", mainFilePath)
	return nil
}

func createDeliveryFile(newModulePath, newModuleName, rootModuleName string) error {

	moduleCamelCase := toCamelCase(newModuleName)
	moduleLowerCamelCase := toLowerCamelCase(newModuleName)

	fileName := moduleLowerCamelCase + ".delivery.go"
	mainFilePath := filepath.Join(newModulePath+"/delivery", fileName)
	content := `package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"kaptan/internal/module/kitchen/domain"
	"kaptan/internal/module/kitchen/dto/kitchen"
	"kaptan/pkg/logger"
	"kaptan/pkg/validators"
)

type KitchenHandler struct {
	kitchenUsecase domain.KitchenUseCase
	validator      *validator.Validate
	logger         logger.ILogger
}

// InitKitchenController will initialize the article's HTTP controller
func InitKitchenController(e *echo.Echo, us domain.KitchenUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &KitchenHandler{
		kitchenUsecase: us,
		validator:      validator,
		logger:         logger,
	}
	dashboard := e.Group("api/v1/admin/kitchen")
	dashboard.POST("", handler.CreateKitchen)
	dashboard.GET("", handler.ListKitchen)
	dashboard.PUT("/:id", handler.UpdateKitchen)
	dashboard.GET("/:id", handler.FindKitchen)
	dashboard.DELETE("/:id", handler.DeleteKitchen)
}
func (a *KitchenHandler) CreateKitchen(c echo.Context) error {
	ctx := c.Request().Context()

	var payload kitchen.StoreKitchenDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.kitchenUsecase.CreateKitchen(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *KitchenHandler) UpdateKitchen(c echo.Context) error {
	ctx := c.Request().Context()

	var payload kitchen.UpdateKitchenDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}
	id := c.Param("id")
	errResp := a.kitchenUsecase.UpdateKitchen(ctx, id, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *KitchenHandler) FindKitchen(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	data, errResp := a.kitchenUsecase.FindKitchen(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"kitchen": data})
}

func (a *KitchenHandler) DeleteKitchen(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	errResp := a.kitchenUsecase.DeleteKitchen(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *KitchenHandler) ListKitchen(c echo.Context) error {
	ctx := c.Request().Context()
	var payload kitchen.ListKitchenDto

	_ = c.Bind(&payload)

	payload.Pagination.SetDefault()

	result, errResp := a.kitchenUsecase.List(&ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, result)
}

	`

	content = strings.Replace(content, "kaptan/", rootModuleName+"/", -1)
	content = strings.Replace(content, "internal/module/retails", newModulePath, -1)
	content = strings.Replace(content, "Kitchen", moduleCamelCase, -1)
	content = strings.Replace(content, "kitchen", moduleLowerCamelCase, -1)

	if err := os.WriteFile(mainFilePath, []byte(content), 0644); err != nil {
		return err
	}
	fmt.Println("Created file:", mainFilePath)
	return nil
}

func createDomainFile(newModulePath, newModuleName, rootModuleName string) error {

	moduleCamelCase := toCamelCase(newModuleName)
	moduleLowerCamelCase := toLowerCamelCase(newModuleName)

	fileName := moduleLowerCamelCase + ".domain.go"
	mainFilePath := filepath.Join(newModulePath+"/domain", fileName)
	content := `package domain

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"kaptan/internal/module/kitchen/dto/kitchen"
	"kaptan/internal/module/kitchen/responses"
	"kaptan/pkg/validators"
	"time"
)

type Kitchen struct {
	mgm.DefaultModel ` + "`bson:\",inline\"`" + `
	Name             Name       ` + "`json:\"name\" bson:\"name\"`" + `
	Email            string     ` + "`json:\"email\" bson:\"email\"`" + `
	Password         string     ` + "`json:\"-\" bson:\"password\"`" + `
	DeletedAt        *time.Time ` + "`json:\"-\" bson:\"deleted_at\"`" + `
}

type Name struct {
	Ar string ` + "`json:\"ar\" bson:\"ar\"`" + `
	En string ` + "`json:\"en\" bson:\"en\"`" + `
}

type KitchenUseCase interface {
	CreateKitchen(ctx context.Context, payload *kitchen.StoreKitchenDto) (err validators.ErrorResponse)
	UpdateKitchen(ctx context.Context, id string, payload *kitchen.UpdateKitchenDto) (err validators.ErrorResponse)
	FindKitchen(ctx context.Context, Id string) (kitchen Kitchen, err validators.ErrorResponse)
	DeleteKitchen(ctx context.Context, Id string) (err validators.ErrorResponse)
	List(ctx *context.Context, dto *kitchen.ListKitchenDto) (*responses.ListResponse, validators.ErrorResponse)
}

type KitchenRepository interface {
	CreateKitchen(kitchen *Kitchen) (err error)
	UpdateKitchen(kitchen *Kitchen) (err error)
	FindKitchen(ctx context.Context, Id primitive.ObjectID) (kitchen *Kitchen, err error)
	DeleteKitchen(ctx context.Context, Id primitive.ObjectID) (err error)
	List(ctx *context.Context, dto *kitchen.ListKitchenDto) (usersRes *[]Kitchen, paginationMeta *PaginationData, err error)
}

`
	content = strings.Replace(content, "kaptan/", rootModuleName+"/", -1)
	content = strings.Replace(content, "internal/module/retails", newModulePath, -1)
	content = strings.Replace(content, "Kitchen", moduleCamelCase, -1)
	content = strings.Replace(content, "kitchen", moduleLowerCamelCase, -1)

	if err := os.WriteFile(mainFilePath, []byte(content), 0644); err != nil {
		return err
	}
	fmt.Println("Created file:", mainFilePath)
	return nil
}

func createUseCaseFile(newModulePath, newModuleName, rootModuleName string) error {

	moduleCamelCase := toCamelCase(newModuleName)
	moduleLowerCamelCase := toLowerCamelCase(newModuleName)

	fileName := moduleLowerCamelCase + ".usecase.go"
	mainFilePath := filepath.Join(newModulePath+"/usecase/"+newModuleName, fileName)
	content := `package kitchen

import (
	"context"
	"kaptan/internal/module/kitchen/domain"
	"kaptan/internal/module/kitchen/dto/kitchen"
	"kaptan/internal/module/kitchen/responses"
	"kaptan/pkg/logger"
	"kaptan/pkg/utils"
	"kaptan/pkg/validators"
	"time"
)

type KitchenUseCase struct {
	repo   domain.KitchenRepository
	logger logger.ILogger
}

const tag = " KitchenUseCase "

func NewKitchenUseCase(repo domain.KitchenRepository, logger logger.ILogger) domain.KitchenUseCase {
	return &KitchenUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (l KitchenUseCase) CreateKitchen(ctx context.Context, payload *kitchen.StoreKitchenDto) (err validators.ErrorResponse) {
	kitchenDomain := domain.Kitchen{}
	kitchenDomain.Name.Ar = payload.Name.Ar
	kitchenDomain.Name.En = payload.Name.En
	kitchenDomain.Email = payload.Email
	password, er := utils.HashPassword(payload.Password)
	if er != nil {
		return validators.GetErrorResponseFromErr(er)
	}
	kitchenDomain.Password = password
	kitchenDomain.CreatedAt = time.Now()
	kitchenDomain.UpdatedAt = time.Now()

	dbErr := l.repo.CreateKitchen(&kitchenDomain)
	if dbErr != nil {
		return validators.GetErrorResponseFromErr(dbErr)
	}
	return
}

func (l KitchenUseCase) UpdateKitchen(ctx context.Context, id string, payload *kitchen.UpdateKitchenDto) (err validators.ErrorResponse) {
	kitchenDomain, dbErr := l.repo.FindKitchen(ctx, utils.ConvertStringIdToObjectId(id))
	if dbErr != nil {
		return validators.GetErrorResponseFromErr(dbErr)
	}
	kitchenDomain.Name.Ar = payload.Name.Ar
	kitchenDomain.Name.En = payload.Name.En
	kitchenDomain.Email = payload.Email

	if payload.Password != "" {
		password, er := utils.HashPassword(payload.Password)
		if er != nil {
			return validators.GetErrorResponseFromErr(er)
		}
		kitchenDomain.Password = password
	}
	kitchenDomain.UpdatedAt = time.Now()

	dbErr = l.repo.UpdateKitchen(kitchenDomain)
	if dbErr != nil {
		return validators.GetErrorResponseFromErr(dbErr)
	}
	return
}
func (l KitchenUseCase) FindKitchen(ctx context.Context, Id string) (kitchen domain.Kitchen, err validators.ErrorResponse) {
	domainKitchen, dbErr := l.repo.FindKitchen(ctx, utils.ConvertStringIdToObjectId(Id))
	if dbErr != nil {
		return *domainKitchen, validators.GetErrorResponseFromErr(dbErr)
	}
	return *domainKitchen, validators.ErrorResponse{}
}

func (l KitchenUseCase) DeleteKitchen(ctx context.Context, Id string) (err validators.ErrorResponse) {

	delErr := l.repo.DeleteKitchen(ctx, utils.ConvertStringIdToObjectId(Id))
	if delErr != nil {
		return validators.GetErrorResponseFromErr(delErr)
	}
	return validators.ErrorResponse{}
}

func (l KitchenUseCase) List(ctx *context.Context, dto *kitchen.ListKitchenDto) (*responses.ListResponse, validators.ErrorResponse) {
	users, paginationMeta, resErr := l.repo.List(ctx, dto)
	if resErr != nil {
		return nil, validators.GetErrorResponseFromErr(resErr)
	}
	return responses.SetListResponse(users, paginationMeta), validators.ErrorResponse{}
}

`
	content = strings.Replace(content, "kaptan/", rootModuleName+"/", -1)
	content = strings.Replace(content, "internal/module/retails", newModulePath, -1)
	content = strings.Replace(content, "Kitchen", moduleCamelCase, -1)
	content = strings.Replace(content, "kitchen", moduleLowerCamelCase, -1)

	if err := os.WriteFile(mainFilePath, []byte(content), 0644); err != nil {
		return err
	}
	fmt.Println("Created file:", mainFilePath)
	return nil
}

func createResponseFile(newModulePath, newModuleName, rootModuleName string) error {

	moduleCamelCase := toCamelCase(newModuleName)
	moduleLowerCamelCase := toLowerCamelCase(newModuleName)

	fileName := "list_dashboard.go"
	mainFilePath := filepath.Join(newModulePath+"/responses/", fileName)
	content := `package responses

import (
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"kaptan/pkg/utils"
)

type ListResponse struct {
	Docs interface{}                ` + "`json:\"docs\" bson:\"docs\"`" + ` 
	Meta *mongopagination.PaginationData ` + "`json:\"meta\" bson:\"meta\"`" + `
}

func SetListResponse(docs interface{}, meta *mongopagination.PaginationData) *ListResponse {
	listResponse := ListResponse{
		Docs: docs,
		Meta: meta,
	}
	if listResponse.Meta == nil {
		listResponse.Meta = &mongopagination.PaginationData{}
	}
	if utils.IsNil(docs) {
		listResponse.Docs = make([]interface{}, 0)
	}
	return &listResponse
}

`
	content = strings.Replace(content, "kaptan/", rootModuleName+"/", -1)
	content = strings.Replace(content, "internal/module/retails", newModulePath, -1)
	content = strings.Replace(content, "Kitchen", moduleCamelCase, -1)
	content = strings.Replace(content, "kitchen", moduleLowerCamelCase, -1)

	if err := os.WriteFile(mainFilePath, []byte(content), 0644); err != nil {
		return err
	}
	fmt.Println("Created file:", mainFilePath)
	return nil
}

func createRepoFile(newModulePath, newModuleName, rootModuleName string) error {

	moduleCamelCase := toCamelCase(newModuleName)
	moduleLowerCamelCase := toLowerCamelCase(newModuleName)

	fileName := moduleLowerCamelCase + ".repo.go"
	mainFilePath := filepath.Join(newModulePath+"/repository/"+newModuleName, fileName)
	content := `package mongodb

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kaptan/internal/module/kitchen/domain"
	"kaptan/internal/module/kitchen/dto/kitchen"
	"kaptan/pkg/logger"
	"time"
)

type KitchenRepository struct {
	kitchenCollection *mgm.Collection
	logger            logger.ILogger
}

const mongoKitchenRepositoryTag = "KitchenMongoRepository"

func NewKitchenMongoRepository(dbs *mongo.Database, log logger.ILogger) domain.KitchenRepository {
	kitchenDbCollection := mgm.Coll(&domain.Kitchen{})

	return &KitchenRepository{
		kitchenCollection: kitchenDbCollection,
		logger:            log,
	}
}

func (l KitchenRepository) CreateKitchen(kitchen *domain.Kitchen) (err error) {
	err = l.kitchenCollection.Create(kitchen)
	if err != nil {
		return err
	}
	return nil

}

func (l KitchenRepository) UpdateKitchen(kitchen *domain.Kitchen) (err error) {
	upsert := true
	opts := options.UpdateOptions{Upsert: &upsert}
	err = l.kitchenCollection.Update(kitchen, &opts)
	return
}
func (l KitchenRepository) FindKitchen(ctx context.Context, Id primitive.ObjectID) (kitchen *domain.Kitchen, err error) {
	domainData := domain.Kitchen{}
	filter := bson.M{"deleted_at": nil, "_id": Id}
	err = l.kitchenCollection.FirstWithCtx(ctx, filter, &domainData)

	return &domainData, err
}

func (l KitchenRepository) DeleteKitchen(ctx context.Context, Id primitive.ObjectID) (err error) {
	kitchenData, err := l.FindKitchen(ctx, Id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	kitchenData.DeletedAt = &now
	kitchenData.UpdatedAt = now
	return l.UpdateKitchen(kitchenData)
}

func (l *KitchenRepository) List(ctx *context.Context, dto *kitchen.ListKitchenDto) (usersRes *[]domain.Kitchen, paginationMeta *PaginationData, err error) {
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.D{{"deleted_at", nil}},
	}}}

	if dto.Query != "" {
		pattern := ".*" + dto.Query + ".*"
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"$or": []bson.M{{"name": bson.M{"$regex": pattern, "$options": "i"}}, {"phone_number": bson.M{"$regex": pattern, "$options": "i"}}}})
	}

	data, err := New(l.kitchenCollection.Collection).Context(*ctx).Limit(dto.Limit).Page(dto.Page).Sort("created_at", -1).Aggregate(matching)

	if data == nil || data.Data == nil {
		return nil, nil, err
	}

	users := make([]domain.Kitchen, 0)
	for _, raw := range data.Data {
		model := domain.Kitchen{}
		err = bson.Unmarshal(raw, &model)
		if err != nil {
			l.logger.Error("kitchen Repo -> List -> ", err)
			break
		}
		users = append(users, model)
	}
	paginationMeta = &data.Pagination
	usersRes = &users

	return
}


`
	content = strings.Replace(content, "kaptan/", rootModuleName+"/", -1)
	content = strings.Replace(content, "internal/module/retails", newModulePath, -1)
	content = strings.Replace(content, "Kitchen", moduleCamelCase, -1)
	content = strings.Replace(content, "kitchen", moduleLowerCamelCase, -1)

	if err := os.WriteFile(mainFilePath, []byte(content), 0644); err != nil {
		return err
	}
	fmt.Println("Created file:", mainFilePath)
	return nil
}

func createDtoFile(newModulePath, newModuleName, rootModuleName string) error {

	moduleCamelCase := toCamelCase(newModuleName)
	moduleLowerCamelCase := toLowerCamelCase(newModuleName)

	fileName := moduleLowerCamelCase + ".dto.go"
	mainFilePath := filepath.Join(newModulePath+"/dto/"+newModuleName, fileName)
	content := `package kitchen

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"kaptan/pkg/validators"
	"kaptan/pkg/utils/dto"
)

type Name struct {
	Ar string ` + "`json:\"ar\" validate:\"required,min=3\"`" + `
	En string ` + "`json:\"en\" validate:\"required,min=3\"`" + `
}

type StoreKitchenDto struct {
	Name     Name   ` + "`json:\"name\" validate:\"required\"`" + `
	Email    string ` + "`json:\"email\" validate:\"required,email\"`" + `
	Password string ` + "`json:\"password\" validate:\"required\"`" + `
}

func (payload *StoreKitchenDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload)
}

type ListKitchenDto struct {
	dto.Pagination
	Query string ` + "`query:\"query\"`" + `
}

type UpdateKitchenDto struct {
	Name     Name   ` + "`json:\"name\" validate:\"required\"`" + `
	Email    string ` + "`json:\"email\" validate:\"required,email\"`" + `
	Password string ` + "`json:\"password\"`" + `
}

func (payload *UpdateKitchenDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload)
}
`
	content = strings.Replace(content, "kaptan/", rootModuleName+"/", -1)
	content = strings.Replace(content, "internal/module/retails", newModulePath, -1)
	content = strings.Replace(content, "Kitchen", moduleCamelCase, -1)
	content = strings.Replace(content, "kitchen", moduleLowerCamelCase, -1)

	if err := os.WriteFile(mainFilePath, []byte(content), 0644); err != nil {
		return err
	}
	fmt.Println("Created file:", mainFilePath)
	return nil
}

// ToCamelCase converts the first letter of each word in a string to camel case
func toCamelCase(input string) string {
	// Split the string into words
	words := strings.Fields(input)

	// Process each word
	for i, word := range words {
		if len(word) > 0 {
			words[i] = string(unicode.ToUpper(rune(word[0]))) + strings.ToLower(word[1:])
		}
	}

	// Join the words back together
	return strings.Join(words, " ")
}

func toLowerCamelCase(s string) string {
	camelCase := toCamelCase(s)
	if len(camelCase) > 0 {
		return strings.ToLower(camelCase[:1]) + camelCase[1:]
	}
	return camelCase
}

// updateModuleFile updates the specified Go module file with new import and provide lines.
func updateModuleFile(newModuleName, newModulePath, rootModuleName string) error {
	moduleCamelCase := toCamelCase(newModuleName)
	moduleLowerCamelCase := toLowerCamelCase(newModuleName)
	useCaseAlias := moduleLowerCamelCase + "_usecase"
	repoAlias := moduleLowerCamelCase + "_repo"
	filePath := newModulePath + "/module.go"
	newImports := []string{
		useCaseAlias + "\"" + rootModuleName + "/" + newModulePath + "/usecase/" + moduleLowerCamelCase + "\"",
		repoAlias + "\"" + rootModuleName + "/" + newModulePath + "/repository/" + moduleLowerCamelCase + "\"",
	}
	newProvides := []string{
		repoAlias + ".New" + moduleCamelCase + "MongoRepository",
		useCaseAlias + ".New" + moduleCamelCase + "UseCase",
	}
	newInvokes := ", delivery." + "Init" + moduleCamelCase + "Controller"

	// Read the original file content
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	lines := strings.Split(string(fileContent), "\n")
	var updatedLines []string

	// Track positions to insert new lines
	importStartPos := -1
	importEndPos := -1
	provideEndPos := -1
	inProvideBlock := false

	for i, line := range lines {
		// Track the fx.Invoke block
		if strings.Contains(line, "fx.Invoke(") {
			line = strings.Replace(line, ")", newInvokes+")", -1)
		}

		updatedLines = append(updatedLines, line)

		// Find the start and end positions of the import block
		if strings.HasPrefix(line, "import (") {
			importStartPos = i
		}
		if importStartPos >= 0 && strings.HasPrefix(line, ")") {
			importEndPos = i
			importStartPos = -1 // Reset to avoid re-entering the block
		}

		// Track the fx.Provide block
		if strings.Contains(line, "fx.Provide(") {
			inProvideBlock = true
		}
		if inProvideBlock {
			if strings.Contains(line, ")") {
				provideEndPos = i
				inProvideBlock = false
			}
		}
	}

	counterAdded := 0
	// Insert new imports
	if importEndPos > 0 {
		for _, newImport := range newImports {
			updatedLines = append(updatedLines[:importEndPos], append([]string{newImport}, updatedLines[importEndPos:]...)...)
			importEndPos++
			counterAdded++
		}
	}

	// Insert new provides
	provideEndPos += counterAdded
	if provideEndPos > 0 {
		for _, newProvide := range newProvides {
			updatedLines = append(updatedLines[:provideEndPos], append([]string{"\t\t" + newProvide + ","}, updatedLines[provideEndPos:]...)...)
			provideEndPos++
			counterAdded++
		}
	}

	// Write back the modified content to the file
	updatedContent := strings.Join(updatedLines, "\n")
	err = os.WriteFile(filePath, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func createModuleFile(newModuleName, newModulePath, rootModuleName string) error {

	moduleCamelCase := toCamelCase(newModuleName)
	moduleLowerCamelCase := toLowerCamelCase(newModuleName)

	mainFilePath := filepath.Join(newModulePath, "module.go")
	content := `package kitchen

import (
	"go.uber.org/fx"
	"kaptan/internal/module/kitchen/delivery"
	kitchen_repo "kaptan/internal/module/kitchen/repository/kitchen"
	kitchen_usecase "kaptan/internal/module/kitchen/usecase/kitchen"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		// App Config
		kitchen_repo.NewKitchenMongoRepository,
		kitchen_usecase.NewKitchenUseCase,
	),
	fx.Invoke(
		delivery.InitKitchenController,
	),
)


`
	content = strings.Replace(content, "kaptan/", rootModuleName+"/", -1)
	content = strings.Replace(content, "internal/module/kitchen", newModulePath, -1)
	content = strings.Replace(content, "Kitchen", moduleCamelCase, -1)
	content = strings.Replace(content, "kitchen", moduleLowerCamelCase, -1)

	if err := os.WriteFile(mainFilePath, []byte(content), 0644); err != nil {
		fmt.Println("ERRRR:", err)
		return err
	}
	fmt.Println("Created file:", mainFilePath)
	return nil
}

func main() {
	// Define the root directory for the module and the folders to be created.
	getNewModuleName := flag.String("new_module", "", "")
	getNewModulePath := flag.String("new_module_path", "", "")
	getRootModuleName := flag.String("root_module", "", "")
	flag.Parse()
	newModuleName := *getNewModuleName
	newModulePath := *getNewModulePath
	rootModuleName := *getRootModuleName
	delivery := newModulePath + "/delivery"
	domain := newModulePath + "/domain"
	dto := newModulePath + "/dto/" + newModuleName
	useCase := newModulePath + "/usecase/" + newModuleName
	responses := newModulePath + "/responses/" + newModuleName
	repository := newModulePath + "/repository/" + newModuleName
	folders := []string{delivery, domain, dto, useCase, responses, repository}

	// Check if module name is provided
	if newModuleName == "" {
		fmt.Println("Error: module name is required")
		flag.Usage()
		return
	}

	// Initialize the module (assuming you have go installed and setup correctly)
	// cmd := exec.Command("go", "mod", "init", newModuleName)

	// if err := cmd.Run(); err != nil {
	// 	fmt.Println("Error initializing module:", err)
	// 	return
	// }
	fmt.Println("Initialized Go module:", newModuleName)

	// Create the folders
	if err := createFolders(".", folders); err != nil {
		fmt.Println("Error creating folders:", err)
		return
	}

	// Create a main.go file
	// if err := createMainFile("."); err != nil {
	// 	fmt.Println("Error creating main.go file:", err)
	// 	return
	// }

	if err := createDeliveryFile(newModulePath, newModuleName, rootModuleName); err != nil {
		fmt.Println("Error createDeliveryFile.go file:", err)
		return
	}
	if err := createDomainFile(newModulePath, newModuleName, rootModuleName); err != nil {
		fmt.Println("Error createDomainFile.go file:", err)
		return
	}

	if err := createUseCaseFile(newModulePath, newModuleName, rootModuleName); err != nil {
		fmt.Println("Error createUseCaseFile.go file:", err)
		return
	}

	if err := createResponseFile(newModulePath, newModuleName, rootModuleName); err != nil {
		fmt.Println("Error createUseCaseFile.go file:", err)
		return
	}

	if err := createRepoFile(newModulePath, newModuleName, rootModuleName); err != nil {
		fmt.Println("Error createUseCaseFile.go file:", err)
		return
	}

	if err := createDtoFile(newModulePath, newModuleName, rootModuleName); err != nil {
		fmt.Println("Error createUseCaseFile.go file:", err)
		return
	}

	err := updateModuleFile(newModuleName, newModulePath, rootModuleName)
	if err != nil {
		fmt.Println("Error updating file:", err)
		crErr := createModuleFile(newModuleName, newModulePath, rootModuleName)
		if crErr != nil {
			fmt.Println("Error creating file:", err)
		}
	} else {
		fmt.Println("File updated successfully")
	}

	fmt.Println("Go module setup completed.")
}

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
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/account"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type AccountHandler struct {
	accountUsecase domain.AccountUseCase
	validator      *validator.Validate
	logger         logger.ILogger
}

// InitAccountController will initialize the article's HTTP controller
func InitAccountController(e *echo.Echo, us domain.AccountUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &AccountHandler{
		accountUsecase: us,
		validator:      validator,
		logger:         logger,
	}
	dashboard := e.Group("api/v1/admin/account")
	dashboard.POST("", handler.StoreAccount)
	dashboard.GET("", handler.ListAccount)
	dashboard.PUT("/:id", handler.UpdateAccount)
	dashboard.GET("/:id", handler.FindAccount)
	dashboard.DELETE("/:id", handler.DeleteAccount)
}
func (a *AccountHandler) StoreAccount(c echo.Context) error {
	ctx := c.Request().Context()

	var payload account.StoreAccountDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.accountUsecase.StoreAccount(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *AccountHandler) UpdateAccount(c echo.Context) error {
	ctx := c.Request().Context()

	var payload account.UpdateAccountDto
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
	errResp := a.accountUsecase.UpdateAccount(ctx, id, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *AccountHandler) FindAccount(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	data, errResp := a.accountUsecase.FindAccount(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"account": data})
}

func (a *AccountHandler) DeleteAccount(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	errResp := a.accountUsecase.DeleteAccount(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *AccountHandler) ListAccount(c echo.Context) error {
	ctx := c.Request().Context()
	var payload account.ListAccountDto

	_ = c.Bind(&payload)

	payload.Pagination.SetDefault()

	result, paginationResult, errResp := a.accountUsecase.ListAccount(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"data": result, "meta": paginationResult})
}	
	`

	content = strings.Replace(content, "samm/", rootModuleName+"/", -1)
	content = strings.Replace(content, "internal/module/retails", newModulePath, -1)
	content = strings.Replace(content, "Account", moduleCamelCase, -1)
	content = strings.Replace(content, "account", moduleLowerCamelCase, -1)

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
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/dto/account"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"time"
)

type Account struct {
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


type AccountUseCase interface {
	StoreAccount(ctx context.Context, payload *account.StoreAccountDto) (err validators.ErrorResponse)
	UpdateAccount(ctx context.Context, id string, payload *account.UpdateAccountDto) (err validators.ErrorResponse)
	FindAccount(ctx context.Context, Id string) (account Account, err validators.ErrorResponse)
	DeleteAccount(ctx context.Context, Id string) (err validators.ErrorResponse)
	ListAccount(ctx context.Context, payload *account.ListAccountDto) (accounts []Account, paginationResult utils.PaginationResult, err validators.ErrorResponse)
}

type AccountRepository interface {
	StoreAccount(ctx context.Context, account *Account) (err error)
	UpdateAccount(ctx context.Context, account *Account) (err error)
	FindAccount(ctx context.Context, Id primitive.ObjectID) (account *Account, err error)
	DeleteAccount(ctx context.Context, Id primitive.ObjectID) (err error)
	ListAccount(ctx context.Context, payload *account.ListAccountDto) (locations []Account, paginationResult utils.PaginationResult, err error)
}

`
	content = strings.Replace(content, "samm/", rootModuleName+"/", -1)
	content = strings.Replace(content, "internal/module/retails", newModulePath, -1)
	content = strings.Replace(content, "Account", moduleCamelCase, -1)
	content = strings.Replace(content, "account", moduleLowerCamelCase, -1)

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
	content := `package account

import (
	"context"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/account"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"time"
)

type AccountUseCase struct {
	repo            domain.AccountRepository
	logger          logger.ILogger
}

const tag = " AccountUseCase "

func NewAccountUseCase(repo domain.AccountRepository, logger logger.ILogger) domain.AccountUseCase {
	return &AccountUseCase{
		repo:            repo,
		logger:          logger,
	}
}

func (l AccountUseCase) StoreAccount(ctx context.Context, payload *account.StoreAccountDto) (err validators.ErrorResponse) {
	accountDomain := domain.Account{}
	accountDomain.Name.Ar = payload.Name.Ar
	accountDomain.Name.En = payload.Name.En
	accountDomain.Email = payload.Email
	password, er := utils.HashPassword(payload.Password)
	if er != nil {
		return validators.GetErrorResponseFromErr(er)
	}
	accountDomain.Password = password
	accountDomain.CreatedAt = time.Now()
	accountDomain.UpdatedAt = time.Now()

	errRe := l.repo.StoreAccount(ctx, &accountDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}

func (l AccountUseCase) UpdateAccount(ctx context.Context, id string, payload *account.UpdateAccountDto) (err validators.ErrorResponse) {
	accountDomain, errRe := l.repo.FindAccount(ctx, utils.ConvertStringIdToObjectId(id))
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	accountDomain.Name.Ar = payload.Name.Ar
	accountDomain.Name.En = payload.Name.En
	accountDomain.Email = payload.Email

	if payload.Password != "" {
		password, er := utils.HashPassword(payload.Password)
		if er != nil {
			return validators.GetErrorResponseFromErr(er)
		}
		accountDomain.Password = password
	}
	accountDomain.UpdatedAt = time.Now()

	errRe = l.repo.UpdateAccount(ctx, accountDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}
func (l AccountUseCase) FindAccount(ctx context.Context, Id string) (account domain.Account, err validators.ErrorResponse) {
	domainAccount, errRe := l.repo.FindAccount(ctx, utils.ConvertStringIdToObjectId(Id))
	if errRe != nil {
		return *domainAccount, validators.GetErrorResponseFromErr(errRe)
	}
	return *domainAccount, validators.ErrorResponse{}
}

func (l AccountUseCase) DeleteAccount(ctx context.Context, Id string) (err validators.ErrorResponse) {

	delErr := l.repo.DeleteAccount(ctx, utils.ConvertStringIdToObjectId(Id))
	if delErr != nil {
		return validators.GetErrorResponseFromErr(delErr)
	}
	return validators.ErrorResponse{}
}

func (l AccountUseCase) ListAccount(ctx context.Context, payload *account.ListAccountDto) (accounts []domain.Account, paginationResult utils.PaginationResult, err validators.ErrorResponse) {
	results, paginationResult, errRe := l.repo.ListAccount(ctx, payload)
	if errRe != nil {
		return results, paginationResult, validators.GetErrorResponseFromErr(errRe)
	}
	return results, paginationResult, validators.ErrorResponse{}

}
`
	content = strings.Replace(content, "samm/", rootModuleName+"/", -1)
	content = strings.Replace(content, "internal/module/retails", newModulePath, -1)
	content = strings.Replace(content, "Account", moduleCamelCase, -1)
	content = strings.Replace(content, "account", moduleLowerCamelCase, -1)

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
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/account"
	"samm/pkg/utils"
	"time"
)

type AccountRepository struct {
	accountCollection *mgm.Collection
}

const mongoAccountRepositoryTag = "AccountMongoRepository"

func NewAccountMongoRepository(dbs *mongo.Database) domain.AccountRepository {
	accountDbCollection := mgm.Coll(&domain.Account{})

	return &AccountRepository{
		accountCollection: accountDbCollection,
	}
}

func (l AccountRepository) StoreAccount(ctx context.Context, account *domain.Account) (err error) {
	_, err = mgm.Coll(&domain.Account{}).InsertOne(ctx, account)
	if err != nil {
		return err
	}
	return nil

}

func (l AccountRepository) UpdateAccount(ctx context.Context, account *domain.Account) (err error) {
	update := bson.M{"$set": account}
	_, err = mgm.Coll(&domain.Account{}).UpdateByID(ctx, account.ID, update)
	return
}
func (l AccountRepository) FindAccount(ctx context.Context, Id primitive.ObjectID) (account *domain.Account, err error) {
	domainData := domain.Account{}
	filter := bson.M{"deleted_at": nil, "_id": Id}
	err = l.accountCollection.FirstWithCtx(ctx, filter, &domainData)

	return &domainData, err
}

func (l AccountRepository) DeleteAccount(ctx context.Context, Id primitive.ObjectID) (err error) {
	accountData, err := l.FindAccount(ctx, Id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	accountData.DeletedAt = &now
	accountData.UpdatedAt = now
	return l.UpdateAccount(ctx, accountData)
}

func (l AccountRepository) ListAccount(ctx context.Context, payload *account.ListAccountDto) (accounts []domain.Account, paginationResult utils.PaginationResult, err error) {

	offset := (payload.Page - 1) * payload.Limit
	findOptions := options.Find().SetLimit(payload.Limit).SetSkip(offset)

	filter := bson.M{}
	match := []bson.M{}
	match = append(match, bson.M{"deleted_at": nil})
	if payload.Query != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name.ar": bson.M{"$regex": payload.Query, "$options": "i"}},
				{"name.en": bson.M{"$regex": payload.Query, "$options": "i"}},
				{"email": bson.M{"$regex": payload.Query, "$options": "i"}},
			},
		}
	}
	filter["$and"] = match

	// Query the collection for the total count of documents
	collection := mgm.Coll(&domain.Account{})
	totalItems, err := collection.CountDocuments(ctx, filter)

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalItems) / float64(payload.Limit)))

	var data []domain.Account
	err = l.accountCollection.SimpleFind(&data, filter, findOptions)

	return data, utils.PaginationResult{Page: payload.Page, TotalPages: int64(totalPages), TotalItems: totalItems}, err

}

`
	content = strings.Replace(content, "samm/", rootModuleName+"/", -1)
	content = strings.Replace(content, "internal/module/retails", newModulePath, -1)
	content = strings.Replace(content, "Account", moduleCamelCase, -1)
	content = strings.Replace(content, "account", moduleLowerCamelCase, -1)

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
	content := `package account

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/validators"
	"samm/pkg/utils/dto"
)

type Name struct {
	Ar string ` + "`json:\"ar\" validate:\"required,min=3\"`" + `
	En string ` + "`json:\"en\" validate:\"required,min=3\"`" + `
}

type StoreAccountDto struct {
	Name     Name   ` + "`json:\"name\" validate:\"required\"`" + `
	Email    string ` + "`json:\"email\" validate:\"required,email\"`" + `
	Password string ` + "`json:\"password\" validate:\"required\"`" + `
}

func (payload *StoreAccountDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload)
}

type ListAccountDto struct {
	dto.Pagination
	Query string ` + "`query:\"query\"`" + `
}

type UpdateAccountDto struct {
	Name     Name   ` + "`json:\"name\" validate:\"required\"`" + `
	Email    string ` + "`json:\"email\" validate:\"required,email\"`" + `
	Password string ` + "`json:\"password\"`" + `
}

func (payload *UpdateAccountDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload)
}
`
	content = strings.Replace(content, "samm/", rootModuleName+"/", -1)
	content = strings.Replace(content, "internal/module/retails", newModulePath, -1)
	content = strings.Replace(content, "Account", moduleCamelCase, -1)
	content = strings.Replace(content, "account", moduleLowerCamelCase, -1)

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
	content := `package account

import (
	"go.uber.org/fx"
	"samm/internal/module/account/delivery"
	account_repo "samm/internal/module/account/repository/account"
	account_usecase "samm/internal/module/account/usecase/account"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		// App Config
		account_repo.NewAccountMongoRepository,
		account_usecase.NewAccountUseCase,
	),
	fx.Invoke(
		delivery.InitAccountController,
	),
)

`
	content = strings.Replace(content, "samm/", rootModuleName+"/", -1)
	content = strings.Replace(content, "internal/module/account", newModulePath, -1)
	content = strings.Replace(content, "Account", moduleCamelCase, -1)
	content = strings.Replace(content, "account", moduleLowerCamelCase, -1)

	if err := os.WriteFile(mainFilePath, []byte(content), 0644); err != nil {
		fmt.Println("ERRRR:", err)
		return err
	}
	fmt.Println("Created file:", mainFilePath)
	return nil
}

func mainTemp() {
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
	repository := newModulePath + "/repository/" + newModuleName
	folders := []string{delivery, domain, dto, useCase, repository}

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

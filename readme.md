### Table of Contents

- [Database Preparation for Relation](#database-preparation-for-relation)
  - [Database Design](#database-design)
  - [Models](#models)
  - [Data Transfer Object (DTO)](#data-transfer-object-dto)
  - [Split Modify for Connection](#modify-repository)

---

# Database Preparation for Relation

## Database Design

![Database Design](./database-design.jpg)

## Models

- Create `user.go` file inside `models` folder, and write this below code

  > File: `models/user.go`

  ```go
  package models

  import "time"

  type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name" gorm:"type: varchar(255)"`
    Email     string    `json:"email" gorm:"type: varchar(255)"`
    Password  string    `json:"-" gorm:"type: varchar(255)"`
    CreatedAt time.Time `json:"-"`
    UpdatedAt time.Time `json:"-"`
  }

  type UsersProfileResponse struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
  }

  func (UsersProfileResponse) TableName() string {
    return "users"
  }
  ```

- Create `profile.go` file inside `models` folder, and write this below code

  > File: `models/profile.go`

  ```go
  package models

  import "time"

  type Profile struct {
    ID        int                  `json:"id" gorm:"primary_key:auto_increment"`
    Phone     string               `json:"phone" gorm:"type: varchar(255)"`
    Gender    string               `json:"gender" gorm:"type: varchar(255)"`
    Address   string               `json:"address" gorm:"type: text"`
    UserID    int                  `json:"user_id"`
    User      UsersProfileResponse `json:"user"`
    CreatedAt time.Time            `json:"-"`
    UpdatedAt time.Time            `json:"-"`
  }

  // for association relation with another table (user)
  type ProfileResponse struct {
    Phone   string `json:"phone"`
    Gender  string `json:"gender"`
    Address string `json:"address"`
    UserID  int    `json:"-"`
  }

  func (ProfileResponse) TableName() string {
    return "profiles"
  }
  ```

- Create `product.go` file inside `models` folder, and write this below code

  > File: `models/product.go`

  ```go
  package models

  import "time"

  type Product struct {
    ID         int                  `json:"id" gorm:"primary_key:auto_increment"`
    Name       string               `json:"name" form:"name" gorm:"type: varchar(255)"`
    Desc       string               `json:"desc" gorm:"type:text" form:"desc"`
    Price      int                  `json:"price" form:"price" gorm:"type: int"`
    Image      string               `json:"image" form:"image" gorm:"type: varchar(255)"`
    Qty        int                  `json:"qty" form:"qty"`
    UserID     int                  `json:"user_id" form:"user_id"`
    User       UsersProfileResponse `json:"user"`
    CreatedAt  time.Time            `json:"-"`
    UpdatedAt  time.Time            `json:"-"`
  }

  type ProductResponse struct {
    ID         int                  `json:"id"`
    Name       string               `json:"name"`
    Desc       string               `json:"desc"`
    Price      int                  `json:"price"`
    Image      string               `json:"image"`
    Qty        int                  `json:"qty"`
    UserID     int                  `json:"-"`
    User       UsersProfileResponse `json:"user"`
  }

  type ProductUserResponse struct {
    ID     int    `json:"id"`
    Name   string `json:"name"`
    Desc   string `json:"desc"`
    Price  int    `json:"price"`
    Image  string `json:"image"`
    Qty    int    `json:"qty"`
    UserID int    `json:"-"`
  }

  func (ProductResponse) TableName() string {
    return "products"
  }

  func (ProductUserResponse) TableName() string {
    return "products"
  }
  ```

## Data Transfer Object (DTO)

- Create `auth` folder, inside it Create `auth_request.go` and write this below code

  > File: `dto/auth/auth_request.go`

  ```go
  package authdto

  type AuthRequest struct {
    Name     string `gorm:"type: varchar(255)" json:"name"`
    Email    string `gorm:"type: varchar(255)" json:"email"`
    Password string `gorm:"type: varchar(255)" json:"password"`
  }
  ```

- Create `product` folder, inside it Create `product_request.go` and write this below code

  > File: `dto/product/product_request.go`

  ```go
  package productdto

  type ProductRequest struct {
    Name       string `json:"name" form:"name" gorm:"type: varchar(255)"`
    Desc       string `json:"desc" gorm:"type:text" form:"desc"`
    Price      int    `json:"price" form:"price" gorm:"type: int"`
    Image      string `json:"image" form:"image" gorm:"type: varchar(255)"`
    Qty        int    `json:"qty" form:"qty" gorm:"type: int"`
    UserID     int    `json:"user_id" gorm:"type: int"`
  }
  ```

- Create `profile` folder, inside it Create `profile_request.go` and write this below code

  > File: `dto/profile/profile_request.go`

  ```go
  package profiledto

  import "golang/models"

  type ProfileResponse struct {
    ID      int                         `json:"id" gorm:"primary_key:auto_increment"`
    Phone   string                      `json:"phone" gorm:"type: varchar(255)"`
    Gender  string                      `json:"gender" gorm:"type: varchar(255)"`
    Address string                      `json:"address" gorm:"type: text"`
    UserID  int                         `json:"user_id"`
    User    models.UsersProfileResponse `json:"user"`
  }
  ```

## Modify Repository

- Inside `repositories` folder, create `repository.go` file, and write this below code

  > File: `repositories/repository.go`

  ```go
  package repositories

  import "gorm.io/gorm"

  type repository struct {
    db *gorm.DB
  }
  ```

  \*`repository` struct move from `user.go` file


  ### Table of Contents

- [GORM Relation belongs to](#gorm-relation-belongs-to)
  - [Handlers](#handlers)
  - [Repository](#repository)
  - [Routes](#routes)

---

# GORM Relation Belongs to

Reference: [Official GORM Website](https://gorm.io/docs/belongs_to.html)

## Relation

For this section, example Belongs To relation:

- `Profile` &rarr; `User`: to get Profile User
- `Product` &rarr; `User`: to get Product Owner

## Handlers

- Inside `handlers` folder, create `profile.go` file, and write this below code

  > File: `handlers/profile.go`

  ```go
  package handlers

  import (
    profiledto "golang/dto/profile"
    dto "golang/dto/result"
    "golang/models"
    "golang/repositories"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
  )

  type handlerProfile struct {
    ProfileRepository repositories.ProfileRepository
  }

  func HandlerProfile(ProfileRepository repositories.ProfileRepository) *handlerProfile {
    return &handlerProfile{ProfileRepository}
  }

  func (h *handlerProfile) GetProfile(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    var profile models.Profile
    profile, err := h.ProfileRepository.GetProfile(id)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
      json.NewEncoder(w).Encode(response)
      return
    }

    w.WriteHeader(http.StatusOK)
    response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseProfile(profile)}
    json.NewEncoder(w).Encode(response)
  }

  func convertResponseProfile(u models.Profile) profiledto.ProfileResponse {
    return profiledto.ProfileResponse{
      ID:      u.ID,
      Phone:   u.Phone,
      Gender:  u.Gender,
      Address: u.Address,
      UserID:  u.UserID,
      User:    u.User,
    }
  }
  ```

- Inside `handlers` folder, create `product.go` file, and write this below code

  > File: `handlers/product.go`

  ```go
  package handlers

  import (
    productdto "golang/dto/product"
    dto "golang/dto/result"
    "golang/models"
    "golang/repositories"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-playground/validator/v10"
    "github.com/gorilla/mux"
  )

  type handlerProduct struct {
    ProductRepository repositories.ProductRepository
  }

  func HandlerProduct(ProductRepository repositories.ProductRepository) *handlerProduct {
    return &handlerProduct{ProductRepository}
  }

  func (h *handlerProduct) FindProducts(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    products, err := h.ProductRepository.FindProducts()
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
      json.NewEncoder(w).Encode(response)
      return
    }

    w.WriteHeader(http.StatusOK)
    response := dto.SuccessResult{Code: http.StatusOK, Data: products}
    json.NewEncoder(w).Encode(response)
  }

  func (h *handlerProduct) GetProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    var product models.Product
    product, err := h.ProductRepository.GetProduct(id)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
      json.NewEncoder(w).Encode(response)
      return
    }

    w.WriteHeader(http.StatusOK)
    response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseProduct(product)}
    json.NewEncoder(w).Encode(response)
  }

  func (h *handlerProduct) CreateProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    request := new(productdto.ProductRequest)
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
      w.WriteHeader(http.StatusBadRequest)
      response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
      json.NewEncoder(w).Encode(response)
      return
    }

    validation := validator.New()
    err := validation.Struct(request)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
      json.NewEncoder(w).Encode(response)
      return
    }

    product := models.Product{
      Name:   request.Name,
      Desc:   request.Desc,
      Price:  request.Price,
      Image:  request.Image,
      Qty:    request.Qty,
      UserID: request.UserID,
    }

    product, err = h.ProductRepository.CreateProduct(product)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
      json.NewEncoder(w).Encode(response)
      return
    }

    product, _ = h.ProductRepository.GetProduct(product.ID)

    w.WriteHeader(http.StatusOK)
    response := dto.SuccessResult{Code: http.StatusOK, Data: product}
    json.NewEncoder(w).Encode(response)
  }

  func convertResponseProduct(u models.Product) models.ProductResponse {
    return models.ProductResponse{
      Name:     u.Name,
      Desc:     u.Desc,
      Price:    u.Price,
      Image:    u.Image,
      Qty:      u.Qty,
      User:     u.User,
    }
  }
  ```

## Repository

- Inside `repositories` folder, create `profile.go` file, and write this below code

  > File: `repositories/profile.go`

  ```go
  package repositories

  import (
    "golang/models"

    "gorm.io/gorm"
  )

  type ProfileRepository interface {
    GetProfile(ID int) (models.Profile, error)
  }

  func RepositoryProfile(db *gorm.DB) *repository {
    return &repository{db}
  }

  func (r *repository) GetProfile(ID int) (models.Profile, error) {
    var profile models.Profile
    err := r.db.Preload("User").First(&profile, ID).Error

    return profile, err
  }
  ```

- Inside `repositories` folder, create `product.go` file, and write this below code

  > File: `repositories/product.go`

  ```go
  package repositories

  import (
    "golang/models"

    "gorm.io/gorm"
  )

  type ProductRepository interface {
    FindProducts() ([]models.Product, error)
    GetProduct(ID int) (models.Product, error)
    CreateProduct(product models.Product) (models.Product, error)
  }

  func RepositoryProduct(db *gorm.DB) *repository {
    return &repository{db}
  }

  func (r *repository) FindProducts() ([]models.Product, error) {
    var products []models.Product
    err := r.db.Preload("User").Find(&products).Error

    return products, err
  }

  func (r *repository) GetProduct(ID int) (models.Product, error) {
    var product models.Product
    // not yet using category relation, cause this step doesnt Belong to Many
    err := r.db.Preload("User").First(&product, ID).Error

    return product, err
  }

  func (r *repository) CreateProduct(product models.Product) (models.Product, error) {
    err := r.db.Create(&product).Error

    return product, err
  }
  ```

## Routes

- Inside `routes` folder, create `profile.go` file, and write this below code

  > File: `routes/profile.go`

  ```go
  package routes

  import (
    "golang/handlers"
    "golang/pkg/mysql"
    "golang/repositories"

    "github.com/gorilla/mux"
  )

  func ProfileRoutes(r *mux.Router) {
    profileRepository := repositories.RepositoryProfile(mysql.DB)
    h := handlers.HandlerProfile(profileRepository)

    r.HandleFunc("/profile/{id}", h.GetProfile).Methods("GET")
  }
  ```

- Inside `routes` folder, create `profile.go` file, and write this below code

  > File: `routes/product.go`

  ```go
  package routes

  import (
    "golang/handlers"
    "golang/pkg/mysql"
    "golang/repositories"

    "github.com/gorilla/mux"
  )

  func ProductRoutes(r *mux.Router) {
    productRepository := repositories.RepositoryProduct(mysql.DB)
    h := handlers.HandlerProduct(productRepository)

    r.HandleFunc("/products", h.FindProducts).Methods("GET")
    r.HandleFunc("/product/{id}", h.GetProduct).Methods("GET")
    r.HandleFunc("/product", h.CreateProduct).Methods("POST")
  }
  ```

- On `routes.go` file, write `ProfileRoutes` and `ProductRoutes`

  > File: `routes/routes.go`

  ```go
  package routes

  import (
    "github.com/gorilla/mux"
  )

  func RouteInit(r *mux.Router) {
    UserRoutes(r)
    ProfileRoutes(r) // Add this code
    ProductRoutes(r) // Add this code
  }
  ```

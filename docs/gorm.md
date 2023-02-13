# gorm

## Setup gorm

    $ go get -u gorm.io/gorm

Install drivers

    $ go get -u gorm.io/driver/sqlite
    $ go get -u gorm.io/driver/mysql

## Basic configuration.

[Connecting to a Database ref](https://gorm.io/docs/connecting_to_the_database.html)

To handle `time.Time` correctly, you need to include `parseTime` as a parameter.
([more parameters](https://github.com/go-sql-driver/mysql#parameters))
To fully support `UTF-8 encoding`, you need to change `charset=utf8` to `charset=utf8mb4`.
See [this article](https://mathiasbynens.be/notes/mysql-utf8mb4) for a detailed explanation

```go
database, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               "<user>:secretpass@/<database>?charset=utf8mb4&parseTime=True&loc=Local",
		DefaultStringSize: 256, // default size for string fields
	}), &gorm.Config{})
```

### Models

```go
type User struct {
	ID             string     `json:"user_id" gorm:"size:11"`
	AccountID      *string    `json:"account_id" gorm:"size:11"`
	Username       string     `json:"user_name" gorm:"size:128"`
	EmailAddress   string     `json:"email_address" gorm:"unique"`
	FirstName      string     `json:"first_name" gorm:"size:128"`
	LastName       string     `json:"last_name" gorm:"size:128"`
	DisplayName    string     `json:"display_name" gorm:"size:128"`
	AvatarImageUrl string     `json:"avatar_image_url"`
	PhoneNumber    string     `json:"phone_number" gorm:"size:64"`
	MobileNumber   string     `json:"mobile_number" gorm:"size:64"`
	Password       []byte     `json:"-" gorm:"size:64"` // don't return password on json
	LastLogin      *time.Time `json:"last_login_date"`  // The LastLogin field takes a pointer to allow setting null value in MySQL
	CreatedAt      time.Time  `json:"created" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `json:"updated" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	Token          string     `json:"token" gorm:"-"`

	Roles   []Role   `json:"roles" gorm:"many2many:user_roles;"`
	Account *Account `json:"-"`
}
```

Some `user` tools

```go
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID short version
	user.ID = util.SortableShortUUID()
}

func (user *User) SetPassword(password string) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
```

### Tricky queries

Create or update, tries to updates a row but if no rows are affected then creates it

```go
	// create or update resetPassword record
	if database.DB.
		Model(&resetPassword).
		Where("email_address = ?", resetPassword.EmailAddress).
		Updates(&resetPassword).RowsAffected == 0 {
		database.DB.Create(&resetPassword)
	}
```

### Dates, timestamps

Having a MySQL table as:

| Field     | Type        | Null | Key | Default           | Extra                                         |
|-----------|-------------|------|-----|-------------------|-----------------------------------------------|
| id        | bigint      | NO   | PRI | NULL              | auto_increment                                |
| user_id   | int         | NO   | MUL | NULL              |                                               |
| lead_id   | bigint      | NO   | MUL | NULL              |                                               |
| name      | varchar(20) | NO   |     | NULL              |                                               |
| type      | varchar(20) | NO   |     | NULL              |                                               |
| note      | text        | YES  |     | NULL              |                                               |
| createdAt | timestamp   | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED                             |
| updatedAt | timestamp   | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED on update CURRENT_TIMESTAMP |

Could potentially receive errors like `error inserting activity: Error 1048 (23000): Column 'createdAt' cannot be null`

To avoid that, add `default:CURRENT_TIMESTAMP` to the struct declaration

```go
// Activity declares the model struct
type Activity struct {
	ID          int64      `json:"id" gorm:"column:id;type:bigint"`
	UserID      int32      `json:"user_id" gorm:"column:user_id;type:int"`
	LeadID      int64      `json:"lead_id" gorm:"column:lead_id;type:bigint"`
	Name        string     `json:"name" gorm:"column:name;type:varchar(20)"`
	Type        string     `json:"type" gorm:"column:type;type:varchar(20)"`
	Note        string     `json:"note" gorm:"column:note;type:text"`
	CreatedAt   *time.Time `json:"createdAt" gorm:"column:createdAt;type:timestamp;default:CURRENT_TIMESTAMP"`
	updatedAt   *time.Time `json:"updatedAt" gorm:"column:updatedAt;type:timestamp;default:CURRENT_TIMESTAMP"`
	User        User
	Lead        Lead
}
```

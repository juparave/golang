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
  	CreatedAt      time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdateAt       time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	Token          string     `json:"token" gorm:"-"`

	Roles   []Role   `json:"roles" gorm:"many2many:user_roles;"`
	Account *Account `json:"-"`
}
```

Reference on `struct tag types` https://gorm.io/docs/models.html#Embedded-Struct

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

Multiple `Or`s

```go
    // apply filter if any
	if params.Filter != "" {
		f := fmt.Sprintf("%%%s%%", params.Filter)
		// need `query` to get group of ORs
		query := database.DB
		db = db.Where(
			// group of ORs
			query.Or("name LIKE ?", f).
				Or("description LIKE ?", f).
				Or("category LIKE ?", f).
				Or("subcategory LIKE ?", f).
				Or("brand LIKE ?", f).
				Or("family LIKE ?", f).
				Or("model LIKE ?", f),
		)
	}
```

Compare with date strings

```go
	var reports []models.QAReport

	database.DB.
		Joins("JOIN qa_samples ON qa_samples.qa_report_id = qa_reports.id").
		Joins("LEFT JOIN species ON species.id = qa_samples.specie_id").
		Joins("LEFT JOIN suppliers ON suppliers.id = qa_reports.supplier_id").
		Where("species.name LIKE ?", "%"+q+"%").
		Or("qa_reports.capture LIKE ?", "%"+q+"%").
		Or("suppliers.name LIKE ?", "%"+q+"%").
		Or("DATE_FORMAT(qa_reports.created_at, '%Y-%m-%d') = ?", q).
		Preload("Supplier").
		Preload("Samples").
		Preload("Samples.Category").
		Preload("Samples.Variety").
		Preload("Samples.Specie").
		Preload("Samples.SubCategory").
		Order("qa_reports.updated_at desc").
		Limit(limit).
		Find(&reports)
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

### Handling short dates

```go
import (
	"database/sql/driver"
	"strings"
	"time"
)

// ShortDate is a custom type to handle dates
type ShortDate struct {
	time.Time
}

// UnmarshalJSON is a custom unmarshaler for ShortDate
// Easier to handle post from JSON input
func (sd *ShortDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		sd.Time = time.Time{}
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	sd.Time = t
	return nil
}

// MarshalJSON is a custom marshaler for ShortDate
// Send JSON date in short format
func (sd *ShortDate) MarshalJSON() ([]byte, error) {
	// returns short date format yyyy-mm-dd
	return []byte(sd.Time.Format("\"2006-01-02\"")), nil
}

// Value implements the driver.Valuer interface
// this allows easy insertion into the db
func (sd ShortDate) Value() (driver.Value, error) {
	return sd.Time.Format("2006-01-02"), nil
}

// Scan implements the sql.Scanner interface
// used to get data from the db
func (sd *ShortDate) Scan(src interface{}) error {
	if src == nil {
		*sd = ShortDate{}
		return nil
	}

	switch t := src.(type) {
	case time.Time:
		*sd = ShortDate{Time: t}
		return nil
	case []byte:
		d, err := time.Parse("2006-01-02", string(t))
		if err != nil {
			return err
		}
		*sd = ShortDate{Time: d}
		return nil
	case string:
		d, err := time.Parse("2006-01-02", t)
		if err != nil {
			return err
		}
		*sd = ShortDate{Time: d}
		return nil
	default:
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into ShortDate", src)
	}
}
```

Example

```go
// PersonalData declares the model struct
type PersonalData struct {
	ID           int64      `json:"id" gorm:"column:personal_id;type:bigint"`
	Paternal     string     `json:"paternal" gorm:"column:personal_paternal;type:varchar(50)"`
	Maternal     string     `json:"maternal" gorm:"column:personal_maternal;type:varchar(50)"`
	Name         string     `json:"name" gorm:"column:personal_name;type:varchar(50)"`
	Email        string     `json:"email" gorm:"column:personal_email;type:varchar(100)"`
	Phone        int64      `json:"phone" gorm:"column:personal_phone;type:bigint"`
	FullNumber   string     `json:"full_number" gorm:"column:full_number;type:varchar(14)"`
	Cp           string     `json:"cp" gorm:"column:personal_cp;type:varchar(20)"`
	Gender       string     `json:"gender" gorm:"column:personal_gender;type:varchar(20)"`
	Birthdate    *ShortDate `json:"birthdate" gorm:"column:personal_birthdate;type:date"`
	CreationDate time.Time  `json:"creation_date" gorm:"column:personal_creation_date;type:timestamp"`
	UpdateDate   time.Time  `json:"update_date" gorm:"column:personal_update_date;type:timestamp"`
}
```

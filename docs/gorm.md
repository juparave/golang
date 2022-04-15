# gorm

Basic configuration.

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

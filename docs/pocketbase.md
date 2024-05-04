# Pocketbase tips

## Extending auth

### Custom login

1. Handler Setup:

* Import necessary packages: net/http, github.com/labstack/echo/v5, and
  github.com/pocketbase/pocketbase/tokens.
* Define an `AuthRequest` struct to hold login credentials: `Identity` (e.g.,
  email or username) and Password.
* Create an ApiAuth function using the echo.Context parameter to handle login
  requests.

2. Request Handling:

* Bind the request body to the Auth struct using c.Bind(&auth). Handle binding
  errors gracefully (e.g., return http.StatusBadRequest).

3. User Record Retrieval:

* Based on your authentication logic, use
  app.PB.Dao().FindFirstRecordByData("users", "email", auth.Identity) to
  retrieve the user record matching the provided Identity.
* For flexibility, consider allowing customization of the field used for
  Identity (e.g., username, phone number).

4. Password Validation:

* Employ `record.ValidatePassword(auth.Password)` to verify the password
  against the retrieved user record.
* Handle invalid credentials gracefully (e.g., return http.StatusUnauthorized
  with an error message).

5. JWT Token Generation:

* Upon successful authentication, create a JWT token using `tokens.NewRecordAuthToken(app.PB, record)`.
* Catch any errors during token generation (e.g., return http.StatusInternalServerError with an informative message).

6. Token Response:

* Return a JSON response (e.g., HTTP status code 200) with the generated token in a field (e.g., "token").


```go
package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tokens"
)

// AuthRequest defines the expected structure of login request data
type AuthRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

// ApiAuth handles user login requests
func ApiAuth(c echo.Context) error {
	// Decode request body into AuthRequest struct
	var req AuthRequest
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Fetch user record by identity (username, email, etc.)
	userRecord, err := findUserRecord(req.Identity)
	if err != nil {
		app.Log.Error("Error finding user record", "error", err)
		// Handle invalid identity gracefully (e.g., "Invalid username or email")
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// Validate password using userRecord.ValidatePassword()
	if !userRecord.ValidatePassword(req.Password) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// Generate JWT token using tokens.NewRecordAuthToken()
	token, err := tokens.NewRecordAuthToken(app.PB, userRecord)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating token")
	}

	// Return successful login response with the JWT token
	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

// findUserRecord fetches the user record based on the provided identity
func findUserRecord(identity string) (*models.Record, error) {
	// Assuming email is used as the identity
	// return app.PB.Dao().FindFirstRecordByData("users", "email", identity)
	return app.PB.Dao().FindAuthRecordByEmail("users", identity)
}
```

### Custom data retrieval

Add custom `handlers` on your main file

```go
...
// serves static files from the provided public dir (if exists)
pb.OnBeforeServe().Add(func(e *core.ServeEvent) error {
    // register new "GET /hello" route
    e.Router.POST("/api/hello", func(c echo.Context) error {
        return c.String(200, "Hello world! from simple handler")
    }, apis.ActivityLogger(pb))

    e.Router.POST("/api/v1/auth", handlers.ApiAuth)
    e.Router.GET("/api/v1/monitor", handlers.ApiMonitor)
    return nil
})

if err := pb.Start(); err != nil {
    log.Fatal(err)
}
...
```


```go
package handlers

import (
	"net/http"
	mymodels "server/internal/models"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
)

type MonitorResponse struct {
	ID         string                 `json:"id"`
	DeviceID   string                 `json:"device_id"`
	DeviceName string                 `json:"device_name"`
	Location   string                 `json:"location"`
	Timestamp  int64                  `json:"timestamp"`
	Sensors    []*mymodels.SensorJSON `json:"sensors"`
}

func ApiMonitor(c echo.Context) error {
	// get authRecord from context
	authRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

	if authRecord == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "unauthorized",
		})
	}

	accountID := authRecord.Get("account_id")

	// get latest sensor from database
	records, err := app.PB.Dao().FindRecordsByFilter(
		"monitors",                   // collection
		"account_id = {:account_id}", // filter
		"-created",                   // sort
		50,                           // limit
		0,                            // offset
		dbx.Params{"account_id": accountID},
	)
	if err != nil {
		app.Log.Error("Error monitors records", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "error finding monitors records",
		})
	}

	// expand the `Device` field
	if errs := app.PB.Dao().ExpandRecords(records, []string{"device_id"}, nil); len(errs) > 0 {
		app.Log.Error("Error expanding device", "errors", errs)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "error expanding device",
		})
	}

	// create slice of SensorDataJSON
	var monitorsResponse []MonitorResponse
	monitorsResponse = make([]MonitorResponse, len(records))

	for i, record := range records {
		// marshal json string to struct
		var sensorData mymodels.SensorDataJSON
		record.UnmarshalJSONField("data", &sensorData)

		// convert string timestamp to int64
		ts, err := strconv.ParseInt(sensorData.Timestamp, 10, 64)
		if err != nil {
			app.Log.Error("Error converting timestamp", err)
		}

		// get the expanded device
		device := record.ExpandedOne("device_id")

		monitorsResponse[i] = MonitorResponse{
			ID:         record.GetString("id"),
			DeviceID:   device.GetString("id"),
			DeviceName: device.GetString("name"),
			Location:   sensorData.LonLat,
			Timestamp:  ts,
			Sensors:    sensorData.Sensors,
		}
	}

	return c.JSON(http.StatusOK, monitorsResponse)
}
```

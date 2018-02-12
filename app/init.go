package app

import (
	"database/sql"
	"fmt"
	"github.com/SpaTmole/gustly/app/models"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	// rgorp "github.com/revel/modules/orm/gorp/app"
	"github.com/revel/revel"
	"gopkg.in/gorp.v2"
	"os"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string

	// DB Manager.
	DB *gorp.DbMap
)

func setColumnSizes(t *gorp.TableMap, colSizes map[string]int) {
	for col, size := range colSizes {
		t.ColMap(col).MaxSize = size
	}
}

func InitDB() {
	var err error

	err = godotenv.Load()
	if err != nil {
		revel.ERROR.Println("Error loading .env file")
		return
	}

	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	connstring := fmt.Sprintf("user=%s password='%s' dbname=%s sslmode=disable", pgUser, pgPassword, "gustly")

	conn, err := sql.Open("postgres", connstring)
	if err != nil {
		revel.ERROR.Println("DB Error", err)
	}
	revel.INFO.Println("DB Connected")

	DB = &gorp.DbMap{Db: conn, Dialect: gorp.PostgresDialect{}}
	t := DB.AddTable(models.User{})
	setColumnSizes(t, map[string]int{
		"Username": 100,
		"Name":     100,
		"Phone":    24,
		"Email":    100,
	})

	// rgorp.Db.TraceOn(revel.AppLog)
	DB.CreateTables()
	revel.INFO.Println("Tables were set up.")
}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}

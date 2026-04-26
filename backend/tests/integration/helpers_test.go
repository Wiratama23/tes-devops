package integration_test

// getTestDatabaseURL returns the URL of a Postgres instance to run integration
// tests against. By default it returns "" so integration tests are skipped.
// Set DATABASE_URL via your own wrapper to opt in (or change this helper to
// read os.Getenv("DATABASE_URL")).
func getTestDatabaseURL() string {
	return ""
}

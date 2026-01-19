data "external_schema" "bun_sqlite" {
  program = [
    "sh",
    "-c",
    "ATLAS_DIALECT=sqlite go run -mod=mod -tags=tools migrations/atlas-schema.go",
  ]
}

data "external_schema" "bun_postgres" {
  program = [
    "sh",
    "-c",
    "ATLAS_DIALECT=postgres go run -mod=mod -tags=tools migrations/atlas-schema.go",
  ]
}

env "sqlite" {
  src = data.external_schema.bun_sqlite.url
  dev = "sqlite://file?mode=memory&_fk=1"
  migration {
    dir = "file://migrations/sqlite"
    format = golang-migrate
  }
}

env "postgres" {
  src = data.external_schema.bun_postgres.url
  dev = "postgres://postgres:postgres@localhost/dev?sslmode=disable"
  migration {
    dir = "file://migrations/postgres"
    format = golang-migrate
  }
}

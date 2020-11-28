module github.com/HellstromIT/go-quickswitch

go 1.15

replace github.com/HellstromIT/go-quickswitch/internal/quickswitch => ../../internal/quickswitch

require (
	github.com/HellstromIT/go-quickswitch/internal/quickswitch v0.0.0-00010101000000-000000000000
	github.com/HellstromIT/go-quickswitch/pkgs/fuzzy v0.0.0-00010101000000-000000000000 // indirect
)

replace github.com/HellstromIT/go-quickswitch/pkgs/fuzzy => ../../pkgs/fuzzy

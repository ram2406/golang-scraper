module github.com/olliefr/docker-gs-ping

go 1.19

replace modules/transform_html v1.0.0 => ./transform_html

require modules/transform_html v1.0.0

require github.com/labstack/echo/v4 v4.10.2

require github.com/akamensky/argparse v1.4.0

require github.com/avast/retry-go v3.0.0+incompatible // indirect

require modules/utils v1.0.0
replace modules/utils v1.0.0 => ./utils


require (
	github.com/PuerkitoBio/goquery v1.8.1
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/antchfx/htmlquery v1.3.0 // indirect
	github.com/antchfx/xmlquery v1.3.15 // indirect
	github.com/antchfx/xpath v1.2.3 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gocolly/colly v1.2.0
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/saintfish/chardet v0.0.0-20230101081208-5e3ef4b5456d // indirect
	github.com/temoto/robotstxt v1.1.2 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/yaml.v3 v3.0.1

)

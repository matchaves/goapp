aws s3 rm s3://store-goapp --recursive
go build citys.go &&
go build name.go &&
go build ip.go && 
go build street.go
aws s3 mv street s3://store-goapp
aws s3 mv ip s3://store-goapp
aws s3 mv name s3://store-goapp
aws s3 mv citys s3://store-goapp

## Add Product
`curl http://localhost:9090/product/addprod --request POST --header 'Productname:prod2' --header 'Productvendor:vendor2' --header 'Productinventory:6' --header 'Productdescription:description is here'`

## Get All Products
`curl http://localhost:9090/product/getprods --request GET`

## Get One Product by Id
`curl http://localhost:9090/product/getprodbyid --request GET --header 'Id:1'`

## Delete Product by Id
`curl http://localhost:9090/product/deletebyid --request DELETE --header 'Id:2'`

## Custom Query Examples
`curl http://localhost:9090/product/customquery --request GET --header 'Type:get' --header 'Query:SELECT name FROM products;'`

`curl http://localhost:9090/product/customquery --request GET --header 'Type:get' --header "Query:UPDATE products SET name = \"prod1\"  WHERE id = 3;"`
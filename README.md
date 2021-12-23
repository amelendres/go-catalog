# GO CATALOG SERVICE

The Catalog of Products

In this gRPC micro service you can see basic features and patterns in Go:
* Hexagonal Architecture
* Read Model  
* SOLID
* Application Service tests
* Integration Tests
* Repository Pattern

## Bounded Contexts

```
  ┌───────────────────┐           ┌───────────────────┐
  │      Catalog      │           │     Pricing       │ 
  └───────────────────┘           └───────────────────┘ 

```
## How can I use it?

### Prerequisites
- Docker

### Installing

Download repository
```sh
git clone https://github.com/amelendres/go-catalog.git
```


Build container
```sh
make build
```

## Run tests

Run the tests
```sh
make test
```


## Try it!
The shopping cart is running on `http://localhost:8050`

List Products
```
curl --location --request GET 'http://localhost:8050/products' \
--header 'Content-Type: application/json'
```




### TODO

* VO validations  
* Separate pricing domain
* Add rdbms

    


## Authors

* **Alfredo Melendres** -  alfredo.melendres@gmail.com

[license](LICENSE.md)

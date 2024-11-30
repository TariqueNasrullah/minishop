# Minishop
A minimalistic online shopping service written in go.

## Features
 - Login
 - Logout
 - Create a new order
 - Listing the orders
 - Cancel order

## How to run this Project

### Build the binary
```shell
  go mod download
  go build -o minishop
  ./minishop --config=config.dev.yaml serve
```

### Using Docker
```shell
  docker compose up
```
By default, the docker compose uses the `config.docker.yml` file.

By Default the backend service runs on port 8080. You can change the configuration from the config file.
Test 

The connection.
```curl
curl --location 'http://loclahost:8080/api/v1/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "01901901901@mailinator.com",
    "password": "321dsa"
}'
```

## Project Architecture

Libraary used
<table>
  <thead>
    <tr>
      <th>Functionality</th>
      <th>Library</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Http Server</td>
      <td>labstack/echo</td>
    </tr>
    <tr>
      <td>Validation</td>
      <td>go-playground/validator</td>
    </tr>
    <tr>
      <td>ORM</td>
      <td>gorm.io/gorm</td>
    </tr>
    <tr>
      <td>Config</td>
      <td>spf13/viper</td>
    </tr>
    <tr>
      <td>CLI</td>
      <td>spf13/cobra</td>
    </tr>
  </tbody>
</table>

## API Doc

### Login
<table>
  <tbody>
    <tr>
      <td>URL</td>
      <td> <code>{{HOST}}/api/v1/login</code></td>
    </tr>
    <tr>
      <td>Method</td>
      <td><code>POST</code></td>
    </tr>
  </tbody>
</table>

```consose
curl --location '{{HOST}}/api/v1/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "01901901901@mailinator.com",
    "password": "321dsa"
}'
```
#### Response
Status code: 200 (Success)
```json
{
    "token_type": "Bearer",
    "expires_in": 5000,
    "access_token": "access_token",
    "refresh_token": "refresh_token"
}
```

Status code: 400 (Bad Request)
```json
{
  "message": "The user credentials were incorrect.",
  "type": "error",
  "code": 400
}
```

### Logout
<table>
  <tbody>
    <tr>
      <td>URL</td>
      <td> <code>{{HOST}}/api/v1/logout</code></td>
    </tr>
    <tr>
      <td>Method</td>
      <td><code>GET</code></td>
    </tr>
  </tbody>
</table>

```consose
curl --location '{{HOST}}/api/v1/logout' \
--header 'authorization: Bearer {{TOKEN}}'
```
#### Response
Status code: 200 (Success)
```json
{
  "message": "Successfully logged out",
  "type": "success",
  "code": 200
}
```

Status code: 401 (Unauthorized)
```json
{
  "message": "Unauthorized",
  "type": "error",
  "code": 401
}
```

### Create New Order
<table>
  <tbody>
    <tr>
      <td>URL</td>
      <td> <code>{{HOST}}/api/v1/orders</code></td>
    </tr>
    <tr>
      <td>Method</td>
      <td><code>POST</code></td>
    </tr>
  </tbody>
</table>

```consose
curl --location '{{HOST}}/api/v1/orders' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {{TOKEN}}' \
--data '{
    "store_id": 131172,
    "merchant_order_id": "AUEDUH1134498",
    "recipient_name": "Tarique",
    "recipient_phone": "01237161343",
    "recipient_address": "banani, gulshan 2, dhaka, bangladesh",
    "recipient_city": 1,
    "recipient_zone": 1,
    "recipient_area": 1,
    "delivery_type": 48,
    "item_type": 2,
    "special_instruction": "Leave At the door",
    "item_quantity": 1,
    "item_weight": 0.5,
    "amount_to_collect": 1300,
    "item_description": "Fragile"
}'
```
#### Response
Status code: 200 (Success)
```json
{
  "consignment_id": "5cecf86e-9a63-4391-9cd6-fdbd9206bbef",
  "merchant_order_id": "AUEDUH1134498",
  "order_status": "Pending",
  "delivery_fee": 60
}
```

Status code: 401 (Unauthorized)
```json
{
  "message": "Unauthorized",
  "type": "error",
  "code": 401
}
```

Status Code: 422 (Error)
```json
{
    "message": "Please fix the given errors",
    "type": "error",
    "code": 422,
    "errors": {
        "amount_to_collect": [
            "The amount to collect field is required."
        ],
        "delivery_type": [
            "The delivery type field is required."
        ],
        "item_quantity": [
            "The item quantity field is required."
        ],
        "item_type": [
            "The item type field is required."
        ],
        "item_weight": [
            "The item weight field is required."
        ],
        "recipient_address": [
            "The recipient address field is required."
        ],
        "recipient_area": [
            "The recipient area field is required."
        ],
        "recipient_city": [
            "The recipient city field is required."
        ],
        "recipient_name": [
            "The recipient name field is required."
        ],
        "recipient_phone": [
            "The recipient phone field is required."
        ],
        "recipient_zone": [
            "The recipient zone field is required."
        ],
        "store_id,omitempty": [
            "The store id,omitempty field is required."
        ]
    }
}
```

### List All Orders
<table>
  <tbody>
    <tr>
      <td>URL</td>
      <td> <code>{{HOST}}/api/v1/orders/all</code></td>
    </tr>
    <tr>
      <td>Method</td>
      <td><code>GET</code></td>
    </tr>
  </tbody>
</table>

```consose
curl --location '{{HOST}}/api/v1/orders/all?limit1&page1&transfer_status=1&archive=0' \
--header 'Authorization: Bearer {{TOKEN}}'
```
#### Response
Status code: 200 (Success blank list)
```json
{
  "message": "Orders successfully fetched.",
  "type": "success",
  "code": 200,
  "data": {
    "data": [],
    "total": 4,
    "current_page": 2,
    "per_page": 4,
    "total_in_page": 0,
    "last_page": 1
  }
}
```

Status code: 200 (Success With Data)
```json
{
  "code": 200,
  "data": {
    "data": [
      {
        "order_consignment_id": "2fdcb4dd-3f79-4d37-b1c5-da9a45a7a093",
        "order_created_at": "2024-11-30T20:20:58.51797+06:00",
        "order_description": "Fragile",
        "merchant_order_id": "AUEDUH1134498",
        "recipient_name": "Tarique",
        "recipient_address": "banani, gulshan 2, dhaka, bangladesh",
        "recipient_phone": "0**3",
        "order_amount": 1300,
        "total_fee": 73,
        "instruction": "Leave At the door",
        "order_type_id": 1,
        "cod_fee": 13,
        "promo_discount": 0,
        "discount": 0,
        "delivery_fee": 60,
        "order_status": "Pending",
        "order_type": "Delivery",
        "item_type": "Parcel",
        "transfer_status": 1,
        "archive": 0,
        "updated_at": "2024-11-30T20:20:58.517971+06:00",
        "created_by": 2,
        "updated_by": 2
      }
    ],
    "total": 1,
    "current_page": 1,
    "per_page": 12,
    "total_in_page": 1,
    "last_page": 1
  },
  "message": "Orders successfully fetched.",
  "type": "success"
}
```

### Cancel Order
<table>
  <tbody>
    <tr>
      <td>URL</td>
      <td> <code>{{HOST}}/api/v1/orders/{{consignment_id}}/cancel</code></td>
    </tr>
    <tr>
      <td>Method</td>
      <td><code>PUT</code></td>
    </tr>
  </tbody>
</table>

```consose
curl --location --request PUT '{{HOST}}/api/v1/orders/{{CONSIGNMENT_ID}}/cancel' \
--header 'Authorization: Bearer {{TOKEN}}'
```
#### Response
Status code: 200 (Success)
```json
{
  "message": "Order Cancelled Successfully",
  "type": "success",
  "code": 200
}
```

Status code: 400 
```json
{
  "message": "Please contact cx to cancel order",
  "type": "error",
  "code": 400
}
```
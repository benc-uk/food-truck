# 'Food Trucks' API Documentation

## OpenAPI Spec

[Link to full Swagger/OpenAPI spec here](./spec.yaml)

## Operations

The API is RESTful and supports the following operations

| Method | Path                                     | Description                                                      | Input | Returns          |
| ------ | ---------------------------------------- | ---------------------------------------------------------------- | ----- | ---------------- |
| GET    | /api/trucks/{lat}/{long}                 | Get at least 5 trucks, close to the provided location            | None  | Array of _Truck_ |
| GET    | /api/trucks/{lat}/{long}?radius={radius} | Get trucks within a given radius in meters, may return no trucks | None  | Array of _Truck_ |

### Status Codes

- **200** - Success
- **400** - Invalid request, e.g. lat and long are not valid
- **500** - Something extra bad happened

JSON error object will be returned on 400 and 500 errors

```json
{
  "message": "Description of the error",
  "component": "Component returning the error",
  "code": "Internal error code"
}
```

## Entity Schema

The API returns and accepts 'trucks' with the following entity structure

```ts
Truck {
  id:            string  // An id of this truck
  address:       string  // Address
  description:   string  // Description of the food served
  name:          string  // Name of the truck or vendor
  lat:           float   // Location latitude
  long:          float   // Location longitude
}
```

## Security

The API is currently unsecured

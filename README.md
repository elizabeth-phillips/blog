# Blog Backend

To run the post API: `go run cmd/postAPI/main.go`

## Routes

| Route        | RESTful method | Example Curl Request | Description           |
|:------------:|:-----------:|:---------------:|:---------------------:|
| `/posts/` | `POST` | `curl --location --request POST 'http://localhost:3001/api/posts/' --header 'Content-Type: application/json' --data-raw '{"title": "My first post", "creator": "Elizabeth", "content": "Here is the content"}'` | Creates a new post    |
| `/posts/` | `GET` | `curl --location --request GET 'http://localhost:3001/api/posts/'` | Gets all of the posts |
| `/posts/:id:` | `GET` | `curl --location --request GET 'http://localhost:3001/api/posts/1'` | Gets one specified post |
| `/posts/:id:` | `PUT` | `curl --location --request POST 'http://localhost:3001/api/posts/' --header 'Content-Type: application/json' --data-raw '{"title": "My first post", "creator": "Elizabeth", "content": "Here is the content"` | Updates the specified post by ID |
| `/posts/:id:` | `DELETE` | `curl --location --request DELETE 'http://localhost:3001/api/posts/1'` | Removes the specified post by ID |

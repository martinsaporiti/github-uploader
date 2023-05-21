## Flow

![image](<./Github%20Authentication%20Flow%402x%20(1).png>)

## How to run this app:

1. Register Application on Github
2. Create CLIENT_ID and CLIENT_SECRET env variables:

```bash
export CLIENT_ID=...
export CLIENT_SECRET=...
```

3. Run `go run main.go`
4. Visit `http://localhost:9000/`

## Refs:

- https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps
- https://clavinjune.dev/en/blogs/serving-embedded-static-file-inside-subdirectory-using-go/

# Iris Backend Template

Handy template to start developing using Iris-framework. This template uses,
among other libraries:

- [Iris go](https://www.iris-go.com)
- [GORM](https://gorm.io)
- [jwt-go](https://github.com/dgrijalva/jwt-go)
- [uuid](https://github.com/google/uuid)
- [godotenv](https://github.com/joho/godotenv)

## Code Structure

```bash
# Directories guide:

| root                # app kicker and routes listing
|---- controllers     # api controllers
|   |---- types       # static definitions for body parsing
|---- middlewares     # middlewares: jwt and user-role filter
|---- models          # GORM models
    |---- user        # static types for user model attributes, such as enum types
|---- utils           # helpers
    |---- context     # Customized context. Extends for quick reply and user extraction
    |---- misc        # various helpers

```

Hope this helps someone :)


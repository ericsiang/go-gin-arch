.
├── README.md
├── api
│   └── v1
│       ├── admins.go
│       └── users.go
├── asset
│   └── markdown
│       ├── md_source
│       │   ├── swagger_1.png
│       │   ├── swagger_2.png
│       │   ├── swagger_3.png
│       │   ├── zap_1.png
│       │   ├── zap_2.png
│       │   └── zap_3.png
│       ├── swagger.md
│       └── zap.md
├── common
│   ├── common_const
│   │   └── common_const.go
│   └── common_msg_id
│       └── common_msg_id.go
├── database
│   ├── migrate
│   │   └── migrate.go
│   └── seeder
│       ├── common_seeder.go
│       └── seeder.go
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── env.yaml
├── go.mod
├── go.sum
├── handler
│   ├── handleError.go
│   └── handleValidate.go
├── initialize
│   ├── config.go
│   ├── env.go
│   ├── initialize.go
│   ├── logger.go
│   ├── mysql.go
│   ├── redis.go
│   └── validateLang.go
├── log
│   ├── error
│   │   ├── error_2023-12-22.log
│   │   └── error_2024-01-02.log
│   └── info
│       ├── info_2023-12-22.log
│       └── info_2024-01-02.log
├── main.go
├── middleware
│   └── jwt_auth.go
├── model
│   ├── gormModel.go
│   ├── model_setting.go
│   └── users.go
├── router
│   └── router.go
├── test
│   ├── bcryptEncap_test
│   │   └── bcryptEncap_test.go
│   ├── jwt_secret_test
│   │   └── jwt_secret_test.go
│   ├── track_time_test
│   │   └── track_time_test.go
│   └── zapLogger_test
│       └── zapLoggger_test.go
├── tree.md
└── util
    ├── bcryptEncap
    │   └── bcrypt.go
    ├── gin_response
    │   └── gin_response.go
    ├── jwt_secret
    │   └── jwt_secret.go
    ├── mysql_manager
    │   └── mysql_err_code.go
    ├── swagger_docs
    │   └── swag_params.go
    ├── track_time
    │   └── track_time.go
    └── zap_logger
        └── zapLogger.go

34 directories, 53 files

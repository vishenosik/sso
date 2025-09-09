# Changelog

All notable changes to this project will be documented in this file.

## [unreleased]

### 🚀 Features

- *(auth)* Make proto-file + taskfile. generate go files.
- Add README
- Update module file
- *(gen code)* Added CHANGELOG, added validate tags
- *(requests_validation)* Added validator pkg/ added registration tests/ UPD Taskfile - added test.
- *(proto)* Add authorization service with Role method. Refactor Authentication service - IDs are now string.
- *(authorization)* Added new method Roles.
- *(env, task)* Added env file. modified taskfile to create dirs (todo).
- *(scripts,taskfile)* Added script with opt (redo to args).
- *(scripts)* New scripts for creating new services. ToDo: add template parsing.
- *(version)* Added version pkg. Will be used to verify requests from microservice client.
- *(grpc)* Added versioning to grpc api.
- *(REST api)* Added server settings for auth functionality.
- *(api)* Add api builder with tests.
- *(errors)* Add errors package to handle incoming errors from service.
- *(tests)* Add test coverage functionality. Add test tasks to taskfile. Add .gitignore with standart ignores + cover files.
- *(Obsidian)* Add obsidian && docs (README).
- *(docs)* Add new docs files.
- *(project)* Shift sso project into single repositiry.
- *(tests)* Add isAdmin & login services tests.
- *(tests)* Add token validation & new fail test.
- *(testing)* Finish covering login method with tests.
- *(codestyle)* Add better service functions err handling.
- *(test)* Add register tests to 95% coverage.
- *(task)* Add swagger task.
- *(swagger)* Client generation experiments.
- *(testing)* Add testing script to exclude untestable directories like mocks etc.
- *(containers)* Add docker-compose files for redis & grafana stack.
- *(paths)* Redo paths
- *(tools)* Add tools dir for postman & docker containers so far.
- *(tools)* Add debug settings into tools dir.
- *(kanban)* Add obsidian driven kanban desk.
- *(docs)* Add codeowners.
- *(context)* Add server context lib.
- *(AI)* Trying to implement ai in my workflow.
- *(logger)* Add pretty logger. Need to configure later.
- *(mutex)* Add mutex trials.
- *(learn)* Add mars game to study mutex.
- *(logger)* Trying to make dev logger work with yaml interpretation.
- *(ai docs)* Add ai generated swagger docs using tabnine.
- *(iter)* Add iterators experiments.
- *(logger)* Add log coloring trials.
- *(iter)* Experiments with iterators. So cool!
- *(iter)* Made more experimental stuff. Need to learn about memory usage in both usecases in func collections.Unique.
- *(testing)* Add benchmarks for unique function.
- *(slices)* Finished benchmarks tests. Left fastest.
- *(pprof)* Experimenting with profiling.
- *(redis)* Add redis cache lib.
- *(cache)* Add cache to stores. Configure redis to handle cache.
- *(cache)* Add Redis config. Start configuring code.
- *(cache)* Configure noop cache. Split app init.
- *(regex\strings\logger)* Add regex precompiled vars. Add strings pkg to move it to toolkit package. Make optional functions for dev logger init.
- *(logger coloring)* Setup coloring to dev logger Handler. Setup regex to search for numbers & keywords.
- *(logger)* Replace std log package to io.Writer.
- *(Auth REST)* Add isAdmin handler into REST server.
- *(collections)* Add iter filters func to test performance of slices through iterators.
- *(collections)* Add more tests to slices filer.
- *(collections)* Add more tests and wrap filter func.
- *(kafka)* Add some trials with kafka.
- *(hashicorp)* Add hashicorp docker-compose & cmd implementations.
- *(kafka/grafana)* Add goroutine producer implementation. Refactor grafana taskfile.
- *(logger)* Add std logger to implement goose logger based on slog.
- *(validator)* Add standalone validator package based on playground validator v10.
- *(sqlx)* Add sqlx pkg trials.
- *(helpers)* Add new fail functions. Add errorsMap helper.
- *(tools)* Add dgraph docker-compose container. Add testing main.go to cmd.
- *(service schema)* Add drawio schemas.
- *(dgraph client)* Start adding dgraph client & config.
- *(init)* Start adding flags work.
- *(rotator)* Add db rotator trial.
- *(build)* Dockerize build & run app operations.
- *(compose)* Configure dgraph & redis containers inside sso compose.
- *(dgraph)* Add dgraph connection into app. Add dgraph configs.
- *(app)* Add debug log to app to watch config during load.
- *(context)* Start adding context to app layer.
- *(git)* Add bin/ to gitkeep directive.
- *(context)* Add app context to rest app init.
- *(context)* Add universal functions to retrieve some values from context.
- *(pkg context)* Add context package - wrapper on std context.
- *(context)* Add request context to pkg.
- *(server context)* Add server context. Setup clean main func.
- *(namings)* Rearrange mod naming.
- *(request context)* Add request context value to gRPC api.
- *(deploy)* Add hmac secret generation to taskfile. Rebind compose volumes accordingly. Update Readme.
- *(dgraph)* Apply dgraph user store implementation.
- *(dgraph)* Add schema migration to dgraph.
- *(versions)* Add migrations versions parser, sorter, filter. Add tests.
- *(migrate)* Add migrator functionality. Add test data. Add test on collect.
- *(migrate fs)* Add migrate schema migrations. Split fs from other functionality. Add iterators to fetching functions.
- *(fs)* Add read up migration func.
- *(migrate versions)* Add collect func. Fix namings. Clear unused funcs.
- *(deploy)* Add hmac secret generation to taskfile. Rebind compose volumes accordingly. Update Readme.
- *(dgraph)* Apply dgraph user store implementation.
- *(dgraph)* Add schema migration to dgraph.
- *(versions)* Add migrations versions parser, sorter, filter. Add tests.
- *(migrate)* Add migrator functionality. Add test data. Add test on collect.
- *(migrate fs)* Add migrate schema migrations. Split fs from other functionality. Add iterators to fetching functions.
- *(fs)* Add read up migration func.
- *(migrate versions)* Add collect func. Fix namings. Clear unused funcs.
- *(sdk)* Add sdk package to manage api separately.
- *(app)* Add build flags implementation.
- *(flags)* Refactor gocherry flags registration.

### 🐛 Bug Fixes

- *(changelog)* Update CHANGELOG.
- *(scripts)* Removed useless autogenerated files.
- *(scripts)* Delete prepare script & "new" task. Remount dir for gen files.
- *(test)* Exclude auth_register from testing.
- *(isAdmin)* Redo isAdmin service method, added error handling, updated tests to 100% coverage.
- *(tests/workflow)* Update deps in tests. Update test func in workflow.
- *(test)* Correct service tests errors.
- *(develop convenience)* Add some generic functions to increase developing convenience.
- *(.env)* Remove .env from index. Move it's contents to Taskfile.
- *(taskfile)* Reimport tasks from sdk.
- *(logger)* Fix inconsistent log output.
- *(namings)* Rename combined functions.
- *(cache)* Remove all fmt logs from cache.
- *(highlight)* Remove fmt logs from highlight pkg.
- *(isAdmin)* Remake error handling in isAdmin gRPC handler.
- *(cache)* Code reorganization.
- *(validator, services)* Add nil uuid handling. Redo login service logging attrs.
- *(mock gen)* Fix a bug when mockery generated mocks for compiled functions.
- *(dgraph)* Add gitkeep to keep /data alive. Rearrange volumes inside docker-compose. Import dgraph Taskfile to root Taskfile.
- *(workflows)* Rename branch name inside changelog workflow.
- *(workflows)* Rename branch name on push step inside changelog workflow.
- *(app config)* Fix panic on Vault fields.
- *(app config)* Fix panic on Vault fields.
- *(compose)* Fix networks inside docker-compose.
- *(service)* Handle errors with changed package namings.

### 🚜 Refactor

- *(grpc exchange validator)* Сделал нормальное тестирование для валидатора. Обновил интерфейс для валидации.
- *(validation)* Update validation error message.
- Перенос папок с proto и gen. Обновление файла
- *(git/trash)* Cleanup gitignore & unused files.
- *(gen)* Remount path to generated files.
- *(namings)* Rename Config interface methods into more reliable ones.
- *(store)* Split authentication store package into users & apps. Fix imports.
- *(jwt)* Deleted error handling in jwt lib. Tests showed that error is never returned from SignedString method when passing []byte param.
- *(operations)* Redo build function.
- *(app init)* Redo rest app init.
- *(git)* .gitignore delete spaces.
- *(sdk)* Separate some packages into sdk repo.
- *(sdk)* Remove sdk from index.
- *(sdk)* Made naming refactoring.
- *(sdk)* Rebase sdk path's.
- *(taskfile/env)* Style taskfile
- *(testing)* Redo test paths.
- *(readme)* Add Readme lint styles.
- *(generated files)* Move ./gen into ./internal. Delete trash files.
- *(docs)* Remove docs into their own directory.
- *(logger)* Rename logger fields.
- *(gen)* Remove authorization from authentication code.
- *(attrs)* Remaster attrs over the code.
- *(dev logger)* Add buffer, mutex, coloring modes to dev logger.
- *(fail results)* Redo fail result return
- *(service auth)* Make methods work as pre-compiled functions.
- *(errors handle)* Add helper map to models/errors. Rename some errors for better understanding.
- *(lib pkg)* Move some lib components to pkg.
- *(config retrieval)* Refactor config fetching from multiple sources.
- *(config)* Refactor configs in service.
- *(compose)* Rebind env variables through env_file directive. Rearrange taskfile dependencies.
- *(app context)* Make better app context. Make config load & logger setup through context.
- *(app stores)* Split stores init from main app init function.
- *(collections)* Make Filter func exported. Fix namings accordingly.
- *(collections)* Make Filter func exported. Fix namings accordingly.
- *(grpc)* Implement gocherry grpc wrapper interface.
- *(pkg)* Replace internal pkg packages with gocherry/pkg packages.
- *(app)* Move app init to main.
- *(stores)* Move packages

### 📚 Documentation

- *(swagger)* Start testing swagger spec & codegen.
- *(api builder)* Add ApiV1 func docs.
- *(colors pkg)* Document colors package.
- *(swagger / env / scripts)* Rearrange swagger storage. Create script to generate root .env.
- *(contributing)* Update contributing document with git regulations.
- *(dgraph)* Update environment example.
- *(changelog)* Add changelog update workflow.
- *(env)* Add env config example.

### ⚡ Performance

- *(collections)* Remake Unique function for better readability.

### 🎨 Styling

- *(app)* Fix codestyle in app layer.

### 🧪 Testing

- *(test settings)* Improve testing by correcting testing script.
- *(sql)* Add test app migration.
- *(sql)* Add test app migration.

### ⚙️ Miscellaneous Tasks

- *(fail)* Make new fail function, that handles nil returns.
- *(delete)* Delete useless functionality.
- *(authentication gRPC api)* Refactor gRPC authentication. Add documentation.
- *(dgraph)* Inport dgraph Taskfile to  root Taskfile.
- *(dgraph cmd)* Refactor the way grpc client to dgraph is configured. Delete deprecated functions.
- *(deploy)* Move docker files to their dirs to clean project's root.
- *(lib)* Return lib dir to improve code logic.
- *(mod)* Update go dependencies to fix security vulnerabilities from github dependabot.
- *(namings)* Rename some functions to improve code understanding.
- *(namings)* Change test namings to beware name collision.
- *(app)* Refactor cache load. Replace method with func using context.
- *(app)* Make some functions prettier.
- *(cmd)* Clean cmd directory from experimental programs.
- *(tools)* Delete containers from tools. Moved them into another repo. Update taskfile accordingly.
- *(scripts)* Delete unused scripts.
- *(integration tests)* Delete unused integration tests.
- *(changelog workflow)* Add release changelog generation.
- *(mocks)* Delete mockery comments from auth api. Delete unused auth mocks. Change ping edp response.
- *(mod)* Upgrade dependencies.
- *(mod)* Update deps.
- *(gen)* Remove old generated files. Update deps.
- *(mod)* Upgrade deps.
- *(actions)* Upgrade github actions.
- *(dgraph)* Delete old dgraph implementation.

### ◀️ Revert

- *(dgraph migrator)* Remove dgraph migrator from repo. Moved to it's own project.
- *(dgraph migrator)* Remove dgraph migrator from repo. Moved to it's own project.

### Build

- *(deploy)* Rearrange namings, move deploy files into single directory.
- *(deploy)* Update ignores, update dockerfile. Update build scripts. Add platform build scripts.

<!-- generated by git-cliff -->

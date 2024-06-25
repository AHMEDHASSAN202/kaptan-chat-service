# SAMM Project General Open Challenge

This document outlines key points regarding the role.

## Tasks

- [ ] **Auto Create Module with Command Line**  
  Command: `go run module_setup.go -new_module='ModuleName' -new_module_path='NewModulePath' -root_module='WholeProjectModuleRoot'`

- [ ] **Restrict Access to Repo Layer Between Modules**

- [ ] **Handle Migration Scripts**

- [ ] **Handle RBAC**  
  Refer to:  
  [Casbin Get Started](https://casbin.org/docs/get-started)  
  [Authorization with Casbin](https://klotzandrew.com/blog/authorization-with-casbin)

## Extra Points

- [ ] **Logging Capabilities**  
  Utilize logrus/zap/...

- [ ] **Add Local Cache Library**  
  Create abstracted method to be replaced with Redis server

- [ ] **Add Event Bus for Internal Communication**  
  Only inside the module

- [ ] **Add Message Queue for Whole Project**

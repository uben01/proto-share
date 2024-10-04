# TODO

- [ ] Create integration tests
- [ ] Handle dependency between modules
    - *option 1: Add `dependencies` field in `module.yml`*
    - *option 2: Read `include` in proto schema*
- [ ] Generate compatibility tests upon schema changes
- [ ] Make templates extendable
    - *use case: maven repository config in build.gradle*
- [ ] Make `module.yml` optional -> Modules defined in config
    - *dependencies would be clearer*
- [ ] Create yml schema validator for module and config
- [ ] Create module/config file generator
- [ ] Accept other formats then yml
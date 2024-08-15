# TODO

- [ ] Test the functionality
- [ ] Create documentation
- [ ] Handle dependency between modules
    - *option 1: Add `dependencies` field in `module.yml`*
    - *option 2: Read `include` in proto schema*
- [ ] Generate compatibility tests upon schema changes
- [ ] Support new languages
- [ ] Make templates extendable
    - *use case: maven repository config in build.gradle*
- [ ] Make `module.yml` optional -> Modules defined in config
    - *dependencies would be clearer*
- [ ] Add `-verbose` flag. Log info only in case of enabled
- [ ] Create yml schema validator for module and config
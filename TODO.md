# TODO

- All system without AI integration
  - [x] POST /auth/sign-up
  - [x] POST /auth/sign-in
  - [x] POST /seasons
  - [ ] POST /seasons/:id/contenders/ingest
  - [ ] POST /seasons/:id/matchmaking
  - [ ] PATCH /seasons/:id
  - [ ] POST /seasons/:id/match/:mid/result
  - [ ] PATCH /seasons/:id/contender/:cid
  - [ ] GET /seasons/:id/matchmaking

- [ ] Codebase improvements

- [ ] Idempotency, retries, resilience

- [ ] AI Integration

- [ ] Usecases with AI Integration
  - [ ] GET /seasons/:id/suggestions?status=pending
  - [ ] PATCH /seasons/:id/suggestions/:sid/accept
  - [ ] PATCH /seasons/:id/suggestions/:sid/rejected

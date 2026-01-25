# Changelog

## [1.4.0](https://github.com/RevoTale/lovely-eye/compare/v1.3.0...v1.4.0) (2026-01-25)


### Features

* add paging to the events count ([e6f77bc](https://github.com/RevoTale/lovely-eye/commit/e6f77bcc001d28f709917f166f69d5b5388061fa))
* show only predefined events in the dahsboard events count card ([af046e6](https://github.com/RevoTale/lovely-eye/commit/af046e6d1872e8ed61823c957dc1ac7199581ea8))
* upgrade and refactor code to use `github.com/oschwald/maxminddb-golang` v2 ([d282207](https://github.com/RevoTale/lovely-eye/commit/d282207875a15c5290973aa3747e6219a38ce79b))


### Bug Fixes

* a bug caused daily chart showing alway jan 1 1970 ([fc5f9d5](https://github.com/RevoTale/lovely-eye/commit/fc5f9d59967c7a310155621045bda1c6a7f67ce1))
* better design priciple ([079c46f](https://github.com/RevoTale/lovely-eye/commit/079c46fe4ab9f7d5c515bbee3536b8362113211a))
* **deps:** update module modernc.org/sqlite to v1.44.3 ([2ac493e](https://github.com/RevoTale/lovely-eye/commit/2ac493ecda397d5e1b882fbbdcd0c39a5e7afdcf))
* **deps:** update module modernc.org/sqlite to v1.44.3 ([9425a78](https://github.com/RevoTale/lovely-eye/commit/9425a78368a6b5420de4de3573d101d5f01b6230))
* performance problem for large chart loading dataset ([41b35ba](https://github.com/RevoTale/lovely-eye/commit/41b35ba6d37cb62151f53016724047fa871ef8a7))
* recent event was not showing the event. Add indicater whther it is event or a pageview ([e4b7d09](https://github.com/RevoTale/lovely-eye/commit/e4b7d09611403af88bf54a1642d9c0c5abc049a6))
* remove useless badge ([dcd9175](https://github.com/RevoTale/lovely-eye/commit/dcd91753791983313540d078f2bd018f4d2318df))
* restore pageview event dedupe ([7abb88f](https://github.com/RevoTale/lovely-eye/commit/7abb88f353fbecf3088ab88de816620f950e3b81))
* unify events/pageviews tracking and schema ([7c70c9c](https://github.com/RevoTale/lovely-eye/commit/7c70c9c0afd0d2a480d9ba19f512a1d9aabb3e63))

## [1.3.0](https://github.com/RevoTale/lovely-eye/compare/v1.2.1...v1.3.0) (2026-01-22)


### Features

* add filtering by the event name and pathname ([3eb9536](https://github.com/RevoTale/lovely-eye/commit/3eb95362688ed438da5c46597e6448e7efc14709))
* enable fragment masking and strict codegen ([e223842](https://github.com/RevoTale/lovely-eye/commit/e223842a01a93dc5cb0ba6be668d99644d71be48))
* paging and paging limits for every graphql method ([944a3dd](https://github.com/RevoTale/lovely-eye/commit/944a3ddb1b497a5e547b1f34cf11d2a500fe8bad))


### Bug Fixes

* chart disappearing on date presets change ([b2557d2](https://github.com/RevoTale/lovely-eye/commit/b2557d25c36623227e740c783e4d724d09f698f1))
* eslint errors ([9bdebbf](https://github.com/RevoTale/lovely-eye/commit/9bdebbff66367332d550f57e13e99d1bf9f969ec))
* prevent evets card flickering ([5f54ea4](https://github.com/RevoTale/lovely-eye/commit/5f54ea45eb4209cee53140cbe55f4b52b453642b))
* prevent more flickering for the eventscard ([da0612e](https://github.com/RevoTale/lovely-eye/commit/da0612e0718f0f57fea41bf00fe6e69e79ca8a5d))
* remove react explicit import ([7b9685d](https://github.com/RevoTale/lovely-eye/commit/7b9685dcbe517802b638245f7a65aa3c52ba1903))
* since graphql paging changes update the dashboard graphql queries ([fcc29d3](https://github.com/RevoTale/lovely-eye/commit/fcc29d30f3be23b6bfa9abc851d18d148aef5192))
* unwanted access to the site creation by other users. ([46e3307](https://github.com/RevoTale/lovely-eye/commit/46e3307ec06bfdf4f640f4ce55fe7cfb40425460))

## [1.2.1](https://github.com/RevoTale/lovely-eye/compare/v1.2.0...v1.2.1) (2026-01-21)


### Bug Fixes

* by download stats I analyzed there is currently no users of this software, so I can do the silent breaking change with migrations to fix migrations issue the hard way. Part of https://github.com/RevoTale/lovely-eye/issues/32. ([9ce44c0](https://github.com/RevoTale/lovely-eye/commit/9ce44c0c3556776a8d623d79a6f95bbe666f54a9))

## [1.2.0](https://github.com/RevoTale/lovely-eye/compare/v1.1.4...v1.2.0) (2026-01-21)


### Features

* batch loading of the chart data to prevent downfall during the large amount traffic ([fb20c85](https://github.com/RevoTale/lovely-eye/commit/fb20c851ce4e9785a354f0ca76297fbe7ffaee0f))
* optimize the event loading, recuce bandwidth usage. ([298c794](https://github.com/RevoTale/lovely-eye/commit/298c794e5a50d318c3339ed6e6253e2a0870974a))

## [1.1.4](https://github.com/RevoTale/lovely-eye/compare/v1.1.3...v1.1.4) (2026-01-20)


### Bug Fixes

* bad group expression for postgres failed data collection ([ecf6734](https://github.com/RevoTale/lovely-eye/commit/ecf67342b89608ab9d8ebd4d280d7804d2411624))
* bad group expression for postgres failed data collection ([9e6a0c4](https://github.com/RevoTale/lovely-eye/commit/9e6a0c4f3f84f12e79e9c19c5c263561fd87f295))

## [1.1.3](https://github.com/RevoTale/lovely-eye/compare/v1.1.2...v1.1.3) (2026-01-20)


### Bug Fixes

* cord error due to misisng credentials header. deduplicate cors header ([02bc51b](https://github.com/RevoTale/lovely-eye/commit/02bc51b6e23bac06862ee06f305c661d271d4def))
* cors credentials header and polish documentation ([261e980](https://github.com/RevoTale/lovely-eye/commit/261e980856aecb4e2a9bf884519a9baa753798a1))
* markdown errors ([2547aef](https://github.com/RevoTale/lovely-eye/commit/2547aef5548b4fdcad63fd7bd4ca65afaa2cda9e))

## [1.1.2](https://github.com/RevoTale/lovely-eye/compare/v1.1.1...v1.1.2) (2026-01-20)


### Bug Fixes

* secure CORS policy and add security headers ([90a9261](https://github.com/RevoTale/lovely-eye/commit/90a92618d4ef9610dbac1bd4d0c3883db3710540))
* secure CORS policy and add security headers ([23e32e5](https://github.com/RevoTale/lovely-eye/commit/23e32e52086178532df3b10a7af7a10c0b00a546))

## [1.1.1](https://github.com/RevoTale/lovely-eye/compare/v1.1.0...v1.1.1) (2026-01-20)


### Bug Fixes

* revert: use `data-site-id` instead of `data-site-key` due to sematics ([7f28da4](https://github.com/RevoTale/lovely-eye/commit/7f28da43c635d5feef1a20bf14d9b5808bb55a1f))
* revert: use `data-site-id` instead of `data-site-key` due to sematics ([fcf8d9f](https://github.com/RevoTale/lovely-eye/commit/fcf8d9f4b3372a0de12da169a1f1d981fafd668e))

## [1.1.0](https://github.com/RevoTale/lovely-eye/compare/v1.0.0...v1.1.0) (2026-01-20)


### Features

* add db connection timeout ([4892f6e](https://github.com/RevoTale/lovely-eye/commit/4892f6e2af4a0581cba5d3397105821b99066e3f))
* use `data-site-id` instead of `data-site-key` ([3a1b7be](https://github.com/RevoTale/lovely-eye/commit/3a1b7bed34b4117d517dfa428597a508abca4b54))


### Bug Fixes

* **deps:** update all non-major dependencies ([1774f44](https://github.com/RevoTale/lovely-eye/commit/1774f44d135b4d08dc6a54995091bbecf1e77b2e))
* **deps:** update all non-major dependencies ([c8603b1](https://github.com/RevoTale/lovely-eye/commit/c8603b1bcfbbd561233c902cc3587ed89ea3f491))
* postgres mount path for v18.1 ([379c244](https://github.com/RevoTale/lovely-eye/commit/379c2444c395840a73bdf59e58d526c99a543b1c))
* renovate auto generate code on package change ([ed02be6](https://github.com/RevoTale/lovely-eye/commit/ed02be6cc72a8456eb9cd99a9775182824a436bf))
* renovate auto generate code on package change ([4f985c8](https://github.com/RevoTale/lovely-eye/commit/4f985c813cd25ceeae09ca5408456e6fea3c6f2b))

## 1.0.0 (2026-01-19)


### Features

* adapt CI to run reusable tests for the PR ([f2bed83](https://github.com/RevoTale/lovely-eye/commit/f2bed8347dd468729ac3d71dbc6ab1af82ec84d0))
* add a CI that checks whether no code forgotten to generate ([66072b3](https://github.com/RevoTale/lovely-eye/commit/66072b37e5b571ee181a31d761de0d62046b6a79))
* add a workarounf for base path serving ([cf114f0](https://github.com/RevoTale/lovely-eye/commit/cf114f007a2a9511f26b0b21f37a94e881494430))
* add ability to block IPs; ability to block IPs by countires. Refactor UI to reduce complexity of some components. ([e67bebf](https://github.com/RevoTale/lovely-eye/commit/e67bebf5c03a8659e4043059ad2121fbfd1fc6df))
* add ability to have insecure cookie for the localhost serving ([304d3d4](https://github.com/RevoTale/lovely-eye/commit/304d3d4b12bdb3af3304e735347e34afa5f361de))
* add bot filtering, daily visitor ID rotation, page view deduplication (10s window), improved IP extraction from proxied requests ([38affaa](https://github.com/RevoTale/lovely-eye/commit/38affaac86b5979f1996c2da1b0cf41a871eee03))
* add configurable dashboard path and consolidate tests ([a16a596](https://github.com/RevoTale/lovely-eye/commit/a16a5965888189c2f9aaf581b95a5432608add0d))
* add CORS headers ([2acbcd6](https://github.com/RevoTale/lovely-eye/commit/2acbcd659f2890c4ced44b937c63b795299a5d3e))
* add event definitions and allowlist validation. Interactive limits for event metadata. Add datetime range filtering UI and docs updates ([162820e](https://github.com/RevoTale/lovely-eye/commit/162820ed4d4394c4563139b95319d09b3afdbbca))
* add example usage for each defined event ([8629eb6](https://github.com/RevoTale/lovely-eye/commit/8629eb614d62b62fbfad0eea88246a90df5ae481))
* add filtering by referer, page path and device. Added the appropriate tests. Refactor old tests to use http coockies in auth ([887ee39](https://github.com/RevoTale/lovely-eye/commit/887ee394c5e97fad49fb034881c48538490a4be9))
* add per-site country tracking with GeoIP downloads. FIltering by country in the dashboard ([9764ca8](https://github.com/RevoTale/lovely-eye/commit/9764ca8b33df406ac2018ee643d107fac8105bf9))
* Add tests for the production setup of docker compose. Make data directory served from `data` in sqlite ([5df376c](https://github.com/RevoTale/lovely-eye/commit/5df376c550ea3fd7d1228428d25d338beaec9d0f))
* allow single site to have multipme domain ([e051769](https://github.com/RevoTale/lovely-eye/commit/e05176928103e8418d5eae863a3db2b7b35eac07))
* auto truncate site domain and validate input of the site name/domain ([b00a8e6](https://github.com/RevoTale/lovely-eye/commit/b00a8e6bf6b18ea7d554df86836791ce314b5900))
* **dashboard:** add server-side paging for top pages/referrers/devices/countries with totals. ([def1f3e](https://github.com/RevoTale/lovely-eye/commit/def1f3e85d4b882608f0e9368edd67eec8495c51))
* **events:** add GraphQL events API with string:string properties (metadata). Add tests for validation and GraphQL retrieval ([25f1e19](https://github.com/RevoTale/lovely-eye/commit/25f1e1944771292b4bd61aef20fdd15afcb82521))
* implement Atlas-based migrations with SQLite and PostgreSQL support. ([c666dcb](https://github.com/RevoTale/lovely-eye/commit/c666dcb5a6e81c2857fe1a364a05b17ba88c8de1))
* implemented basic dashboard with Bounce Rate, Total Visitors, chart, Top Referrers, Device Types etc. Fix bugs with the tracker script not collecting events. ([7d11ac9](https://github.com/RevoTale/lovely-eye/commit/7d11ac9e4eb582b634719867ce1c618a275bcd08))
* loggin level configuration ([42ab3e2](https://github.com/RevoTale/lovely-eye/commit/42ab3e244438f9632f8678a12dbc49c1a57b9667))
* make a strict eslint rules via `eslint-config-love` and fix all related errors ([9b8bb61](https://github.com/RevoTale/lovely-eye/commit/9b8bb6165d79ed102065453825570c38a154b01c))
* make testing of all database migrations on both sqlite ans postgresql before each release with the real-world dcker compose setup ([4dc40d5](https://github.com/RevoTale/lovely-eye/commit/4dc40d556b0f849c910fc35e5a9bbd3b5b7c2dcf))
* polish mobile support: handle date input overflow, tighten header on small screen ([30aa5da](https://github.com/RevoTale/lovely-eye/commit/30aa5da1ee84a2c46e0c2b83f0cc254e544e4ff7))
* prefill site name from domain and reorder fields so the user has to enter less ([6b0a234](https://github.com/RevoTale/lovely-eye/commit/6b0a23474fe8d9eeefccbb2153bb4bc24042cfd8))
* remove the hacky ways to manage migrations. Make a valid prebuilt environment for devcontainess to make both sqlite and postgresql migrations work out of box, be straightforward and no code mess ([286d7b8](https://github.com/RevoTale/lovely-eye/commit/286d7b8db31aaf4e123bf0a0f780d0e9d693d257))
* reuse dashboard list patterns and refactor datetime picker to avoid react related runtime page crash ([0a3e4d5](https://github.com/RevoTale/lovely-eye/commit/0a3e4d58b0a3a2d924b1f94084b233239cf72509))
* sunset light/dark mode color palette. Fix dark theme flickering. Adjust soem compoentns ([b4e3cb3](https://github.com/RevoTale/lovely-eye/commit/b4e3cb362d1be32d50f564628239366036461724))
* support devcontainer overrides for the customization of the per user dev environment ([737a6c9](https://github.com/RevoTale/lovely-eye/commit/737a6c99777661ffea4de2d47e736b192e1461c9))
* support multi-value dashboard filters and composable links. normalize country unknown filter ([d0c99c9](https://github.com/RevoTale/lovely-eye/commit/d0c99c9324602e9e90ee68257551538048e6ad84))
* testing the dashboard availablity in the production docker compose setup. ([7147ee9](https://github.com/RevoTale/lovely-eye/commit/7147ee9c1d7f67c8df11e03344f9f6fff2e06303))


### Bug Fixes

* `Scan error on column index 0, name "COALESCE(AVG(duration * 1.0), 0)": bun: can't scan 0 (int64) into float64, want nil` ([367dd13](https://github.com/RevoTale/lovely-eye/commit/367dd13df6005fceefaa67c9da2198becf64269f))
* add manual workflow dispatch for the release-please PRs ([c56b0f0](https://github.com/RevoTale/lovely-eye/commit/c56b0f073bb2d93df6ba1f321dcd83b7c382783f))
* add manual workflow dispatch for the release-please PRs ([429c920](https://github.com/RevoTale/lovely-eye/commit/429c92072ce1f4396fd6fe7287af52fc43822d17))
* auto run Ci tests for the release PR ([9335f61](https://github.com/RevoTale/lovely-eye/commit/9335f6199d7e0bc5f759c43a00bb2956dcce2ad7))
* auto run Ci tests for the release PR ([51348a3](https://github.com/RevoTale/lovely-eye/commit/51348a3838712081a6f3f0449d503631cd40bcdd))
* bug with db data types conversion ([b0231c3](https://github.com/RevoTale/lovely-eye/commit/b0231c3e178ed9d8606bac087c93faec1d21d24c))
* bump golangci/golangci-lint-action to support latest go ([fceb15b](https://github.com/RevoTale/lovely-eye/commit/fceb15b3621bd8456d67d0443c2380c0ca1b4c71))
* cast gloat value ([3568cdb](https://github.com/RevoTale/lovely-eye/commit/3568cdbc71ccd9f655da8f656d5123a3da5e04f7))
* CI tests not triggered on release action ([b0db8dc](https://github.com/RevoTale/lovely-eye/commit/b0db8dca257a451db87acca3d3cf006a8eeaac16))
* CI tests not triggered on release action ([f6e5095](https://github.com/RevoTale/lovely-eye/commit/f6e50953e8dd605fddb0f17414f53dd286d0fb3f))
* crtf token auth ([707e879](https://github.com/RevoTale/lovely-eye/commit/707e8795c5bd013d45f24b9454f5deac00540533))
* dashboard CI ([7693fea](https://github.com/RevoTale/lovely-eye/commit/7693fea52fc66a21657a735e3898958b16a4e97c))
* date filter layout on mobile screen ([4e8bedd](https://github.com/RevoTale/lovely-eye/commit/4e8beddd12fea24557545ac6658abeb6aeae8ad0))
* date range apply indicator syncing ([c5339ca](https://github.com/RevoTale/lovely-eye/commit/c5339ca86fc4179b9d065646f5429366f755e5ac))
* header layout on ultra small screen ([e93dedb](https://github.com/RevoTale/lovely-eye/commit/e93dedb5372852c6aac9ad226186be86b7dbe5b5))
* make a single for all tests passing ([e09918e](https://github.com/RevoTale/lovely-eye/commit/e09918e28828a900ce1013102937587d6ba029a3))
* more accurate readme ([2919108](https://github.com/RevoTale/lovely-eye/commit/2919108c5d8255652c06104f2dba880e61a7bd48))
* optimize the list of devcontain extensions installed ([5f16c3d](https://github.com/RevoTale/lovely-eye/commit/5f16c3d99c99a387621f23045e7f508b6b301895))
* redirection to the login page in case 401 error returned by one of api call ([65a3e76](https://github.com/RevoTale/lovely-eye/commit/65a3e7687ff08886c75aff4d9528310d68e01839))
* release-type and project structure mismatch ([8e1e31e](https://github.com/RevoTale/lovely-eye/commit/8e1e31ec79cea18332df185b5c69780e4332a226))
* release-type and project structure mismatch ([94e9e66](https://github.com/RevoTale/lovely-eye/commit/94e9e66764e7d808fd53e059251dd99fb17d6229))
* release-type and project structure mismatch ([86b6b4e](https://github.com/RevoTale/lovely-eye/commit/86b6b4ea93f0760a55cde5c25ed15555790c80df))
* release-type and project structure mismatch ([f68cbf0](https://github.com/RevoTale/lovely-eye/commit/f68cbf06a958ecbe4dc67380165917d60f7b00f4))
* remove csrf protection because brower features has replaced such protection according ot ther following thread https://www.reddit.com/r/node/comments/1im7yj0/jwt_csrf_a_good_security_practice/ ([004ff07](https://github.com/RevoTale/lovely-eye/commit/004ff077c57a8637449a581c8b5a1d0737ba5eb7))
* Scan error on column index 0, name "COALESCE(AVG(duration * 1.0), 0)": bun: can't scan 0 (int64) into float64, want nil ([a8246f8](https://github.com/RevoTale/lovely-eye/commit/a8246f881e3ffafc3e5e5781fb099db7976110b7))
* some go lint errors ([4f73c00](https://github.com/RevoTale/lovely-eye/commit/4f73c0008be5c87a22702536f396af9efd72ae4e))
* tests failing due to cors validation ([45236bf](https://github.com/RevoTale/lovely-eye/commit/45236bfe3f29a62f5371f1adf5cb059039873086))
* tracking script snippet ([e4ee8d8](https://github.com/RevoTale/lovely-eye/commit/e4ee8d863714144bd0d38fa4f8f37afc64d1c01b))
* unblock CI by installing golangci-lint ([d98351e](https://github.com/RevoTale/lovely-eye/commit/d98351efba5254ebaa3d7acc6c63916d9a0261ed))
* use fork to work around https://github.com/googleapis/release-pl… ([5c9d221](https://github.com/RevoTale/lovely-eye/commit/5c9d22161803ad251befd141422f428794c6deda))
* use fork to work around https://github.com/googleapis/release-please/issues/2265 ([dca17aa](https://github.com/RevoTale/lovely-eye/commit/dca17aaa167938aca87533b8f4dbd211828926f3))
* use PAT to fix CI run on publishing tag via release-please ([318b39c](https://github.com/RevoTale/lovely-eye/commit/318b39c6454eb0bf1e78af19ca7de6dbd29fe0b7))
* use PAT to fix CI run on publishing tag via release-please ([8738efa](https://github.com/RevoTale/lovely-eye/commit/8738efac825b233a3489ef2967b64b307067457d))
* validation of the public collection domain origin ([2c665af](https://github.com/RevoTale/lovely-eye/commit/2c665af3bdfcfe59ded19767373d09bfc893ec05))
* wrong build path after refactoring ([7341985](https://github.com/RevoTale/lovely-eye/commit/7341985887c91cddf66fe3c6f874fd17c177171b))
* wrong linter version in CI ([84e6a18](https://github.com/RevoTale/lovely-eye/commit/84e6a1886863ae5be120095cda9f0195658e8ab9))

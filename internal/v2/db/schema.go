package db

// ~ application 基础应用表，来自云端，本地如果有克隆，也往这个表插入数据
func (dm *dbModule) initApplicationSchema() error {
	// ~ 普通应用    app_id         = app_raw_id     // 比如 md5(testapp1)
	// ~            app_name       = app_raw_name   // 比如 testapp1
	// ~ 克隆应用    app_id         = md5(testapp1XXXXXX)
	// ~            app_name       = testapp1XXXXXX
	// ~            app_raw_id     = md5(testapp1)
	// ~            app_raw_name   = testapp1
	createApplicationSchema := `
	CREATE TABLE IF NOT EXISTS application (
		id                   BIGSERIAL     PRIMARY KEY,
		source_id            VARCHAR(50)   NOT NULL,
		app_id               VARCHAR(32)   NOT NULL,
		app_raw_id           VARCHAR(32)   NOT NULL,
		app_name             VARCHAR(100)  NOT NULL,
		app_raw_name         VARCHAR(100)  NOT NULL,
		app_version          VARCHAR(16)   NOT NULL,
		app_type             VARCHAR(32)   NOT NULL,  # app or middleware?
		app_entry            JSONB,
		app_image_analysis   JSONB,
		installed_type       VARCHAR(10)   NOT NULL,    # 这个字段存安装方式：完整安装，Server 端，Client 端
		is_cloned            BOOLEAN       NOT NULL DEFAULT false,
		created_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
		updated_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
		CONSTRAINT ck_application_app_entry_json
			CHECK (app_entry IS NULL OR jsonb_typeof(app_entry) = 'object'),
		CONSTRAINT ck_application_app_image_analysis_json
			CHECK (app_image_analysis IS NULL OR jsonb_typeof(app_image_analysis) = 'object')
	);`

	_, err := dm.db.Exec(createApplicationSchema)
	if err != nil {
		return err
	}

	indexApplicationSchema := `
	CREATE INDEX IF NOT EXISTS idx_application_app_id ON application (app_id);
	CREATE INDEX IF NOT EXISTS idx_application_app_raw_name ON application (app_raw_name);
	`

	_, err = dm.db.Exec(indexApplicationSchema)
	if err != nil {
		return err
	}

	return nil
}

// ~ user-application
// ~ 已渲染的 app：渲染完后更新插入这个表。克隆的应用也往这个表插入一条记录
func (dm *dbModule) initUserApplicationSchema() error {
	createUserApplicationSchema := `
	CREATE TABLE IF NOT EXISTS user_application (
		id                   BIGSERIAL     PRIMARY KEY,
		user_id              VARCHAR(120)  NOT NULL,       # 这个明天跟前端核对下，看看应该用多长
		application_id       BIGSERIAL,    # 表 application.id
		app_raw_data         JSONB,
		app_image_analysis   JSONB,
		is_upgrade           BOOLEAN       NOT NULL DEFAULT false,
		created_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
		updated_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
		CONSTRAINT ck_ua_app_raw_json
			CHECK (app_raw_data IS NULL OR jsonb_typeof(app_raw_data) = 'object'),
		CONSTRAINT ck_ua_app_image_analysis_json
			CHECK (app_image_analysis IS NULL OR jsonb_typeof(app_image_analysis) = 'object')
	);`

	_, err := dm.db.Exec(createUserApplicationSchema)
	if err != nil {
		return err
	}

	indexUserApplicationSchema := `
	CREATE INDEX IF NOT EXISTS idx_ua_app_id ON user_application (app_id);
	CREATE INDEX IF NOT EXISTS idx_ua_app_raw_id ON user_application (app_raw_id);
	CREATE INDEX IF NOT EXISTS idx_ua_app_name ON user_application (app_name);
	CREATE INDEX IF NOT EXISTS idx_ua_app_raw_name ON user_application (app_raw_name);
	`

	_, err = dm.db.Exec(indexUserApplicationSchema)
	if err != nil {
		return err
	}

	return nil
}

// ~ app-state
// ~ 已安装的应用表，靠 user_id + app_id 来查询
func (dm *dbModule) initApplicationStateSchema() error {
	createUserApplicationStateSchema := `
	CREATE TABLE IF NOT EXISTS user_application_state (
		id                   BIGSERIAL      PRIMARY KEY,
		user_application_id  BIGSERIAL,     # 表 user_application.id
		app_version          VARCHAR(16)    NOT NULL,   # app 升级，会先更新上面两个表，检查已安装的版本，告知前端；如果用户点击更新后，如果最后更新完后再更新这个字段到最新版本
		state                VARCHAR(64)    NOT NULL DEFAULT '',
		reason               VARCHAR(200)   NOT NULL DEFAULT '',
		message              VARCHAR(200)   NOT NULL DEFAULT '',
		progress             VARCHAR(10)    NOT NULL DEFAULT '',
		spec                 JSONB,
		status               JSONB,
		created_at           TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
		updated_at           TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
		CONSTRAINT ck_user_application_state_spec_object
			CHECK (spec IS NULL OR jsonb_typeof(spec) = 'object'),
		CONSTRAINT ck_user_application_state_status_object
			CHECK (status IS NULL OR jsonb_typeof(status) = 'object')
	);`

	_, err := dm.db.Exec(createUserApplicationStateSchema)
	if err != nil {
		return err
	}

	indexUserApplicationStateSchema := `
	CREATE INDEX IF NOT EXISTS idx_uas_user_id ON user_application_state (user_id);
	CREATE INDEX IF NOT EXISTS idx_uas_app_id ON user_application_state (app_id);
	`

	_, err = dm.db.Exec(indexUserApplicationStateSchema)
	if err != nil {
		return err
	}

	return nil
}

// ~ 存储市场源信息，比如 tag，topic；市场、upload、studio、cli 都会存到这个表
func (dm *dbModule) initSourceSchema() error {
	createMarketSourceSchema := `
	CREATE TABLE IF NOT EXISTS market_source (
		id             BIGSERIAL     PRIMARY KEY,
		source_id      VARCHAR(50)   NOT NULL,
		source_title   VARCHAR(50)   NOT NULL,
		source_url     TEXT          NOT NULL,
		source_type    VARCHAR(16)   NOT NULL CHECK (source_type IN ('local', 'remote')),
		description    TEXT          NOT NULL DEFAULT '',
		priority       INTEGER       NOT NULL DEFAULT 100 CHECK (priority >= 0),  # 展示顺序
		others         JSONB,
		created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
		updated_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
		CONSTRAINT ck_ms_others_object
			CHECK (others IS NULL OR jsonb_typeof(others) = 'object')
	);`

	_, err := dm.db.Exec(createMarketSourceSchema)
	if err != nil {
		return err
	}

	indexMarketSourceSchema := `
	CREATE INDEX IF NOT EXISTS idx_ms_source_id ON market_source (source_id);
	`

	_, err = dm.db.Exec(indexMarketSourceSchema)
	if err != nil {
		return err
	}

	return nil
}
